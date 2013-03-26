/*
	Package config implements simple TOML-based configuration variables, based on
	the flag package in the standard Go library (In fact, it's just a simple
	wrapper around flag.FlagSet). It is used in a similar manner, minus the usage
	strings and other command-line specific bits.

	Usage:

	Given the following TOML file:

		country = "USA"

		[atlanta]
		enabled = true
		population = 432427
		temperature = 99.6

	Define your config variables and give them defaults:

		import config "github.com/stvp/tomlcfg"

		var (
			country            = config.String("country", "Unknown")
			atlantaEnabled     = config.Bool("atlanta.enabled", false)
			alantaPopulation   = config.Int("atlanta.population", 0)
			atlantaTemperature = config.Float("atlanta.population", 0)
		)

	After all the config variables are defined, load the config file to overwrite
	the default values with the user-supplied config settings:

		if err := config.Parse("/path/to/myconfig.conf"); err != nil {
			panic(err)
		}

	You can also create separate ConfigSets for different config files:

		networkConfig = config.NewConfigSet("network settings", config.ExitOnError)
*/
package config

import (
	"errors"
	"flag"
	"fmt"
	toml "github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
	"time"
)

type ConfigSet struct {
	*flag.FlagSet
}

// Bool defines a bool config variable with a given name and default value for
// a ConfigSet.
func (c *ConfigSet) Bool(name string, value bool) *bool {
	return c.FlagSet.Bool(name, value, "")
}

// Int defines a int config variable with a given name and default value for a
// ConfigSet.
func (c *ConfigSet) Int(name string, value int) *int {
	return c.FlagSet.Int(name, value, "")
}

// Int64 defines a int64 config variable with a given name and default value
// for a ConfigSet.
func (c *ConfigSet) Int64(name string, value int64) *int64 {
	return c.FlagSet.Int64(name, value, "")
}

// Uint defines a uint config variable with a given name and default value for
// a ConfigSet.
func (c *ConfigSet) Uint(name string, value uint) *uint {
	return c.FlagSet.Uint(name, value, "")
}

// Uint64 defines a uint64 config variable with a given name and default value
// for a ConfigSet.
func (c *ConfigSet) Uint64(name string, value uint64) *uint64 {
	return c.FlagSet.Uint64(name, value, "")
}

// String defines a string config variable with a given name and default value
// for a ConfigSet.
func (c *ConfigSet) String(name string, value string) *string {
	return c.FlagSet.String(name, value, "")
}

// Float64 defines a float64 config variable with a given name and default
// value for a ConfigSet.
func (c *ConfigSet) Float64(name string, value float64) *float64 {
	return c.FlagSet.Float64(name, value, "")
}

// Duration defines a time.Duration config variable with a given name and
// default value.
func (c *ConfigSet) Duration(name string, value time.Duration) *time.Duration {
	return globalConfig.FlagSet.Duration(name, value, "")
}

// Parse takes a path to a TOML file and loads it. This must be called after
// all the config flags in the ConfigSet have been defined but before the flags
// are accessed by the program.
func (c *ConfigSet) Parse(path string) error {
	configBytes, err := ioutil.ReadFile(path)
	if err != nil {
		errorString := fmt.Sprintf("Error loading %s: %s", path, err.Error())
		return errors.New(errorString)
	}

	tomlTree, err := toml.Load(string(configBytes))
	if err != nil {
		errorString := fmt.Sprintf("Error parsing config: %s", err.Error())
		return errors.New(errorString)
	}

	// TODO: need to get *all* keys "section.name" etc.
	for _, key := range tomlTree.Keys() {
		c.Set(key, fmt.Sprintf("%v", tomlTree.Get(key)))
	}

	return nil
}

const (
	ContinueOnError flag.ErrorHandling = flag.ContinueOnError
	ExitOnError     flag.ErrorHandling = flag.ExitOnError
	PanicOnError    flag.ErrorHandling = flag.PanicOnError
)

// NewConfigSet returns a new ConfigSet with the given name and error handling
// policy. The three valid error handling policies are: flag.ContinueOnError,
// flag.ExitOnError, and flag.PanicOnError.
func NewConfigSet(name string, errorHandling flag.ErrorHandling) *ConfigSet {
	return &ConfigSet{
		flag.NewFlagSet(name, errorHandling),
	}
}

// -- globalConfig

var globalConfig = NewConfigSet(os.Args[0], flag.ExitOnError)

// Bool defines a bool config variable with a given name and default value.
func Bool(name string, value bool) *bool {
	return globalConfig.Bool(name, value)
}

// Int defines a int config variable with a given name and default value.
func Int(name string, value int) *int {
	return globalConfig.Int(name, value)
}

// Int64 defines a int64 config variable with a given name and default value.
func Int64(name string, value int64) *int64 {
	return globalConfig.Int64(name, value)
}

// Uint defines a uint config variable with a given name and default value.
func Uint(name string, value uint) *uint {
	return globalConfig.Uint(name, value)
}

// Uint64 defines a uint64 config variable with a given name and default value.
func Uint64(name string, value uint64) *uint64 {
	return globalConfig.Uint64(name, value)
}

// String defines a string config variable with a given name and default value.
func String(name string, value string) *string {
	return globalConfig.String(name, value)
}

// Float64 defines a float64 config variable with a given name and default
// value.
func Float64(name string, value float64) *float64 {
	return globalConfig.Float64(name, value)
}

// Duration defines a time.Duration config variable with a given name and
// default value.
func Duration(name string, value time.Duration) *time.Duration {
	return globalConfig.Duration(name, value)
}

func Parse(path string) error {
	return globalConfig.Parse(path)
}
