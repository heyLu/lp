const std = @import("std");

const audio_format_pcm = 1;
const audio_format_float = 3;

pub fn main() !void {
    // var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    // defer arena.deinit();

    // const allocator = arena.allocator();

    std.debug.print("hi.\n", .{});

    const file = try std.fs.cwd().openFileZ(std.os.argv[1], .{});
    defer file.close();

    const stat = try file.stat();

    var header_buffer: [36]u8 = undefined;
    var bytes_read = try file.readAll(&header_buffer);
    std.debug.print("read {d} of {d} bytes.\n", .{ bytes_read, stat.size });
    // TODO: return error values instead
    std.debug.assert(bytes_read == 36);

    std.debug.print("{c}\n", .{header_buffer[0..bytes_read]});

    if (!std.mem.eql(u8, header_buffer[0..4], "RIFF")) {
        std.debug.print("not a .wav\n", .{});
        std.process.exit(1);
    }

    const size = std.mem.readVarInt(u32, header_buffer[4..8], .little) + 8;
    std.debug.assert(stat.size == size);

    if (!std.mem.eql(u8, header_buffer[8..12], "WAVE")) {
        std.debug.print("not a .wav\n", .{});
        std.process.exit(1);
    }

    if (!std.mem.eql(u8, header_buffer[12..16], "fmt ")) {
        std.debug.print("not a .wav\n", .{});
        std.process.exit(1);
    }

    const block_size = std.mem.readVarInt(u32, header_buffer[16..20], .little);
    const audio_format = std.mem.readVarInt(u16, header_buffer[20..22], .little);
    std.debug.assert(audio_format == audio_format_pcm or audio_format == audio_format_float);
    const num_channels = std.mem.readVarInt(u16, header_buffer[22..24], .little);
    const sample_rate = std.mem.readVarInt(u32, header_buffer[24..28], .little);
    const bytes_per_sec = std.mem.readVarInt(u32, header_buffer[28..32], .little);
    const bytes_per_block = std.mem.readVarInt(u16, header_buffer[32..34], .little);
    const bits_per_sample = std.mem.readVarInt(u16, header_buffer[34..36], .little);
    std.debug.print(
        \\blockSize     = {d: >10}
        \\audioFormat   = {d: >10}
        \\numChannels   = {d: >10}
        \\sampleRate    = {d: >10}
        \\bytesPerSec   = {d: >10}
        \\bytesPerBlock = {d: >10}
        \\bitsPerSample = {d: >10}
        \\--------------------------------------------------
        \\
    , .{ block_size, audio_format, num_channels, sample_rate, bytes_per_sec, bytes_per_block, bits_per_sample });

    var chunk_header: [8]u8 = undefined;
    while (bytes_read > 0) {
        bytes_read = try file.readAll(&chunk_header);
        std.debug.assert(bytes_read == 8);

        // FIXME: this fails on second read because we don't actually read the chunk
        std.debug.assert(std.mem.eql(u8, chunk_header[0..4], "data"));

        const data_size = std.mem.readVarInt(u32, chunk_header[4..8], .little);
        std.debug.print("chunkSize     = {d: >10}\n", .{data_size});

        // FIXME: read the chunk!  (in bytes_per_sec increments? 🤔)
    }
}
