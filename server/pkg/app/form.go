package app

import (
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"zt-server/pkg/e"
)

// BindAndValid bind and valid
func BindAndValid(c *gin.Context, form interface{}) int {
	err := c.Bind(form)
	if err != nil {
		fmt.Println("err0:", err)
		return e.InvalidParams
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		fmt.Println("err1:", err)
		return e.ERROR
	}

	fmt.Println(form)

	if !check {
		fmt.Println("err2:", valid.Errors)
		return e.InvalidParams
	}
	return e.SUCCESS
}
