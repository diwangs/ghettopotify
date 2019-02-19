package main

import (
	"net"
	"bufio"
	"fmt"
	"os"
	"time"
	"strings"
	"bytes"
	"strconv"

	"github.com/hajimehoshi/oto"
)

// Declare the connection object as a global variable
var subList []string 
var conn *net.UDPConn

func main() {
	subList = []string{"127.0.0.1:6969"}
	address, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:6969")
	conn, _ =  net.DialUDP("udp", nil, address)
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("> ")
		scanner.Scan()
		handleCmd(strings.Fields(scanner.Text()))
	}
}

func handleCmd(cmd []string) {
	switch cmd[0] {
		case "play":
			if len(cmd) < 2 {
				fmt.Println("What song?")
				return
			} 
			var buf []byte
			temp := make([]byte, 2)
			conn.Write([]byte(fmt.Sprintf("stream %s", cmd[1])))
			n, _ := conn.Read(temp) // Check whether the song exists
			if n == 1 {
				fmt.Println("Playing...")
				go fillBuffer(&buf)
				for len(buf) == 0 {
					time.Sleep(1000 * time.Millisecond)
				}
				play(&buf)
			} else { 
				fmt.Println("Track not found!")
			}
		case "ls":
			conn.Write([]byte("ls"))
			var rawName []byte
			for {
				rawName = make([]byte, 20)
				n, _ := conn.Read(rawName)
				if n == 1 {break}
				fmt.Println(string(bytes.Trim(rawName, "\x00")))
			}
		case "lschan":
			for _, ch := range subList {
				fmt.Println(ch)
			}
		case "chchan":
			if len(cmd) < 2 {
				fmt.Println("What channel?")
				return
			}
			sellen, _ := strconv.Atoi(cmd[1])
			if len(subList) <= sellen {
				fmt.Println("Channel not found on sublist")
				return
			}
			address, _ := net.ResolveUDPAddr("udp4", subList[sellen])
			conn, _ =  net.DialUDP("udp", nil, address)
		case "sub":
			if len(cmd) < 2 {
				fmt.Println("What to sub?")
				return
			}
			subList = append(subList, cmd[1])
		case "exit":
			fmt.Println("Peace out")
			os.Exit(0)
		case "help":
			fmt.Printf("play <file>\tPlay a song file\nexit\t\tExit Ghettopotify\nls\t\tList songs\n")
		default:
			fmt.Printf("%s is not a valid command, maybe you need `help'?\n", cmd[0])
	}
}

func fillBuffer(buffer *[]byte) {
	packetBuf := make([]byte, 4608)
	for {
		n, _ := conn.Read(packetBuf)
		// 1 = EOF
		if n == 1 {break}
		*buffer = append(*buffer, packetBuf...)
	}
	fmt.Println("Buffering finished!")
}

func play(buffer *[]byte) {
	p, err := oto.NewPlayer(44100, 2, 2, 8192)
	if err != nil {panic(err)}
	defer p.Close()

	// Pausing mechanic
	pauser := false
	
	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for {
			scanner.Scan()
			if scanner.Text() == "pause" {
				fmt.Println("Pausing...")
				pauser = true
			} else if (scanner.Text() == "resume") {
				fmt.Println("Resuming...")
				pauser = false
			} else {break} // Bug: must enter arbitrary character to enter a command again
		}
	}()
			
	// Play
	p.Write(buffer, &pauser) 
}
