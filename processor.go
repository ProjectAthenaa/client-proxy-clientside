package main

import (
	"context"
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/fasttls"
	"github.com/ProjectAthenaa/sonic-core/fasttls/tls"
	client_proxy "github.com/ProjectAthenaa/sonic-core/protos/clientProxy"
	"github.com/prometheus/common/log"
)

type server struct {
	stream client_proxy.Proxy_RegisterClient
	ctx    context.Context
	cancel context.CancelFunc
	client *fasttls.Client
}

func NewServer(stream client_proxy.Proxy_RegisterClient) *server {
	s := &server{
		stream: stream,
		client: fasttls.NewClient(tls.HelloChrome_91, nil),
	}
	s.ctx, s.cancel = context.WithCancel(stream.Context())

	return s
}

func (s *server) listen() {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			newReq, err := s.stream.Recv()
			if err != nil {
				log.Error("receive request: ", err)
				continue
			}
			log.Info("Received Request")


			if _, ok := newReq.Headers["STOP"]; ok {
				s.cancel()
				continue
			}

			go s.process(newReq)
		}
	}
}

func (s *server) process(clientReq *client_proxy.Request) {
	req := convertToRequest(clientReq)

	resp, err := s.client.DoCtx(s.ctx, req)
	if err != nil {
		_ = s.stream.Send(&client_proxy.Response{Headers: map[string]string{"ERROR": fmt.Sprint(err)}})
		return
	}

	clientResp := convertToResponse(resp, clientReq.TaskID)

	_ = s.stream.Send(clientResp)
}

func (s *server) Start() {
	log.Info("Proxy started")
	s.listen()
}
