package stuff

import (
	"fmt"
	"sync"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

type Server struct {
	Key []byte

	Server   *socketio.Server
	present  map[string]bool
	socketOf map[string]*socketio.Conn

	queueLock sync.Mutex
	Queue     []string

	playerIsIn map[string]string

	games map[string]*BattleShips
	//addUserHandler func()
}

func (server *Server) Init() {
	server.present = make(map[string]bool)
	server.socketOf = make(map[string]*socketio.Conn)
	server.games = make(map[string]*BattleShips)

	server.Queue = make([]string, 0)

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				// for , y := range server.socketOf {
				// 	fmt.Println(x, (*y).Context())
				// }
			case <-quit:
				ticker.Stop()
				return
			}
		}

	}()

	ticker2 := time.NewTicker(5 * time.Second)
	quit2 := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker2.C:
				for _, k := range server.Queue {
					fmt.Println(k)
				}
			case <-quit2:
				ticker2.Stop()
				return
			}
		}

	}()

	server.Server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.Server.OnEvent("/", "updateSocket", server.JoinGameHandler)

	server.Server.OnEvent("/", "boardMade", server.BoardMadeHandler)
	server.Server.OnEvent("/", "makeMove", server.MakeMoveHandler)

	server.Server.OnEvent("/", "join", server.JoinGameHandler)
	server.Server.OnEvent("/", "addUser", server.AddUserHandler)

	server.Server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.Server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("Internal Server Error:", e.Error())
	})

	server.Server.OnDisconnect("/", server.DisconnectHandler)
}
