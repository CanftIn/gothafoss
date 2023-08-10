package server

import (
	"context"

	"github.com/CanftIn/gothafoss/lib/im/http"
	"github.com/CanftIn/gothafoss/lib/log"
)

type Server struct {
	ctx      *context.Context
	imHttp   *http.IMHttp
	addr     string
	sslAddr  string
	grpcAddr string
	log.TLog
}

func New(ctx *context.Context) *Server {
	imHttp := http.New()
	imHttp.Use(http.CORSMiddleware())
	s := &Server{
		ctx:      ctx,
		imHttp:   imHttp,
		addr:     ctx.GetConfig().Addr,
		sslAddr:  ctx.GetConfig().SSLAddr,
		grpcAddr: ctx.GetConfig().GRPCAddr,
	}
	return s
}

func (s *Server) Init() error {
	return nil
}

func (s *Server) run(sslAddr string, addr ...string) error {
	s.imHttp.Static("/web", "./asserts/web")
	s.imHttp.Any("/v1/ping", func(c *http.Context) {
		c.ResponseOK()
	})

}
