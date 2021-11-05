package model

import "github.com/jinzhu/gorm"

type Email struct {
	Model

	Server   string `json:"server" gorm:"column:server;not null;unique"`
	Port     int    `json:"port" gorm:"column:port;not null"`
	Email    string `json:"email" gorm:"column:email;not null"`
	Password string `json:"password" gorm:"column:password"`

	Remark string `json:"remark" gorm:"column:remark"`
}

func GetEmail() (*Email, error) {
	var email Email
	err := db.First(&email).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &email, nil
}

func AddEmail(data map[string]interface{}) error {
	email := &Email{
		Server:   data["server"].(string),
		Port:     data["port"].(int),
		Email:    data["email"].(string),
		Password: data["password"].(string),
		Remark:   data["remark"].(string),
	}
	err := db.Create(&email).Error
	if err != nil {
		return err
	}

	return nil
}

func UpdateEmail(id uint64, data map[string]interface{}) error {
	email := Email{}
	email.Model.ID = id

	err := db.Model(&email).Update(data).Error
	if err != nil {
		return err
	}

	return nil
}
