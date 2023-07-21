package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type DatabaseInitConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Net      string `json:"net"`
}

func main() {
	// Capture connection properties.
	mySQLConfig, configErr := loadMySQLConfig()
	if configErr != nil {
		log.Fatal(configErr)
		return
	}
	cfg := convertToMySQLConfig(mySQLConfig)
	// Get a database handle.
	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the database
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS myAppDatabase")
	if err != nil {
		log.Fatal(err)
	}

	// Use the database
	_, err = db.Exec("USE myAppDatabase")
	if err != nil {
		log.Fatal(err)
	}

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS promotions (
			id int NOT NULL AUTO_INCREMENT,
			promotion_id varchar(255) NOT NULL,
			price decimal(5,2) NOT NULL,
			expiration_date datetime DEFAULT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY promotion_id_UNIQUE (promotion_id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database initialized successfully!")
}

func loadMySQLConfig() (DatabaseInitConfig, error) {
	file, err := os.Open("mysql_db_init_config.json")
	if err != nil {
		return DatabaseInitConfig{}, err
	}
	defer file.Close()

	var config DatabaseInitConfig
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return DatabaseInitConfig{}, err
	}

	return config, nil
}

func convertToMySQLConfig(config DatabaseInitConfig) mysql.Config {
	return mysql.Config{
		User:   config.User,
		Passwd: config.Password,
		Net:    config.Net,
		Addr:   fmt.Sprintf("%s:%s", config.Host, config.Port),
	}
}
