#!/bin/bash
# Argo CD Version Endpoint Debugging Script
# Tests all three API methods: REST API, gRPC-Web
# Use this script to debug connectivity and protocol issues with Argo CD

set -e

# Configuration
HOST="${1:-172.16.0.203}"
PORT="${2:-443}"
HTTP="${3:-https}"
ENDPOINT="${HOST}:${PORT}"

# Color codes for better readability
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo ""
    echo "========================================"
    echo -e "${BLUE}$1${NC}"
    echo "========================================"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_info() {
    echo "  $1"
}

# Main header
clear
echo "========================================"
echo "  Argo CD Version Endpoint Debugger"
echo "========================================"
echo ""
echo "Target: $ENDPOINT"
echo "Time: $(date)"
echo ""

# Check prerequisites
print_header "Checking Prerequisites"

CURL_VERSION=$(curl --version | head -n1)
print_info "curl: $CURL_VERSION"

if command -v jq &> /dev/null; then
    JQ_VERSION=$(jq --version)
    print_success "jq: $JQ_VERSION"
    HAS_JQ=true
else
    print_warning "jq not found (JSON formatting will be limited)"
    HAS_JQ=false
fi

if command -v xxd &> /dev/null; then
    print_success "xxd: available"
    HAS_XXD=true
else
    print_warning "xxd not found (hex dump will be limited)"
    HAS_XXD=false
fi


# Test 1: REST API (HTTP/1.1 + JSON)
print_header "Test 1: REST API (HTTP/1.1 + JSON)"
print_info "Method: GET /api/version"
print_info "Protocol: HTTP/1.1"
print_info "Content-Type: application/json"
echo ""

REST_OUTPUT=$(mktemp)
REST_RESPONSE=$(mktemp)

curl -k -s -w "\nHTTP_CODE:%{http_code}\nTIME_TOTAL:%{time_total}s\n" \
     -H 'Accept: application/json' \
     "$HTTP://${ENDPOINT}/api/version" \
     -o "$REST_RESPONSE" 2>&1 | tee "$REST_OUTPUT" > /dev/null

HTTP_CODE=$(grep "HTTP_CODE:" "$REST_OUTPUT" | cut -d: -f2)
TIME_TOTAL=$(grep "TIME_TOTAL:" "$REST_OUTPUT" | cut -d: -f2)

if [ "$HTTP_CODE" = "200" ]; then
    print_success "REST API successful (HTTP $HTTP_CODE, ${TIME_TOTAL})"
    echo ""
    print_info "Response:"
    if [ "$HAS_JQ" = true ]; then
        cat "$REST_RESPONSE" | jq '.' | head -20
    else
        cat "$REST_RESPONSE"
    fi
    REST_SUCCESS=true
else
    print_error "REST API failed (HTTP $HTTP_CODE)"
    cat "$REST_RESPONSE"
    REST_SUCCESS=false
fi

# Test 2: gRPC-Web (HTTP/1.1 + Protobuf)
print_header "Test 2: gRPC-Web (HTTP/1.1 + Protobuf)"
print_info "Method: POST /version.VersionService/Version"
print_info "Protocol: HTTP/1.1"
print_info "Content-Type: application/grpc-web+proto"
echo ""

# Create empty gRPC request
GRPCWEB_REQUEST=$(mktemp)
{
    printf "\x00"              # Compression flag (0 = no compression)
    printf "\x00\x00\x00\x00"  # Message length = 0 (4 bytes, big-endian)
} > "$GRPCWEB_REQUEST"

GRPCWEB_RESPONSE=$(mktemp)
GRPCWEB_OUTPUT=$(mktemp)

curl -k -s -w "\nHTTP_CODE:%{http_code}\nTIME_TOTAL:%{time_total}s\n" \
     -X POST \
     -H 'Content-Type: application/grpc-web+proto' \
     -H 'Accept: application/grpc-web+proto' \
     --data-binary "@${GRPCWEB_REQUEST}" \
     "$HTTP://${ENDPOINT}/version.VersionService/Version" \
     -o "$GRPCWEB_RESPONSE" 2>&1 | tee "$GRPCWEB_OUTPUT" > /dev/null

HTTP_CODE=$(grep "HTTP_CODE:" "$GRPCWEB_OUTPUT" | cut -d: -f2)
TIME_TOTAL=$(grep "TIME_TOTAL:" "$GRPCWEB_OUTPUT" | cut -d: -f2)

if [ "$HTTP_CODE" = "200" ]; then
    RESP_SIZE=$(wc -c < "$GRPCWEB_RESPONSE" | tr -d ' ')
    print_success "gRPC-Web successful (HTTP $HTTP_CODE, ${TIME_TOTAL}, ${RESP_SIZE} bytes)"
    echo ""
    
    print_info "Response (protobuf binary):"
    if [ "$HAS_XXD" = true ]; then
        xxd "$GRPCWEB_RESPONSE" | head -10
    else
        print_warning "xxd not available for hex dump"
    fi
    
    echo ""
    print_info "Extracted strings from protobuf:"
    strings "$GRPCWEB_RESPONSE" 2>/dev/null | grep -E '.' | while read -r line; do
        print_info "  → $line"
    done | head -15
    GRPCWEB_SUCCESS=true
else
    print_error "gRPC-Web failed (HTTP $HTTP_CODE)"
    cat "$GRPCWEB_RESPONSE"
    GRPCWEB_SUCCESS=false
fi

# Summary
print_header "Summary"
echo ""

print_info "Test Results for $ENDPOINT:"
echo ""

if [ "$REST_SUCCESS" = true ]; then
    print_success "REST API (HTTP/1.1 + JSON)"
else
    print_error "REST API (HTTP/1.1 + JSON)"
fi

if [ "$GRPCWEB_SUCCESS" = true ]; then
    print_success "gRPC-Web (HTTP/1.1 + Protobuf)"
else
    print_error "gRPC-Web (HTTP/1.1 + Protobuf)"
fi


echo ""
print_header "Troubleshooting Tips"
echo ""

if [ "$REST_SUCCESS" = false ]; then
    print_error "REST API Failed:"
    print_info "  • Check if Argo CD server is running"
    print_info "  • Verify host:port is correct"
    print_info "  • Check firewall rules"
fi

if [ "$GRPCWEB_SUCCESS" = false ] && [ "$REST_SUCCESS" = true ]; then
    print_error "gRPC-Web Failed but REST API Works:"
    print_info "  • gRPC-Web might not be enabled"
    print_info "  • Check Content-Type header routing"
fi


if [ "$REST_SUCCESS" = true ] && [ "$GRPCWEB_SUCCESS" = true ]; then
    echo ""
    print_success "All protocols working correctly!"
fi

echo ""
print_info "Generated: $(date)"
echo ""

# Cleanup
rm -f "$REST_OUTPUT" "$REST_RESPONSE" \
      "$GRPCWEB_REQUEST" "$GRPCWEB_RESPONSE" "$GRPCWEB_OUTPUT" \
      "$GRPC_REQUEST" "$GRPC_RESPONSE" "$GRPC_OUTPUT" 2>/dev/null

exit 0

