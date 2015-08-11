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

		import "github.com/stvp/go-toml-config"

		var (
			country            = config.String("country", "Unknown")
			atlantaEnabled     = config.Bool("atlanta.enabled", false)
			alantaPopulation   = config.Int("atlanta.population", 0)
			atlantaTemperature = config.Float("atlanta.temperature", 0)
		)

	After all the config variables are defined, load the config file to overwrite
	the default values with the user-supplied config settings:

		if err := config.Parse("/path/to/myconfig.conf"); err != nil {
			panic(err)
		}

	You can also create separate ConfigSets for different config files:

		networkConfig = config.NewConfigSet("network settings", config.ExitOnError)
		networkConfig.String("host", "localhost")
		networkConfig.Int("port", 8080)
		networkConfig.Parse("/path/to/network.conf")
*/
package config

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/pelletier/go-toml"
)

// -- ConfigSet

type ConfigSet struct {
	*flag.FlagSet
}

// BoolVar defines a bool config with a given name and default value for a ConfigSet.
// The argument p points to a bool variable in which to store the value of the config.
func (c *ConfigSet) BoolVar(p *bool, name string, value bool) {
	c.FlagSet.BoolVar(p, name, value, "")
}

// Bool defines a bool config variable with a given name and default value for
// a ConfigSet.
func (c *ConfigSet) Bool(name string, value bool) *bool {
	return c.FlagSet.Bool(name, value, "")
}

// IntVar defines a int config with a given name and default value for a ConfigSet.
// The argument p points to a int variable in which to store the value of the config.
func (c *ConfigSet) IntVar(p *int, name string, value int) {
	c.FlagSet.IntVar(p, name, value, "")
}

// Int defines a int config variable with a given name and default value for a
// ConfigSet.
func (c *ConfigSet) Int(name string, value int) *int {
	return c.FlagSet.Int(name, value, "")
}

// Int64Var defines a int64 config with a given name and default value for a ConfigSet.
// The argument p points to a int64 variable in which to store the value of the config.
func (c *ConfigSet) Int64Var(p *int64, name string, value int64) {
	c.FlagSet.Int64Var(p, name, value, "")
}

// Int64 defines a int64 config variable with a given name and default value
// for a ConfigSet.
func (c *ConfigSet) Int64(name string, value int64) *int64 {
	return c.FlagSet.Int64(name, value, "")
}

// UintVar defines a uint config with a given name and default value for a ConfigSet.
// The argument p points to a uint variable in which to store the value of the config.
func (c *ConfigSet) UintVar(p *uint, name string, value uint) {
	c.FlagSet.UintVar(p, name, value, "")
}

// Uint defines a uint config variable with a given name and default value for
// a ConfigSet.
func (c *ConfigSet) Uint(name string, value uint) *uint {
	return c.FlagSet.Uint(name, value, "")
}

// Uint64Var defines a uint64 config with a given name and default value for a ConfigSet.
// The argument p points to a uint64 variable in which to store the value of the config.
func (c *ConfigSet) Uint64Var(p *uint64, name string, value uint64) {
	c.FlagSet.Uint64Var(p, name, value, "")
}

// Uint64 defines a uint64 config variable with a given name and default value
// for a ConfigSet.
func (c *ConfigSet) Uint64(name string, value uint64) *uint64 {
	return c.FlagSet.Uint64(name, value, "")
}

// StringVar defines a string config with a given name and default value for a ConfigSet.
// The argument p points to a string variable in which to store the value of the config.
func (c *ConfigSet) StringVar(p *string, name string, value string) {
	c.FlagSet.StringVar(p, name, value, "")
}

// String defines a string config variable with a given name and default value
// for a ConfigSet.
func (c *ConfigSet) String(name string, value string) *string {
	return c.FlagSet.String(name, value, "")
}

// Float64Var defines a float64 config with a given name and default value for a ConfigSet.
// The argument p points to a float64 variable in which to store the value of the config.
func (c *ConfigSet) Float64Var(p *float64, name string, value float64) {
	c.FlagSet.Float64Var(p, name, value, "")
}

// Float64 defines a float64 config variable with a given name and default
// value for a ConfigSet.
func (c *ConfigSet) Float64(name string, value float64) *float64 {
	return c.FlagSet.Float64(name, value, "")
}

// DurationVar defines a time.Duration config with a given name and default value for a ConfigSet.
// The argument p points to a time.Duration variable in which to store the value of the config.
func (c *ConfigSet) DurationVar(p *time.Duration, name string, value time.Duration) {
	c.FlagSet.DurationVar(p, name, value, "")
}

// Duration defines a time.Duration config variable with a given name and
// default value.
func (c *ConfigSet) Duration(name string, value time.Duration) *time.Duration {
	return c.FlagSet.Duration(name, value, "")
}

// Parse takes a path to a TOML file and loads it. This must be called after
// all the config flags in the ConfigSet have been defined but before the flags
// are accessed by the program.
func (c *ConfigSet) Parse(path string) error {
	configBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	tomlTree, err := toml.Load(string(configBytes))
	if err != nil {
		errorString := fmt.Sprintf("%s is not a valid TOML file. See https://github.com/mojombo/toml", path)
		return errors.New(errorString)
	}

	err = c.loadTomlTree(tomlTree, []string{})
	if err != nil {
		return err
	}

	return nil
}

// loadTomlTree recursively loads a TomlTree into this ConfigSet's config
// variables.
func (c *ConfigSet) loadTomlTree(tree *toml.TomlTree, path []string) error {
	for _, key := range tree.Keys() {
		fullPath := append(path, key)
		value := tree.Get(key)
		if subtree, isTree := value.(*toml.TomlTree); isTree {
			err := c.loadTomlTree(subtree, fullPath)
			if err != nil {
				return err
			}
		} else {
			fullPath := strings.Join(append(path, key), ".")
			err := c.Set(fullPath, fmt.Sprintf("%v", value))
			if err != nil {
				return buildLoadError(fullPath, err)
			}
		}
	}
	return nil
}

// buildLoadError takes an error from flag.FlagSet#Set and makes it a bit more
// readable, if it recognizes the format.
func buildLoadError(path string, err error) error {
	missingFlag := regexp.MustCompile(`^no such flag -([^\s]+)`)
	invalidSyntax := regexp.MustCompile(`^.+ parsing "(.+)": invalid syntax$`)
	errorString := err.Error()

	if missingFlag.MatchString(errorString) {
		errorString = missingFlag.ReplaceAllString(errorString, "$1 is not a valid config setting")
	} else if invalidSyntax.MatchString(errorString) {
		errorString = "The value for " + path + " is invalid"
	}

	return errors.New(errorString)
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

// BoolVar defines a bool config with a given name and default value.
// The argument p points to a bool variable in which to store the value of the config.
func BoolVar(p *bool, name string, value bool) {
	globalConfig.BoolVar(p, name, value)
}

// Bool defines a bool config variable with a given name and default value.
func Bool(name string, value bool) *bool {
	return globalConfig.Bool(name, value)
}

// IntVar defines a int config with a given name and default value.
// The argument p points to a int variable in which to store the value of the config.
func IntVar(p *int, name string, value int) {
	globalConfig.IntVar(p, name, value)
}

// Int defines a int config variable with a given name and default value.
func Int(name string, value int) *int {
	return globalConfig.Int(name, value)
}

// Int64Var defines a int64 config with a given name and default value.
// The argument p points to a int64 variable in which to store the value of the config.
func Int64Var(p *int64, name string, value int64) {
	globalConfig.Int64Var(p, name, value)
}

// Int64 defines a int64 config variable with a given name and default value.
func Int64(name string, value int64) *int64 {
	return globalConfig.Int64(name, value)
}

// UintVar defines a uint config with a given name and default value.
// The argument p points to a uint variable in which to store the value of the config.
func UintVar(p *uint, name string, value uint) {
	globalConfig.UintVar(p, name, value)
}

// Uint defines a uint config variable with a given name and default value.
func Uint(name string, value uint) *uint {
	return globalConfig.Uint(name, value)
}

// Uint64Var defines a uint64 config with a given name and default value.
// The argument p points to a uint64 variable in which to store the value of the config.
func Uint64Var(p *uint64, name string, value uint64) {
	globalConfig.Uint64Var(p, name, value)
}

// Uint64 defines a uint64 config variable with a given name and default value.
func Uint64(name string, value uint64) *uint64 {
	return globalConfig.Uint64(name, value)
}

// StringVar defines a string config with a given name and default value.
// The argument p points to a string variable in which to store the value of the config.
func StringVar(p *string, name string, value string) {
	globalConfig.StringVar(p, name, value)
}

// String defines a string config variable with a given name and default value.
func String(name string, value string) *string {
	return globalConfig.String(name, value)
}

// Float64Var defines a float64 config with a given name and default value.
// The argument p points to a float64 variable in which to store the value of the config.
func Float64Var(p *float64, name string, value float64) {
	globalConfig.Float64Var(p, name, value)
}

// Float64 defines a float64 config variable with a given name and default
// value.
func Float64(name string, value float64) *float64 {
	return globalConfig.Float64(name, value)
}

// DurationVar defines a time.Duration config with a given name and default value.
// The argument p points to a time.Duration variable in which to store the value of the config.
func DurationVar(p *time.Duration, name string, value time.Duration) {
	globalConfig.DurationVar(p, name, value)
}

// Duration defines a time.Duration config variable with a given name and
// default value.
func Duration(name string, value time.Duration) *time.Duration {
	return globalConfig.Duration(name, value)
}

// Parse takes a path to a TOML file and loads it into the global ConfigSet.
// This must be called after all config flags have been defined but before the
// flags are accessed by the program.
func Parse(path string) error {
	return globalConfig.Parse(path)
}
