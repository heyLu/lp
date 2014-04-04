var ws = require('ws');

var wss = new ws.Server({port: 8080});

wss.broadcast = function(data) {
	var msg = typeof data === "string" ? data : JSON.stringify(data);
	this.clients.forEach(function(client) {
		client.send(msg);
	});
};

wss.on('connection', function(ws) {
	var name = randomName();

	ws.on('message', function(msg) {
		var msg = JSON.parse(msg);
		msg.author = name;
		msg.timestamp = Date.now();

		wss.broadcast(msg);
	});

	ws.on('close', function() {
		var msg = {
			type: "disconnect",
			author: name,
			timestamp: Date.now()
		};
		wss.broadcast(msg);
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
