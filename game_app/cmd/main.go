package main

import (
	"github.com/coda-payments/game_app/internal/config"
	"github.com/coda-payments/game_app/internal/routes"
)

func main() {
	routes.Launch(config.GetPort())
}
