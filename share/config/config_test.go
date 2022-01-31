package config

import (
	"testing"
)


func TestConfig(t *testing.T)  {
	var Cfg = struct {
		Logger struct{Name string}
	}{}
	err := Load(&Cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(Cfg)
}
