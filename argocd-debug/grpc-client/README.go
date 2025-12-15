# ArgoCD gRPC Client Example

A simple Go gRPC client that demonstrates how to call the ArgoCD version service to validate grpc connection to ArgoCD server.

## Features

- ✓ TLS and plaintext connections
- ✓ Bearer token authentication
- ✓ Configurable timeout
- ✓ Skip TLS verification (for self-signed certificates)
- ✓ Clean, well-documented code

## Prerequisites

- Go 1.25 or later
- Access to an ArgoCD server

## Build

```bash
cd examples/grpc-client
go build -o grpc-client

```

## Usage

### Basic usage (with TLS, skip verification)

```bash
./grpc-client -server 172.16.0.204:443
```

### With authentication

```bash
# Get token first
TOKEN=$(curl -k -s https://172.16.0.204/api/v1/session \
  -d '{"username":"admin","password":"your-password"}' \
  | jq -r .token)

# Call with token
./grpc-client -server 172.16.0.204:443 -token "$TOKEN"
```

### Plaintext connection (no TLS)

```bash
./grpc-client -server 172.16.0.204:80 -plaintext
```

### Custom timeout

```bash
./grpc-client -server 172.16.0.204:443 -timeout 30s
```

### All options

```bash
./grpc-client \
  -server 172.16.0.204:443 \
  -token "your-token-here" \
  -insecure \
  -timeout 10s
```

## Command Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `-server` | `172.16.0.204:443` | ArgoCD server address (host:port) |
| `-token` | `""` | Bearer token for authentication |
| `-insecure` | `true` | Skip TLS certificate verification |
| `-plaintext` | `false` | Use plaintext connection (no TLS) |
| `-timeout` | `10s` | Request timeout duration |

## Example Output

```
Connecting to ArgoCD server: 172.16.0.204:443
Plaintext: false, Insecure TLS: true

Calling version.VersionService/Version...

✓ SUCCESS! Version information:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Version:          v2.13.2+c4b283e
Build Date:       2025-01-15T10:30:00Z
Git Commit:       c4b283e1234567890abcdef
Git Tag:          v2.13.2
Git Tree State:   clean
Go Version:       go1.23.1
Compiler:         gc
Platform:         linux/amd64
Kustomize:        v5.4.3
Helm:             v3.16.3
Kubectl:          v1.32.0
Jsonnet:          v0.20.0
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## How It Works

1. **Create gRPC Connection**: Establishes a gRPC connection with TLS or plaintext
2. **Create Client**: Uses the generated `NewVersionServiceClient` from protobuf
3. **Add Authentication**: Adds bearer token to context metadata if provided
4. **Make Request**: Calls `Version()` with an empty protobuf message
5. **Display Results**: Prints the version information returned by the server

## Code Structure

- `main.go` - Main client implementation
- `go.mod` - Go module definition with local ArgoCD dependency
- `README.md` - This file

## Extending This Example

To call other ArgoCD services, follow the same pattern:

```go
// For application service
import applicationpkg "github.com/argoproj/argo-cd/v3/pkg/apiclient/application"

client := applicationpkg.NewApplicationServiceClient(conn)
apps, err := client.List(ctx, &application.ApplicationQuery{})

// For cluster service
import clusterpkg "github.com/argoproj/argo-cd/v3/pkg/apiclient/cluster"

client := clusterpkg.NewClusterServiceClient(conn)
clusters, err := client.List(ctx, &cluster.ClusterQuery{})
```

## Troubleshooting

### Connection refused
- Check if the server address and port are correct
- Verify the ArgoCD server is running and accessible

### TLS errors
- Use `-insecure` flag to skip certificate verification
- Or use `-plaintext` for non-TLS connections

### Permission denied
- The version endpoint may require authentication
- Use `-token` flag with a valid bearer token



