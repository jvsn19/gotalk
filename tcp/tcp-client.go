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
	if debug {
		output, _ := os.OpenFile("./output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		output.Write([]byte(timeElapsed.String() + "\n"))
		output.Close()
	}
	return timeElapsed.String()
}

var sw stopwatch
var debug bool

// Thread to send messages to another client
func sender(nickname string, conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		message := nickname + ": " + text
		sw.start = time.Now()
		conn.Write([]byte("1\n"))
		fmt.Fprint(conn, message+"\n")
	}
}

// Thread to receive messages and print on screen
func receiver(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, _ := reader.ReadString('\n')
		if message == "1\n" {
			message, _ := reader.ReadString('\n')
			conn.Write([]byte("2\n"))
			fmt.Print(message)
		} else if message == "2\n" {
			sw.end = time.Now()
			sw.getTime()
		}
	}
}

// Open a socket for connection
func openTCPConnection(nickname, port string, conn net.Conn, listener net.Listener) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err = listener.Accept()
	if err != nil {
		return
	}
	go sender(nickname, conn)
	go receiver(conn)
}

// Sucessful connect to another client
func startTCPConnection(nickname, destiny string, conn net.Conn) {
	conn, err := net.Dial("tcp", destiny)
	if err != nil {
		return
	}
	go sender(nickname, conn)
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

	// Connection variable
	var conn net.Conn
	var listener net.Listener

	go openTCPConnection(*nicknamePtr, *portPtr, conn, listener)
	go startTCPConnection(*nicknamePtr, *destinyPtr, conn)
	for {
	}
}
