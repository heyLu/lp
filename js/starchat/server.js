import { Peer } from "@earthstar/earthstar";
import { ExtensionSyncWebsocket, RuntimeDriverDeno, Server, getStorageDriverFilesystem } from "@earthstar/earthstar/deno";

const extensions = [
  new ExtensionSyncWebsocket("sync"),
];

const server = new Server(
  extensions,
  {
    peer: new Peer({
      password: "myextremelygoodlongpassword",
      runtime: new RuntimeDriverDeno(),
      storage: await getStorageDriverFilesystem("star-storage"),
    }),
    port: "8080",
  }
);
