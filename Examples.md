# tcping Usage Examples

## Basic Usage

### Test if a port is open

```shell
tcping google.com 443
# Output:    google.com port 443 open
# Exit code: 0
```

### Test with timeout

```shell
tcping -t 5 example.com 80
# Output:    example.com port 80 open
# Exit code: 0
```

### Quiet mode (for scripts)

```shell
tcping -q google.com 443
# Output: <no output, only exit code>
echo $?
# Output: 0
```

## Advanced Usage

### Force IPv4 or IPv6

```shell
# Force IPv4
tcping -f 4 google.com 443
# Force IPv6
tcping -f 6 google.com 443
```

### Microsecond precision timeout

```shell
# 500ms timeout (500ms = 500000 microseconds)
tcping -u 500000 example.com 443
```

### Combine options

```shell
tcping -q -f 4 -t 3 127.0.0.1 22
```
