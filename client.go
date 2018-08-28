package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
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

// Thread to send messages to another client
func sender(nickname string, conn net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		message := nickname + ": " + text
		sw.start = time.Now()
		conn.Write([]byte("1\n"))
		fmt.Fprint(conn, message+"\n")
	}
}

// Thread to receive messages and print on screen
func receiver(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
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
func openConnection(nickname, port string, conn net.Conn, listener net.Listener) {
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
func startConnection(nickname, destiny string, conn net.Conn) {
	conn, err := net.Dial("tcp", destiny)
	if err != nil {
		return
	}
	go sender(nickname, conn)
	go receiver(conn)
}

func main() {
	portPtr := flag.String("p", "8080", "Local Port")
	destinyPtr := flag.String("dest", "127.0.0.1:8080", "IP:Port to connect")
	nicknamePtr := flag.String("u", "Picolino", "Nickname")
	flag.Parse()

	fmt.Println("Local Port: ", *portPtr)
	fmt.Println("Destiny: ", *destinyPtr)

	// Connection variable
	var conn net.Conn
	var listener net.Listener

	go openConnection(*nicknamePtr, *portPtr, conn, listener)
	go startConnection(*nicknamePtr, *destinyPtr, conn)
	for {
	}
}
