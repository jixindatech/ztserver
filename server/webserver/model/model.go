package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
	"zt-server/settings"
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

func OpenDatabase(cfg *settings.DataBase) error {
	var err error
	if cfg.Type == "mysql" {
		db, err = gorm.Open(cfg.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Name))
		if err != nil {
			return err
		}
	}

	if len(cfg.TablePrefix) > 0 {
		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			return cfg.TablePrefix + defaultTableName
		}
	}

	db.SingularTable(true)
	db.LogMode(false)

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.SingularTable(true)
	/*
		db.LogMode(true)
		db.DB().SetMaxIdleConns(10)
		db.DB().SetMaxOpenConns(100)
	*/
	db.AutoMigrate(User{})
	db.AutoMigrate(Resource{})
	db.AutoMigrate(Email{})
	db.AutoMigrate(SSL{})
	db.AutoMigrate(Upstream{})
	db.AutoMigrate(Router{})

	return nil
}
