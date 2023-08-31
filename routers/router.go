package routers

import (
	"fiber/controller"
	handler "fiber/pkg"
	"fiber/pkg/logs"
	"fiber/pkg/setting"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func InitRouter() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:       false, //是否开启多进程
		CaseSensitive: false, //路由定义大小写问题的匹配
		StrictRouting: true,
		ServerHeader:  "sdboon", //定义响应头中的Server的标记头
		//ErrorHandler:
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	//抛出异常全局捕获
	//app.Use(recover.New())
	app.Use(handler.Recover)
	// 全局中间件件，对所有的路由生效
	app.Use(func(c *fiber.Ctx) error {
		sec, err := setting.Cfg.GetSection("ips")
		if err != nil {
			panic(fmt.Sprintf("Fail to get section 'ftp': %v", err))
		}
		ips := sec.Key("ips").MustString("127.0.0.1")
		isBool := inSlice(strings.Split(ips, ","), c.IP())
		if isBool {
			return c.Next()
		}
		logs.Logger.Error("IP不匹配:" + c.IP())
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"code": 500,
			"msg":  "IP不匹配",
			"data": "",
		})
	})

	apiV1 := app.Group("/v1")
	{
		{
			apiV1.Post("/GetCpuInfo", controller.GetCpuInfo)
			apiV1.Post("/GetMemInfo", controller.GetMemInfo)
			apiV1.Post("/GetSysLoad", controller.GetSysLoad)
			apiV1.Post("/GetDiskInfo", controller.GetDiskInfo)
		}
	}
	return app
}

// 是否存在切片中
func inSlice(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
