package stuff

import (
	"encoding/json"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	socketio "github.com/googollee/go-socket.io"
	"github.com/rs/zerolog/log"
)

type myCustomClaims2 struct {
	ID  string `json:"id"`
	Exp int64  `json:"exp"`
	Nbf int64  `json:"nbf"`
	jwt.StandardClaims
}

func (server *Server) verifyGame(gameID, token string) bool {
	token2, err := jwt.ParseWithClaims(token, &myCustomClaims2{}, func(token *jwt.Token) (interface{}, error) {
		return server.Key, nil
	})

	if claims, ok := token2.Claims.(*myCustomClaims2); ok && token2.Valid {
		if claims.ID == gameID {
			return true
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
	}
	fmt.Println(token2.Claims, err)
	return false
}

type thisData struct {
	Player        string
	UserToken     string
	GameToken     string
	GameID        string
	ShipPlacement shipPlacement `json:"shipPlacement"`
}

// BoardMadeHandler is
func (server *Server) BoardMadeHandler(s socketio.Conn, msg string) {
	//log.Debug().Str("service", "BoardMadeHandler").Msgf("Data : %s", msg)
	//fmt.Println("BoardMadeHandler", msg)

	var dat *thisData
	if err := json.Unmarshal([]byte(msg), &dat); err != nil {
		//fmt.Println("BoardMadeHandler", err)
		panic(fmt.Errorf("JSON UNMARSHAL ERROR %v", err))
	}

	if !server.vetify(dat.Player, dat.UserToken) {
		return
	}

	if !server.verifyGame(dat.GameID, dat.GameToken) {
		return
	}

	fmt.Println(dat.ShipPlacement)

	gg := server.games[dat.GameID]

	aa := gg.PlayerReady(dat.Player, dat.ShipPlacement)

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
type shipPlacement map[string](map[int]BoardPoint)
