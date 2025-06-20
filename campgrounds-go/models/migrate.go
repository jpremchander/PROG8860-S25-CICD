package models

import (
	"yelpcamp-go/config"
	"log"
)

func AutoMigrate() {
	db := config.GetDB()
	
	err := db.AutoMigrate(&User{}, &Campground{}, &Review{}, &Image{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	
	log.Println("Database migration completed")
}
