# tcping-go

> `tcping-go` is ported from [mkirchner/tcping](https://github.com/mkirchner/tcping)

Check if a desired port is reachable via TCP.

## Features

- TCP connection testing with customizable timeout
- IPv4/IPv6 support with auto-detection
- Quite mode for scripting
- Standard exit codes for automation
- Cross-platform (linux, macOS, Windows)
- Zero external dependencies
- Small binary size

## Installation

### Prebuilt Releases

See the latest [Releases](https://github.com/AyakuraYuki/tcping-go/releases).

### Build from source code

```shell
go build -o tcping main.go

# Or use Makefile
make build
```

### Cross-Compilation

```shell
# linux-amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tcping main.go
# or
GOOS=linux GOARCH=amd64 make cross-build

# linux-arm64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o tcping main.go
# or
GOOS=linux GOARCH=arm64 make cross-build

# macOS (Apple Silicon)
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o tcping main.go
# or
GOOS=darwin GOARCH=arm64 make cross-build

# macOS (Intel)
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o tcping main.go
# or
GOOS=darwin GOARCH=amd64 make cross-build

# Windows AMD64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o tcping.exe main.go
# or
GOOS=windows GOARCH=amd64 make cross-build
```

For more platforms such as BSD, feel free to modify `GOOS` and `GOARCH` for
building. Some platforms may require compilation within a physical system
environment, where using a Docker environment could be helpful.

## Usage

```shell
tcping [-q] [-f <4|6>] [-t timeout_sec] [-u timeout_usec] <host> <port>
```

### Arguments

- `<host>`: Target hostname or IP address
- `<port>`: Target port number or service name

### Options

- `-q`: Quite mode - suppress output
- `-f <4|6>`: IP family: 4 for IPv4, 6 for IPv6
- `-t <seconds>`: Timeout in seconds
- `-u <microseconds>`: Timeout in microseconds

### Examples

```shell
tcping google.com 443
tcping -f 4 localhost 80
tcping -t 3 10.10.10.100 22
```

For more examples, see [Examples.md](https://github.com/AyakuraYuki/tcping-go/blob/main/Examples.md)

## Exit Codes

- `0`: Connection successful
- `1`: Connection refused (port closed)
- `2`: Connection timeout
- `255`: Error occurred (invalid arguments, DNS failure, etc.)

## License

This Go implementation maintains compatibility with the original tcping project.

- Original C version: Copyright (C) 2003-2019 Marc Kirchner
- Go implementation: Copyright (C) 2026 Ayakura Yuki

## Contributing

Contributes are welcome, please ensure:

- Code follows Go conventions (`go fmt`, `go vet`)
- Maintains compatibility with original tcping behavior
- Includes tests for new features

## Acknowledgments

- Original tcping by Marc Kirchner
- Go standard library for excellent network primitives
