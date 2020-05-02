package test

import (
	"fmt"
	"testing"
	"vortex"
)

func TestGetConfig(t *testing.T) {
	ls := vortex.GetConfig()
	fmt.Println(ls.Name)
}

func TestSetConfig(t *testing.T) {
	ls := vortex.GetConfig()
	ls.Name = "修改名称"
	_ = ls.SetConfig()
	ls = vortex.GetConfig()
	if ls.Name != "修改名称" {
		t.Fail()
	}
	ls.Name = "节点实例名"
	_ = ls.SetConfig()
}

func TestDeadline(t *testing.T) {
	ls := vortex.GetConfig()
	println(ls.Name)

	ls = vortex.GetConfig()
	println(ls.Name)
}

func TestGetAttach(t *testing.T) {
	addr := "192.168.1.223"
	ls := vortex.GetConfig()
	k := ls.GetSensorSetByAttach(addr)
	if k[0].Attach != addr {
		t.Fail()
	}
}
