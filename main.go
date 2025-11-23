package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/calimapp/snmp-simulator/internal/config"
	"github.com/calimapp/snmp-simulator/internal/snmp"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal(err.Error())
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	if cfg.Polling.Enabled {
		go snmp.StartPollingSimulator(cfg)
	}
	if cfg.Trap.Enabled {
		if err := snmp.StartTrapSimulator(cfg); err != nil {
			log.Fatal(err)
		}
	}
	<-stop
}
