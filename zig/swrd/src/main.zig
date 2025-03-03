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

const Swrd = struct {
    x: f32,
    y: f32,
    size: f32,
    dir_x: f32,
    dir_y: f32,

    fn draw(self: Swrd, renderer: *c.SDL_Renderer) !void {
        try errify(c.SDL_RenderPoint(renderer, self.x, self.y));
        try errify(c.SDL_RenderRect(renderer, &c.SDL_FRect{ .x = self.x - self.size / 2, .y = self.y - self.size / 2, .w = self.size, .h = self.size }));

        try errify(c.SDL_RenderPoint(renderer, self.x + self.dir_x * 20, self.y + self.dir_y * 20));
        // try errify(c.SDL_RenderLine(renderer, self.x, self.y, self.x + self.dir_x * 20, self.y + self.dir_y * 20));
    }

    fn move(self: *Swrd, x: f32, y: f32) void {
        self.x += x;
        self.y += y;
    }
};

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

    var quit = false;
    var paused = false;

    const sample_rate = 44100;
    const sounds_spec = c.SDL_AudioSpec{ .format = c.SDL_AUDIO_F32, .channels = 1, .freq = sample_rate };
    const audio_stream = try errify(c.SDL_OpenAudioDeviceStream(c.SDL_AUDIO_DEVICE_DEFAULT_PLAYBACK, &sounds_spec, null, null));
    defer c.SDL_DestroyAudioStream(audio_stream);
    try errify(c.SDL_ResumeAudioStreamDevice(audio_stream));

    const audio_data = try allocator.alloc(f32, 1024);
    var freq: f32 = 440;
    var volume: f32 = 0.5;
    const audio_thread = try std.Thread.spawn(.{}, doAudio, .{ &quit, sounds_spec, audio_stream, audio_data, &freq, &volume });
    audio_thread.detach();

    var rrnd = std.Random.DefaultPrng.init(0);
    const rnd = std.Random.DefaultPrng.random(&rrnd);

    var swrd = Swrd{ .x = @as(f32, @floatFromInt(window_w)) / 2, .y = @as(f32, @floatFromInt(window_h)) / 2, .size = 10, .dir_x = 0, .dir_y = -1 };

    try errify(c.SDL_SetRenderDrawColor(renderer, 255, 255, 255, 255));

    const render_budget = std.time.us_per_s / 120;
    var last_render = std.time.microTimestamp();

    var move_x: f32 = 0;
    var move_y: f32 = 0;
    while (!quit) {
        try errify(c.SDL_SetRenderDrawColor(renderer, 0, 0, 0, 255));
        try errify(c.SDL_RenderClear(renderer));
        try errify(c.SDL_SetRenderDrawColor(renderer, 255, 255, 255, 255));

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

                    const max_volume = 1.0;
                    volume = event.motion.y / @as(f32, @floatFromInt(window_h)) * max_volume;
                    volume = max_volume - @min(volume, max_volume);

                    try errify(c.SDL_RenderLine(renderer, event.motion.x, 0, event.motion.x, @floatFromInt(window_h)));
                },
                c.SDL_EVENT_WINDOW_FOCUS_GAINED => paused = false,
                c.SDL_EVENT_WINDOW_FOCUS_LOST => paused = true,
                c.SDL_EVENT_KEY_DOWN => {
                    std.log.debug("key: {}", .{event.key.key});
                    switch (event.key.key) {
                        c.SDLK_SPACE => freq = 400,
                        c.SDLK_ESCAPE => paused = !paused,
                        c.SDLK_UP, c.SDLK_W => move_y = -1,
                        c.SDLK_DOWN, c.SDLK_S => move_y = 1,
                        c.SDLK_LEFT, c.SDLK_A => move_x = -1,
                        c.SDLK_RIGHT, c.SDLK_D => move_x = 1,
                        else => {},
                    }
                },
                c.SDL_EVENT_KEY_UP => {
                    switch (event.key.key) {
                        c.SDLK_UP, c.SDLK_DOWN, c.SDLK_W, c.SDLK_S => move_y = 0,
                        c.SDLK_LEFT, c.SDLK_RIGHT, c.SDLK_A, c.SDLK_D => move_x = 0,
                        else => {},
                    }
                },
                c.SDL_EVENT_GAMEPAD_ADDED => {
                    const pad = try errify(c.SDL_OpenGamepad(event.gdevice.which));
                    std.log.debug("{}", .{pad});
                },
                c.SDL_EVENT_GAMEPAD_AXIS_MOTION => {
                    switch (event.gaxis.axis) {
                        c.SDL_GAMEPAD_AXIS_LEFTX => {
                            move_x = @as(f32, @floatFromInt(event.gaxis.value)) / std.math.maxInt(i16);
                            swrd.dir_x = move_x;
                        },
                        c.SDL_GAMEPAD_AXIS_LEFTY => {
                            move_y = @as(f32, @floatFromInt(event.gaxis.value)) / std.math.maxInt(i16);
                            swrd.dir_y = move_y;
                        },
                        c.SDL_GAMEPAD_AXIS_RIGHTX => swrd.dir_x = @as(f32, @floatFromInt(event.gaxis.value)) / std.math.maxInt(i16),
                        c.SDL_GAMEPAD_AXIS_RIGHTY => swrd.dir_y = @as(f32, @floatFromInt(event.gaxis.value)) / std.math.maxInt(i16),
                        else => {},
                    }
                    if (@abs(move_x) < 0.05) move_x = 0;
                    if (@abs(move_y) < 0.05) move_y = 0;
                },
                c.SDL_EVENT_GAMEPAD_BUTTON_DOWN => {
                    switch (event.gbutton.button) {
                        c.SDL_GAMEPAD_BUTTON_START => paused = !paused,
                        else => {},
                    }
                },
                c.SDL_EVENT_GAMEPAD_UPDATE_COMPLETE => {},
                c.SDL_EVENT_JOYSTICK_AXIS_MOTION, c.SDL_EVENT_JOYSTICK_UPDATE_COMPLETE, c.SDL_EVENT_JOYSTICK_BUTTON_DOWN, c.SDL_EVENT_JOYSTICK_BUTTON_UP => {},
                c.SDL_EVENT_GAMEPAD_BUTTON_UP => {},
                c.SDL_EVENT_WINDOW_RESIZED => {
                    window_w = event.window.data1;
                    window_h = event.window.data2;
                },
                else => std.log.debug("unhandled event {} ({})", .{ event.type, event }),
            }
        }
        // try errify(c.SDL_RenderDebugTextFormat(renderer, @as(f32, @floatFromInt(window_w)) - 130, 5, "%.5f %.5f", swrd.dir_x, swrd.dir_y));

        if (paused) {
            try errify(c.SDL_SetAudioStreamGain(audio_stream, 0.0));
            std.time.sleep(10 * 1000 * 1000);
            continue;
        }
        try errify(c.SDL_SetAudioStreamGain(audio_stream, 1.0));

        try errify(c.SDL_RenderPoint(renderer, std.Random.float(rnd, f32) * @as(f32, @floatFromInt(window_w)), std.Random.float(rnd, f32) * @as(f32, @floatFromInt(window_h))));

        const wave_height = 20;
        for (0..audio_data.len) |i| {
            try errify(c.SDL_RenderPoint(renderer, @as(f32, @floatFromInt(i)) / @as(f32, @floatFromInt(audio_data.len)) * @as(f32, @floatFromInt(window_w)), @as(f32, @floatFromInt(window_h - 1)) - (wave_height / 2 + audio_data[i] * wave_height)));
        }

        swrd.move(move_x, move_y);
        try swrd.draw(renderer);

        // make the window appear
        try errify(c.SDL_RenderPresent(renderer));

        const now = std.time.microTimestamp();
        const render_time = now - last_render;
        last_render = now;
        // std.log.debug("render_time: {}, sleep_time: {}", .{ std.fmt.fmtDurationSigned(render_time * 1000), std.fmt.fmtDurationSigned((render_budget - render_time) * 1000) });
        std.time.sleep(@intCast(@max(0, (render_budget - render_time) * 1000)));
    }
}

fn doAudio(quit: *bool, sounds_spec: c.SDL_AudioSpec, audio_stream: *c.SDL_AudioStream, audio_data: []f32, freq: *f32, volume: *f32) !void {
    errdefer |err| if (err == error.SdlError) {
        std.log.err("SDL error: {s}", .{c.SDL_GetError()});
        std.process.exit(1);
    };

    const sample_rate: f32 = @floatFromInt(sounds_spec.freq);

    // const minimum_audio = sample_rate * @sizeOf(f32) / 2;
    const minimum_audio = audio_data.len * @sizeOf(f32) * 2;

    var current_sine_sample: i32 = 0;
    var last_freq = freq.*;
    while (!quit.*) {
        const queued = c.SDL_GetAudioStreamQueued(audio_stream);
        if (queued < minimum_audio) {
            const current_freq = freq.*;

            var start: usize = 0;
            if (@abs(last_freq - current_freq) > 0.01) {
                audio_data[0] = audio_data[audio_data.len - 1];
                const numSteps = 300;
                const step = audio_data[0] / numSteps;
                for (1..numSteps) |i| {
                    start = i;

                    if (@abs(audio_data[i - 1]) < step * 2) {
                        audio_data[i] = 0;
                        break;
                    }

                    // audio_data[i] = audio_data[i - 1] - step;
                    audio_data[i] = audio_data[0] * (1 - @as(f32, @floatFromInt(i)) / @as(f32, @floatFromInt(numSteps)));
                }
                std.log.debug("different! {d:.5} -> {d:.5}; {}*{d:.5}, -1={d:.5}, 0={d:.5}, {}={d:.5}", .{ last_freq, current_freq, numSteps, step, audio_data[0], audio_data[1], start, audio_data[start] });
                // start += 1;

                last_freq = current_freq;
                current_sine_sample = 0;
            }

            // std.log.debug("audio {} {} {} {}", .{ queued, minimum_audio, current_sine_sample, (audio_data.len * @sizeOf(f32)) });
            for (start..audio_data.len) |i| {
                const phase = @as(f32, @floatFromInt(current_sine_sample)) * current_freq / sample_rate;
                audio_data[i] = c.SDL_sinf(phase * 2 * c.SDL_PI_F) * volume.*;
                current_sine_sample += 1;
            }
            current_sine_sample = @mod(current_sine_sample, sounds_spec.freq);

            try errify(c.SDL_PutAudioStreamData(audio_stream, audio_data.ptr, @intCast(audio_data.len * @sizeOf(f32))));
            // try errify(c.SDL_FlushAudioStream(audio_stream));
        }
        std.time.sleep(1 * 1000 * 1000);
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
