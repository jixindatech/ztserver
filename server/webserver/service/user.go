package service

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"mime"
	"time"
	"zt-server/cache"
	"zt-server/pkg/core/golog"
	"zt-server/webserver/model"
	"zt-server/webserver/utils"
)

type User struct {
	ID uint64

	Name        string
	Phone       string
	Email       string
	Status      int
	Remark      string
	ResourceIds []uint64

	Page     int
	PageSize int
}

func (u *User) Save() error {
	data := make(map[string]interface{})
	data["name"] = u.Name
	data["email"] = u.Email
	data["phone"] = u.Phone
	data["status"] = u.Status
	data["remark"] = u.Remark

	if u.ID > 0 {
		return model.PutUser(u.ID, data)
	}
	return model.AddUser(data)
}

func (u *User) Get() (*model.User, error) {
	return model.GetUser(u.ID)
}

func (u *User) GetList() ([]*model.User, int, error) {
	data := make(map[string]interface{})
	data["name"] = u.Name
	data["page"] = u.Page
	data["pagesize"] = u.PageSize

	return model.GetUsers(data)
}

func (u *User) Delete() error {
	user, err := model.GetUser(u.ID)
	if err != nil {
		return err
	}
	email := user.Email
	err = model.DeleteUser(u.ID)
	if err != nil {
		return err
	}

	return cache.Del(email)
}

func (u *User) SaveUserResource() error {
	data := make(map[string]interface{})
	data["ids"] = u.ResourceIds
	err := model.SaveUserResource(u.ID, data)
	if err != nil {
		return err
	}

	user, err := model.GetUser(u.ID)
	if err != nil {
		return err
	}

	/* skip users without send mail */
	if user.Timestamp == 0 {
		return nil
	}

	resouces, err := model.GetUserResource(user.ID)
	if err != nil {
		return err
	}
	userResource := make(map[string]int)
	for _, resource := range resouces {
		userResource[resource.Server] = 1
	}
	err = setupUserCache(user.Email, user.ID, userResource, user.Timestamp)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetUserResource() ([]uint64, error) {
	data := make(map[string]interface{})
	data["ids"] = u.ResourceIds

	resources, err := model.GetUserResource(u.ID)
	if err != nil {
		return nil, err
	}

	var res []uint64
	for _, resource := range resources {
		res = append(res, resource.ID)
	}

	if len(res) == 0 {
		res = []uint64{}
	}

	return res, nil
}

func setupUserCache(email string, id uint64, resources map[string]int, timestamp int64) error {
	data := cache.UserCache{}
	data.ID = id
	data.Resource = resources
	data.Timestamp = timestamp

	userData, err := json.Marshal(data)
	if err != nil {
		golog.Error("user", zap.String("err", err.Error()))
		return err
	}

	err = cache.Set(email, string(userData))
	if err != nil {
		golog.Error("user", zap.String("err", err.Error()))
		return err
	}

	return nil
}

func (u *User) SendMail() error {
	email, err := model.GetEmail()
	if err != nil {
		return err
	}

	subject := "零信任客户端信息"
	//body := "零信任客户端信息"

	m := gomail.NewMessage()
	m.SetHeader("From", mime.QEncoding.Encode("UTF-8", "Support")+"<"+email.Email+">")
	m.SetHeader("To", u.Email)
	m.SetHeader("Subject", subject)

	data := make(map[string]interface{})
	data["timestamp"] = time.Now().Unix()

	if u.ID > 0 {
		// store timestamp to database
		err = model.PutUser(u.ID, data)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("%s", "invalid user id")
	}

	data["email"] = u.Email
	data["id"] = u.ID
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}

	encryptedData, err := utils.CBCEncryptWithPKCS7(string(jsonStr))
	if err != nil {
		return err
	}
	m.SetBody("text/html", encryptedData)
	d := gomail.NewDialer(email.Server, email.Port, email.Email, email.Password)

	go func() {
		err := d.DialAndSend(m)
		if err != nil {
			golog.Error("email", zap.String("send", err.Error()))
		}
	}()

	user, err := model.GetUser(u.ID)
	if err != nil {
		return err
	}

	resources, err := model.GetResourceByUser(user.ID)
	if err != nil {
		return err
	}

	userResource := make(map[string]int)
	for _, resource := range resources {
		userResource[resource.Server] = 1
	}

	err = setupUserCache(user.Email, user.ID, userResource, user.Timestamp)
	if err != nil {
		golog.Error("user", zap.String("err", err.Error()))
		return err
	}

	return nil
}

func SetupUser() (err error) {
	userSrv := User{}
	users, count, err := userSrv.GetList()
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}

	for _, user := range users {
		resources, err := model.GetResourceByUser(user.ID)
		if err != nil {
			return err
		}

		if user.Timestamp == 0 {
			continue
		}

		userResource := make(map[string]int)
		for _, resource := range resources {
			userResource[resource.Server] = 1
		}

		err = setupUserCache(user.Email, user.ID, userResource, user.Timestamp)
		if err != nil {
			golog.Error("user", zap.String("err", err.Error()))
			return err
		}
	}

	return err
}
