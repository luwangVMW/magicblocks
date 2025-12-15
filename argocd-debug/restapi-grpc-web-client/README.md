# Argo CD Connection Debugger

A comprehensive bash script for testing and debugging Argo CD API connectivity and protocol support.

## Overview

`connection.sh` tests multiple Argo CD API methods to diagnose connectivity and protocol issues. It provides detailed output with color-coded results for easy troubleshooting.

## Features

- **REST API Testing**: HTTP/1.1 + JSON (`/api/version`)
- **gRPC-Web Testing**: HTTP/1.1 + Protobuf (`/version.VersionService/Version`)
- **Color-coded Output**: Visual feedback for success/failure states
- **Detailed Diagnostics**: Response times, payload inspection, and troubleshooting tips
- **Prerequisite Checks**: Validates required tools (curl, jq, xxd)

## Usage

```bash
./connection.sh [HOST] [PORT] [PROTOCOL]
```

### Parameters

- `HOST` - Target Argo CD server (default: `172.16.0.203`)
- `PORT` - Port number (default: `443`)
- `PROTOCOL` - http or https (default: `https`)

### Examples

```bash
# Use default settings
./connection.sh

# Test specific host
./connection.sh argocd.example.com

# Custom host and port
./connection.sh 192.168.1.100 8080

# HTTP (non-TLS) connection
./connection.sh localhost 8080 http
```

## Prerequisites

### Required
- `curl` - HTTP client (pre-installed on most systems)

### Optional (for enhanced output)
- `jq` - JSON formatting
- `xxd` - Hex dump viewer

Install on macOS:
```bash
brew install jq
```

## Output

The script provides:
1. **Prerequisites Check**: Verifies available tools
2. **REST API Test**: JSON response with HTTP status codes
3. **gRPC-Web Test**: Protobuf binary response with extracted strings
4. **Summary**: Overview of all test results
5. **Troubleshooting Tips**: Context-specific guidance for failures

## Exit Codes

- `0` - Script completed (check summary for individual test results)

## Notes

- The script uses `-k` flag to skip SSL certificate verification (useful for self-signed certs)
- Temporary files are automatically cleaned up on exit
- All tests use the version endpoint which doesn't require authentication

