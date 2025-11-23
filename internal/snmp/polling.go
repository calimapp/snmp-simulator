package snmp

import (
	"fmt"
	"log"

	"github.com/calimapp/snmp-simulator/internal/config"
	"github.com/gosnmp/gosnmp"
	"github.com/slayercat/GoSNMPServer"
)

func StartPollingSimulator(cfg *config.Config) {
	master := GoSNMPServer.MasterAgent{
		Logger: GoSNMPServer.NewDefaultLogger(),
		SecurityConfig: GoSNMPServer.SecurityConfig{
			AuthoritativeEngineID:    GoSNMPServer.DefaultAuthoritativeEngineID(),
			AuthoritativeEngineBoots: 1,
			Users: []gosnmp.UsmSecurityParameters{
				{
					UserName:                 cfg.Security.V3.Username,
					AuthenticationProtocol:   cfg.Security.V3.GetAuthProtocol(),
					AuthenticationPassphrase: cfg.Security.V3.AuthPassword,
					PrivacyProtocol:          cfg.Security.V3.GetPrivProtocol(),
					PrivacyPassphrase:        cfg.Security.V3.PrivPassword,
				},
			},
		},
		SubAgents: []*GoSNMPServer.SubAgent{
			{
				CommunityIDs: []string{cfg.Security.V2c.Community},
				OIDs:         cfg.Polling.GetPDUs(),
			},
			{
				CommunityIDs: []string{""},
				OIDs:         cfg.Polling.GetPDUs(),
			},
		},
	}

	server := GoSNMPServer.NewSNMPServer(master)
	if err := server.ListenUDP("udp", fmt.Sprintf("0.0.0.0:%d", cfg.Polling.Port)); err != nil {
		log.Fatalf("error listening: %v", err)
	}
	server.ServeForever()
}
