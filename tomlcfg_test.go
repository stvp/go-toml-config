package config

import (
	"flag"
	"testing"
)

const (
	GOOD_CONFIG_PATH    = "examples/good.conf"
	INVALID_CONFIG_PATH = "examples/invalid.conf"
	MISSING_CONFIG_PATH = "examples/nope.conf"
)

func testBadParse(t *testing.T, c *ConfigSet) {
	err := c.Parse(MISSING_CONFIG_PATH)
	if err == nil || err.Error() != "open examples/nope.conf: no such file or directory" {
		t.Error("Expected error when loading missing TOML file, got", err)
	}

	err = c.Parse(INVALID_CONFIG_PATH)
	if err == nil || err.Error() != "examples/invalid.conf is not a valid TOML file. See https://github.com/mojombo/toml for syntax help." {
		t.Error("Expected error when loading missing TOML file, got", err)
	}
}

func testGoodParse(t *testing.T, c *ConfigSet) {
	boolSetting := c.Bool("my_bool", false)
	intSetting := c.Int("my_int", 0)
	int64Setting := c.Int64("my_bigint", 0)
	uintSetting := c.Uint("my_uint", 0)
	uint64Setting := c.Uint64("my_biguint", 0)
	stringSetting := c.String("my_string", "nope")
	float64Setting := c.Float64("my_bigfloat", 0)
	nestedSetting := c.String("section.name", "")
	deepNestedSetting := c.String("places.california.name", "")

	err := c.Parse(GOOD_CONFIG_PATH)
	if err != nil {
		t.Fatal(err)
	}

	if *boolSetting != true {
		t.Error("bool setting should be true, is", *boolSetting)
	}
	if *intSetting != 22 {
		t.Error("int setting should be 22, is", *intSetting)
	}
	if *int64Setting != int64(-23) {
		t.Error("int64 setting should be -23, is", *int64Setting)
	}
	if *uintSetting != 24 {
		t.Error("uint setting should be 24, is", *uintSetting)
	}
	if *uint64Setting != uint64(25) {
		t.Error("uint64 setting should be 25, is", *uint64Setting)
	}
	if *stringSetting != "ok" {
		t.Error("string setting should be \"ok\", is", *stringSetting)
	}
	if *float64Setting != float64(26.1) {
		t.Error("float64 setting should be 26.1, is", *float64Setting)
	}
	if *nestedSetting != "cool dude" {
		t.Error("nested setting should be \"cool dude\", is", *nestedSetting)
	}
	if *deepNestedSetting != "neat dude" {
		t.Error("deepNested setting should be \"neat dude\", is", *deepNestedSetting)
	}
}

func TestParse(t *testing.T) {
	testBadParse(t, globalConfig)
	testBadParse(t, NewConfigSet("App Config", flag.ExitOnError))
	testGoodParse(t, globalConfig)
	testGoodParse(t, NewConfigSet("App Config", flag.ExitOnError))
}
