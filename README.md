# Build

Run `make build` to build the application. Executable will be built to `build/alphabill`. 

### Build dependencies

* `golang` version 1.20. (https://go.dev/doc/install)
* in order to rebuild everything including protobuf definitions (`make` or `make all`):
  * `protoc` version 3.21.9+ (https://grpc.io/docs/protoc-installation)
  * `protoc-gen-go` (https://grpc.io/docs/languages/go/quickstart/)

# Money Partition

1. Run script `./setup-testab.sh -m 3 -t 0 -d 0` to generate configuration for a root chain and 3 money partition nodes.
    The script generates rootchain and partition node keys, genesis files.
    Node configuration files are located in `testab` directory.
2. Run script `./start.sh -r -p money -b money` to start rootchain and 3 money partition nodes and money backend

3. Run script `stop.sh -a` to stop the root chain and partition nodes.
   
   Alternatively, use `stop.sh` to stop any partition or root and `start.sh` to resume. See command help for more details. 

## Configuration

It's possible to define the configuration values from (in the order of precedence):

* Command line flags (e.g. `--address="/ip4/127.0.0.1/tcp/26652"`)
* Environment (Prefix 'AB' must be used. E.g. `AB_ADDRESS="/ip4/127.0.0.1/tcp/26652"`)
* Configuration file (properties file) (E.g. `address="/ip4/127.0.0.1/tcp/26652"`)
* Default values

The default location of configuration file is `$AB_HOME/config.props`

The default `$AB_HOME` is `$HOME/.alphabill`

# Verifiable Data Partition
1. Run script `./setup-testab.sh -m 0 -t 0 -d 3` to generate configuration for a root chain and 3 vd partition nodes.
   The script generates rootchain and partition node keys, genesis files.
   Node configuration files are located in `testab` directory.
2. Run script `./start.sh -r -p vd` to start rootchain and 3 vd partition nodes
3. Run script `stop.sh -a` to stop the root chain and partition nodes.

# User Token Partition
1. Run script `./setup-testab.sh -m 0 -t 3 -d 0` to generate configuration for a root chain and 3 token partition nodes.
   The script generates rootchain and partition node keys, genesis files.
   Node configuration files are located in `testab` directory.
2. Run script `./start.sh -r -p tokens -b tokens` to start rootchain and 3 token partition nodes and token backend
3. Run script `stop.sh -a` to stop the root chain and partition nodes.

# Start all partitions at once
1. Run script `./setup-testab.sh` to generate genesis for root, and 3 money, vd and tokens nodes.
2. Run `start.sh -r -p money -p tokens -p vd -b money -b tokens` to start everything
3. Run `stop.sh -a` to stop everything

# Logging configuration

Logging can be configured through a yaml configuration file. See `cli/alphabill/config/logger-config.yaml` for example.

Default location of the logger configuration file is `$AB_HOME/logger-config.yaml`

The location can be changed through `--logger-config` configuration key. If it's relative URL, then it's relative
to `$AB_HOME`.

# Wallet Logging Configuration

Wallet logging can be configured only through CLI parameters. 

`./alphabill wallet --log-file=<path/to/my/file> --log-level=INFO`

Default log output is `stderr` and default log level is `INFO`. 

Possible log level values: `ERROR, WARNING, NOTICE, INFO, DEBUG`

# Set up autocompletion

To use autocompletion (supported with `bash`, `fish`, `powershell` and `zsh`), run the following commands after
building (this is `bash` example):

* `./alphabill completion bash > /tmp/completion`
* `source /tmp/completion`

# CI setup

See gitlab-ci.yml for details.

GitLab runs the CI job inside docker container defined in `alphabill/gitlab-ci-image`.
