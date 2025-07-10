package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Role1776/agent/internal/app"
	"github.com/Role1776/agent/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}

	defer client.CloseIdleConnections()

	agent := app.NewApp(cfg, client)
	agent.Run()
}
