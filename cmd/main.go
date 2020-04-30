package main

import (
	"vortex"
)

func main() {
	vortex.InitLogger()
	go vortex.RunDTUConnector()
	select {}
}
