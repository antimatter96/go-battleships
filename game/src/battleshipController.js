const Battleships = require('./game');
const GameServer = require("./gameController");

const r = require('rethinkdb');

class BattleshipsServer extends GameServer {
  constructor(server, keys, dbName) {
    super(server, keys, dbName, Battleships);
  }

  connect(socket) {
    super.connect(socket);

    socket.on('boardMade',
      this.rejectIfTokenInvalid.bind(this,
        this.rejectIfGameMissing.bind(this),
        this.boardMade.bind(this),
        socket
      )
    );
    socket.on('makeMove',
      this.rejectIfTokenInvalid.bind(this,
        this.rejectIfGameMissing.bind(this),
        this.move.bind(this),
        socket
      )
    );
  }

  async boardMade(socket, player, game, data) {
    console.log("boardMade");
    //console.log("boardMade", data);
    let shipPlacement = data.shipPlacement;
    if (shipPlacement == undefined) {
      console.log("missing shipPlacement");
      return;
    }

    let res = game.playerReady(player, shipPlacement);
    //console.log(res);
    let status = await r.table("games").filter({id: game.id}).update({content: JSON.stringify(game)}).run(this.db);
    //console.log("Replaced", status["replaced"]);
    let otherPlayer = game.otherPlayer(player);
    this.sendStuff(socket, otherPlayer, res);
  }

  async move(socket, player, game, data) {
    console.log("Move");
    //console.log("Move", data);
    let move = data.move;
    if (move == undefined) {
      return;
    }

    let res = game.makeMove(player, move);
    //console.log(res);
    let status = await r.table("games").filter({id: game.id}).update({content: JSON.stringify(game)}).run(this.db);
    //console.log("Replaced", status["replaced"]);
    let otherPlayer = game.otherPlayer(player);
    this.sendStuff(socket, otherPlayer, res);
  }
}

module.exports = BattleshipsServer;
