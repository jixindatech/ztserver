package app

import (
	"fmt"
	"github.com/astaxie/beego/validation"
)

// MarkErrors mark error log info
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		// log.Println(err.Key, err.Message)
		fmt.Println(err.Message)
	}

	return
}
