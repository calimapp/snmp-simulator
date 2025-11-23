package snmp

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/calimapp/snmp-simulator/internal/config"
	"github.com/gosnmp/gosnmp"
)

func StartTrapSimulator(cfg *config.Config) error {
	g := &gosnmp.GoSNMP{
		Target:    cfg.Trap.TargetHost,
		Port:      uint16(cfg.Trap.TargetPort),
		Transport: "udp",
		Timeout:   2 * time.Second,
		Retries:   1,
	}
	if cfg.Security.V2c.Enabled {
		g.Version = gosnmp.Version2c
		g.Community = cfg.Security.V2c.Community
	} else if cfg.Security.V3.Enabled {
		g.Version = gosnmp.Version3
		g.SecurityModel = gosnmp.UserSecurityModel
		g.MsgFlags = gosnmp.AuthPriv
		g.SecurityParameters = &gosnmp.UsmSecurityParameters{
			UserName:                 cfg.Security.V3.Username,
			AuthenticationProtocol:   cfg.Security.V3.GetAuthProtocol(),
			AuthenticationPassphrase: cfg.Security.V3.AuthPassword,
			PrivacyProtocol:          cfg.Security.V3.GetPrivProtocol(),
			PrivacyPassphrase:        cfg.Security.V3.PrivPassword,
		}
	}
	if err := g.Connect(); err != nil {
		return err
	}

	for _, trap := range cfg.Trap.Traps {
		go scheduleTrapEvent(g, trap)
		slog.Info(fmt.Sprintf("Schedule trap every %s -> %s:%d", trap.Interval, cfg.Trap.TargetHost, cfg.Trap.TargetPort))
	}
	return nil
}

const sysUpTimeTrapOID string = "1.3.6.1.2.1.1.3.0"
const snmpTrapOID string = "1.3.6.1.6.3.1.1.4.1.0"

func buildTrap(trap config.Trap) gosnmp.SnmpTrap {
	snmpVars := []gosnmp.SnmpPDU{
		{Name: sysUpTimeTrapOID, Type: gosnmp.TimeTicks, Value: uint32(12345)},
		{Name: snmpTrapOID, Type: gosnmp.ObjectIdentifier, Value: trap.TrapOID},
	}
	for _, snmpVar := range trap.VarBinds {
		snmpVars = append(snmpVars, buildPDU(snmpVar.Oid, snmpVar.Type, snmpVar.Value))
	}
	return gosnmp.SnmpTrap{Variables: snmpVars}
}

func scheduleTrapEvent(g *gosnmp.GoSNMP, trap config.Trap) {
	interval, _ := time.ParseDuration(trap.Interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	snmpTrap := buildTrap(trap)
	for range ticker.C {
		_, err := g.SendTrap(snmpTrap)
		if err != nil {
			fmt.Printf("SendTrap error: %v\n", err)
		} else {
			slog.Debug(fmt.Sprintf("Send Trap %s", trap.TrapOID))
		}
	}
}
