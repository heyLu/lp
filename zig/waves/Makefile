all: waves wav.wasm

waves: waves.zig wav.zig
	zig build-exe waves.zig

wav.wasm: wav.zig
	zig build-exe wav.zig -target wasm32-freestanding-musl -rdynamic -fno-entry --import-memory --stack 8192 --initial-memory=65536 --max-memory=65536000 --global-base=6560 --export=parse_header_c -O ReleaseSmall
