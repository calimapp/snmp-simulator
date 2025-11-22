package config

import "github.com/gosnmp/gosnmp"

type Security struct {
	V2c V2cSecurity `yaml:"v2c" required:"true"`
	V3  V3Security  `yaml:"v3" required:"true"`
}
type V2cSecurity struct {
	Enabled   bool
	Community string
}

type V3Security struct {
	Enabled      bool
	Username     string `yaml:"username"`
	AuthProtocol string `yaml:"authProtocol"`
	AuthPassword string `yaml:"authPassword"`
	PrivProtocol string `yaml:"privProtocol"`
	PrivPassword string `yaml:"privPassword"`
}

func (s *V3Security) GetAuthProtocol() gosnmp.SnmpV3AuthProtocol {
	switch s.AuthProtocol {
	case "md5", "MD5":
		return gosnmp.MD5
	case "sha", "SHA":
		return gosnmp.SHA
	default:
		return gosnmp.NoAuth
	}
}

func (s *V3Security) GetPrivProtocol() gosnmp.SnmpV3PrivProtocol {
	switch s.PrivProtocol {
	case "des", "DES":
		return gosnmp.DES
	case "aes", "AES":
		return gosnmp.AES
	default:
		return gosnmp.NoPriv
	}
}
