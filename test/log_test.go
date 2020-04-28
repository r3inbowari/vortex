package test

import (
	"testing"
	"vortex"
)

func TestInfo(t *testing.T) {
	vortex.InitLogger()
	vortex.Info("hello", map[string]interface{}{"param0:": "hi", "param1": 0})
}

func TestWarn(t *testing.T) {
	vortex.InitLogger()
	vortex.Warn("hello", map[string]interface{}{"param0:": "hi", "param1": 0})
}

func TestError(t *testing.T) {
	vortex.InitLogger()
	vortex.Error("hello", map[string]interface{}{"param0:": "hi", "param1": 0})
}

func TestFatal(t *testing.T) {
	vortex.InitLogger()
	vortex.Fatal("hello", map[string]interface{}{"param0:": "hi", "param1": 0})
}

func TestPanic(t *testing.T) {
	vortex.InitLogger()
	vortex.Panic("hello", map[string]interface{}{"param0:": "hi", "param1": 0})
}

func TestTrace(t *testing.T) {
	vortex.InitLogger()
	vortex.Trace("hello", map[string]interface{}{"param0:": "hi", "param1": 0})
}
