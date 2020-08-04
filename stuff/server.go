package stuff

import (
	"fmt"

	socketio "github.com/googollee/go-socket.io"
)

type Server struct {
	Key []byte

	Server  *socketio.Server
	present map[string]bool
	//addUserHandler func()
}

func (server *Server) Init() {

	server.present = make(map[string]bool)

	server.Server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.Server.OnEvent("/", "updateSocket", server.JoinGameHandler)

	server.Server.OnEvent("/", "boardMade", server.JoinGameHandler)
	server.Server.OnEvent("/", "makeMove", server.JoinGameHandler)

	server.Server.OnEvent("/", "join", server.JoinGameHandler)
	server.Server.OnEvent("/", "addUser", server.AddUserHandler)

	server.Server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})

	server.Server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.Server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e.Error(), s.Context())
	})

	server.Server.OnDisconnect("/", server.DisconnectHandler)
}
