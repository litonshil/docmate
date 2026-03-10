package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"docmate/client/conn"
	"docmate/internal/model"

	"github.com/spf13/cobra"
	"gorm.io/gorm/clause"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Database seeder commands",
}

var seedMedicinesCmd = &cobra.Command{
	Use:   "medicines",
	Short: "Seed medicines from CSV file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return seedMedicines()
	},
}

func RegisterSeedCommand() {
	seedCmd.AddCommand(seedMedicinesCmd)
	RegisterSubCommand(seedCmd)
}

func seedMedicines() error {
	db := conn.Db()
	if db == nil {
		return fmt.Errorf("DB connection not initialized")
	}

	// Determine the CSV file path. Check variations of locations
	// for running locally vs running inside docker root.
	candidatePaths := []string{
		filepath.Join(".", "docs", "medicines.csv"),
		filepath.Join("..", "docs", "medicines.csv"),
		filepath.Join("/project", "docs", "medicines.csv"),
	}

	var csvFilePath string
	var file *os.File
	var err error

	for _, path := range candidatePaths {
		file, err = os.Open(path)
		if err == nil {
			csvFilePath = path
			fmt.Printf("Found CSV at %s\n", csvFilePath)

			break
		}
	}

	if file == nil {
		return fmt.Errorf("failed to open medicines.csv in any candidate paths")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Some descriptions might have unquoted newlines/HTML or other weird characters
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	// Read all records (this assumes the file fits in memory, which is ~19MB, perfectly fine)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV records: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file is empty or only contains headers")
	}

	// Skip header
	header := records[0]
	fmt.Printf("Parsed header: %v\n", header)

	var medicines []model.Medicine
	now := time.Now()

	fmt.Printf("Processing %d records...\n", len(records)-1)

	for i, row := range records[1:] {
		// brand_name,generic_name,type,form,strength,manufacturer,description
		if len(row) < 7 {
			fmt.Printf("Skipping row %d due to insufficient columns: %v\n", i+2, row)

			continue
		}

		brandName := strings.TrimSpace(row[0])
		genericName := strings.TrimSpace(row[1])
		// skip type row[2] as we don't have it in model
		formStr := strings.TrimSpace(strings.ToLower(row[3]))
		strength := strings.TrimSpace(row[4])
		manufacturer := strings.TrimSpace(row[5])
		description := strings.TrimSpace(row[6])

		// Map string to MedicineFormType
		var form model.MedicineFormType
		switch formStr {
		case "tablet":
			form = model.MedicineFormTablet
		case "capsule":
			form = model.MedicineFormCapsule
		case "syrup":
			form = model.MedicineFormSyrup
		case "suspension":
			form = model.MedicineFormSuspension
		case "injection":
			form = model.MedicineFormInjection
		case "inhaler":
			form = model.MedicineFormInhaler
		case "drops":
			form = model.MedicineFormDrops
		case "cream":
			form = model.MedicineFormCream
		case "ointment":
			form = model.MedicineFormOintment
		case "gel":
			form = model.MedicineFormGel
		case "patch":
			form = model.MedicineFormPatch
		case "suppository":
			form = model.MedicineFormSuppository
		case "powder":
			form = model.MedicineFormPowder
		case "sachet":
			form = model.MedicineFormSachet
		default:
			form = model.MedicineFormOther
		}

		med := model.Medicine{
			CreatedBy:    1, // System user or admin
			BrandName:    brandName,
			GenericName:  genericName,
			Form:         form,
			Strength:     strength,
			Manufacturer: manufacturer,
			Description:  description,
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		medicines = append(medicines, med)
	}

	fmt.Printf("Successfully mapped %d medicines. Starting database insertion in batches...\n", len(medicines))

	// Insert in batches of 1000
	batchSize := 1000

	// CreateInBatches using on conflict do nothing in case a duplicate is found if we re-run
	// Usually seeding requires handling conflicts. Assuming generic_name is not unique per se,
	// but combination of brand, generic, strength forms could be unique.
	// If no unique constraints exist on these, this might just duplicate if ran twice.
	result := db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(medicines, batchSize)
	if result.Error != nil {
		return fmt.Errorf("failed to bulk insert medicines: %w", result.Error)
	}

	fmt.Printf("Successfully seeded %d medicines.\n", result.RowsAffected)

	return nil
}
