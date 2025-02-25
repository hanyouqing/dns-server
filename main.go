package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <domain>")
		os.Exit(1)
	}

	domain := os.Args[1]

	ips, err := net.LookupIP(domain)
	if err != nil {
		fmt.Printf("Could not get IPs: %v\n", err)
		os.Exit(1)
	}

	for _, ip := range ips {
		fmt.Println(ip)
	}
}
