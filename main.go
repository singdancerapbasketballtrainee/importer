package main

import (
	"importer/app/config"
	"importer/app/di"
	"importer/app/log"
)

func main() {
	config.ConfigureInit()
	log.InitLog()
	log.Log.Info("Start up")
	app, cleanup, err := di.InitApp()
	if err != nil {
		log.Log.Fatal(err.Error())
	}
	err = app.Start()
	if err != nil {
		log.Log.Fatal(err.Error())
	}
	cleanup()
}
