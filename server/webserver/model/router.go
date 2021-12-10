package model

import (
	"github.com/jinzhu/gorm"
)

type Router struct {
	Model

	Name    string  `json:"name" gorm:"column:name;not null;unique"`
	Host    string  `json:"host" gorm:"column:host;not null;"`
	Path    string  `json:"path" gorm:"column:path;not null"`
	//Proxy   datatypes.JSON `json:"proxy" gorm:"column:proxy;not null"`
	Remark  string `json:"remark" gorm:"column:remark"`

	Upstreams []*Upstream `json:"upstreamRef" gorm:"many2many:router_upstream;"`
}

func AddRouter(data map[string]interface{}) error {
	router := &Router{
		Name:  data["name"].(string),
		Host:  data["host"].(string),
		Path:  data["path"].(string),
	}

	err := db.Create(&router).Error
	if err != nil {
		return err
	}

	var upstreams []*Upstream
	temp := Upstream{}
	temp.Model.ID = data["upstreamRef"].(uint64)
	upstreams = append(upstreams, &temp)

	err = db.Debug().Model(&router).Association("Upstreams").Replace(upstreams).Error
	if err != nil {
		return err
	}

	return nil
}

func DeleteRouter(id uint64) error {
	router := Router{}
	router.Model.ID = id
	/* clear Associations with upstreams*/
	db.Model(&router).Association("Upstreams").Clear()

	err := db.Delete(&router).Error
	if err != nil {
		return err
	}
	return nil
}


func PutRouter(id uint64, data map[string]interface{}) error {
	router := Router{}
	router.Model.ID = id

	var upstreams []*Upstream
	temp := Upstream{}
	temp.Model.ID = data["upstreamRef"].(uint64)
	upstreams = append(upstreams, &temp)

	delete(data, "upstreamRef")

	err := db.Model(&router).Update(data).Error
	if err != nil {
		return err
	}

	err = db.Debug().Model(&router).Association("Upstreams").Replace(upstreams).Error
	if err != nil {
		return err
	}

	return nil
}

func GetRouter(id uint64) (*Router, error) {
	var router Router

	err := db.Preload("Upstreams").Where("id = ?", id).Find(&router).Error
	if err != nil {
		return &router, err
	}

	return &router, nil
}

func GetRouters(data map[string]interface{}) ([]*Router, int, error) {
	var routers []*Router
	name := data["name"].(string)
	page := data["page"].(int)
	pageSize := data["pagesize"].(int)

	var count int
	if page > 0 {
		offset := (page - 1) * pageSize
		if len(name) > 0 {
			name = "%" + name + "%"
			err := db.Preload("Upstreams").Where("name LIKE ?", name).Offset(offset).Limit(pageSize).Find(&routers).Count(&count).Error
			if err != nil {
				return nil, 0, err
			}
		} else {
			err := db.Preload("Upstreams").Offset(offset).Limit(pageSize).Find(&routers).Count(&count).Error
			if err != nil {
				return nil, 0, err
			}
		}
	} else {
		err := db.Preload("Upstreams").Find(&routers).Count(&count).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, 0, err
		}
	}

	return routers, count, nil
}
