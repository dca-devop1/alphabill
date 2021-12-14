# Build

Run `make` to test and build the application.

### Build dependencies

* `golang` version 1.17 or newer. (https://go.dev/doc/install)
* `protoc` version 3 or newer. (https://grpc.io/docs/protoc-installation)

# Run

Executable will be built to `build/alphabill`

Run the executable `alphabill shard` to start shard node with default configuration. To see possible configuration options run with `--help` flag.

## Configuration

It's possible to define the configuration values from (in the order of precedence):
* Command line flags (e.g. `--initial-bill-value=1000`)
* Environment (Prefix 'AB' must be used. E.g. `AB_INITIAL_BILL_VALUE=1000`)
* Configuration file (properties file) (E.g. `initial-bill-value=1000`)
* Default values

# CI setup

See gitlab-ci.yml for details.

GitLab runs the CI job inside docker container defined in `deployments/gitlab`.
