package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

type Model struct {
	ID        uint64    `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updateAt"`
	// DeletedAt *time.Time `json:"deletedAt"`
}

type Database struct {
	DB string
}

const path = "db"

var db *gorm.DB

func OpenDatabase(database string) error {
	var err error
	database = path + "/" + database
	db, err = gorm.Open("sqlite3", database)
	if err != nil {
		return err
	}

	db.SingularTable(true)
	/*
		db.LogMode(true)
		db.DB().SetMaxIdleConns(10)
		db.DB().SetMaxOpenConns(100)
	*/
	db.AutoMigrate(User{})
	db.AutoMigrate(Resource{})
	db.AutoMigrate(Email{})
	db.AutoMigrate(Cert{})
	db.AutoMigrate(Proxy{})
	return nil
}
