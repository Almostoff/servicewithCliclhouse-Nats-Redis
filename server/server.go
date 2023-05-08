package main

import (
	"taskFive/server/app"
	"taskFive/server/config"
)

func main() {
	cfg := new(config.Config)
	cfg.InitFile()
	app := app.InitApp(*cfg)
	app.Run()
}
