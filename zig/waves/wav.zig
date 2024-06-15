const std = @import("std");

pub const audio_format_pcm = 1;
pub const audio_format_float = 3;

pub const FormatInfo = struct {
    size: u64,
    block_size: u32,
    audio_format: u16,
    num_channels: u16,
    sample_rate: u32,
    bytes_per_sec: u32,
    bytes_per_block: u16,
    bits_per_sample: u16,
};

pub const Error = error{
    NotAWAV,
    InvalidHeader,
    UnknownFormat,
};

pub fn parse_header(header_bytes: []const u8) Error!FormatInfo {
    if (header_bytes.len < 36) {
        return Error.InvalidHeader;
    }

    if (!std.mem.eql(u8, header_bytes[0..4], "RIFF")) {
        return Error.NotAWAV;
    }

    const size = std.mem.readVarInt(u32, header_bytes[4..8], .little) + 8;

    if (!std.mem.eql(u8, header_bytes[8..12], "WAVE")) {
        return Error.NotAWAV;
    }

    if (!std.mem.eql(u8, header_bytes[12..16], "fmt ")) {
        return Error.NotAWAV;
    }

    const block_size = std.mem.readVarInt(u32, header_bytes[16..20], .little);
    const audio_format = std.mem.readVarInt(u16, header_bytes[20..22], .little);
    if (!(audio_format == audio_format_pcm or audio_format == audio_format_float)) {
        return Error.UnknownFormat;
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
