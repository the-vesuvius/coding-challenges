package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/spf13/cobra"
)

const (
	headerSizeBytes = 12
)

func Parse(buffer []byte) Message {
	header := ParseHeader(buffer[0:headerSizeBytes])
	q, a := parseTheRest(buffer[headerSizeBytes:])
	return Message{
		Header:   header,
		Question: q,
		Answer:   a,
	}
}

func parseTheRest(buffer []byte) (Question, ResourceRecord) {
	fmt.Println(buffer)
	q := Question{}
	a := ResourceRecord{}
	i := 0
	nameParts := []string{}
	for buffer[i] != 0 {
		n := int(buffer[i])
		nameParts = append(nameParts, string(buffer[i+1:i+1+n]))
		i = i + 1 + n
	}
	i++
	q.Qname = strings.Join(nameParts, ".")
	q.Qtype = binary.BigEndian.Uint16(buffer[i : i+2])

	i += 2
	q.Qclass = binary.BigEndian.Uint16(buffer[i : i+2])
	i += 2

	fmt.Printf("%08b\n", buffer[i])
	fmt.Printf("%08b\n", buffer[i+1])
	fmt.Println(buffer[i:])

	return q, a
}

type Message struct {
	Header   Header
	Question Question
	Answer   ResourceRecord
}

func (m Message) Encode() []byte {
	result := []byte{}
	result = append(result, m.Header.Encode()...)
	result = append(result, m.Question.Encode()...)
	return result
}

type Header struct {
	Id uint16
	Qr bool // 1 bit. 0 query, 1 response

	// 4 bit
	//0               a standard query (QUERY)
	//1               an inverse query (IQUERY)
	//2               a server status request (STATUS)
	//3-15            reserved for future use
	Opcode uint8

	// 1 bit
	//Authoritative Answer - this bit is valid in responses,
	//and specifies that the responding name server is an
	//authority for the domain name in question section.
	//Note that the contents of the answer section may have
	//multiple owner names because of aliases.  The AA bit
	//corresponds to the name which matches the query name, or
	//the first owner name in the answer section.
	Aa bool

	// 1 bit
	//TrunCation - specifies that this message was truncated
	//due to length greater than that permitted on the
	//transmission channel.
	Tc bool

	// 1 bit
	//Recursion Desired - this bit may be set in a query and
	//is copied into the response.  If RD is set, it directs
	//the name server to pursue the query recursively.
	//Recursive query support is optional.
	Rd bool

	// 1 bit
	//Recursion Available - this be is set or cleared in a
	//response, and denotes whether recursive query support is
	//available in the name server.
	Ra bool

	// 3 bit
	//Reserved for future use.  Must be zero in all queries
	//and responses.
	Z uint8

	// 4 bit
	//Response code - this 4 bit field is set as part of
	//responses.  The values have the following
	//interpretation:
	//
	//0               No error condition
	//
	//1               Format error - The name server was
	//unable to interpret the query.
	//
	//2               Server failure - The name server was
	//unable to process this query due to a
	//problem with the name server.
	//
	//3               Name Error - Meaningful only for
	//responses from an authoritative name
	//server, this code signifies that the
	//domain name referenced in the query does
	//not exist.
	//
	//4               Not Implemented - The name server does
	//not support the requested kind of query.
	//
	//5               Refused - The name server refuses to
	//perform the specified operation for
	//policy reasons.  For example, a name
	//server may not wish to provide the
	//information to the particular requester,
	//or a name server may not wish to perform
	//a particular operation (e.g., zone
	//transfer) for particular data.
	//
	//6-15            Reserved for future use.
	Rcode uint8

	//an unsigned 16 bit integer specifying the number of
	//entries in the question section.
	QdCount uint16

	//an unsigned 16 bit integer specifying the number of
	//entries in the question section.
	AnCount uint16

	//an unsigned 16 bit integer specifying the number of name
	//server resource records in the authority records
	//section.
	NsCount uint16

	//an unsigned 16 bit integer specifying the number of
	//resource records in the additional records section.
	ArCount uint16
}

func (h Header) Encode() []byte {
	result := make([]byte, 0, 12)

	// ID
	result = binary.BigEndian.AppendUint16(result, h.Id)

	byt := byte(0)
	// QR
	if h.Qr {
		byt ^= 0b10000000
	}
	// Opcode
	byt ^= h.Opcode << 4

	if h.Aa {
		byt ^= 0b00000100
	}

	if h.Tc {
		byt ^= 0b00000010
	}

	if h.Rd {
		byt ^= 0b00000001
	}

	result = append(result, byt)
	byt = byte(0)
	if h.Ra {
		byt ^= 0b10000000
	}

	byt ^= h.Z << 4
	byt ^= h.Rcode
	result = append(result, byt)

	result = binary.BigEndian.AppendUint16(result, h.QdCount)
	result = binary.BigEndian.AppendUint16(result, h.AnCount)
	result = binary.BigEndian.AppendUint16(result, h.NsCount)
	result = binary.BigEndian.AppendUint16(result, h.ArCount)

	return result
}

func ParseHeader(buffer []byte) Header {
	h := Header{}
	h.Id = binary.BigEndian.Uint16(buffer[0:2])
	byt := buffer[2]

	h.Qr = byt&0b10000000 > 0

	h.Opcode = (byt >> 4) & 0b00001111
	h.Aa = byt&0b00000100 > 0
	h.Tc = byt&0b00000010 > 0
	h.Rd = byt&0b00000001 > 0

	byt = buffer[3]
	h.Ra = byt&0b10000000 > 0

	h.Z = (byt >> 4) & 0b00000111
	h.Rcode = byt & 0b00001111

	h.QdCount = binary.BigEndian.Uint16(buffer[4:6])
	h.AnCount = binary.BigEndian.Uint16(buffer[6:8])
	h.NsCount = binary.BigEndian.Uint16(buffer[8:10])
	h.ArCount = binary.BigEndian.Uint16(buffer[10:12])

	return h
}

type Question struct {
	//a domain name represented as a sequence of labels, where
	//each label consists of a length octet followed by that
	//number of octets.  The domain name terminates with the
	//zero length octet for the null label of the root.  Note
	//that this field may be an odd number of octets; no
	//padding is used.
	Qname string

	//a two octet code which specifies the type of the query.
	//The values for this field include all codes valid for a
	//TYPE field, together with some more general codes which
	//can match more than one type of RR.
	Qtype uint16

	//a two octet code that specifies the class of the query.
	//For example, the QCLASS field is IN for the Internet.
	Qclass uint16
}

func (q Question) Encode() []byte {
	result := make([]byte, 0)

	result = append(result, encodeName(q.Qname)...)

	result = binary.BigEndian.AppendUint16(result, q.Qtype)
	result = binary.BigEndian.AppendUint16(result, q.Qclass)

	return result
}

type ResourceRecord struct {
	//a domain name to which this resource record pertains.
	Name string

	//two octets containing one of the RR type codes.  This
	//field specifies the meaning of the data in the RDATA
	//field.
	Type uint16

	//two octets which specify the class of the data in the
	//RDATA field.
	Class uint16

	//a 32 bit unsigned integer that specifies the time
	//interval (in seconds) that the resource record may be
	//cached before it should be discarded.  Zero values are
	//interpreted to mean that the RR can only be used for the
	//transaction in progress, and should not be cached.
	Ttl uint32

	//an unsigned 16 bit integer that specifies the length in
	//octets of the RDATA field.
	RdLength uint16

	//a variable length string of octets that describes the
	//resource.  The format of this information varies
	//according to the TYPE and CLASS of the resource record.
	//For example, the if the TYPE is A and the CLASS is IN,
	//the RDATA field is a 4 octet ARPA Internet address.
	RdData string
}

func main() {
	cmd := cobra.Command{
		Use: "dnsr",
		Run: func(cmd *cobra.Command, args []string) {
			m := Message{
				Header: Header{
					Id:      22,
					Qr:      false,
					Opcode:  0,
					Aa:      false,
					Tc:      false,
					Rd:      true,
					Ra:      false,
					Z:       0,
					Rcode:   0,
					QdCount: 1,
					AnCount: 0,
					NsCount: 0,
					ArCount: 0,
				},
				Question: Question{
					Qname:  "dns.google.com",
					Qtype:  1,
					Qclass: 1,
				},
			}

			udpAddr, err := net.ResolveUDPAddr("udp", "8.8.8.8:53")

			if err != nil {
				log.Fatalf("resolving UDP address: %v", err)
			}

			// Dial to the address with UDP
			conn, err := net.DialUDP("udp", nil, udpAddr)
			if err != nil {
				log.Fatalf("dialing up UDP address: %v", err)
			}
			defer conn.Close()
			_, err = conn.Write(m.Encode())
			if err != nil {
				log.Fatalf("writing message: %v", err)
			}

			buffer := make([]byte, 1024)
			n, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				log.Fatalf("reading from UDP : %v", err)
			}
			
			for _, b := range buffer {
				fmt.Printf("%08b ", b)
			}
			respMessage := Parse(buffer[0:n])
			fmt.Printf("n: %d, %+v\n", n, respMessage)
		},
	}

	err := cmd.Execute()
	if err != nil {
		log.Fatalf("Something went wrong: %v\n", err)
	}
}

func encodeName(name string) []byte {
	result := []byte{}
	parts := strings.Split(name, ".")
	for _, part := range parts {
		result = append(result, byte(len(part)))
		result = append(result, []byte(part)...)
	}
	result = append(result, byte(0))
	return result
}
