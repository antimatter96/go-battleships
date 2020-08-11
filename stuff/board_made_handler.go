package stuff

import (
	"encoding/json"

	socketio "github.com/googollee/go-socket.io"
	"github.com/rs/zerolog/log"
)

type thisData struct {
	commonData
	ShipPlacement shipPlacement `json:"shipPlacement"`
}

func (server *Server) boardMadeHandler(s socketio.Conn, msg string) {
	sublogger := log.With().Str("service", "boardMadeHandler").Logger()

	//sublogger.Debug().Msgf("Data : %s", msg)

	var dat *thisData
	if err := json.Unmarshal([]byte(msg), &dat); err != nil {
		sublogger.Error().AnErr("JSON marshalling error", err)
		panic(err)
	}

	if !server.verify(dat.Player, dat.UserToken) {
		return
	}

	if !server.verify(dat.GameID, dat.GameToken) {
		return
	}

	sublogger.Debug().Msgf("Data : %+v", dat.ShipPlacement)

	gg := server.games[dat.GameID]

	aa := gg.PlayerReady(dat.Player, dat.ShipPlacement)

	sublogger.Debug().Msgf("User %-10s - Response: %+v", dat.Player, aa.thisPlayerRes)
	sublogger.Debug().Msgf("User %-10s - Response: %+v", gg.OtherPlayer(dat.Player), aa.otherPlayerRes)

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
