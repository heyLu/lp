// Playing around with SDL + TTF.
//
// Based on `hello_sdl.c`, zig + SDL code from https://github.com/andrewrk/sdl-zig-demo.
//
// Resources:
// - http://wiki.libsdl.org/CategoryAPI
// - https://www.libsdl.org/projects/SDL_ttf/docs/SDL_ttf.html
//
// Unrelatedly, https://ziglang.org/learn/samples/ lead me to
// raylib, which looks like a really neat library to get started with
// game programming without Godot or Unity:
// https://github.com/raysan5/raylib
const c = @cImport({
    @cInclude("SDL2/SDL.h");
    @cInclude("SDL2/SDL_ttf.h");
});
const std = @import("std");

pub fn main() !void {
    var general_purpose_allocator = std.heap.GeneralPurposeAllocator(.{}){};
    const gpa = &general_purpose_allocator.allocator;
    const args = try std.process.argsAlloc(gpa);
    defer std.process.argsFree(gpa, args);

    if (c.SDL_Init(c.SDL_INIT_VIDEO) != 0) {
        c.SDL_Log("Unable to initialize SDL: %s", c.SDL_GetError());
        return error.SDLInitializationFailed;
    }
    defer c.SDL_Quit();

    if (c.TTF_Init() != 0) {
        c.SDL_Log("Unable to initialize SDL_ttf: %s", c.TTF_GetError());
        return error.TTFInitializationFailed;
    }
    defer c.TTF_Quit();

    var font_file = if (args.len > 1) args[1] else "./FantasqueSansMono-Regular.ttf";
    const font = c.TTF_OpenFont(font_file, 20) orelse {
        c.SDL_Log("Unable to load font: %s", c.TTF_GetError());
        return error.TTFInitializationFailed;
    };
    defer c.TTF_CloseFont(font);
    c.SDL_Log("Using font %s", font_file.ptr);

    const screen = c.SDL_CreateWindow("hello fonts", c.SDL_WINDOWPOS_CENTERED, c.SDL_WINDOWPOS_CENTERED, 600, 100, c.SDL_WINDOW_BORDERLESS) orelse {
        c.SDL_Log("Unable to create window: %s", c.SDL_GetError());
        return error.SDLInitializationFailed;
    };
    defer c.SDL_DestroyWindow(screen);

    var surface = c.SDL_GetWindowSurface(screen);

    // assume monospace font
    var glyph_width: c_int = 0;
    if (c.TTF_GlyphMetrics(font, 'g', null, null, null, null, &glyph_width) != 0) {
        c.SDL_Log("Unable to measure glyph: %s", c.TTF_GetError());
        return error.TTFInitializationFailed;
    }
    var glyph_height = c.TTF_FontLineSkip(font);

    var msg = "howdy there, enby! 🐘                                          ".*;
    var pos: usize = 0;
    var max_chars = std.math.min(@divTrunc(@intCast(usize, surface.*.w), @intCast(usize, glyph_width)), msg.len);

    var result: ?[*:0]u8 = null;
    result = "";

    const keyboardState = c.SDL_GetKeyboardState(null);

    c.SDL_StartTextInput();

    var quit = false;
    while (!quit) {
        var event: c.SDL_Event = undefined;
        while (c.SDL_PollEvent(&event) != 0) {
            const ctrlPressed = (keyboardState[c.SDL_SCANCODE_LCTRL] != 0);
            switch (event.@"type") {
                c.SDL_QUIT => {
                    quit = true;
                },
                c.SDL_WINDOWEVENT => {
                    switch (event.window.event) {
                        c.SDL_WINDOWEVENT_SIZE_CHANGED => {
                            surface = c.SDL_GetWindowSurface(screen);
                        },
                        else => {},
                    }
                },
                c.SDL_KEYDOWN => {
                    if (ctrlPressed) {
                        switch (event.key.keysym.sym) {
                            c.SDLK_a => {
                                pos = 0;
                                msg[pos] = '_';
                            },
                            c.SDLK_k => {
                                var i: usize = 0;
                                while (i < max_chars) : (i += 1) {
                                    msg[i] = ' ';
                                }
                                msg[max_chars] = 0;
                                pos = 0;
                            },
                            else => {},
                        }
                    } else {
                        switch (event.key.keysym.sym) {
                            c.SDLK_ESCAPE => {
                                quit = true;
                            },
                            c.SDLK_BACKSPACE => {
                                pos = if (pos == 0) max_chars - 1 else (pos - 1) % (max_chars - 1);
                                msg[pos] = '_';
                            },
                            c.SDLK_RETURN => {
                                result = try runCommand(&msg, gpa);
                                var i: usize = 0;
                                while (i < max_chars) : (i += 1) {
                                    msg[i] = ' ';
                                }
                                msg[max_chars] = 0;
                                pos = 0;
                            },
                            else => {},
                        }
                    }
                },
                c.SDL_TEXTINPUT => {
                    if (!ctrlPressed and event.text.text.len > 0) {
                        c.SDL_Log("input: '%s' at %d", event.text.text, pos);
                        msg[pos] = event.text.text[0];
                        pos = (pos + 1) % (max_chars - 1);
                    }
                },
                else => {},
            }
        }

        // thanks to https://stackoverflow.com/questions/22886500/how-to-render-text-in-sdl2 for some actually useful code here
        const white: c.SDL_Color = c.SDL_Color{ .r = 255, .g = 255, .b = 255, .a = 255 };
        const black: c.SDL_Color = c.SDL_Color{ .r = 0, .g = 0, .b = 0, .a = 255 };
        // Shaded vs Solid gives a nicer output, with solid the output
        // was squiggly and not aligned with a baseline.
        const text = c.TTF_RenderUTF8_Shaded(font, &msg, white, black);
        _ = c.SDL_BlitSurface(text, null, surface, null);

        const result_text = c.TTF_RenderUTF8_Shaded(font, result, white, black);
        _ = c.SDL_BlitSurface(result_text, null, surface, &c.SDL_Rect{ .x = 0, .y = glyph_height, .w = surface.*.w, .h = surface.*.h - glyph_height });

        _ = c.SDL_UpdateWindowSurface(screen);

        c.SDL_Delay(16);
    }
}

fn runCommand(raw_cmd: []const u8, allocator: *std.mem.Allocator) !?[*:0]u8 {
    const cmd = std.mem.trim(u8, std.mem.sliceTo(raw_cmd, 0), &std.ascii.spaces);
    const argv = if (std.mem.startsWith(u8, cmd, "go "))
        &[_][]const u8{ "go", "doc", cmd[3..] }
    else
        &[_][]const u8{ "/usr/bin/qalc", "-terse", cmd };
    for (argv) |arg| {
        std.debug.print("'{s}' ", .{arg});
    }
    const result = try std.ChildProcess.exec(.{ .allocator = allocator, .argv = argv });
    const buf = try allocator.allocSentinel(u8, 100, 0);
    std.mem.copy(u8, buf, result.stdout[0..std.math.min(100, result.stdout.len)]);
    var i: usize = result.stdout.len;
    while (i < buf.len) : (i += 1) {
        buf[i] = ' ';
    }
    buf[buf.len - 1] = 0;
    std.debug.print("stderr: '{s}'\n", .{result.stderr});
    defer {
        allocator.free(result.stdout);
        allocator.free(result.stderr);
    }
    return buf;
}

// tests

test "trim []const u8" {
    const untrimmed: []const u8 = "   hey there   ";
    const trimmed = std.mem.trim(u8, untrimmed, &std.ascii.spaces);
    try std.testing.expect(trimmed.len < untrimmed.len);
    try std.testing.expect(trimmed.len == 9);
    try std.testing.expect(std.mem.eql(u8, trimmed, "hey there"));
}

test "trim [*:0]const u8" {
    const untrimmed: [*:0]const u8 = "   hey there   ";
    const to_trim: [*:0]const u8 = " ";
    const trimmed = std.mem.trim(u8, std.mem.sliceTo(untrimmed, 0), std.mem.sliceTo(to_trim, 0));
    try std.testing.expect(std.mem.len(trimmed) < std.mem.len(untrimmed));
    try std.testing.expect(trimmed.len == 9);
    try std.testing.expect(std.mem.eql(u8, trimmed, "hey there"));
}
