go-toml-config
==============

[![Build Status](https://travis-ci.org/stvp/go-toml-config.png?branch=master)](https://travis-ci.org/stvp/go-toml-config)

go-toml-config is a simple [TOML](https://github.com/mojombo/toml)-based
configuration package for Golang apps that allows you to easily load
configuration files and set defaults. It's a simple wrapper around
[`flag.FlagSet`](http://golang.org/pkg/flag/), so you can use it in pretty much
the same exact way.

[API documentation](http://godoc.org/github.com/stvp/go-toml-config)

Example
--------

With `my_app.conf`:

```toml
country = "USA"

[atlanta]
enabled = true
population = 432427
temperature = 99.6
```

Use:

```go
import "github.com/stvp/go-toml-config"

var (
  country            = config.String("country", "Unknown")
  atlantaEnabled     = config.Bool("atlanta.enabled", false)
  alantaPopulation   = config.Int("atlanta.population", 0)
  atlantaTemperature = config.Float64("atlanta.temperature", 0)
)

func main() {
  config.Parse("/path/to/my_app.conf")
}
```

You can also create different ConfigSets to manage different logical groupings
of config variables:

```go
networkConfig = config.NewConfigSet("network settings", config.ExitOnError)
networkConfig.String("host", "localhost")
networkConfig.Int("port", 8080)
networkConfig.Parse("/path/to/network.conf")
```

Contributors
------------

Thanks all!

* @tysonmote
* @matrixik
* @fwang2002
* @shanks
* @xboston
* @tgulacsi


