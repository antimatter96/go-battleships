const http = require('http');

const BattleshipsController = require("./battleshipController");
const Server = require("./server");

/**
 * 
*/

const config = require("../config");
const app = Server.getExpressApp(config);

const server = http.createServer(app);
server.listen(process.env.PORT || 8000);
server.on("error", function(err) {
  console.log(err);
});

const gameController = new BattleshipsController(server, config.keys, "battleships");
gameController.Start();
