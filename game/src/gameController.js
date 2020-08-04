const { List } = require('immutable');

const SocketIO = require('socket.io');
const r = require('rethinkdb');
const jwt = require("jsonwebtoken");

class GameServer {
  constructor(server, keys, dbName, gameClass) {
    if (!server || typeof (server.listeners) !== "function") {
      throw new Error("Server not present");
    }
    this.io = SocketIO(server);

    this.privateKey = keys.privateKey;
    this.publicKey = keys.publicKey;
    /*
    =====================================
      Currently storing these in memory;
      Might use a datastore
    ======================================
    */

    this.UsersInQueue = List();
    this.Users = new Set();
    this.socketOfUser = [];

    this.playerIsIn = {};

    r.connect({
      db: dbName,
    }, this.rethinkDBConnectionCallback.bind(this));

    this.Game = gameClass;
  }

  rethinkDBConnectionCallback(err, conn) {
    if (err) {
      throw new Error("Can not connect to DB");
    }

    this.db = conn;
  }

  Start() {
    this.io.on('connect', this.connect.bind(this));
  }

  connect(socket) {
    console.log("_____client connected_____");

    socket.on('disconnect', this.disconnect.bind(this, socket));
    socket.on('addUser', this.addUser.bind(this, socket));
    socket.on('updateSocket', this.updateSocket.bind(this, socket));
    socket.on('join', this.join.bind(this, socket));
  }

  disconnect(_socket, _data) {
    console.log("_____client disconnected_____");
  }

  addUser(socket, data) {
    console.log("Add user", data.username);
    //add gaurd
    if (this.Users.has(data.name)) {
      socket.emit('userAdded', {
        msg: 'Username taken'
      });
      return;
    }

    this.Users.add(data.name);
    socket.username = data.name;
    this.socketOfUser[data.name] = socket.id;
    socket.emit('userAdded', {
      msg: 'OK',
      name: data.name,
      userToken: jwt.sign({ name: data.name }, this.privateKey, {
        "algorithm": "RS256",
        "expiresIn": "12h",
      })
    });

  }

  updateSocket(socket, data) {
    //console.log("Updating socket", data);
    if (!data.userToken) {
      socket.emit('updateFailed');
      return;
    }

    var decoded;
    try {
      decoded = jwt.verify(data.userToken, this.publicKey, {
        "algorithm": "R256",
      });
    } catch (error) {
      // Better to rate limit this user
      socket.emit('updateFailed');
      return;
    }

    if (decoded.name != data.player) {
      socket.emit('updateFailed');
      return;
    }

    this.socketOfUser[data.player] = socket.id;
    socket.username = data.player;

    socket.emit('updateSuccess');
  }

  async join(socket, data) {
    console.dir("Joining", data.player);
    let player1 = socket.username;
    if (player1 !== data.player) {
      this.updateSocket(socket, data);
      player1 = data.player;
    }

    if (this.UsersInQueue.includes(player1)) {
      socket.emit('lockJoin');
      return;
    }

    if (this.UsersInQueue.size <= 0) {
      this.UsersInQueue = this.UsersInQueue.push(player1);
      return;
    }

    let player2 = this.UsersInQueue.first();
    if (player2 === player1) {
      // Da actual faq
      return;
    }

    this.UsersInQueue = this.UsersInQueue.shift();

    let newGame = new this.Game(player1, player2);

    await r.table("games").insert({
      id: newGame.id,
      player1: player1,
      player2: player2,
      content: JSON.stringify(newGame),
    }).run(this.db);

    this.playerIsIn[player1] = newGame.id;
    this.playerIsIn[player2] = newGame.id;

    let gameToken = jwt.sign({ gameId: newGame.id }, this.privateKey, {
      "algorithm": "RS256",
      "expiresIn": "12h",
    });

    socket.emit('startGame', {
      'otherPlayer': player2,
      'gameId': newGame.id,
      'gameToken': gameToken,
    });
    socket.to(this.socketOfUser[player2]).emit('startGame', {
      'otherPlayer': player1,
      'gameId': newGame.id,
      'gameToken': gameToken,
    });

  }

  sendStuff(currentSocket, otherPlayer, res) {
    for (let message of res.thisPlayer) {
      currentSocket.emit(message.message, message.data);
    }

    if (res.otherPlayer) {
      for (let message of res.otherPlayer) {
        currentSocket.to(this.socketOfUser[otherPlayer]).emit(message.message, message.data);
      }
    }
  }

  async rejectIfGameMissing(callback, socket, data) {
    //console.log("rejectIfGameMissing");
    if (data == null || typeof (data) !== "object" || Object.keys(data).length === 0) {
      return new Error("missing data");
    }

    let player = data.player;
    if (typeof (player) !== "string" || player.trim() === "") {
      return new Error("missing playerId");
    }

    let gameId = this.playerIsIn[player];
    if (typeof (gameId) !== "string" || gameId.trim() === "") {
      return new Error("missing gameId");
    }

    let games = await r.table("games").filter({ id: gameId }).run(this.db);
    if (games.constructor.name != "Cursor") {
      throw new Error("db error");
    }

    let storedGame = null;
    try {
      storedGame = await games.next();
    } catch (error) {
      if (error.name != "ReqlDriverError") {
        return new Error("missing game");
      }
    }

    if (storedGame == null) {
      return new Error("missing game");
    }

    let game = this.Game.gameFromString(storedGame.content);

    if (game == null || typeof (game) !== "object" || Object.keys(game).length === 0) {
      return new Error("missing game");
    }

    callback(socket, player, game, data);
  }

  async rejectIfTokenInvalid(callThis, callWith, socket, data) {
    if (!data.userToken) {
      socket.emit('updateFailed');
      return;
    }
    if (!data.gameToken) {
      socket.emit('updateFailed');
      return;
    }

    var decodedUser;
    var decodedGame;
    try {
      decodedUser = jwt.verify(data.userToken, this.publicKey, {
        "algorithm": "R256",
      });
      decodedGame = jwt.verify(data.gameToken, this.publicKey, {
        "algorithm": "R256",
      });
    } catch (error) {
      // Better to rate limit this user
      socket.emit('updateFailed');
      return;
    }

    if (decodedUser.name != data.player) {
      socket.emit('updateFailed');
      return;
    }

    if (socket.username != data.player) {
      console.log(socket.username, data.player);
      socket.emit('updateFailed');
      return;
    }

    if (decodedGame.gameId != data.gameId) {
      console.log(decodedGame, data.gameId);
      socket.emit('updateFailed');
      return;
    }

    callThis(callWith, socket, data);
  }
}

module.exports = GameServer;
