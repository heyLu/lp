# babl

a simple, web-based chat server using websockets.

## running it

    $ npm install
    $ PORT=10001 node server

The server is now running on <http://localhost:10001>.

You can also archive all messages ever sent:

    $ HOST_URL=ws://localhost:10001 node archiver

This allows clients to load the previous messages on initial load.

## ideas

* starting out should be as easy as possible
    - start with a purely client-side bot/script
        * who will answer on the client? the logged in user or a
          user-specific bot?
    - use same api on server
        * still, some things aren't possible on the client (unrestricted
          http requests, file writing or custom protocols)
* would be cool to support images & youtube videos via the same
  (client-side) api, possibly even audio or canvas?
* don't want to reinvent security, protocols... is it possible to use
  hubot? (setup is quite complex, need to understand quite a few
  different parts)

## wishlist

* how do we websocket efficiently? should we batch (both client and
  server), should/could we send binary?

## api playground

here are some ideas how the api could look.

    // client-side
    var chat = document.querySelector("#chat");

    babl.on('connection', function(conn) {
        conn.on('message', function(msg) {
            var display_msg = msg.user + "(" + msg.timestamp + "): " + msg.text;
            chat.textContent += display_msg + "\n";
        });
    });
