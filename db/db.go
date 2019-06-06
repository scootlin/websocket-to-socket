package db

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var db *gorm.DB

// InitDB initializes a database connection
func InitDB() {
	var err error
	db, err = gorm.Open("postgres", "")
	db.LogMode(true)
	db.SingularTable(true)
	if err != nil {
		log.Panic(err)
	}
}

func Get() *gorm.DB {
	return db
}
