package grpc

import (
	"app/src/app_settings"
	"net"

	ggrpc "google.golang.org/grpc"
)

type Server struct {
	server   *ggrpc.Server
	listener net.Listener
}

func NewServer(settings *app_settings.GrpcSettings) (*Server, error) {
	listener, err := net.Listen("tcp", settings.ListenAddress)
	if err != nil {
		return nil, err
	}

	return &Server{
		server:   ggrpc.NewServer(),
		listener: listener,
	}, nil
}

func (s *Server) GRPC() *ggrpc.Server {
	return s.server
}

func (s *Server) Address() string {
	return s.listener.Addr().String()
}

func (s *Server) Serve() error {
	return s.server.Serve(s.listener)
}

func (s *Server) GracefulStop() {
	s.server.GracefulStop()
}
