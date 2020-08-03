package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http"
)

func main() {
	log := logrus.New()
	log.Infoln("Starting Command-Handler process...")

	cfg, err := configuration.LoadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("unable to load app configuration: %s", err.Error()))
	}

	// Starting gateway http API
	api := http.NewApi(log)
	api.Start("0.0.0.0", cfg.API.Port)
}
