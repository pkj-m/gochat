package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"strings"
)

// list of supported commands

// !update-name <new_name>

// note: all commands should end with a space to ensure that the input is not a prefix
const (
	// UpdateNameCommand lets clients update entity
	UpdateNameCommand = "!update-name"
)

// processMessage accepts client input and attempts to process it as either a
// server command or an ordinary text message
func (s *Server) processMessage(ws *websocket.Conn, b []byte) error {
	msg := string(b)
	username := s.usernames[ws]
	msgs := strings.Split(msg, " ")

	if len(msgs) == 2 && msgs[0] == UpdateNameCommand {
		// validate the username
		invalid := isUsernameValid(msgs[1])
		if invalid != nil {
			_ = s.replyToClient(ws, []byte(fmt.Sprintf("-- error: %s", invalid.Error())))
			return invalid
		}

		// update the connection name with the new name
		s.usernames[ws] = msgs[1]
		_ = s.replyToClient(ws, []byte(fmt.Sprintf("-- username updated to %s", msgs[1])))
		return nil
	}

	formattedMsg := fmt.Sprintf("[%s] %s", username, msg)
	s.broadcast(ws, []byte(formattedMsg))
	return nil
}
