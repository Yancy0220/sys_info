package main

import (
	"fiber/pkg/setting"
	"fiber/routers"
	"fmt"
)

func main() {
	app := routers.InitRouter()

	app.Listen(fmt.Sprintf(":%d", setting.HTTPPort))
}
