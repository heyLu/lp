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
    const font = c.TTF_OpenFont(font_file, 16) orelse {
        c.SDL_Log("Unable to load font: %s", c.TTF_GetError());
        return error.TTFInitializationFailed;
    };
    defer c.TTF_CloseFont(font);
    c.SDL_Log("Using font %s", font_file.ptr);

    // assume monospace font
    var glyph_width: c_int = 0;
    if (c.TTF_GlyphMetrics(font, 'g', null, null, null, null, &glyph_width) != 0) {
        c.SDL_Log("Unable to measure glyph: %s", c.TTF_GetError());
        return error.TTFInitializationFailed;
    }
    var glyph_height = c.TTF_FontLineSkip(font);

    var window_width = glyph_width * 100;
    var window_height = glyph_height * 20;
    const screen = c.SDL_CreateWindow("hello fonts", c.SDL_WINDOWPOS_CENTERED, c.SDL_WINDOWPOS_CENTERED, window_width, window_height, c.SDL_WINDOW_BORDERLESS) orelse {
        c.SDL_Log("Unable to create window: %s", c.SDL_GetError());
        return error.SDLInitializationFailed;
    };
    defer c.SDL_DestroyWindow(screen);

    const renderer = c.SDL_CreateRenderer(screen, -1, c.SDL_RENDERER_ACCELERATED);

    var msg = "                                                                                                    ".*;
    var pos: usize = 0;
    var max_chars = std.math.min(@divTrunc(@intCast(usize, window_width), @intCast(usize, glyph_width)), msg.len);

    var result: []const u8 = try gpa.alloc(u8, 0);

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
                            window_width = event.window.data1;
                            window_height = event.window.data2;
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
                                gpa.free(result);
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

        _ = c.SDL_RenderClear(renderer);

        // thanks to https://stackoverflow.com/questions/22886500/how-to-render-text-in-sdl2 for some actually useful code here
        const white: c.SDL_Color = c.SDL_Color{ .r = 255, .g = 255, .b = 255, .a = 255 };
        const black: c.SDL_Color = c.SDL_Color{ .r = 0, .g = 0, .b = 0, .a = 255 };
        // Shaded vs Solid gives a nicer output, with solid the output
        // was squiggly and not aligned with a baseline.
        const text = c.TTF_RenderUTF8_Shaded(font, &msg, white, black);
        const texture = c.SDL_CreateTextureFromSurface(renderer, text);
        c.SDL_FreeSurface(text);
        _ = c.SDL_RenderCopy(renderer, texture, null, &c.SDL_Rect{ .x = 0, .y = 0, .w = @intCast(c_int, msg.len) * glyph_width, .h = glyph_height });

        var i: c_int = 1;
        var lines = std.mem.split(result, "\n");
        var line = lines.next();
        while (line != null) {
            const line_c = try gpa.dupeZ(u8, line.?);
            const result_text = c.TTF_RenderUTF8_Shaded(font, line_c, white, black);
            gpa.free(line_c);
            const result_texture = c.SDL_CreateTextureFromSurface(renderer, result_text);
            _ = c.SDL_RenderCopy(renderer, result_texture, null, &c.SDL_Rect{ .x = 0, .y = i * glyph_height, .w = @intCast(c_int, line.?.len) * glyph_width, .h = glyph_height });
            c.SDL_FreeSurface(result_text);
            c.SDL_DestroyTexture(result_texture);

            i += 1;
            line = lines.next();
        }

        _ = c.SDL_RenderPresent(renderer);

        c.SDL_Delay(16);
    }

    _ = general_purpose_allocator.detectLeaks();
}

fn runCommand(raw_cmd: []const u8, allocator: *std.mem.Allocator) ![]const u8 {
    const cmd = std.mem.trim(u8, std.mem.sliceTo(raw_cmd, 0), &std.ascii.spaces);
    const argv = if (std.mem.startsWith(u8, cmd, "go "))
        &[_][]const u8{ "go", "doc", cmd[3..] }
    else if (std.mem.startsWith(u8, cmd, "py "))
        &[_][]const u8{ "python", "-c", try std.fmt.allocPrint(allocator, "import {s}; help({s});", .{ std.mem.sliceTo(cmd["py ".len..], '.'), cmd["py ".len..] }) }
    else if (std.mem.endsWith(u8, cmd, " --help"))
        // TODO: handle --help that outputs on stderr
        &[_][]const u8{ cmd[0..(cmd.len - " --help".len)], "--help" }
    else if (std.mem.startsWith(u8, cmd, "man "))
        // TODO: handle `man 5 sway`
        &[_][]const u8{ "man", cmd["man ".len..] }
    else
        &[_][]const u8{ "/usr/bin/qalc", "-terse", cmd };
    for (argv) |arg| {
        std.debug.print("'{s}' ", .{arg});
    }
    const result = try std.ChildProcess.exec(.{ .allocator = allocator, .argv = argv });
    std.debug.print("stderr: '{s}'\n", .{result.stderr});
    defer {
        // FIXME: allocator.free(result.stdout);
        allocator.free(result.stderr);
    }

    return result.stdout;
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
