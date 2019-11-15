package main
import(
	"fmt"
	"strings"
)
func main(){
	host:="www.abcd.efg.higklm.com"
	name:=[]byte{}
	hostParts:=strings.Split(host,`.`)
	for _,part:=range hostParts{
		name=append(name,byte(len(part)))
		name=append(name,[]byte(part)...)
	}
	name=append(name,byte(0))

	queryMsg=MSG{
		Header:HeaderSection{
			ID:uint16(12345678),
			QR:true,
			Opcode:StandardQuery,
			AA:false,
			TC:false,
			RD:false,
			RA:false,
			Z:false,
			RCODE:uint8(0),
			QDCOUNT:uint16(1),
			ANCOUNT:uint16(0),
			NSCOUNT:uint16(0),
			ARCOUNT:uint16(0),
		},
		Question:QuerySection{
			QNAME:name,
			QTYPE:A,
			QCLASS:uint16(1)
		}
	}
}

type HeaderSection struct{
	ID uint16,
	QR bool,
	Opcode Op,
	AA bool,
	TC bool,
	RD bool,
	RA bool,
	Z uint8,
	RCODE ResponseCode,
	QDCOUNT uint16,
	ANCOUNT uint16,
	NSCOUNT uint16,
	ARCOUNT uint16
}

func (header *HeaderSection)rawBytes()[]byte{
	res:=[]byte{}
	res=
}


type ResponseCode uint8
const(
	var NoError ResponseCode=iota
	var FormatError ResponseCode
	var ServerFailure ResponseCode
	var NameError ResponseCode
	var NotImplemented ResponseCode
	var Refused ResponseCode
)

type Op uint8
const(
	var StandardQuery=iota
	var InverseQuery
	var Status
)

type QuerySection struct{
	QNAME []byte
	QTYPE RRType
	QCLASS uint16
}



type RRType uint16

const(
	var _ RRType=iota
	var A RRType
	var NS 
	var MD
	var MF
	var CNAME
	var SOA
	var MB
	var MG
	var MR
	var NULL
	var WKS
	var PTR
	var HINFO
	var MINFO
	var MX
	var TXT
)

type RRClass uint16

const(
	_=iota
	var IN 
	var CS
	var CH
	var HS
)

type RR struct{
	NAME []byte
	TYPE RRType
	CLASS RRClass
	TTL uint16
	RDLENGTH uint16
	RDATA uint16
}

type MSG struct{
	Header HeaderSection
	Question QuerySection
	Answer []RR
	Authority []RR
	Additional []RR
}