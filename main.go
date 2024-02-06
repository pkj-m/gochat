package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log/slog"
	"net/http"
)

const UsernameDefaultLength = 6

type Server struct {
	conns     map[*websocket.Conn]bool
	usernames map[*websocket.Conn]string
}

func NewServer() *Server {
	return &Server{
		conns:     make(map[*websocket.Conn]bool),
		usernames: make(map[*websocket.Conn]string),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	username := generateRandomString(UsernameDefaultLength)
	slog.Info(fmt.Sprintf("new incoming connection: %s - %s", ws.RemoteAddr(), username))
	s.conns[ws] = true
	s.usernames[ws] = username
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)

	// send welcome msg to client with their username
	username := s.usernames[ws]
	_, _ = ws.Write([]byte(fmt.Sprintf("Hello. Your username is: %s", username)))

	// read all subsequent messages
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			slog.Error(fmt.Sprintf("[%s] read error: %s", username, err.Error()))
			continue
		}
		msg := buf[:n]

		err = s.processMessage(ws, msg)
		if err != nil {
			slog.Error(fmt.Sprintf("[%s] process message error: %s", username, err.Error()))
		}
	}
}

// broadcast accepts the sender connection and the byte message to be delivered to all other clients
func (s *Server) broadcast(sender *websocket.Conn, b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			// do not send msg back to sender
			username := s.usernames[ws]
			if ws == sender {
				return
			}
			if _, err := ws.Write(b); err != nil {
				slog.Error(fmt.Sprintf("[%s] write error: %s", username, err.Error()))
			}
		}(ws)
	}
}

// replyToClient allows server to send messages to a single client. Auto-logs errors.
func (s *Server) replyToClient(receiver *websocket.Conn, b []byte) error {
	_, err := receiver.Write(b)
	if err != nil {
		slog.Error(fmt.Sprintf("error in replyToClient: %s", err.Error()))
	}
	return err
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	slog.Info("starting web server at port:3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		slog.Error(fmt.Sprintf("server err: %s", err.Error()))
	}
}
