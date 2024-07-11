package main

import (
	"net/http"

	"github.com/kmyokoyama/go-template/internal/components"
	"github.com/kmyokoyama/go-template/internal/diplomat/inbound"
	"github.com/kmyokoyama/go-template/internal/diplomat/outbound"
	"go.uber.org/fx"
)

type App struct {
	Server *http.Server
}

func NewApp(srv *http.Server) *App {
	return &App{Server: srv}
}

func main() {
	fx.New(
		fx.Provide(NewApp),
		fx.Provide(components.NewComponents),
		fx.Provide(inbound.NewRouter),
		fx.Provide(inbound.NewHttpServer),
		fx.Provide(outbound.NewDatabase),
		fx.Provide(components.NewLogger),
		fx.Invoke(func(*App) {}),
	).Run()
}
