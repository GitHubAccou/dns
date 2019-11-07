package main
import(
	"fmt"
)
func main(){

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
	
}