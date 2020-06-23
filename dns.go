package main
import(
	"fmt"
	"net"
	"strings"
)
func main(){
	conn,_:=net.ListenUDP("udp",&net.UDPAddr{Port:53})
	for{
		buf:=make([]byte,1024*8)
		len,remote,_:=conn.ReadFrom(buf)	
		go ServeDNS(conn,buf[:len],remote)
	}
}
func ServeDNS(conn *net.UDPConn,data []byte,client net.Addr){
	if dns,err:=packet2DNS(data);err==nil{
		fmt.Println("request "+queryName2Host(dns.Query.Name))
		if dns.Query.QueryType==QueryType_PTR{
			fmt.Println("----------------------------PTR------------------------")
			rData:=makeupPTRResponseData(dns.ID,queryName2Host(dns.Query.Name),[]byte{192,168,20,48})
			conn.WriteTo(rData,client)
		}else{
			rData:=makeupResponseData(dns.ID,queryName2Host(dns.Query.Name),[]byte{192,168,20,48})
			conn.WriteTo(rData,client)
		}
	}else{
		fmt.Println(err.Error())
	}
}
func queryName2Host(data []byte)string{
	res:=make([]byte,0)
	for i:=0;true;{
		len:=int(data[i])
		if len==0{
			break
		}
		res=append(res,data[i+1:i+1+len]...)
		res=append(res,byte('.'))
		i+=len+1
	}
	return string(res[:len(res)-1])
}

func packet2DNS(data []byte)(dns DNS,err error){
	dns.ID=uint16(data[0])<<8+uint16(data[1])
	qr:=data[2]>>7
	if qr==1{
		dns.Flag.QR=true
	}
	opcode:=data[2]<<1>>4
	dns.Flag.OpCode=OpCode(opcode)
	aa:=data[2]<<5>>7
	if aa==1{
		dns.Flag.AA=true
	}
	tc:=data[2]<<6>>7
	if tc==1{
		dns.Flag.TC=true
	}
	rd:=data[2]<<7>>7
	if rd==1{
		dns.Flag.RD=true
	}
        ra:=data[3]>>7
	if ra==1{
		dns.Flag.RA=true
	}
	rcode:=data[3]&0x0F
	dns.Flag.RCode=RCode(rcode)
	dns.QC=uint16(data[4])<<8+uint16(data[5])
	dns.RC=uint16(data[6])<<8+uint16(data[6])
	dns.AC=uint16(data[8])<<8+uint16(data[9])
	dns.OC=uint16(data[10])<<8+uint16(data[11])
	i:=12
	for ;i<len(data);{
		len:=int(data[i])
		if len==0{
			break
		}else{
			i+=len+1
		}
	}
	dns.Query.Name=data[12:i+1]
	dns.Query.QueryType=QueryType(uint16(data[i+1])<<8+uint16(data[i+2]))
	return 
}
func QueryIP(domain string){
	labels:=strings.Split(domain,".")
	nameD:=make([]byte,0)
	for _,label:=range labels{
		nameD=append(nameD,byte(len(label)))
		nameD=append(nameD,[]byte(label)...)
	}
	dns:=DNS{
		ID:uint16(1),
		Flag:Flag{
				QR:false,
				OpCode:OpCode(0),
				AA:false,
				TC:false,
				RD:true,
				RA:false,
				RCode:RCode(0),
			},
		QC:1,
		RC:0,
		AC:0,
		OC:0,
		Query:Query{
			Name:nameD,
			QueryType:QueryType_A,
			QueryClass:QueryClass_IN,
		},
	}
	if conn,err:=net.DialUDP("udp",nil,&net.UDPAddr{IP:net.IPv4(192,168,20,8),Port:53});err==nil{
		defer conn.Close()
		conn.Write(dns.Data())
	}else{
		fmt.Println(err.Error())
	}
}
func makeupResponseData(id uint16,host string,ip []byte)[]byte{
	labels:=strings.Split(host,".")
	nameD:=make([]byte,0)
	for _,label:=range labels{
		nameD=append(nameD,byte(len(label)))
		nameD=append(nameD,[]byte(label)...)
	}
	nameD=append(nameD,0)
	dns:=DNS{
		ID:id,
		Flag:Flag{
				QR:true,
				OpCode:OpCode(0),
				AA:false,
				TC:false,
				RD:true,
				RA:true,
				RCode:RCode(0),
			},
		QC:1,
		RC:1,
		AC:0,
		OC:0,
		Query:Query{
			Name:nameD,
			QueryType:QueryType_A,
			QueryClass:QueryClass_IN,
		},
		Answer:ResourceRecord{
				Name :[]byte{byte(0xc0),byte(0x0c)},
				QueryType: QueryType_A,
				QueryClass: QueryClass_IN,
				TTL: uint32(60*60*24),
				RRLen: uint16(4),
				DataS:ip,
		},
	}
	return dns.Data()
}

func makeupPTRResponseData(id uint16,host string,ip []byte)[]byte{
	labels:=strings.Split(host,".")
	nameD:=make([]byte,0)
	for _,label:=range labels{
		nameD=append(nameD,byte(len(label)))
		nameD=append(nameD,[]byte(label)...)
	}
	nameD=append(nameD,0)
	hostName:=make([]byte,0)
	hostName=append(hostName,byte(6))
	hostName=append(hostName,[]byte("easter")...)
	hostName=append(hostName,byte(0))
	dns:=DNS{
		ID:id,
		Flag:Flag{
				QR:true,
				OpCode:OpCode(0),
				AA:false,
				TC:false,
				RD:true,
				RA:true,
				RCode:RCode(0),
			},
		QC:1,
		RC:1,
		AC:0,
		OC:0,
		Query:Query{
			Name:nameD,
			QueryType:QueryType_PTR,
			QueryClass:QueryClass_IN,
		},
		Answer:ResourceRecord{
				Name :[]byte{byte(0xc0),byte(0x0c)},
				QueryType: QueryType_PTR,
				QueryClass: QueryClass_IN,
				TTL: uint32(60*60*24),
				RRLen: uint16(len(nameD)),
				DataS:hostName,
		},
	}
	return dns.Data()
}

type DNS struct{
	ID uint16
	Flag Flag
	QC uint16
	RC uint16
	AC uint16
	OC uint16
	Query Query
	Answer ResourceRecord
	Auth ResourceRecord
	Other ResourceRecord
}

func (dns *DNS)Data()[]byte{
	data:=make([]byte,0)
	data=append(data,byte(dns.ID>>8))
	data=append(data,byte(dns.ID<<8>>8))
	data=append(data,dns.Flag.Data()...)
	data=append(data,byte(dns.QC>>8))
	data=append(data,byte(dns.QC<<8>>8))
	data=append(data,byte(dns.RC>>8))
	data=append(data,byte(dns.RC<<8>>8))
	data=append(data,byte(dns.AC>>8))
	data=append(data,byte(dns.AC<<8>>8))
	data=append(data,byte(dns.OC>>8))
	data=append(data,byte(dns.OC<<8>>8))
	data=append(data,dns.Query.Data()...)
	data=append(data,dns.Answer.Data()...)
//	data=append(data,dns.Auth.Data()...)
//	data=append(data,dns.Other.Data()...)
	return data;
}
type Flag struct{
	QR bool 
	OpCode OpCode
	AA bool
	TC bool
	RD bool
	RA bool
	RCode RCode
}
func (flag *Flag)Data()[]byte{
	var x uint8=0
	var y uint8=0
	if flag.QR {
		x+=(1<<7)
	}else{
		x+=(0<<7)
	}
	x+=uint8(flag.OpCode<<3)&0x78
	if flag.AA {
		x+=(1<<2)
	}else{
		x+=(0<<2)
	}
	if flag.TC {
		x+=(1<<1)
	}else{
		x+=(0<<1)
	}
	if flag.RD {
		x+=1
	}else{
		x+=0
	}
	if flag.RA {
		y+=(1<<7)
	}else{
		y+=(0<<7)
	}
	y+=byte(flag.RCode<<4>>4)
	return []byte{byte(x),byte(y)}
}
type OpCode uint8
type RCode uint8
type Query struct{
	Name []byte
	QueryType QueryType
	QueryClass QueryClass
}
func (query *Query)Data()[]byte{
	data:=make([]byte,0)
	data=append(data,query.Name...)
	data=append(data,byte(query.QueryType>>8))
	data=append(data,byte(query.QueryType<<8>>8))
	data=append(data,byte(query.QueryClass>>8))
	data=append(data,byte(query.QueryClass<<8>>8))
	return data;
}
type ResourceRecord struct{
	Name [] byte
	QueryType QueryType
	QueryClass QueryClass
	TTL uint32
	RRLen uint16
	DataS []byte

}
func (resourceRecord *ResourceRecord)Data()[]byte{
	data:=make([]byte,0)
	data=append(data,resourceRecord.Name...)
	data=append(data,byte(resourceRecord.QueryType>>8))
	data=append(data,byte(resourceRecord.QueryType<<8>>8))
	data=append(data,byte(resourceRecord.QueryClass>>8))
	data=append(data,byte(resourceRecord.QueryClass<<8>>8))
	data=append(data,byte(resourceRecord.TTL>>24))
	data=append(data,byte(resourceRecord.TTL<<8>>24))
	data=append(data,byte(resourceRecord.TTL<<16>>24))
	data=append(data,byte(resourceRecord.TTL<<24>>24))
	data=append(data,byte(resourceRecord.RRLen>>8))
	data=append(data,byte(resourceRecord.RRLen<<8>>8))
	data=append(data,resourceRecord.DataS...)
	return data;
}

type QueryType uint16 
type QueryClass uint16

const QueryType_A QueryType=1
const QueryType_NS QueryType=2
const QueryType_CNAME QueryType=5
const QueryType_PTR QueryType=12
const QueryType_HINFO QueryType=13
const QueryType_MX QueryType=15
const QueryType_AXFR QueryType=252
const QueryType_ANY QueryType=255
const QueryClass_IN QueryClass=1
