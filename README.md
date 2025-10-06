# goscanr

A fast, concurrent TCP port scanner written in Go.

## Description

`goscanr` is a simple yet efficient port scanner that uses Go's concurrency features to scan multiple ports simultaneously. It can quickly identify open TCP ports on a target host within a specified range.

## Features

- **Fast concurrent scanning**: Uses goroutines to scan multiple ports simultaneously
- **Configurable port range**: Specify start and end ports for scanning
- **Adjustable timeout**: Set connection timeout for port probes
- **Simple CLI interface**: Easy-to-use command-line flags
- **Lightweight**: Minimal dependencies, pure Go implementation

## Installation

### Prerequisites

- Go 1.24.7 or later

### Build from source

```bash
git clone https://github.com/bvdwalt/goscanr.git
cd goscanr
go build -o goscanr main.go
```

## Usage

### Basic Usage

```bash
./goscanr -target <hostname_or_ip>
```

### Advanced Usage

```bash
./goscanr -target <hostname_or_ip> -start <start_port> -end <end_port> -timeout <timeout_ms>
```

### Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-target` | Target hostname or IP address (required) | - |
| `-start` | Starting port number | 1 |
| `-end` | Ending port number | 1024 |
| `-timeout` | Connection timeout in milliseconds | 300 |

### Examples

#### Scan common ports on localhost
```bash
./goscanr -target localhost
```

#### Scan a specific port range
```bash
./goscanr -target example.com -start 80 -end 443
```

#### Scan with custom timeout
```bash
./goscanr -target 192.168.1.1 -start 1 -end 65535 -timeout 1000
```

#### Quick scan of well-known ports
```bash
./goscanr -target target.example.com -start 20 -end 1024 -timeout 200
```

## Sample Output

```
Scanning example.com from port 1 to 1024...
Port 22 is open
Port 53 is open
Port 80 is open
Port 443 is open
Scan complete.
```

## How It Works

The scanner uses Go's `net.DialTimeout()` function to attempt TCP connections to each port in the specified range. Each port is scanned in a separate goroutine for maximum concurrency, with a `sync.WaitGroup` ensuring all scans complete before the program exits.

## Performance Considerations

- **Concurrency**: The scanner creates one goroutine per port, which provides excellent performance but may be resource-intensive for very large port ranges
- **Timeout**: Lower timeout values increase scan speed but may miss slower-responding services
- **Network limits**: Very aggressive scanning may trigger rate limiting or be flagged by intrusion detection systems

## Legal and Ethical Use

This tool is intended for:
- Testing your own systems and networks
- Authorized penetration testing and security assessments
- Educational purposes

**Important**: Only scan systems you own or have explicit permission to test. Unauthorized port scanning may violate terms of service, local laws, or regulations.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source. Please check the repository for license details.

## Author

Created by [bvdwalt](https://github.com/bvdwalt)