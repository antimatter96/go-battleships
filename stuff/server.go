package stuff

import (
	"fmt"
	"sync"

	socketio "github.com/googollee/go-socket.io"
)

type Server struct {
	Key []byte

	Server   *socketio.Server
	present  map[string]bool
	socketOf map[string]*socketio.Conn

	queueLock sync.Mutex
	Queue     []string

	games map[string]*BattleShips
	//addUserHandler func()
}

// Init is used to start everything
func (server *Server) Init() {
	server.present = make(map[string]bool)
	server.socketOf = make(map[string]*socketio.Conn)
	server.games = make(map[string]*BattleShips)

	server.Queue = make([]string, 0)

	server.Server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.Server.OnEvent("/", "addUser", server.addUserHandler)
	server.Server.OnEvent("/", "updateSocket", server.joinGameHandler)

	server.Server.OnEvent("/", "join", server.joinGameHandler)

	server.Server.OnEvent("/", "boardMade", server.boardMadeHandler)
	server.Server.OnEvent("/", "makeMove", server.makeMoveHandler)

	server.Server.OnDisconnect("/", server.DisconnectHandler)

	server.Server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("Internal Server Error:", e.Error())
	})
}
