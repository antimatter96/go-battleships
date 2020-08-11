package stuff

import (
	"encoding/json"

	socketio "github.com/googollee/go-socket.io"
	"github.com/rs/zerolog/log"
)

type thisData2 struct {
	commonData
	Move boardPoint `json:"move"`
}

func (server *Server) makeMoveHandler(s socketio.Conn, msg string) {
	sublogger := log.With().Str("service", "makeMoveHandler").Logger()

	sublogger.Debug().Msgf("Data : %s", msg)
	//fmt.Println("MakeMoveHandler", msg)

	var dat *thisData2
	if err := json.Unmarshal([]byte(msg), &dat); err != nil {
		sublogger.Error().AnErr("JSON marshalling error", err)
		panic(err)
	}

	sublogger.Debug().Msgf("Parsed Data : %+v", dat.Move)

	if !server.verify(dat.Player, dat.UserToken) {
		return
	}

	if !server.verify(dat.GameID, dat.GameToken) {
		return
	}
	gg := server.games[dat.GameID]

	aa := gg.MakeMove(dat.Player, dat.Move)

	sublogger.Debug().Msgf("User %-10s - Response: %+v    User %-10s - Response: %+v", dat.Player, aa.thisPlayerRes, gg.OtherPlayer(dat.Player), aa.otherPlayerRes)

	for _, v := range aa.thisPlayerRes {
		//fmt.Println("this", v.data)
		b, err := json.Marshal(v.data)
		if err != nil {
			sublogger.Error().AnErr("JSON marshalling error", err)
			return
		}
		s.Emit(v.message, string(b))
	}

	if len(aa.otherPlayerRes) != 0 {
		otherPlayerSocket := server.socketOf[gg.OtherPlayer(dat.Player)]
		for _, v := range aa.otherPlayerRes {
			//fmt.Println("other", v.data)
			b, err := json.Marshal(v.data)
			if err != nil {
				sublogger.Error().AnErr("JSON marshalling error", err)
				return
			}
			(*otherPlayerSocket).Emit(v.message, string(b))
		}
	}

}
