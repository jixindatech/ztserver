package service

import "zt-server/webserver/model"

type Resource struct {
	ID     uint64
	Name   string
	Server string
	Port   int
	Path   string
	Method string
	Remark string

	Page     int
	PageSize int
}

func (r *Resource) Save() error {
	data := make(map[string]interface{})
	data["name"] = r.Name
	data["server"] = r.Server
	data["remark"] = r.Remark

	if r.ID > 0 {
		err :=  model.PutResource(r.ID, data)
		if err != nil {
			return err
		}

		return SetupUser()
	}
	return model.AddResource(data)
}

func (u *Resource) Get() (*model.Resource, error) {
	data := make(map[string]interface{})
	data["id"] = u.ID

	return model.GetResource(data)
}

func (u *Resource) GetList() ([]*model.Resource, int, error) {
	data := make(map[string]interface{})
	data["name"] = u.Name
	data["page"] = u.Page
	data["pagesize"] = u.PageSize

	return model.GetResources(data)
}

func (u *Resource) Delete() error {
	return model.DeleteResource(u.ID)
}
