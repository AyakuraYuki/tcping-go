package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var version = "0.1.0"

const (
	TcpingOpen    = 0   // Connection successful
	TcpingClosed  = 1   // Connection refused (port closed)
	TcpingTimeout = 2   // Connection timeout
	TcpingError   = 255 // Error occurred
)

type Config struct {
	Host         string
	Port         string
	Quite        bool
	IPFamily     string // "4", "6", or "" for auto
	TimeoutSec   uint64
	TimeoutMicro uint64
}

func main() {
	config := parseArgs()

	timeout := time.Duration(config.TimeoutSec)*time.Second +
		time.Duration(config.TimeoutMicro)*time.Microsecond

	if timeout == 0 {
		timeout = 10 * time.Second
	}

	exitCode := tcping(config.Host, config.Port, config.IPFamily, timeout, config.Quite)
	os.Exit(exitCode)
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [-q] [-f <4|6>] [-t timeout_sec] [-u timeout_usec] <host> <port>\n", os.Args[0])
	_, _ = fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
	_, _ = fmt.Fprintf(os.Stderr, "\nReturn codes:\n")
	_, _ = fmt.Fprintf(os.Stderr, "  %-3d - Connection successful\n", TcpingOpen)
	_, _ = fmt.Fprintf(os.Stderr, "  %-3d - Connection refused (port closed)\n", TcpingClosed)
	_, _ = fmt.Fprintf(os.Stderr, "  %-3d - Connection timeout\n", TcpingTimeout)
	_, _ = fmt.Fprintf(os.Stderr, "  %-3d - Error occurred\n", TcpingError)
}

func parseArgs() Config {
	var config Config

	flag.BoolVar(&config.Quite, "q", false, "Quite mode - suppress output")
	flag.StringVar(&config.IPFamily, "f", "", "IP family: 4 for IPv4, 6 for IPv6")
	flag.Uint64Var(&config.TimeoutSec, "t", 0, "Timeout in seconds")
	flag.Uint64Var(&config.TimeoutMicro, "u", 0, "Timeout in microseconds")

	flag.Usage = usage

	flag.Parse()

	if config.IPFamily != "" && config.IPFamily != "4" && config.IPFamily != "6" {
		_, _ = fmt.Fprintf(os.Stderr, "error: invalid IP family, use 4 or 6.\n\n")
		flag.Usage()
		os.Exit(TcpingError)
	}

	args := flag.Args()
	if len(args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "error: missing required arguments <host> and <port>\n\n")
		flag.Usage()
		os.Exit(TcpingError)
	}

	config.Host = args[0]
	config.Port = args[1]

	if _, err := strconv.Atoi(config.Port); err != nil {
		if _, err = net.LookupPort("tcp", config.Port); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error: invalid port: %s\n\n", config.Port)
			os.Exit(TcpingError)
		}
	}

	return config
}

func tcping(host, port, ipFamily string, timeout time.Duration, quite bool) int {
	var network string

	// determine network type based on IP family
	switch ipFamily {
	case "4":
		network = "tcp4"
	case "6":
		network = "tcp6"
	default:
		network = "tcp"
	}

	// resolve address
	address := net.JoinHostPort(host, port)

	// resolve to get the actual IP and port we're connecting to
	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		var (
			netErr net.Error
			opErr  *net.OpError
		)

		if errors.As(err, &netErr) && netErr.Timeout() {
			logf(quite, "%s port %s user timeout.\n", host, port)
			return TcpingTimeout
		}

		if errors.As(err, &opErr) && opErr != nil && opErr.Err.Error() == "connect: connection refused" {
			logf(quite, "%s port %s closed.\n", host, port)
			return TcpingClosed
		}

		logf(quite, "error: %s port %s: %v\n", host, port, err)
		return TcpingError
	}

	// connection successful
	defer func(conn net.Conn) { _ = conn.Close() }(conn)

	// get the actual resolved IP
	remoteAddr := conn.RemoteAddr().String()
	logf(quite, "%s port %s open\n", host, port)

	if !quite {
		if resolvedHost, resolvedPort, err := net.SplitHostPort(remoteAddr); err == nil {
			if resolvedHost != host {
				fmt.Printf("  (resolved to %s:%s)\n", resolvedHost, resolvedPort)
			}
		}
	}

	return TcpingOpen
}

func logf(quite bool, format string, args ...any) {
	if !quite {
		fmt.Printf(format, args...)
	}
}
