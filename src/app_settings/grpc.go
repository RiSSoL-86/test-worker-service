package app_settings

import (
	"errors"
	"os"
	"strings"
)

type GrpcSettings struct {
	ListenAddress string
}

func NewGrpcSettings() (*GrpcSettings, error) {
	listenAddress := strings.TrimSpace(os.Getenv("GRPC_LISTEN_ADDRESS"))
	if listenAddress == "" {
		return nil, errors.New("GRPC_LISTEN_ADDRESS is not set")
	}

	return &GrpcSettings{ListenAddress: listenAddress}, nil
}
