// see also: https://en.wikipedia.org/wiki/WAV
// see also: http://www.tactilemedia.com/info/MCI_Control_Info.html

const std = @import("std");

const audio_format_pcm = 1;
const audio_format_float = 3;

pub fn main() !void {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();

    const allocator = arena.allocator();

    std.debug.print("hi.\n", .{});

    const file = try std.fs.cwd().openFileZ(std.os.argv[1], .{});
    defer file.close();

    const stat = try file.stat();

    var header_buffer: [36]u8 = undefined;
    var bytes_read = try file.readAll(&header_buffer);

    const info = try parse_header(&header_buffer);

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

            for (0..chunk_data[0..bytes_read].len / info.bytes_per_block) |block_idx| {
                const sample_left = std.mem.readVarInt(u32, chunk_data[block_idx * info.bytes_per_block .. block_idx * info.bytes_per_block + sample_size], .little);
                const sample_right = std.mem.readVarInt(u32, chunk_data[block_idx * info.bytes_per_block + sample_size .. block_idx * info.bytes_per_block + info.bytes_per_block], .little);
                std.debug.assert(sample_left < max_sample);
                std.debug.assert(sample_right < max_sample);
                // if (block_idx % 1000 == 0) {
                //     std.debug.print("{d} {d}\n", .{ sample_left, sample_right });
                // }

                samples_total += 1;
            }
        }
    }

    std.debug.print("done?  {d}sec of delicious audio, {d} samples\n", .{ data_total / info.bytes_per_sec, samples_total });
}

const FormatInfo = struct {
    size: u64,
    block_size: u32,
    audio_format: u16,
    num_channels: u16,
    sample_rate: u32,
    bytes_per_sec: u32,
    bytes_per_block: u16,
    bits_per_sample: u16,
};

const WavError = error{
    NotAWAV,
    InvalidHeader,
    UnknownFormat,
};

fn parse_header(header_bytes: []u8) WavError!FormatInfo {
    if (header_bytes.len < 36) {
        return WavError.InvalidHeader;
    }

    if (!std.mem.eql(u8, header_bytes[0..4], "RIFF")) {
        return WavError.NotAWAV;
    }

    const size = std.mem.readVarInt(u32, header_bytes[4..8], .little) + 8;

    if (!std.mem.eql(u8, header_bytes[8..12], "WAVE")) {
        return WavError.NotAWAV;
    }

    if (!std.mem.eql(u8, header_bytes[12..16], "fmt ")) {
        return WavError.NotAWAV;
    }

    const block_size = std.mem.readVarInt(u32, header_bytes[16..20], .little);
    const audio_format = std.mem.readVarInt(u16, header_bytes[20..22], .little);
    if (!(audio_format == audio_format_pcm or audio_format == audio_format_float)) {
        return WavError.UnknownFormat;
    }

    const num_channels = std.mem.readVarInt(u16, header_bytes[22..24], .little);
    const sample_rate = std.mem.readVarInt(u32, header_bytes[24..28], .little);
    const bytes_per_sec = std.mem.readVarInt(u32, header_bytes[28..32], .little);
    const bytes_per_block = std.mem.readVarInt(u16, header_bytes[32..34], .little);
    const bits_per_sample = std.mem.readVarInt(u16, header_bytes[34..36], .little);

    return FormatInfo{
        .size = size,
        .block_size = block_size,
        .audio_format = audio_format,
        .num_channels = num_channels,
        .sample_rate = sample_rate,
        .bytes_per_sec = bytes_per_sec,
        .bytes_per_block = bytes_per_block,
        .bits_per_sample = bits_per_sample,
    };
}
