const std = @import("std");

const c = @cImport({
    @cDefine("SDL_DISABLE_OLD_NAMES", {});
    @cInclude("SDL3/SDL.h");
    @cInclude("SDL3/SDL_revision.h");
    // For programs that provide their own entry points instead of relying on SDL's main function
    // macro magic, 'SDL_MAIN_HANDLED' should be defined before including 'SDL_main.h'.
    @cDefine("SDL_MAIN_HANDLED", {});
    @cInclude("SDL3/SDL_main.h");
});

pub fn main() !void {
    errdefer |err| if (err == error.SdlError) std.log.err("SDL error: {s}", .{c.SDL_GetError()});

    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();

    c.SDL_SetMainReady();

    try errify(c.SDL_SetAppMetadata("swrd", "0.0.1", "org.papill0n.swrd"));
    try errify(c.SDL_Init(c.SDL_INIT_VIDEO | c.SDL_INIT_AUDIO | c.SDL_INIT_GAMEPAD));
    defer c.SDL_Quit();

    var window_w: i32 = 640;
    var window_h: i32 = 480;
    const window: *c.SDL_Window, const renderer: *c.SDL_Renderer = create_window_and_renderer: {
        var window: ?*c.SDL_Window = null;
        var renderer: ?*c.SDL_Renderer = null;
        try errify(c.SDL_CreateWindowAndRenderer("swrd", window_w, window_h, 0, &window, &renderer));
        errdefer comptime unreachable;

        break :create_window_and_renderer .{ window.?, renderer.? };
    };

    defer c.SDL_DestroyRenderer(renderer);
    defer c.SDL_DestroyWindow(window);

    std.log.debug("hello?", .{});

    var freq: f32 = 440;
    const audio_thread = try std.Thread.spawn(.{}, doAudio, .{ allocator, &freq });
    audio_thread.detach();

    var rrnd = std.Random.DefaultPrng.init(0);
    const rnd = std.Random.DefaultPrng.random(&rrnd);

    try errify(c.SDL_SetRenderDrawColor(renderer, 255, 255, 255, 255));

    var quit = false;
    while (!quit) {
        var event: c.SDL_Event = undefined;
        while (c.SDL_PollEvent(&event)) {
            switch (event.type) {
                c.SDL_EVENT_QUIT => quit = true,
                c.SDL_EVENT_MOUSE_MOTION => {
                    const max_freq = 500;
                    freq = round_frequency(event.motion.x / @as(f32, @floatFromInt(window_w)) * max_freq);
                    freq = @min(freq, max_freq);
                    // rrnd = std.Random.DefaultPrng.init(@intFromFloat(event.motion.x));
                    std.log.debug("freq: {}", .{freq});

                    try errify(c.SDL_SetRenderDrawColor(renderer, 0, 0, 0, 255));
                    try errify(c.SDL_RenderClear(renderer));
                    try errify(c.SDL_SetRenderDrawColor(renderer, 255, 255, 255, 255));
                    try errify(c.SDL_RenderLine(renderer, event.motion.x, 0, event.motion.x, @floatFromInt(window_h)));
                },
                c.SDL_EVENT_KEY_DOWN => {
                    std.log.debug("key: {}", .{event.key.key});
                    if (event.key.key == c.SDLK_SPACE) {
                        freq = 400;
                    }
                },
                c.SDL_EVENT_WINDOW_RESIZED => {
                    window_w = event.window.data1;
                    window_h = event.window.data2;
                },
                else => std.log.debug("unhandled event {}", .{event.type}),
            }
        }

        try errify(c.SDL_RenderPoint(renderer, std.Random.float(rnd, f32) * @as(f32, @floatFromInt(window_w)), std.Random.float(rnd, f32) * @as(f32, @floatFromInt(window_h))));

        // make the window appear
        try errify(c.SDL_RenderPresent(renderer));

        std.time.sleep(1 * 1000 * 1000);
    }
}

fn doAudio(allocator: std.mem.Allocator, _: *f32) !void {
    errdefer |err| if (err == error.SdlError) {
        std.log.err("SDL error: {s}", .{c.SDL_GetError()});
        std.process.exit(1);
    };

    const sample_rate = 2000;
    const sounds_spec = c.SDL_AudioSpec{ .format = c.SDL_AUDIO_F32, .channels = 1, .freq = sample_rate };
    const audio_stream = try errify(c.SDL_OpenAudioDeviceStream(c.SDL_AUDIO_DEVICE_DEFAULT_PLAYBACK, &sounds_spec, null, null));
    defer c.SDL_DestroyAudioStream(audio_stream);

    try errify(c.SDL_ResumeAudioStreamDevice(audio_stream));

    const minimum_audio = sample_rate / @sizeOf(f32) / 2;
    var audio_data = try allocator.alloc(f32, 512);

    const current_freq: f32 = 123;
    var current_sine_sample: i32 = 0;
    while (true) {
        const queued = c.SDL_GetAudioStreamQueued(audio_stream);
        if (queued < minimum_audio) {
            // current_sine_sample = 0;
            // current_freq = freq.*;

            // std.log.debug("audio {} {} {} ", .{ queued, minimum_audio, current_sine_sample });
            for (0..audio_data.len) |i| {
                const phase: f32 = @as(f32, @floatFromInt(current_sine_sample)) * current_freq / sample_rate;
                audio_data[i] = c.SDL_sinf(phase * 2 * c.SDL_PI_F);
                current_sine_sample += 1;
            }
            current_sine_sample = @mod(current_sine_sample, sample_rate);

            try errify(c.SDL_PutAudioStreamData(audio_stream, audio_data.ptr, @intCast(audio_data.len)));
            // try errify(c.SDL_FlushAudioStream(audio_stream));

            // try errify(c.SDL_SetRenderDrawColor(renderer, 0, 0, 0, 255));
            // try errify(c.SDL_RenderClear(renderer));
            // try errify(c.SDL_SetRenderDrawColor(renderer, 255, 255, 255, 255));
            // for (0..audio_data.len) |i| {
            //     try errify(c.SDL_RenderPoint(renderer, @floatFromInt(i), audio_data[i] * 100));
            // }
        }
        std.time.sleep(1000 * 1000);
    }
}

fn round_frequency(f: f32) f32 {
    return @as(f32, @floatFromInt(@as(i32, @intFromFloat(f * 100)))) / 100;
}

inline fn errify(value: anytype) error{SdlError}!switch (@typeInfo(@TypeOf(value))) {
    .Bool => void,
    .Pointer, .Optional => @TypeOf(value.?),
    .Int => |info| switch (info.signedness) {
        .signed => @TypeOf(@max(0, value)),
        .unsigned => @TypeOf(value),
    },
    else => @compileError("unhandled type: " ++ @typeName(@TypeOf(value))),
} {
    return switch (@typeInfo(@TypeOf(value))) {
        .Bool => if (!value) error.SdlError,
        .Pointer, .Optional => value orelse error.SdlError,
        .Int => |info| switch (info.signedness) {
            .signed => if (value >= 0) @max(0, value) else error.SdlError,
            .unsigned => if (value != 0) value else error.SdlError,
        },
        else => comptime unreachable,
    };
}
