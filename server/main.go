package main

import (
	"bufio"
	"log"
	"net"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}

}

var (
	openConnection = make(map[net.Conn]bool)
	newConnection  = make(chan net.Conn)
	deadConnection = make(chan net.Conn)
)

func main() {
	ln, err := net.Listen("tcp", ":3000")
	logFatal(err)

	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			logFatal(err)

			openConnection[conn] = true
			newConnection <- conn
		}

	}()
	for {
		select {
		case conn := <-newConnection:
			//Invoke broadcase Message (broadcaset to the other connection)
			go broadcastMessage(conn)
		case conn := <-deadConnection:
			//remove or delete the connection

			for item := range openConnection {
				if item == conn {
					break
				}
			}
			delete(openConnection, conn)
		}
	}
}
func broadcastMessage(conn net.Conn) {
	for {

		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')

		if err != nil {
			break
		}
		//loop through all the open connection
		//and send message to these connection
		//except the connection thAT SEND the message
		for item := range openConnection {
			if item != conn {
				item.Write([]byte(message))
			}
		}

	}
}
