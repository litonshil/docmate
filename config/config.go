package config

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"

	_ "github.com/spf13/viper/remote"
)

const (
	LOCAL = "LOCAL"
)

var c Config

type Config struct {
	App   *AppConfig
	Db    *DbClient
	Cache *CacheClient
	Queue *QueueClient
}

// AppConfig application specific config
type AppConfig struct {
	Name         string
	Port         int
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
	AppKeyHeader string
	AppKey       string
	ENV          string
	LogLevel     string
	JWTSecret    string
}

type ServiceConfig struct {
	ServiceURL   string
	AppKeyHeader string
	AppKey       string
	Timeout      time.Duration
}

type DbClient struct {
	Db *DBConfig
}

type CacheClient struct {
	Redis *RedisConfig
}

type QueueClient struct {
	Asynq *AsynqConfig
}

// DBConfig database specific config
type DBConfig struct {
	Host        string
	Port        int
	Name        string
	Username    string
	Password    string
	MaxLifeTime time.Duration
	MaxIdleConn int
	MaxOpenConn int
	Debug       bool
	PrepareStmt bool
}

type RedisConfig struct {
	Host            string
	Port            int
	Name            string
	Username        string
	Password        string
	ValueExpiredIn  int // seconds
	MaxIdleConn     int
	MaxOpenConn     int
	Database        int
	MandatoryPrefix string
}

type AsynqConfig struct {
	RedisAddr        string
	DB               int
	Password         string
	Concurrency      int
	Queue            string
	SyncTaskQueue    string
	Retention        time.Duration // in hours
	UniquenessTTL    time.Duration
	RetryCount       int
	TaskExecTimeUnit string
}

// Get returns all configurations
func Get() Config {
	return c
}

func App() *AppConfig {
	return c.App
}

func DB() *DbClient {
	return c.Db
}

func Cache() *CacheClient {
	return c.Cache
}

func Queue() *QueueClient {
	return c.Queue
}

// Load the config
func Load() error {
	setDefaultConfig()

	if err := bindEnvVariables(); err != nil {
		return err
	}

	consulURL := viper.GetString("CONSUL_URL")
	consulPath := viper.GetString("CONSUL_PATH")
	configFilePath := viper.GetString("CONFIG_FILE_PATH")

	// Load from file if CONFIG_FILE_PATH is provided
	if configFilePath != "" {
		return loadConfigFromFile(configFilePath)
	}

	// Load from Consul if CONSUL_URL and CONSUL_PATH are provided
	if consulURL != "" && consulPath != "" {
		return loadConfigFromConsul(consulURL, consulPath)
	}

	// Log missing configuration variables
	log.Println("CONFIG_FILE_PATH or CONSUL_URL and CONSUL_PATH are missing from ENV")
	return nil
}

// bindEnvVariables binds the environment variables to Viper
func bindEnvVariables() error {
	envVars := []string{"env", "CONSUL_URL", "CONSUL_PATH", "CONFIG_FILE_PATH"}
	for _, v := range envVars {
		if err := viper.BindEnv(v); err != nil {
			log.Printf("Error binding env variable %s: %v", v, err)
			return err
		}
	}
	return nil
}

// loadConfigFromFile loads config from a local file
func loadConfigFromFile(configFilePath string) error {
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Failed to read config from file: %v", err)
		return err
	}

	log.Printf("Loaded configuration from %s", configFilePath)
	return unmarshalConfig()
}

// loadConfigFromConsul loads config from Consul
func loadConfigFromConsul(consulURL, consulPath string) error {
	viper.SetConfigType("json")

	if err := viper.AddRemoteProvider("consul", consulURL, consulPath); err != nil {
		log.Printf("Failed to add remote provider: %v", err)
		return err
	}

	if err := viper.ReadRemoteConfig(); err != nil {
		log.Printf("Failed to read remote config from Consul: %v", err)
		return err
	}

	log.Printf("Loaded configuration from Consul URL: %v, Path: %v", consulURL, consulPath)
	return unmarshalConfig()
}

// unmarshalConfig unmarshal the configuration into the Config struct
func unmarshalConfig() error {
	c = Config{}
	if err := viper.Unmarshal(&c); err != nil {
		log.Printf("Unable to decode config: %v", err)
		return err
	}

	// Optionally, print the configuration for verification
	if r, err := json.MarshalIndent(&c, "", "  "); err == nil {
		fmt.Println(string(r))
	}

	return nil
}

// ReadDotENV reads the environment from the .env file
func ReadDotENV() string {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading .env file: %v", err)
		return ""
	}

	env := viper.GetString("ENV")
	return env
}

// setDefaultConfig sets the default configurations
func setDefaultConfig() {
	env := ReadDotENV()

	c.App = &AppConfig{
		Name:         "docmate",
		Port:         8080,
		ReadTimeout:  30,
		WriteTimeout: 30,
		IdleTimeout:  30,
		AppKeyHeader: "app-key",
		AppKey:       "appkey",
		ENV:          env,
	}

	c.Db = &DbClient{
		Db: &DBConfig{
			Host:        "127.0.0.1",
			Name:        "docmate",
			Port:        3307,
			Username:    "root",
			Password:    "12345678",
			MaxLifeTime: 30,
			MaxIdleConn: 1,
			MaxOpenConn: 2,
			Debug:       true,
			PrepareStmt: true,
		},
	}

	c.Cache = &CacheClient{
		Redis: &RedisConfig{
			Host:            "127.0.0.1",
			Port:            6379,
			Username:        "",
			Password:        "",
			Database:        4,
			ValueExpiredIn:  0,
			MandatoryPrefix: "docmate_",
		},
	}

	c.Queue = &QueueClient{
		Asynq: &AsynqConfig{
			RedisAddr:        "127.0.0.1:6379",
			DB:               15,
			Concurrency:      10,
			Queue:            "docmate",
			Retention:        168,
			RetryCount:       25,
			UniquenessTTL:    1,
			TaskExecTimeUnit: "SECOND",
		},
	}
}
