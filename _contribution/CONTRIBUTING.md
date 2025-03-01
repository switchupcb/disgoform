# Contributing

## License

You agree to license any contribution to this library under the [Apache License 2.0](#license).

## Pull Requests

Pull requests must pass [static code analysis](#static-code-analysis) and work with all [test cases](#test).

## Project Structure

### Static Code Analysis

Disgo uses [golangci-lint](https://github.com/golangci/golangci-lint) to statically analyze code.

You can install golangci-lint with `go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5`.

Perform static code analysis using `golangci-lint run`.

#### Runtime Errors

1. If you receive a `diff` error, add a `diff` tool in your PATH: There is one located in the `Git/bin` directory.
2. If you receive `File is not ... with -...`, use `golangci-lint run --disable-all --no-config -Egofmt --fix` or ignore it.

### Test

#### Unit Tests

Unit tests test logic.

#### Integration Tests

Integration tests prove functionality between the underlying [API Wrapper (`disgo`)](https://github.com/switchupcb/disgo) and Discord.

#### Running Tests

Use `go test ./tests` to run the tests from the current directory.

Use [Github Action Workflow Files](/.github/workflows/) to find the correct test command and environment variables for a module.

# Roadmap

You can implement the following additional features.

1. Reverse Sync which generates a configuration file from the Discord Bot's current state.