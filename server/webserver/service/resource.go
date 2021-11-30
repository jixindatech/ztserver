package service

import "zt-server/webserver/model"

type Resource struct {
	ID     uint64
	Name   string
	Host string
	Port   int
	Path   string
	Method []byte
	Remark string

	Page     int
	PageSize int
}

func (r *Resource) Save() error {
	data := make(map[string]interface{})
	data["name"] = r.Name
	data["host"] = r.Host
	data["path"] = r.Path
	data["method"] = r.Method
	data["remark"] = r.Remark

	if r.ID > 0 {
		err :=  model.PutResource(r.ID, data)
		if err != nil {
			return err
		}
		// resource update
		return SetupUserResource()
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
	err := model.DeleteResource(u.ID)
	if err != nil {
		return err
	}

	return SetupUserResource()
}
