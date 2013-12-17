var ws = require('ws');
var http = require('http');
var express = require('express');

var app = express();
app.use(express.static(__dirname + "/public"));

app.get('/stats', function(req, res) {
	var stats = {
		users: wss.clients.length,
		pixls: Object.keys(world).length
	};
	res.setHeader('Content-Type', 'text/plain');
	res.send(JSON.stringify(stats, null, "  "));
});

app.get('/world', function(req, res) {
	res.setHeader('Content-Type', 'text/plain');
	res.send(JSON.stringify(world));
});

app.get('/reset', function(req, res) {
	world = {};
	wss.clients.forEach(function(client) {
		client.send("{}");
	});
	res.setHeader('Content-Type', 'text/plain');
	res.send(JSON.stringify({status: "ok"}));
});

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
