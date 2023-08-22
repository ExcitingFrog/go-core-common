package provider

import (
	"os"
	"os/signal"
)

type IProvider interface {
	Run() error
	Close() error
}

type Providers struct {
	providers []IProvider
}

func (s *Providers) handleInterrupt() {
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		s.Close()
		close(cleanupDone)
	}()
	<-cleanupDone
}

func NewProviders() *Providers {
	return &Providers{
		providers: make([]IProvider, 0),
	}
}

func (s *Providers) AddProvider(provider IProvider) {
	s.providers = append(s.providers, provider)
}

func (s *Providers) Run() {
	for _, provider := range s.providers {
		go provider.Run()
	}
	s.handleInterrupt()
}

func (s *Providers) Close() {
	for _, provider := range s.providers {
		if err := provider.Close(); err != nil {
			panic(err)
		}
	}
}
