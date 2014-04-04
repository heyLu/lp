var ws = require('ws');
var http = require('http');
var express = require('express');

var app = express();
app.use(express.static(__dirname + "/public"));

var server = http.createServer(app);

app.get('/msgs', function(req, res) {
	res.sendfile('msgs.json');
});

var wss = new ws.Server({server: server});

wss.broadcast = function(data) {
	var msg = typeof data === "string" ? data : JSON.stringify(data);
	this.clients.forEach(function(client) {
		client.send(msg);
	});
};

wss.on('connection', function(ws) {
	var name = randomName();

	wss.broadcast({type: "connect", author: name, timestamp: Date.now()});

	ws.on('message', function(msg) {
		var msg = JSON.parse(msg);
		msg.author = name;
		msg.timestamp = Date.now();

		wss.broadcast(msg);
	});

	ws.on('close', function() {
		wss.broadcast({type: "disconnect", author: name, timestamp: Date.now()});
	});
});

function randomName() {
	var name = "";
	var a = "a".charCodeAt(0);
	for (var i = 0; i < 8; i++) {
		name += String.fromCharCode(a + Math.random() * 26);
	}
	return name;
}

var port = parseInt(process.env.PORT || "8080");
server.listen(port);
