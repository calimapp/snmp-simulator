package config

import (
	"fmt"

	"github.com/gosnmp/gosnmp"
	"github.com/slayercat/GoSNMPServer"
)

type Config struct {
	Security Security      `yaml:"security"`
	Polling  PollingConfig `yaml:"polling"`
}

type PollingConfig struct {
	Port    int     `yaml:"port"`
	Scalars []OID   `yaml:"scalars"`
	Tables  []Table `yaml:"tables"`
}

type Table struct {
	Name    string              `yaml:"name"`
	RootOID string              `yaml:"rootOID"`
	Columns map[int]TableColumn `yaml:"columns"`
}

func (t *Table) buildTablePDUs() []*GoSNMPServer.PDUValueControlItem {
	var pdus []*GoSNMPServer.PDUValueControlItem
	for columnID, column := range t.Columns {
		for rowID, value := range column.Rows {
			oid := fmt.Sprintf("%s.%d.%d", t.RootOID, columnID, rowID)
			alias := fmt.Sprintf("%s.%s.%d", t.Name, column.Name, rowID)
			pdus = append(pdus, buildPDU(oid, alias, column.Type, value))
		}
	}
	return pdus
}

type TableColumn struct {
	Name string      `yaml:"name"`
	Type string      `yaml:"type"`
	Rows map[int]any `yaml:"rows"`
}

func (p PollingConfig) GetPDUs() []*GoSNMPServer.PDUValueControlItem {
	var pdus []*GoSNMPServer.PDUValueControlItem
	for _, oid := range p.Scalars {
		pdus = append(pdus, oid.toAgentPDU())
	}
	for _, table := range p.Tables {
		pdus = append(pdus, table.buildTablePDUs()...)
	}
	return pdus
}

type OID struct {
	Oid   string `yaml:"oid"`
	Alias string `yaml:"alias"`
	Value any    `yaml:"value"`
	Type  string `yaml:"type"`
}

func (oid *OID) toAgentPDU() *GoSNMPServer.PDUValueControlItem {
	return buildPDU(oid.Oid, oid.Alias, oid.Type, oid.Value)
}

func buildPDU(oid string, alias string, typ string, value any) *GoSNMPServer.PDUValueControlItem {
	pdu := &GoSNMPServer.PDUValueControlItem{
		OID:      oid,
		Document: alias,
	}
	switch typ {
	case "OCTET_STRING", "octetstring", "STRING", "string":
		pdu.Type = gosnmp.OctetString
		pdu.OnGet = func() (v any, err error) {
			return GoSNMPServer.Asn1OctetStringWrap(value.(string)), nil
		}
	case "INTEGER", "integer", "INT", "int":
		pdu.Type = gosnmp.Integer
		pdu.OnGet = func() (v any, err error) {
			return GoSNMPServer.Asn1IntegerWrap(int(value.(uint64))), nil
		}
	case "OID", "oid", "OBJECT_IDENTIFER", "objectidentifier":
		pdu.Type = gosnmp.ObjectIdentifier
		pdu.OnGet = func() (v any, err error) {
			return GoSNMPServer.Asn1OctetStringWrap(value.(string)), nil
		}
	case "TIMETICKS", "timeticks":
		pdu.Type = gosnmp.TimeTicks
		pdu.OnGet = func() (v any, err error) {
			return GoSNMPServer.Asn1TimeTicksWrap(uint32(value.(uint64))), nil
		}
	default:
		pdu.Type = gosnmp.UnknownType
		pdu.OnGet = func() (v any, err error) {
			return GoSNMPServer.Asn1OctetStringWrap(fmt.Sprintf("%v", value)), nil
		}
	}
	return pdu
}
