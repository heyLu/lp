var ws = require('ws');
var http = require('http');
var express = require('express');
var fs = require('fs');
var path = require('path');

var app = express();
app.use(express.static(__dirname + "/public"));

app.get('/stats', function(req, res) {
	var stats = {
		users: wss.clients.length,
		pixls: Object.keys(world).length,
		last_active: new Date(last_active).toISOString()
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

app.get('/save/:name', function(req, res) {
	fs.writeFile(path.join('./data/', req.params.name + '.json'), JSON.stringify(world), function(err) {
		res.setHeader('Content-Type', 'text/plain');
		if (err) {
			res.statusCode = 404;
			res.send(JSON.stringify({error: err}));
		} else {
			res.send(JSON.stringify({name: req.params.name}));
		}
	});
});

app.get('/load/:name', function(req, res) {
	fs.readFile(path.join('./data', req.params.name + '.json'), function(err, data) {
		res.setHeader('Content-Type', 'text/plain');
		if (err) {
			res.statusCode = 404;
			res.send(JSON.stringify({error: err}));
		} else {
			world = JSON.parse(data);
			wss.clients.forEach(function(client) {
				client.send(JSON.stringify(world));
			});
			res.send(JSON.stringify({pixls: Object.keys(world).length}));
		}
	})
});

var server = http.createServer(app);
server.listen(8001);

var wss = new ws.Server({server: server});

var last_active = 0;
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

		last_active = Date.now();
	});
});
