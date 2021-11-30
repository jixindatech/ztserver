package app

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"zt-server/pkg/e"
)

// BindAndValid bind and valid
func BindAndValid(c *gin.Context, form interface{}) int {
	err := c.Bind(form)
	if err != nil {
		//fmt.Println("bind error:", err)
		return e.InvalidParams
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		//fmt.Println("valid error:", err)
		return e.ERROR
	}

	if !check {
		//fmt.Println("check error:", valid.Errors)
		return e.InvalidParams
	}
	return e.SUCCESS
}
