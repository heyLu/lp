import { ExtensionSyncWebsocket, Server } from "@earthstar/earthstar/deno";

const extensions = [
  new ExtensionSyncWebsocket("sync"),
];

const server = new Server(
  extensions,
  {
    peer: {
      password: "myextremelygoodlongpassword"
    },
    port: "8080",
  }
);
