# Overview
This project is intended to be used to clone every project accuweather has

# Config
Configuration can be placed in the following locations:

1. `/etc/git-dr/git-dr.yml`
2. `~/.git-dr/git-dr.yml`
3. `./git-dr.yml`

[Configuration can be "JSON, TOML, YML, HCL, or Java properties formats"](https://github.com/spf13/viper)

Check git-dr.example.yml for configuration options.

Config options can also be specified via environment variables.

# Running
Simply run the compiled binary. If there's an error, it'll angrily panic and return an exit code of 1 and log its error. Otherwise it'll return 0.
