var ws = require('ws');

var server = new ws.Server({port: 8001});

var world = {};

server.on('connection', function(socket) {
	socket.send(JSON.stringify(world));

	socket.on('message', function(msg) {
		var pixls = JSON.parse(msg);
		pixls.forEach(function(pixl) {
			world[pixl.x + "," + pixl.y] = {color: pixl.color};
			console.log(pixl);
		});

		server.clients.forEach(function(client) {
			if (socket != client) {
				client.send(msg);
			}
		});
	});
});