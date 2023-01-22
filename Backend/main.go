package main

import (
	"social/configs"
	"social/routes"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	configs.ConnectDB()
	routes.Run()
}
