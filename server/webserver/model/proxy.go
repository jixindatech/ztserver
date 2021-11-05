package model

import (
	"github.com/jinzhu/gorm"
	"gorm.io/datatypes"
)

type Proxy struct {
	Model

	Server  string         `json:"server" gorm:"column:server;not null;unique"`
	Lb      string         `json:"lb" gorm:"column:lb;not null"`
	Backend datatypes.JSON `json:"backend" gorm:"column:backend;not null"`
	Remark  string         `json:"remark" gorm:"column:remark;"`
}

func AddProxy(data map[string]interface{}) error {
	proxy := &Proxy{
		Server:  data["server"].(string),
		Lb:      data["lb"].(string),
		Backend: data["backend"].([]byte),
		Remark:  data["remark"].(string),
	}
	err := db.Create(&proxy).Error
	if err != nil {
		return err
	}

	return nil
}

func PutProxy(id uint64, data map[string]interface{}) error {
	proxy := Proxy{}
	proxy.Model.ID = id

	err := db.Model(&proxy).Update(data).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteProxy(id uint64) error {
	proxy := Proxy{}
	proxy.Model.ID = id

	err := db.Delete(&proxy).Error
	if err != nil {
		return err
	}
	return nil
}

func GetProxy(id uint64) (*Proxy, error) {
	var proxy Proxy

	err := db.Where("id = ?", id).Find(&proxy).Error
	if err != nil {
		return &proxy, err
	}

	return &proxy, nil
}

func GetProxys(data map[string]interface{}) ([]*Proxy, int, error) {
	var proxys []*Proxy
	server := data["server"].(string)
	page := data["page"].(int)
	pageSize := data["pagesize"].(int)

	var count int
	if page > 0 {
		offset := (page - 1) * pageSize
		if len(server) > 0 {
			server = "%" + server + "%"
			err := db.Where("server LIKE ?", server).Offset(offset).Limit(pageSize).Find(&proxys).Count(&count).Error
			if err != nil {
				return nil, 0, err
			}
		} else {
			err := db.Offset(offset).Limit(pageSize).Find(&proxys).Count(&count).Error
			if err != nil {
				return nil, 0, err
			}
		}
	} else { // All of proxys
		err := db.Find(&proxys).Count(&count).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, 0, err
		}
	}

	return proxys, count, nil
}
