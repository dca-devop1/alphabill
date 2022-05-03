# Build

Run `make` to test and build the application.

### Build dependencies

* `golang` version 1.18. (https://go.dev/doc/install)
* `protoc` version 3 or newer. (https://grpc.io/docs/protoc-installation)
* `protoc-gen-go` (https://grpc.io/docs/languages/go/quickstart/)

# Run

Executable will be built to `build/alphabill`

Run the executable `alphabill shard` to start shard node with default configuration. To see possible configuration options run with `--help` flag.

## Configuration

It's possible to define the configuration values from (in the order of precedence):
* Command line flags (e.g. `--initial-bill-value=1000`)
* Environment (Prefix 'AB' must be used. E.g. `AB_INITIAL_BILL_VALUE=1000`)
* Configuration file (properties file) (E.g. `initial-bill-value=1000`)
* Default values

The default location of configuration file is `$AB_HOME/config.props`

The default $AB_HOME is `$HOME/.alphabill`

### Logging configuration

Logging can be configured through a yaml configuration file. See `cli/alphabill/config/logger-config.yaml` for example.

Default location of the logger configuration file is `$AB_HOME/logger-config.yaml`

The location can be changed through `--logger-config` configuration key. If it's relative URL, then it's relative to `$AB_HOME`. 

### Wallet Logging Configuration

Wallet logging can be configured only through CLI parameters. 

`./alphabill wallet --log-file=<path/to/my/file> --log-level=INFO`

Default log file is `$AB_HOME/wallet/wallet.log`. Default log level is `INFO`

Possible log level values: `NONE, ERROR, WARNING, NOTICE, INFO, DEBUG`

## Set up autocompletion

To use autocompletion (supported with `bash`, `fish`, `powershell` and `zsh`), run the following commands after building (this is `bash` example):
* `./alphabill completion bash > /tmp/completion`
* `source /tmp/completion`

# CI setup

See gitlab-ci.yml for details.

GitLab runs the CI job inside docker container defined in `deployments/gitlab`.
