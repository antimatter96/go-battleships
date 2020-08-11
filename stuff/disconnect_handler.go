package stuff

import (
	"fmt"

	socketio "github.com/googollee/go-socket.io"
)

// DisconnectHandler is
func (server *Server) DisconnectHandler(s socketio.Conn, msg string) {
	fmt.Println("DisconnectHandler", msg, s.Context())

	if s.Context() != nil {
		name := s.Context().(string)

		delete(server.socketOf, name)
	}
}
