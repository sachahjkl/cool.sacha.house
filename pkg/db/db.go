package db

import (
	"log"
	"sachahjkl/htmx_go/pkg/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init(url string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Migrating db @ %v", url)
	db.AutoMigrate(&model.User{}, &model.Todo{})
	log.Printf("Finished migration")

	return db
}
