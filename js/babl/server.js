var ws = require('ws');

var wss = new ws.Server({port: 8080});

wss.on('connection', function(ws) {
	var name = randomName();

	ws.on('message', function(msg) {
		var msg = JSON.parse(msg);
		msg.author = name;
		msg.timestamp = Date.now();

		wss.clients.forEach(function(sock) {
			sock.send(JSON.stringify(msg));
		});
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
