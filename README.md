# goscanr

A fast, concurrent TCP port scanner written in Go.

## Description

`goscanr` is a simple yet efficient port scanner that uses Go's concurrency features to scan multiple ports simultaneously. It can quickly identify open TCP ports on a target host within a specified range.

## Features

- **Fast concurrent scanning**: Uses goroutines to scan multiple ports simultaneously
- **Adaptive concurrency**: Automatically adjusts concurrency based on network conditions
- **nmap integration**: Pipes open ports to nmap for service detection when available
- **Banner grabbing**: Captures service banners on open ports
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
| `-concurrency` | Initial concurrent scans (adapts automatically) | 500 |

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
example.com (1.2.3.4) — ports 1-1024
┌─────────┬───────┬─────────┬────────────────────────┐
│ PORT    │ STATE │ SERVICE │ BANNER                 │
├─────────┼───────┼─────────┼────────────────────────┤
│ 22/tcp  │ open  │ ssh     │ SSH-2.0-OpenSSH_9.3    │
│ 80/tcp  │ open  │ http    │                        │
│ 443/tcp │ open  │ https   │                        │
└─────────┴───────┴─────────┴────────────────────────┘
Done in 365ms
```

## How It Works

Scanning happens in two phases:

1. **Port discovery**: goscanr attempts TCP connections across the specified port range using goroutines. Ports are processed in batches and concurrency adapts automatically — increasing when the network is healthy, backing off when timeouts are detected.
2. **Service detection**: Open ports are passed to nmap (if installed) for service and version identification. Results are displayed in a table with any banners captured during the initial connection.

## Performance Considerations

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

MIT — see [LICENSE](LICENSE) for details.

## Author

Created by [bvdwalt](https://github.com/bvdwalt)