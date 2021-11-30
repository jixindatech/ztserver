package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type User struct {
	Model

	Name      string `json:"name" gorm:"column:name;not null"`
	Phone     string `json:"phone" gorm:"column:phone"`
	Secret    string `json:"-" gorm:"column:secret"`
	Email     string `json:"email" gorm:"column:email;not null;unique"`
	Timestamp int64  `json:"-" gorm:"column:timestamp;"`
	Status    int    `json:"status" gorm:"column:status"`
	Remark    string `json:"remark" gorm:"column:remark"`

	Resources []Resource `gorm:"many2many:user_resource;"`
}

func AddUser(data map[string]interface{}) error {
	user := User{
		Name:   data["name"].(string),
		Phone:  data["phone"].(string),
		Email:  data["email"].(string),
		Status: data["status"].(int),
		Remark: data["remark"].(string),
	}

	err := db.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUser(id uint64) (*User, error) {
	var user User
	err := db.Where("id = ?", id).Find(&user).Error
	if err != nil {
		return &user, err
	}

	return &user, nil
}

func GetUsers(data map[string]interface{}) ([]*User, int, error) {
	var users []*User
	name := data["name"].(string)
	page := data["page"].(int)
	pageSize := data["pagesize"].(int)

	var count int
	if page > 0 {
		offset := (page - 1) * pageSize
		if len(name) > 0 {
			name = "%" + name + "%"
			err := db.Where("name LIKE ?", name).Offset(offset).Limit(pageSize).Find(&users).Count(&count).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, 0, err
			}
		} else {
			err := db.Offset(offset).Limit(pageSize).Find(&users).Count(&count).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, 0, err
			}
		}
	} else {
		err := db.Find(&users).Count(&count).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, 0, err
		}
	}

	return users, count, nil
}

func GetResourceByUser(id uint64) ([]*Resource, error) {
	var user User
	err := db.Where("id = ?", id).Find(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if user.ID > 0 {
		var resources []*Resource
		db.Model(&user).Association("Resources").Find(&resources)
		return resources, nil
	}

	return nil, fmt.Errorf("%s", "unknown error")
}

func GetResourceByUserEmail(email string) ([]*Resource, error) {
	var user User
	err := db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return nil, err
	}

	if user.ID > 0 {
		var resources []*Resource
		db.Model(&user).Association("Resources").Find(&resources)
		return resources, nil
	}

	return nil, fmt.Errorf("%s", "unknown error")
}

func PutUser(id uint64, data map[string]interface{}) error {
	user := User{}
	user.Model.ID = id

	err := db.Model(&user).Update(data).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(id uint64) error {
	user := User{}
	user.Model.ID = id
	/* delete Associations with resources*/
	db.Model(&user).Association("Resources").Clear()

	err := db.Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func SaveUserResource(id uint64, data map[string]interface{}) error {
	user := User{}
	user.Model.ID = id

	ids := data["ids"].([]uint64)
	var resources []*Resource
	for _, id := range ids {
		temp := Resource{}
		temp.Model.ID = id
		resources = append(resources, &temp)
	}

	err := db.Debug().Model(&user).Association("Resources").Replace(resources).Error
	if err != nil {
		return err
	}

	return nil
}

func GetUserResource(id uint64) ([]*Resource, error) {
	user := User{}
	user.Model.ID = id
	var resources []*Resource
	err := db.Model(&user).Association("Resources").Find(&resources).Error
	if err != nil {
		return nil, err
	}

	return resources, nil
}
