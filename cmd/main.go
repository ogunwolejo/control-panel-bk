package main

import (
	"control-panel-bk/internal"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	internal.ControlPanelServer()
}
