// see also: https://en.wikipedia.org/wiki/WAV
// see also: http://www.tactilemedia.com/info/MCI_Control_Info.html

const std = @import("std");

const wav = @import("wav.zig");

pub fn main() !void {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();

    const allocator = arena.allocator();

    const file = try std.fs.cwd().openFileZ(std.os.argv[1], .{});
    defer file.close();

    const stat = try file.stat();

    var header_buffer: [36]u8 = undefined;
    var bytes_read = try file.readAll(&header_buffer);

    const info = try wav.parse_header(&header_buffer);

    std.debug.assert(info.size == stat.size);

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
    , .{ info.block_size, info.audio_format, info.num_channels, info.sample_rate, info.bytes_per_sec, info.bytes_per_block, info.bits_per_sample });

    const dataStart = try file.getPos();

    const stdout = std.io.getStdOut().writer();

    var chunk_header: [8]u8 = undefined;
    // var chunk_data: [1024]u8 = undefined;
    var chunk_data = try allocator.alloc(u8, info.bytes_per_sec);
    var data_total: u64 = 0;
    var samples_total: u64 = 0;
    while (bytes_read > 0) {
        bytes_read = try file.readAll(&chunk_header);
        if (bytes_read == 0) {
            break;
        }

        if (bytes_read != 8 or !std.mem.eql(u8, chunk_header[0..4], "data")) {
            std.debug.print("unknown chunk header {c} {d} bytes\n", .{ chunk_header, bytes_read });
            break;
        }

        std.debug.assert(bytes_read == 8);

        // FIXME: this fails on second read because we don't actually read the chunk
        std.debug.assert(std.mem.eql(u8, chunk_header[0..4], "data"));

        const data_size = std.mem.readVarInt(u32, chunk_header[4..8], .little);
        std.debug.print("chunkSize     = {d: >10}\n", .{data_size});
        data_total += data_size;

        // std.debug.print("{}\n", .{@as(f32, @floatFromInt(data_size)) / @as(f32, @floatFromInt(chunk_data.len))});
        const sample_size = info.bytes_per_block / info.num_channels;
        const max_sample = 16777216; // 2 << bits_per_sample - 1
        for (0..@divFloor(data_size, chunk_data.len) + 1) |_| {
            bytes_read = try file.readAll(chunk_data);
            // std.debug.assert(bytes_read == chunk_data.len);

            // std.debug.print("read {d}\n", .{bytes_read});

            var local_max: u32 = 0;
            for (0..chunk_data[0..bytes_read].len / info.bytes_per_block) |block_idx| {
                const sample_left = std.mem.readVarInt(u32, chunk_data[block_idx * info.bytes_per_block .. block_idx * info.bytes_per_block + sample_size], .little);
                const sample_right = std.mem.readVarInt(u32, chunk_data[block_idx * info.bytes_per_block + sample_size .. block_idx * info.bytes_per_block + info.bytes_per_block], .little);
                std.debug.assert(sample_left < max_sample);
                std.debug.assert(sample_right < max_sample);
                // if (block_idx % 1000 == 0) {
                //     std.debug.print("{d} {d}\n", .{ sample_left, sample_right });
                // }

                local_max = @max(local_max, sample_left);
                if (block_idx % (info.size / 80) == 0) {
                    try stdout.print("{d} ", .{local_max});
                }

                samples_total += 1;
            }
        }
    }
    try stdout.print("\n", .{});

    std.debug.print("done?  {d}sec of delicious audio, {d} samples\n", .{ data_total / info.bytes_per_sec, samples_total });

    try file.seekTo(dataStart);
}
