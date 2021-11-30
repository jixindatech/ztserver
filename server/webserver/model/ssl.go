package model

import (
	"github.com/jinzhu/gorm"
	"gorm.io/datatypes"
)

type SSL struct {
	Model

	Name   string  `json:"name" gorm:"column:name;not null;unique"`
	Server datatypes.JSON `json:"server" gorm:"column:server;not null;unique"`
	Pub    string `json:"-" gorm:"type:varchar(5120);column:pub;not null;"`
	Pri    string `json:"-" gorm:"type:varchar(5120);column:pri;not null;"`
	Remark string `json:"remark" gorm:"column:remark;"`
}

func AddSSL(data map[string]interface{}) error {
	ssl := &SSL{
		Name:   data["name"].(string),
		Server: data["server"].([]byte),
		Pub:    data["pub"].(string),
		Pri:    data["pri"].(string),
		Remark: data["remark"].(string),
	}
	err := db.Create(&ssl).Error
	if err != nil {
		return err
	}

	return nil
}

func GetSSL(id uint64) (*SSL, error) {
	var ssl SSL

	err := db.Where("id = ?", id).Find(&ssl).Error
	if err != nil {
		return &ssl, err
	}

	return &ssl, nil
}

func GetSSLs(data map[string]interface{}) ([]*SSL, int, error) {
	var ssls []*SSL
	name := data["name"].(string)
	page := data["page"].(int)
	pageSize := data["pagesize"].(int)

	var count int
	if page > 0 {
		offset := (page - 1) * pageSize
		if len(name) > 0 {
			name = "%" + name + "%"
			err := db.Where("server LIKE ?", name).Offset(offset).Limit(pageSize).Find(&ssls).Count(&count).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, 0, err
			}
		} else {
			err := db.Offset(offset).Limit(pageSize).Find(&ssls).Count(&count).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, 0, err
			}
		}
	} else { // All of caches
		err := db.Find(&ssls).Count(&count).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, 0, err
		}
	}

	return ssls, count, nil
}

func PutSSL(id uint64, data map[string]interface{}) error {
	ssl := SSL{}
	ssl.Model.ID = id

	err := db.Model(&ssl).Update(data).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteSSL(id uint64) error {
	ssl := SSL{}
	ssl.Model.ID = id

	err := db.Delete(&ssl).Error
	if err != nil {
		return err
	}
	return nil
}
