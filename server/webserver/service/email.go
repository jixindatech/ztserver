package service

import "zt-server/webserver/model"

type Email struct {
	ID       uint64
	Server   string
	Port     int
	Email    string
	Password string

	Remark string
}

func (e *Email) Get() (*model.Email, error) {
	return model.GetEmail()
}

func (e *Email) Save() error {
	data := make(map[string]interface{})
	data["server"] = e.Server
	data["port"] = e.Port
	data["email"] = e.Email
	data["password"] = e.Password
	data["remark"] = e.Remark

	if e.ID > 0 {
		return model.UpdateEmail(e.ID, data)
	}

	return model.AddEmail(data)
}
