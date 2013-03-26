package config

import (
	"flag"
	"testing"
)

const (
	GOOD_CONFIG_PATH = "examples/good.conf"
)

func TestParse(t *testing.T) {
	cfg := NewConfigSet("App Config", flag.ExitOnError)

	// Set in cfg file
	var name = cfg.String("name", "foo")
	var count = cfg.Int("section.count", 5)

	// Not set in cfg file
	var foo = cfg.String("foo", "default1")
	var nested_foo = cfg.String("nope.foo", "default2")

	err := cfg.Parse(GOOD_CONFIG_PATH)
	if err != nil {
		t.Error(err)
	}

	if *name != "tomlcfg" {
		t.Error(*name)
	}
	if *count != 10 {
		t.Error(*count)
	}
	if *foo != "default1" {
		t.Error(*foo)
	}
	if *nested_foo != "default2" {
		t.Error(*nested_foo)
	}
}
