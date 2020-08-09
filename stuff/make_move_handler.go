package stuff

import (
	"encoding/json"
	"fmt"

	socketio "github.com/googollee/go-socket.io"
	"github.com/rs/zerolog/log"
)

type thisData2 struct {
	Player    string
	UserToken string
	GameToken string
	GameID    string
	Move      BoardPoint `json:"move"`
}

// BoardMadeHandler is
func (server *Server) MakeMoveHandler(s socketio.Conn, msg string) {
	//log.Debug().Str("service", "MakeMoveHandler").Msgf("Data : %s", msg)
	//fmt.Println("MakeMoveHandler", msg)

	var dat *thisData2
	if err := json.Unmarshal([]byte(msg), &dat); err != nil {
		//fmt.Println("BoardMadeHandler", err)
		panic(fmt.Errorf("JSON UNMARSHAL ERROR %v", err))
	}

	fmt.Printf("MakeMoveHandler %+v\n", dat.Move)

	if !server.vetify(dat.Player, dat.UserToken) {
		return
	}

	if !server.verifyGame(dat.GameID, dat.GameToken) {
		return
	}
	gg := server.games[dat.GameID]

	aa := gg.MakeMove(dat.Player, dat.Move)

	fmt.Printf("%+v\n", aa)

	for _, v := range aa.thisPlayerRes {
		//fmt.Println("this", v.data)
		b, err := json.Marshal(v.data)
		if err != nil {
			log.Fatal().Err(err)
		}
		s.Emit(v.message, string(b))
	}

	if len(aa.otherPlayerRes) != 0 {
		otherPlayerSocket := server.socketOf[gg.OtherPlayer(dat.Player)]
		for _, v := range aa.otherPlayerRes {
			//fmt.Println("other", v.data)
			b, err := json.Marshal(v.data)
			if err != nil {
				log.Fatal().Err(err)
			}
			(*otherPlayerSocket).Emit(v.message, string(b))
		}
	}

}
