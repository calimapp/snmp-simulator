package snmp

import (
	"fmt"

	"github.com/gosnmp/gosnmp"
)

func buildPDU(oid string, typ string, value any) gosnmp.SnmpPDU {
	pdu := gosnmp.SnmpPDU{
		Name: oid,
	}
	switch typ {
	case "OCTET_STRING", "octetstring", "STRING", "string":
		pdu.Type = gosnmp.OctetString
		pdu.Value = value.(string)
	case "INTEGER", "integer", "INT", "int":
		pdu.Type = gosnmp.Integer
		pdu.Value = int(value.(uint64))
	case "OID", "oid", "OBJECT_IDENTIFER", "objectidentifier":
		pdu.Type = gosnmp.ObjectIdentifier
		pdu.Value = value.(string)
	case "TIMETICKS", "timeticks":
		pdu.Type = gosnmp.TimeTicks
		pdu.Value = uint32(value.(uint64))
	default:
		pdu.Type = gosnmp.UnknownType
		pdu.Value = fmt.Sprintf("%v", value)
	}
	return pdu
}
