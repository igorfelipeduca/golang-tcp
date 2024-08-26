package main

import (
	"fmt"
	"net"
)

type Server struct {
	listenAddr string // address to listen on
	ln         net.Listener // listener for incoming connections
	quitch     chan struct{} // quitch is the channel to quit the server
	msgch      chan Message // msgch is the channel to send messages to the server
}

type Message struct {
	from string // the address of the sender
	payload  []byte // the message
}

func NewServer(listenAddr string) *Server {
	// creating a new server and appending the listenAddr to it
	return &Server{
		listenAddr: listenAddr,
		quitch: make(chan struct{}),
		msgch: make(chan Message, 10), // we're defaulting the chan size to 10 to avoid it being empty and blocking the read loop
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.acceptLoop()

	fmt.Println("Server started, listening on", s.listenAddr)

	// Start the message handling goroutine
	go s.handleMessages()

	// Block until quitch is closed
	<-s.quitch

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		go s.readLoop(conn)
	}
}

func (s *Server) readLoop (conn net.Conn) error {
	buf := make([]byte, 2048) // create a buffer to store the message

	for {
		n, err := conn.Read(buf) // read the message buffer from the connection

		if err != nil {
			fmt.Println("read error:", err)
			continue // we're not stopping the server in case of errors, we're just not going to read from this connection anymore
		}

		msg := buf[:n]

		s.msgch <- Message {
			from: conn.RemoteAddr().String(), // the address of the sender
			payload: msg, // the message
		}

		incommingCon := conn.RemoteAddr().String()

		s.writeMessage(conn, Message{
			from: incommingCon,
			payload: []byte(fmt.Sprintf("Hey %s. I got your message: %s", incommingCon, string(msg))),
		})
	}
}

func (s *Server) writeMessage(conn net.Conn, msg Message) error {
	_, err := conn.Write(msg.payload)

	if err != nil {
		fmt.Println("Error while writing to the connection:", err)
		return err;
	}

	return nil
}

func (s *Server) handleMessages() {
	for msg := range s.msgch {
		fmt.Printf("Received message from connection (%s): %s\n", msg.from, string(msg.payload))
	}
}

func main() {
	server := NewServer(":8080")

	go func() {
		err := server.Start()

		if err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()

	// Keep the main goroutine alive
	select {}
}