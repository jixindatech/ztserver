package model

type Resource struct {
	Model

	Name   string `json:"name" gorm:"column:name;unique;not null"`
	Server string `json:"server" gorm:"column:server;unique;not null"`
	Port   int    `json:"port" gorm:"column:port"`
	Path   string `json:"path" gorm:"column:path"`
	Method string `json:"method" gorm:"column:method"`
	Remark string `json:"remark" gorm:"column:remark"`

	Users []User `gorm:"many2many:user_resource;"`
}

func AddResource(data map[string]interface{}) error {
	resource := Resource{
		Name:   data["name"].(string),
		Server: data["server"].(string),
		Remark: data["remark"].(string),
	}

	err := db.Create(&resource).Error
	if err != nil {
		return err
	}
	return nil
}

func GetResource(data map[string]interface{}) (*Resource, error) {
	var resource Resource
	id := data["id"].(uint64)

	err := db.Where("id = ?", id).Find(&resource).Error
	if err != nil {
		return &resource, err
	}

	return &resource, nil
}

func GetResources(data map[string]interface{}) ([]*Resource, int, error) {
	var resources []*Resource
	name := data["name"].(string)
	page := data["page"].(int)
	pageSize := data["pagesize"].(int)

	var count int
	offset := (page - 1) * pageSize
	if len(name) > 0 {
		name = "%" + name + "%"
		err := db.Where("name LIKE ?", name).Offset(offset).Limit(pageSize).Find(&resources).Count(&count).Error
		if err != nil {
			return nil, 0, err
		}
	} else {
		err := db.Offset(offset).Limit(pageSize).Find(&resources).Count(&count).Error
		if err != nil {
			return nil, 0, err
		}
	}

	return resources, count, nil
}

/*
func GetResourceByUser(email string) ([]*Resource, error) {
	var user User
	err := db.Where("email = ？", email).Find(&user).Error
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
*/
func PutResource(id uint64, data map[string]interface{}) error {
	resource := Resource{}
	resource.Model.ID = id

	err := db.Model(&resource).Update(data).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteResource(id uint64) error {
	resource := Resource{}
	resource.Model.ID = id
	/* delete Associations with users*/
	db.Model(&resource).Association("Users").Clear()

	err := db.Delete(&resource).Error
	if err != nil {
		return err
	}
	return nil
}
