package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

type stopwatch struct {
	start time.Time
	end   time.Time
}

func (sw stopwatch) getTime() string {
	timeElapsed := sw.end.Sub(sw.start)
	fmt.Println(timeElapsed)
	return timeElapsed.String()
}

var sw stopwatch
var debug bool

// Thread to send messages to another client
func sender(nickname string, conn *net.UDPConn, addr *net.UDPAddr) {
	if addr == nil {
		buffer := make([]byte, 1024)
		_, addr, _ = conn.ReadFromUDP(buffer)
		fmt.Println("ok")
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		message := nickname + ": " + text
		sw.start = time.Now()
		_, err := conn.WriteToUDP([]byte("1\n"), addr)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprint(conn, message)
	}
}

// Thread to receive messages and print on screen
func receiver(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
		}
		message := string(buffer[:n])
		if message == "1\n" {
			n, addr, err = conn.ReadFromUDP(buffer)
			if err != nil {
				fmt.Println(err)
			}
			message = string(buffer[:n])
			conn.WriteToUDP([]byte("2\n"), addr)
			fmt.Print(message)
		} else if message == "2\n" {
			sw.end = time.Now()
			sw.getTime()
		}
	}
}

// Open a socket for connection
func openUDPConnection(nickname, destiny, port string) {
	addr, _ := net.ResolveUDPAddr("udp4", "localhost:"+port)
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	go sender(nickname, conn, nil)
	go receiver(conn)
}

// Sucessful connect to another client
func startUDPConnection(nickname, destiny string) {
	addr, _ := net.ResolveUDPAddr("udp4", destiny)
	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return
	}
	time.Sleep(time.Second * 5)
	_, err = conn.WriteToUDP([]byte("1\n"), addr)
	if err != nil {
		fmt.Println(err)
	}
	go sender(nickname, conn, addr)
	go receiver(conn)
}

func readFile(path string) []string {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	toSend := strings.Split(string(fileContent), "\n")
	return toSend
}

func main() {
	portPtr := flag.String("p", "8080", "Local Port")
	destinyPtr := flag.String("dest", "127.0.0.1:8080", "IP:Port to connect")
	nicknamePtr := flag.String("u", "Picolino", "Nickname")
	debugPtr := flag.Bool("d", false, "If true, a log file will be created with message and rtt")
	flag.Parse()

	fmt.Println("Local Port: ", *portPtr)
	fmt.Println("Destiny: ", *destinyPtr)
	debug = *debugPtr

	go openUDPConnection(*nicknamePtr, *destinyPtr, *portPtr)
	go startUDPConnection(*nicknamePtr, *destinyPtr)
	for {
	}
}
