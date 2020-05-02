package test

import (
	"net"
	"testing"
	"time"
	"vortex"
)

func TestGetTask(t *testing.T) {
	vortex.InitLogger()
	go vortex.RunDTUConnector()

	Connect()

	dds := vortex.GetDTUSessions()
	a, _ := dds.Load("127.0.0.1")
	b := a.(vortex.DTUSession)
	c := b.GetTask("da12-abba")
	println(c)

}

func Connect() {
	net.Dial("tcp", "127.0.0.1:6564")
	time.Sleep(time.Second * 1)
}
