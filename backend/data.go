package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

type User struct {
	ID     uint   `gorm:"primary_key" json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Rights string `json:"rights"`
}

var db *gorm.DB

func initDB() {
	connectString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		Config.DB.User, Config.DB.Password, Config.DB.Host, Config.DB.Database)

	db, err := gorm.Open("mysql", connectString)
	if err != nil {
		log.Print(connectString)
		log.Print(err)
		log.Fatal("failed to connect database")
	}
	defer db.Close()
	log.Print("DB connection open\n")
	//db.LogMode(true)
	db.AutoMigrate(&User{})
}
