package model

import "github.com/jinzhu/gorm"

type Cert struct {
	Model

	Server string `json:"server" gorm:"column:server;not null;unique"`
	Pub    string `json:"-" gorm:"type:varchar(5120);column:pub;not null;"`
	Pri    string `json:"-" gorm:"type:varchar(5120);column:pri;not null;"`
	Remark string `json:"remark" gorm:"column:remark;"`
}

func AddCert(data map[string]interface{}) error {
	cert := &Cert{
		Server: data["server"].(string),
		Pub:    data["pub"].(string),
		Pri:    data["pri"].(string),
		Remark: data["remark"].(string),
	}
	err := db.Create(&cert).Error
	if err != nil {
		return err
	}

	return nil
}

func GetCert(id uint64) (*Cert, error) {
	var cert Cert

	err := db.Where("id = ?", id).Find(&cert).Error
	if err != nil {
		return &cert, err
	}

	return &cert, nil
}

func GetCerts(data map[string]interface{}) ([]*Cert, int, error) {
	var certs []*Cert
	server := data["server"].(string)
	page := data["page"].(int)
	pageSize := data["pagesize"].(int)

	var count int
	if page > 0 {
		offset := (page - 1) * pageSize
		if len(server) > 0 {
			server = "%" + server + "%"
			err := db.Where("server LIKE ?", server).Offset(offset).Limit(pageSize).Find(&certs).Count(&count).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, 0, err
			}
		} else {
			err := db.Offset(offset).Limit(pageSize).Find(&certs).Count(&count).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, 0, err
			}
		}
	} else { // All of caches
		err := db.Find(&certs).Count(&count).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, 0, err
		}
	}

	return certs, count, nil
}

func PutCert(id uint64, data map[string]interface{}) error {
	cert := Cert{}
	cert.Model.ID = id

	err := db.Model(&cert).Update(data).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteCert(id uint64) error {
	cert := Cert{}
	cert.Model.ID = id

	err := db.Delete(&cert).Error
	if err != nil {
		return err
	}
	return nil
}
