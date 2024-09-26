package main

import (
	"github.com/reinaldo-silva/savina-stock/config"
	"github.com/reinaldo-silva/savina-stock/internal/app"
)

func main() {

	cfg := config.LoadConfig()

	a := app.App{}

	a.Initialize(cfg)

	a.Run(cfg)
}
