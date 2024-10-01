package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/LevisThors/traefik-auth-plugin/mapper"
	v1 "github.com/LevisThors/traefik-auth-plugin/pb"

	"github.com/http-wasm/http-wasm-guest-tinygo/handler"
	"github.com/http-wasm/http-wasm-guest-tinygo/handler/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
}

func init() {
	var config Config
	err := json.Unmarshal(handler.Host.GetConfig(), &config)
	if err != nil {
		handler.Host.Log(api.LogLevelError, fmt.Sprintf("Could not load config %v", err))
		os.Exit(1)
	}

	mw, err := New(config)
	if err != nil {
		handler.Host.Log(api.LogLevelError, fmt.Sprintf("Could not create middleware %v", err))
		os.Exit(1)
	}

	handler.HandleRequestFn = mw.handleRequest
}

// Config holds the plugin configuration.
type Config struct {
	AuthServiceURL string `json:"authServiceURL"`
}

// AuthMiddleware represents the plugin.
type AuthMiddleware struct {
	authClient *Client
}

// New creates a new AuthMiddleware plugin.
func New(config Config) (*AuthMiddleware, error) {
	client, err := NewAuthServiceClient(context.Background(), config.AuthServiceURL)
	if err != nil {
		return nil, fmt.Errorf("could not create auth service client: %w", err)
	}

	return &AuthMiddleware{
		authClient: client,
	}, nil
}

func (a *AuthMiddleware) handleRequest(req api.Request, resp api.Response) (next bool, reqCtx uint32) {
	authHeader, ok := req.Headers().Get("Authorization")
	if authHeader == "" || !ok {
		resp.SetStatusCode(http.StatusUnauthorized)
		resp.Body().Write([]byte("Authorization header missing"))
		return false, 0
	}

	token := strings.TrimPrefix(authHeader, "JWT ")
	if token == authHeader {
		resp.SetStatusCode(http.StatusUnauthorized)
		resp.Body().Write([]byte("Invalid Authorization header format"))
		return false, 0
	}

	payload, isValid, err := a.authClient.CheckToken(context.Background(), token)
	if err != nil {
		resp.SetStatusCode(http.StatusInternalServerError)
		resp.Body().Write([]byte("Error validating token"))
		return false, 0
	}

	if !isValid {
		resp.SetStatusCode(http.StatusUnauthorized)
		resp.Body().Write([]byte("Invalid token"))
		return false, 0
	}

	for key, value := range payload {
		req.Headers().Set(key, fmt.Sprintf("%v", value))
	}

	return true, 0
}

type Client struct {
	connection *grpc.ClientConn
	client     v1.AuthServiceClient
}

func NewAuthServiceClient(ctx context.Context, address string) (*Client, error) {
	connection, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("creating gRPC client: %w", err)
	}

	return &Client{
		connection: connection,
		client:     v1.NewAuthServiceClient(connection),
	}, nil
}

func (c *Client) CheckToken(ctx context.Context, token string) (map[string]interface{}, bool, error) {
	req := &v1.CheckTokenRequest{Token: token}
	resp, err := c.client.CheckToken(ctx, req)
	if err != nil {
		return nil, false, fmt.Errorf("error calling CheckToken: %w", err)
	}

	response, err := mapper.ToPayload(resp.GetPayload())
	if err != nil {
		return nil, false, fmt.Errorf("error mapping payload: %w", err)
	}

	isValid := resp.GetStatus() == v1.TokenStatus_VALID
	return response, isValid, nil
}
