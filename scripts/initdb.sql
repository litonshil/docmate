-- Create the database
CREATE DATABASE IF NOT EXISTS `docmate`;

-- Create the user 'docmate_user' with the specified password
CREATE USER IF NOT EXISTS 'docmate_user'@'%' IDENTIFIED BY '12345678';

-- Grant all privileges to 'docmate_user' on the 'docmate' database
GRANT ALL PRIVILEGES ON docmate.* TO 'docmate_user'@'%';

-- Apply changes
FLUSH PRIVILEGES;



CREATE TABLE users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    type VARCHAR(50) NOT NULL,
    user_name VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL,
    phone VARCHAR(20),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    PRIMARY KEY (id),
    UNIQUE KEY ux_users_email (email),
    UNIQUE KEY ux_users_username (user_name),
    INDEX ix_users_type (type),
    INDEX ix_users_created_at (created_at),
    INDEX ix_users_deleted_at (deleted_at)
);
