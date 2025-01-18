import { Path, Peer, RuntimeDriverUniversal, isErr, notErr } from "@earthstar/earthstar";
import { StorageDriverIndexedDB } from "@earthstar/earthstar/browser";

const peer = new Peer({
  password: "password1234",
  runtime: new RuntimeDriverUniversal(),
  storage: new StorageDriverIndexedDB(),
});

console.log(peer.shares());

const authorKeypair = await peer.createIdentity("test");
// const authorKeypair = { tag: "@test.bslpgo6hu2epr3xusywnrhvb7yczr2hfka6tqbo3i5qtnpehbyhua", secretKey: "bnwewxf3oekevgvei7fze6gvgknkr6r56ivbzpjgvdjk77x5bwasq" };
const shareKeypair = await peer.createShare("chatting");
// const shareKeypair = { tag: "-chatting.bfptozwv5fte7tngdywldklbkgflsctrogh7mhdie57m2gyedj3tq", secretKey: "bpssslrdwy4ns5izkgyyhd5g7erlxcassqkjkxtxrbwa6rzacaemq" };

if (notErr(shareKeypair) && notErr(authorKeypair)) {
	console.group("Share keypair");
	console.log(shareKeypair);
	console.groupEnd();

	console.group("Author keypair");
	console.log(authorKeypair);
	console.groupEnd();
} else if (isErr(shareKeypair)) {
	console.error(shareKeypair);
} else if (isErr(authorKeypair)) {
	console.error(authorKeypair);
}

const readCap = await peer.mintCap(shareKeypair.tag, authorKeypair.tag, "read");
if (isErr(readCap)) {
  console.error("mint read", readCap);
}
const writeCap = await peer.mintCap(shareKeypair.tag, authorKeypair.tag, "write");
if (isErr(writeCap)) {
  console.error("mint write", writeCap);
}

const store = await peer.getStore(shareKeypair.tag);
if (isErr(store)) {
  console.error("get store", store);
}

const form = document.getElementById("message-form");
const input = document.querySelector("input");

form.addEventListener("submit", async (event) => {
  // prevent page from reloading
  event.preventDefault();

  let res = await store.set({
    path: Path.fromStrings("chat", `~${authorKeypair.tag}`, `${Date.now()}`),
    identity: authorKeypair.tag,
    payload: new TextEncoder().encode(input.value),
  });
  if (isErr(res)) {
    console.error("set", res);
  }

  input.value = "";

  renderMessages();
});


// Read messages from chat.
const messages = document.getElementById("messages");

async function renderMessages() {
	messages.innerHTML = "";

  for await (const doc of store.queryDocs({
    pathPrefix: Path.fromStrings("chat"),
    // timestampGte: lastWeek,
    // limit: 10
  })) {
    const message = document.createElement("li");
    message.textContent = `${doc.identity}: ${new TextDecoder().decode(await doc.payload.bytes())} (${new Date(Number(doc.timestamp / 1000n))})`
    messages.append(message);
  }
}

// cache.onCacheUpdated(() => {
// 	renderMessages();
// });

renderMessages();

const res = await peer.syncHttp("http://localhost:8080");
if (isErr(res)) {
  console.error(res);
}
