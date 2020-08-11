package stuff

import (
	"encoding/json"
	"fmt"

	socketio "github.com/googollee/go-socket.io"
	"github.com/rs/zerolog/log"
)

type thisData struct {
	Player        string
	UserToken     string
	GameToken     string
	GameID        string
	ShipPlacement shipPlacement `json:"shipPlacement"`
}

func (server *Server) boardMadeHandler(s socketio.Conn, msg string) {
	//log.Debug().Str("service", "BoardMadeHandler").Msgf("Data : %s", msg)
	//fmt.Println("BoardMadeHandler", msg)

	var dat *thisData
	if err := json.Unmarshal([]byte(msg), &dat); err != nil {
		//fmt.Println("BoardMadeHandler", err)
		panic(fmt.Errorf("JSON UNMARSHAL ERROR %v", err))
	}

	if !server.verify(dat.Player, dat.UserToken) {
		return
	}

	if !server.verify(dat.GameID, dat.GameToken) {
		return
	}

	fmt.Println(dat.ShipPlacement)

	gg := server.games[dat.GameID]

	aa := gg.PlayerReady(dat.Player, dat.ShipPlacement)

	fmt.Println(dat.Player, aa.thisPlayerRes)
	fmt.Println(gg.OtherPlayer(dat.Player), aa.otherPlayerRes)

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

//
type shipPlacement map[string](map[int]boardPoint)
