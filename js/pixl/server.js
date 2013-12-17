var ws = require('ws');
var http = require('http');
var express = require('express');

var app = express();
app.use(express.static(__dirname + "/public"));

var server = http.createServer(app);
server.listen(8001);

var wss = new ws.Server({server: server});

var world = {};

wss.on('connection', function(socket) {
	socket.send(JSON.stringify(world));

	socket.on('message', function(msg) {
		var pixls = JSON.parse(msg);
		pixls.forEach(function(pixl) {
			world[pixl.x + "," + pixl.y] = {color: pixl.color};
			console.log(pixl);
		});

		wss.clients.forEach(function(client) {
			if (socket != client) {
				client.send(msg);
			}
		});
	});
});