package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/work-acc/Wildberries-L2/dev11/internal/api"
	"github.com/work-acc/Wildberries-L2/dev11/internal/config"
	"github.com/work-acc/Wildberries-L2/dev11/internal/service"
)

func main() {
	data, err := os.ReadFile("./configs/config.json")
	if err != nil {
		log.Fatal(err)
	}

	cfg := new(config.Config)
	if err := json.Unmarshal(data, cfg); err != nil {
		log.Fatal(err)
	}

	service := service.New()
	router := api.New(*cfg, service)

	if err := router.Start(); err != nil {
		log.Fatal(err)
	}
}
