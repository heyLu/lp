#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <strings.h>
#include <time.h>
#include <unistd.h>

#include <sys/socket.h>
#include <sys/un.h>

#define MSG_LEN 100
#define SOCKET_PATH "/tmp/test-socket.sock"

int sock;

void cleanup() {
	if (close(sock)) {
		perror("close");
	}

	if (unlink(SOCKET_PATH)) {
		perror("unlink");
	}
}

static void handle_signal(int signum) {
	switch(signum) {
	case SIGINT:
		cleanup();
		exit(0);
	default:
		fprintf(stderr, "unhandled signal: %d", signum);
		exit(1);
	}
}

int main(int argc, char **argv) {
	sock = socket(AF_UNIX, SOCK_STREAM, 0);

	struct sockaddr_un addr;
	addr.sun_family = AF_UNIX;
	strncpy(addr.sun_path, SOCKET_PATH, sizeof(addr.sun_path) - 1);
	if (bind(sock, (struct sockaddr*) &addr, sizeof(addr))) {
		perror("bind");
		goto err;
	}

	struct sigaction sa;
	sa.sa_handler = handle_signal;
	sigemptyset(&sa.sa_mask);
	sa.sa_flags = SA_RESTART;
	if (sigaction(SIGINT, &sa, NULL) < 0) {
		perror("sigaction");
	}

	printf("listening on %s\n", SOCKET_PATH);
	if (listen(sock, 0) < 0) {
		perror("listen");
		goto err;
	}

	for (;;) {
		int conn = accept(sock, NULL, NULL);
		if (conn < 0) {
			perror("accept");
			goto err;
		}

		char msg[MSG_LEN];
		time_t now = time(NULL);
		int len = strftime(msg, MSG_LEN, "%Y-%m-%d %H:%M:%S %z\n", localtime(&now));
		if (write(conn, msg, len) < 0) {
			perror("write");
		}

		if (close(conn) < 0) {
			perror("close conn");
		}
	}

err:
	cleanup();
	return 0;
}
