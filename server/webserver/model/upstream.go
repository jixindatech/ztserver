package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gorm.io/datatypes"
)

type Upstream struct {
	Model

	Name  string         `json:"name" gorm:"column:name;not null;unique"`
	Lb      string         `json:"lb" gorm:"column:lb;not null"`
	Key     string         `json:"key" gorm:"column:key;not null"`
	Backend datatypes.JSON `json:"backend" gorm:"column:backend;not null"`
	Retry int       `json:"retry" gorm:"column:retry;not null"`
	TimeoutConnect int `json:"timeoutConnect" gorm:"column:timeout_connect;not null"`
	TimeoutSend int `json:"timeoutSend" gorm:"column:timeout_send;not null"`
	TimeoutReceive int `json:"timeoutReceive" gorm:"column:timeout_receive;not null"`

	Remark  string         `json:"remark" gorm:"column:remark;"`

	Routers []Router `gorm:"many2many:router_upstream;"`
}

func AddUpstream(data map[string]interface{}) error {
	upstream := &Upstream{
		Name:  data["name"].(string),
		Lb:      data["lb"].(string),
		Key: data["key"].(string),
		Backend: data["backend"].([]byte),
		Retry: data["retry"].(int),
		TimeoutConnect: data["timeoutConnect"].(int),
		TimeoutSend: data["timeoutSend"].(int),
		TimeoutReceive: data["timeoutReceive"].(int),
		Remark:  data["remark"].(string),

	}
	err := db.Create(&upstream).Error
	if err != nil {
		return err
	}

	return nil
}

func PutUpstream(id uint64, data map[string]interface{}) error {
	upstream := Upstream{}
	upstream.Model.ID = id

	err := db.Model(&upstream).Update(data).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteUpstream(id uint64) error {
	upstream := Upstream{}
	upstream.Model.ID = id

	count := db.Model(&upstream).Association("Routers").Count()
	if count > 0 {
		return fmt.Errorf("%s", "exist router link")
	}

	err := db.Delete(&upstream).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUpstream(id uint64) (*Upstream, error) {
	var upstream Upstream

	err := db.Where("id = ?", id).Find(&upstream).Error
	if err != nil {
		return &upstream, err
	}

	return &upstream, nil
}

func GetUpstreams(data map[string]interface{}) ([]*Upstream, int, error) {
	var upstreams []*Upstream
	name := data["name"].(string)
	page := data["page"].(int)
	pageSize := data["pagesize"].(int)

	var count int
	if page > 0 {
		offset := (page - 1) * pageSize
		if len(name) > 0 {
			name = "%" + name + "%"
			err := db.Where("name LIKE ?", name).Offset(offset).Limit(pageSize).Find(&upstreams).Count(&count).Error
			if err != nil {
				return nil, 0, err
			}
		} else {
			err := db.Offset(offset).Limit(pageSize).Find(&upstreams).Count(&count).Error
			if err != nil {
				return nil, 0, err
			}
		}
	} else { // All of upstreams
		err := db.Find(&upstreams).Count(&count).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, 0, err
		}
	}

	return upstreams, count, nil
}
