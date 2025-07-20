package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

//{type: "init", userId: "xyz"} or {type: "message", to:"xyz", body:"xqyzhjoaig"}

type Server struct {
	muWsConns sync.Mutex
	muUsers   sync.Mutex
	conns     map[*websocket.Conn]bool
	users     map[string]*websocket.Conn
}

type BaseMessage struct {
	Type string `json:"type"`
}

type InitMessage struct {
	Type string `json:"type"`
	User string `json:"user"`
}

type ChatMessage struct {
	Type     string `json:"type"`
	Receiver string `json:"to"`
	Body     string `json:"body"`
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
		users: make(map[string]*websocket.Conn),
	}
}

func (s *Server) handleWs(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client: ", ws.RemoteAddr())

	s.muWsConns.Lock()
	s.conns[ws] = true
	s.muWsConns.Unlock()

	s.ReadLoop(ws)
	ws.Close()

	s.muWsConns.Lock()
	delete(s.conns, ws)
	s.muWsConns.Unlock()
}

func (s *Server) ReadLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		numOfBytesRead, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error: ", err)
			continue
		}
		msg := buf[:numOfBytesRead]
		var baseMessage BaseMessage
		err = json.Unmarshal(msg, &baseMessage)
		if err != nil {
			fmt.Println("Error parsing to base Message: ", err)
			continue
		}
		switch baseMessage.Type {
		case "init":
			var initMessage InitMessage
			err = json.Unmarshal(msg, &initMessage)
			if err != nil {
				fmt.Println("Error parsing to init Message: ", err)
				continue
			}
			s.muUsers.Lock()
			s.users[initMessage.User] = ws
			s.muUsers.Unlock()
		case "message":
			var chatMessage ChatMessage
			err = json.Unmarshal(msg, &chatMessage)
			if err != nil {
				fmt.Println("Error parsing to chatMessage: ", err)
				continue
			}
			err := s.sendMessage(chatMessage, ws)
			if err != nil {
				fmt.Println("err while sending messages: ", err)
			}
		default:
			fmt.Println("Cant find type of this message.")
		}
		ws.Write([]byte("thank you for the message"))
	}
}

func (s *Server) sendMessage(message ChatMessage, senderWs *websocket.Conn) error {
	s.muUsers.Lock()
	receiverWsConn := s.users[message.Receiver]
	s.muUsers.Unlock()
	if receiverWsConn == nil {
		senderWs.Write([]byte("no user found with the provided username"))
		return errors.New("value for this key not found")
	}
	_, err := receiverWsConn.Write([]byte("Message received from User: " + senderWs.RemoteAddr().String() + "\n Message: " + message.Body))
	if err != nil {
		return errors.New("error writing to receiver Ws Conn")
	}
	return nil
}

func main() {
	server := NewServer()
	//http.Handle("/ws", corsMiddleware(websocket.Handler(server.handleWs)))
	http.Handle("/ws", websocket.Handler(server.handleWs))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
