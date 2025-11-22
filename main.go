package main

import (
	"fmt"
	"log"

	"github.com/calimapp/snmp-simulator/internal/config"
	"github.com/calimapp/snmp-simulator/internal/snmp"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("=> %+v\n", cfg)
	snmp.NewAgent(cfg)
}
