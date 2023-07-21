package persistance

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var conn_pool *sql.DB

type DatabaseConfig struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Database  string `json:"database"`
	Net       string `json:"net"`
	ParseTime bool   `json:"parseTime"`
}

// Get
func GetMySqlConnection() (*sql.DB, error) {
	// Capture connection properties.
	mySQLConfig, configErr := loadMySQLConfig()
	if configErr != nil {
		return nil, configErr
	}
	cfg := convertToMySQLConfig(mySQLConfig)

	// Get a database handle.
	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
		return nil, pingErr
	}
	log.Println("Connected to db!")
	return db, nil
}

func InitMySqlConnectionPool(maxConnections int) error {
	// Capture connection properties.
	mySQLConfig, configErr := loadMySQLConfig()
	if configErr != nil {
		return configErr
	}
	cfg := convertToMySQLConfig(mySQLConfig)
	// Get a database handle.
	var err error
	conn_pool, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}

	conn_pool.SetMaxOpenConns(maxConnections)
	conn_pool.SetMaxIdleConns(maxConnections / 2)

	pingErr := conn_pool.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
		return pingErr
	}
	fmt.Println("Connected!")
	return nil
}

func GetConnection() *sql.DB {
	return conn_pool
}

func CloseDB(db *sql.DB) {
	db.Close()
	log.Println("DB closed!")
}

func loadMySQLConfig() (DatabaseConfig, error) {
	file, err := os.Open("../resources/mysql_db_config.json")
	if err != nil {
		return DatabaseConfig{}, err
	}
	defer file.Close()

	var config DatabaseConfig
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return DatabaseConfig{}, err
	}

	return config, nil
}

func convertToMySQLConfig(config DatabaseConfig) mysql.Config {
	return mysql.Config{
		User:      config.User,
		Passwd:    config.Password,
		Net:       config.Net,
		Addr:      fmt.Sprintf("%s:%s", config.Host, config.Port),
		DBName:    config.Database,
		ParseTime: config.ParseTime,
	}
}
