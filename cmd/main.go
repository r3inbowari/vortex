package main

import (
	"vortex"
)

func main() {
	vortex.InitLogger()
	vortex.GetMQTTInstance()
	go vortex.RunDTUService()
	select {}
}
