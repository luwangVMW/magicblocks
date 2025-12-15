package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	versionpkg "github.com/argoproj/argo-cd/v3/pkg/apiclient/version"
)

func main() {
	// Command line flags
	var (
		server       = flag.String("server", "172.16.0.204:443", "ArgoCD server address (host:port)")
		token        = flag.String("token", "", "Bearer token for authentication (optional)")
		insecureFlag = flag.Bool("insecure", true, "Skip TLS verification")
		plaintext    = flag.Bool("plaintext", false, "Use plaintext connection (no TLS)")
		timeout      = flag.Duration("timeout", 10*time.Second, "Request timeout")
		verbose      = flag.Bool("verbose", false, "Verbose output")
	)
	flag.Parse()

	if *verbose {
		fmt.Printf("Connecting to ArgoCD server: %s\n", *server)
		fmt.Printf("Plaintext: %v, Insecure TLS: %v\n", *plaintext, *insecureFlag)
		fmt.Println()
	}

	// Create gRPC connection
	conn, err := createGRPCConnection(*server, *plaintext, *insecureFlag)
	if err != nil {
		log.Fatalf("Failed to create gRPC connection: %v", err)
	}
	defer conn.Close()

	// Create version service client
	client := versionpkg.NewVersionServiceClient(conn)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// Add authentication token to context if provided
	if *token != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+*token)
	}

	// Call Version endpoint
	if *verbose {
		fmt.Println("Calling version.VersionService/Version...")
	}

	version, err := client.Version(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("Failed to get version: %v", err)
	}

	// Print version information
	fmt.Println("\n✓ SUCCESS! Version information:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Version:          %s\n", version.Version)
	fmt.Printf("Build Date:       %s\n", version.BuildDate)
	fmt.Printf("Git Commit:       %s\n", version.GitCommit)
	fmt.Printf("Git Tag:          %s\n", version.GitTag)
	fmt.Printf("Git Tree State:   %s\n", version.GitTreeState)
	fmt.Printf("Go Version:       %s\n", version.GoVersion)
	fmt.Printf("Compiler:         %s\n", version.Compiler)
	fmt.Printf("Platform:         %s\n", version.Platform)

	if *verbose {
		fmt.Printf("Kustomize:        %s\n", version.KustomizeVersion)
		fmt.Printf("Helm:             %s\n", version.HelmVersion)
		fmt.Printf("Kubectl:          %s\n", version.KubectlVersion)
		fmt.Printf("Jsonnet:          %s\n", version.JsonnetVersion)
		if version.ExtraBuildInfo != "" {
			fmt.Printf("Extra Build Info: %s\n", version.ExtraBuildInfo)
		}
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// createGRPCConnection creates a gRPC client connection with the specified options
func createGRPCConnection(server string, plaintext, insecureSkipVerify bool) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	// Configure TLS/credentials
	if plaintext {
		// Use insecure credentials (no TLS)
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// Use TLS
		tlsConfig := &tls.Config{
			InsecureSkipVerify: insecureSkipVerify,
		}
		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	// Set max message size (200MB default in ArgoCD)
	opts = append(opts,
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(200*1024*1024),
			grpc.MaxCallSendMsgSize(200*1024*1024),
		),
	)

	// Create connection
	conn, err := grpc.NewClient(server, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	return conn, nil
}
