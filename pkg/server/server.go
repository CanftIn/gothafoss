package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/CanftIn/gothafoss/pkg/im/config"
	"github.com/CanftIn/gothafoss/pkg/im/imhttp"
	"github.com/CanftIn/gothafoss/pkg/im/module"
	"github.com/CanftIn/gothafoss/pkg/im/register"
	"github.com/CanftIn/gothafoss/pkg/log"

	"github.com/unrolled/secure"
)

type Server struct {
	ctx      *config.Context
	imHttp   *imhttp.IMHttp
	addr     string
	sslAddr  string
	grpcAddr string
	log.TLog
}

func New(ctx *config.Context) *Server {
	imHttp := imhttp.New()
	imHttp.Use(imhttp.CORSMiddleware())
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
	s.imHttp.Any("/v1/ping", func(c *imhttp.Context) {
		c.ResponseOK()
	})

	s.imHttp.Any("/swagger/:module", func(c *imhttp.Context) {
		m := c.Param("module")
		module := register.GetModuleByName(m, s.ctx)
		if strings.TrimSpace(module.Swagger) == "" {
			c.Status(http.StatusNotFound)
			return
		}
		c.String(http.StatusOK, module.Swagger)

	})

	if len(addr) != 0 {
		if sslAddr != "" {
			go func() {
				err := s.imHttp.Run(addr...)
				if err != nil {
					panic(err)
				}
			}()
		} else {
			err := s.imHttp.Run(addr...)
			if err != nil {
				return err
			}
		}

	}

	// https 服务
	if sslAddr != "" {
		s.imHttp.Use(TlsHandler(sslAddr))
		currDir, _ := os.Getwd()
		return s.imHttp.RunTLS(sslAddr, currDir+"/assets/ssl/ssl.pem", currDir+"/assets/ssl/ssl.key")
	}

	return nil
}

func (s *Server) Start() error {
	go func() {
		err := s.run(s.sslAddr, s.addr)
		if err != nil {
			panic(err)
		}
	}()

	err := module.Start(s.ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() error {

	return module.Stop(s.ctx)
}

func TlsHandler(sslAddr string) imhttp.HandlerFunc {
	return func(c *imhttp.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     sslAddr,
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}

// GetRoute 获取路由
func (s *Server) GetRoute() *imhttp.IMHttp {
	return s.imHttp
}
