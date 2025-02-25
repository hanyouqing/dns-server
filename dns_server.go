package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:53")
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("DNS server listening on 127.0.0.1:53")

	for {
		buf := make([]byte, 512)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		go handleRequest(conn, addr, buf[:n])
	}
}

func handleRequest(conn *net.UDPConn, addr *net.UDPAddr, buf []byte) {
	// Simple parsing of DNS query (very basic)
	queryName := parseQuery(buf)
	if queryName == "" {
		return
	}

	ips, err := net.LookupIP(queryName)
	if err != nil || len(ips) == 0 {
		fmt.Println("Lookup failed for:", queryName)
		return
	}
	ip := ips[0].To4()
	if ip == nil {
		fmt.Println("No IPv4 address found for:", queryName)
		return
	}

	response := buildResponse(buf, ip)
	conn.WriteToUDP(response, addr)
}

func parseQuery(buf []byte) string {
	// Very simplified parsing (for demonstration)
	queryStart := 12
	queryEnd := 0
	for i := queryStart; i < len(buf); i++ {
		if buf[i] == 0 {
			queryEnd = i
			break
		}
	}
	if queryEnd == 0 {
		return ""
	}
	return string(buf[queryStart:queryEnd])
}

func buildResponse(query []byte, ip net.IP) []byte {
	// Very simplified response building (for demonstration)
	response := make([]byte, 512)
	copy(response, query[:12])                     // Copy header
	response[2] |= 0x80                            // Set response flag
	response[3] |= 0x80                            // Set recursion available
	binary.BigEndian.PutUint16(response[6:8], 1)   // Answer count: 1
	binary.BigEndian.PutUint16(response[10:12], 1) // Additional count: 1

	queryEnd := 0
	for i := 12; i < len(query); i++ {
		if query[i] == 0 {
			queryEnd = i + 1
			break
		}
	}

	copy(response[12:], query[12:queryEnd]) // Copy query
	answerStart := 12 + len(query[12:queryEnd])

	binary.BigEndian.PutUint16(response[answerStart:answerStart+2], 0xc00c) // pointer to name
	binary.BigEndian.PutUint16(response[answerStart+2:answerStart+4], 1)    // Type A
	binary.BigEndian.PutUint16(response[answerStart+4:answerStart+6], 1)    // Class IN
	binary.BigEndian.PutUint32(response[answerStart+6:answerStart+10], 60)  // TTL 60
	binary.BigEndian.PutUint16(response[answerStart+10:answerStart+12], 4)  // Data length: 4
	copy(response[answerStart+12:], ip)                                     // IP address

	return response[:answerStart+16]
}
