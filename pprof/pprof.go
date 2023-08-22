package pprof

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/ExcitingFrog/go-core-common/provider"
	"github.com/sirupsen/logrus"
)

type PProf struct {
	provider.IProvider

	Config *Config
	server *http.Server
	addr   string
}

func NewPprof(config *Config) *PProf {
	if config == nil {
		config = NewConfig()
	}

	return &PProf{
		Config: config,
	}
}

func (p *PProf) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(p.Config.endpoint+"/", pprof.Index)
	mux.HandleFunc(p.Config.endpoint+"/cmdline", pprof.Cmdline)
	mux.HandleFunc(p.Config.endpoint+"/profile", pprof.Profile)
	mux.HandleFunc(p.Config.endpoint+"/symbol", pprof.Symbol)
	mux.HandleFunc(p.Config.endpoint+"/trace", pprof.Trace)

	p.addr = fmt.Sprintf(":%d", p.Config.port)
	p.server = &http.Server{Addr: p.addr, Handler: mux}

	logrus.Info("pprof server listen on ", p.addr)
	if err := p.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (p *PProf) Close() error {
	return p.server.Close()
}
