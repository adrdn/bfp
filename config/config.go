package config

import (
	"os"
	"database/sql"
	"encoding/json"
)

// ConfigFile specifies the JSON file
const ConfigFile = "config/config.json"

// Config is about fetching the path from the JSON file
type Config struct {
	DbDriver 	string 	`json:"dbDriver"`
	DbUser		string	`json:"dbUser"`
	DbPass		string	`json:"dbPass"`
	DbName		string	`json:"dbName"`
}

// LoadConfig sets the configuration
func LoadConfig(path string) (Config, error) {
	var conf Config

	dir, err := os.Open(path)
	defer dir.Close()
	if err != nil {
		return conf, err
	}
	jsonParser := json.NewDecoder(dir)
	jsonParser.Decode(&conf)
	return conf, err
}

func getConfig() Config {
	conf, err := LoadConfig(ConfigFile)
	if err != nil {
		panic(err)
	}
	return conf
}

// DbConn is the Database Configuration
func DbConn() (db *sql.DB) {
	dbDriver	:= 	getConfig().DbDriver
	dbUser		:= 	getConfig().DbUser
	dbPass		:=	getConfig().DbPass
	dbName		:=	getConfig().DbName
	db, err := sql.Open(dbDriver, dbUser + ":" + dbPass + "@/" + dbName)
	if err != nil {
		panic(err)
	}
	return db
}