package main

import (
	"net"
	"os"
	"fmt"
	"strings"
	"bytes"
	"io/ioutil"

	"github.com/hajimehoshi/go-mp3"
)

func main() {
	address, _ := net.ResolveUDPAddr("udp4", ":6969")
	conn, err :=  net.ListenUDP("udp", address)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Listener
	var rawReq []byte
	var req []string
	for {
		rawReq = make([]byte, 256)
		req = make([]string, 3)
		_, address, _ := conn.ReadFromUDP(rawReq)
		req = strings.Fields(string(bytes.Trim(rawReq, "\x00")))
		go handleReq(conn, *address, req)
	}
}

func handleReq(UDPHandle *net.UDPConn, addr net.UDPAddr, req []string) {
	fmt.Printf("%s request from %s\n", req[0], addr.String())
	switch req[0] {
		case "stream":
			err := stream(UDPHandle, &addr, req[1])
			if err != nil {
				UDPHandle.WriteToUDP([]byte{0, 0}, &addr) // error, send 2 bytes
			}
		case "ls":
			files, err := ioutil.ReadDir("./")
   			if err != nil {panic(err)}

    		for _, f := range files {
				if strings.HasSuffix(f.Name(), ".mp3") {
					UDPHandle.WriteToUDP([]byte(f.Name()), &addr)
				}
			}
			UDPHandle.WriteToUDP([]byte{0}, &addr) // confirm, send 1 byte
	}
}

func stream(UDPHandle *net.UDPConn, addr *net.UDPAddr, song string) error {
	f, err := os.Open(song)
	defer f.Close()
	if err != nil {return err}
	UDPHandle.WriteToUDP([]byte{0}, addr) // confirm, send 1 byte

	fmt.Printf("Begin streaming %s to %s\n", song, addr.String())
	// Decodes file
	d, err := mp3.NewDecoder(f)
	defer d.Close()
	if err != nil {return err}

	// Streams file
	var songRawBuf []byte = make([]byte, 4608)
	for {
		_, err := d.Read(songRawBuf)
		if err != nil {break}
		UDPHandle.WriteToUDP(songRawBuf, addr)
	}
	UDPHandle.WriteToUDP([]byte{0}, addr)
	fmt.Println("Finished streaming!")
	return nil
}
