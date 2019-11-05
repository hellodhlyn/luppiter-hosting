package main

import (
	"os"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func init() {
	newDB, err := gorm.Open("postgres", "host="+os.Getenv("DB_HOST")+" port="+os.Getenv("DB_PORT")+" user="+os.Getenv("DB_USERNAME")+" password="+os.Getenv("DB_PASSWORD")+" dbname="+os.Getenv("DB_NAME")+" sslmode=disable")
	if err != nil {
		panic(err)
	}

	db = newDB
}
