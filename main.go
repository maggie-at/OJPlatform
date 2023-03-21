package main

import "OJPlatform/router"

func main() {
	ginServer := router.Router()
	ginServer.Run(":8888")
}
