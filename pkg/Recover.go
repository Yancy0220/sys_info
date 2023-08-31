package handler

import (
	"fiber/pkg/logs"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func Recover(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			//打印错误堆栈信息
			logs.Error(fmt.Sprintf("panic: %v\n", r))
			//debug.PrintStack()
			//封装通用json返回
			c.Status(fiber.StatusOK).JSON(fiber.Map{"code": 500, "msg": errorToString(r), "data": nil})
			return
		}
	}()
	//加载完 defer recover，继续后续接口调用
	return c.Next()
}

// recover错误，转string
func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return r.(string)
	}
}
