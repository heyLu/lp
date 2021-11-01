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

// TODO: instant output (for some commands like `py ...`, `go ...`, qcalc)

// commands wishlist:
// - search (e.g. default current dir + /usr/include)
// - launch with logs (default launcher, use systemd-run --user --unit=name name?)
// - view (the above logs)
// - switch to window
// - open url
// - open shortcuts (logs -> ..., tickets)
// - history (could be another command + some special keybindings)

// output line-by-line -> saved by caller?
// output can be reset
// incremental output vs. final output/action

const ProcessWithOutput = struct {
    process: *std.ChildProcess,
    stdout_buf: std.ArrayList(u8),
    stderr_buf: std.ArrayList(u8),

    dead_fds: usize = 0,
    max_output_bytes: usize = 50 * 1024,

    cleanup_done: bool = false,

    fn spawn(allocator: *std.mem.Allocator, argv: []const []const u8, max_output_bytes: usize) !ProcessWithOutput {
        const child = try std.ChildProcess.init(argv, allocator);
        child.stdin_behavior = std.ChildProcess.StdIo.Ignore;
        child.stdout_behavior = std.ChildProcess.StdIo.Pipe;
        child.stderr_behavior = std.ChildProcess.StdIo.Pipe;
        try child.spawn();

        return ProcessWithOutput{ .process = child, .stdout_buf = std.ArrayList(u8).init(allocator), .stderr_buf = std.ArrayList(u8).init(allocator), .dead_fds = 0, .max_output_bytes = max_output_bytes };
    }

    fn is_running(self: *ProcessWithOutput) bool {
        if (self.process.term) |_| {
            return false;
        } else {
            return true;
        }
    }

    fn stdout(self: ProcessWithOutput) []u8 {
        return self.stdout_buf.items;
    }

    fn stderr(self: ProcessWithOutput) []u8 {
        return self.stderr_buf.items;
    }

    // poll: https://github.com/ziglang/zig/blob/master/lib/std/child_process.zig#L206
    //   basically do one iteration with no blocking each time it runs and thus get the output incrementally?
    fn poll(self: *ProcessWithOutput) !void {
        if (!self.is_running()) {
            return;
        }

        var poll_fds = [_]std.os.pollfd{
            .{ .fd = self.process.stdout.?.handle, .events = std.os.POLL.IN, .revents = undefined },
            .{ .fd = self.process.stderr.?.handle, .events = std.os.POLL.IN, .revents = undefined },
        };

        // We ask for ensureTotalCapacity with this much extra space. This has more of an
        // effect on small reads because once the reads start to get larger the amount
        // of space an ArrayList will allocate grows exponentially.
        const bump_amt = 512;

        const err_mask = std.os.POLL.ERR | std.os.POLL.NVAL | std.os.POLL.HUP;

        if (self.dead_fds >= poll_fds.len) {
            return;
        }

        const events = try std.os.poll(&poll_fds, 0);
        if (events == 0) {
            return;
        }

        var remove_stdout = false;
        var remove_stderr = false;
        // Try reading whatever is available before checking the error
        // conditions.
        // It's still pstd.ossible to read after a POLL.HUP is received, always
        // check if there's some data waiting to be read first.
        if (poll_fds[0].revents & std.os.POLL.IN != 0) {
            // stdout is ready.
            const new_capacity = std.math.min(self.stdout_buf.items.len + bump_amt, self.max_output_bytes);
            try self.stdout_buf.ensureTotalCapacity(new_capacity);
            const buf = self.stdout_buf.unusedCapacitySlice();
            if (buf.len == 0) return error.StdoutStreamTooLong;
            const nread = try std.os.read(poll_fds[0].fd, buf);
            self.stdout_buf.items.len += nread;

            // Remove the fd when the EOF condition is met.
            remove_stdout = nread == 0;
        } else {
            remove_stdout = poll_fds[0].revents & err_mask != 0;
        }

        if (poll_fds[1].revents & std.os.POLL.IN != 0) {
            // stderr is ready.
            const new_capacity = std.math.min(self.stderr_buf.items.len + bump_amt, self.max_output_bytes);
            try self.stderr_buf.ensureTotalCapacity(new_capacity);
            const buf = self.stderr_buf.unusedCapacitySlice();
            if (buf.len == 0) return error.StderrStreamTooLong;
            const nread = try std.os.read(poll_fds[1].fd, buf);
            self.stderr_buf.items.len += nread;

            // Remove the fd when the EOF condition is met.
            remove_stderr = nread == 0;
        } else {
            remove_stderr = poll_fds[1].revents & err_mask != 0;
        }

        // Exclude the fds that signaled an error.
        if (remove_stdout) {
            poll_fds[0].fd = -1;
            self.dead_fds += 1;
        }
        if (remove_stderr) {
            poll_fds[1].fd = -1;
            self.dead_fds += 1;
        }
    }

    fn deinit(self: *ProcessWithOutput) void {
        self.stdout_buf.deinit();
        self.stderr_buf.deinit();
        self.process.deinit();
    }
};

const RegexRunner = struct {
    run_always: bool,
    process: ProcessWithOutput,

    command_to_argv: fn (cmd: []u8) [][]u8,

    fn run(self: RegexRunner, cmd: []u8, is_confirmed: true) !void {
        if (!self.run_always and !is_confirmed) {
            return;
        }

        if (self.is_running()) {
            self.process.kill();
        }
        const argv = self.command_to_argv(cmd);
        self.process = ProcessWithOutput(argv);
    }

    fn output(self: RegexRunner) []u8 {
        if (!self.is_running()) {
            // TODO: return buffer
            return "<done>";
        }

        // get more input?

        // TODO: return buffer
    }
};

pub fn main() !void {
    var general_purpose_allocator = std.heap.GeneralPurposeAllocator(.{}){};
    defer {
        _ = general_purpose_allocator.detectLeaks();
    }
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
    const screen = c.SDL_CreateWindow("hello fonts", c.SDL_WINDOWPOS_CENTERED, c.SDL_WINDOWPOS_CENTERED, window_width, window_height, c.SDL_WINDOW_BORDERLESS | c.SDL_WINDOW_OPENGL) orelse {
        c.SDL_Log("Unable to create window: %s", c.SDL_GetError());
        return error.SDLInitializationFailed;
    };
    defer c.SDL_DestroyWindow(screen);

    const op: f32 = 0.5;
    if (c.SDL_SetWindowOpacity(screen, op) != 0) {
        c.SDL_Log("Unable to make window transparent: %s", c.SDL_GetError());
    }
    var opacity: f32 = 10.0;
    _ = c.SDL_GetWindowOpacity(screen, &opacity);
    c.SDL_Log("opacity: %f", opacity);

    const renderer = c.SDL_CreateRenderer(screen, -1, c.SDL_RENDERER_ACCELERATED);

    var msg = "                                                                                                    ".*;
    var pos: usize = 0;
    var max_chars = std.math.min(@divTrunc(@intCast(usize, window_width), @intCast(usize, glyph_width)), msg.len);

    var result: []const u8 = try gpa.alloc(u8, 0);
    defer gpa.free(result);

    const keyboardState = c.SDL_GetKeyboardState(null);

    c.SDL_StartTextInput();

    var process = &try ProcessWithOutput.spawn(gpa, &[_][]const u8{ "ag", "\\bshit\\b", "/usr/include" }, 1024 * 1024);

    var quit = false;
    var skip: i32 = 0;
    var num_lines: i32 = 0;
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
                            c.SDLK_c => {
                                const clipboard_text = try gpa.dupeZ(u8, result);
                                if (c.SDL_SetClipboardText(clipboard_text) != 0) {
                                    c.SDL_Log("Could not set clipboard text: %s", c.SDL_GetError());
                                }
                                gpa.free(clipboard_text);
                            },
                            c.SDLK_v => {
                                const clipboard_text = c.SDL_GetClipboardText();
                                if (std.mem.len(clipboard_text) == 0) {
                                    c.SDL_Log("Could not get clipboard: %s", c.SDL_GetError());
                                } else {
                                    const initial_pos = pos;
                                    while (pos < max_chars and pos - initial_pos < std.mem.len(clipboard_text)) : (pos += 1) {
                                        msg[pos] = clipboard_text[pos - initial_pos];
                                    }
                                    msg[pos] = ' ';
                                    msg[max_chars] = 0;
                                }
                                c.SDL_free(clipboard_text);
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
                                skip = 0;
                                gpa.free(result);
                                result = try runCommand(&msg, gpa);
                                var i: usize = 0;
                                while (i < max_chars) : (i += 1) {
                                    msg[i] = ' ';
                                }
                                msg[max_chars] = 0;
                                pos = 0;

                                num_lines = 0;
                                var lines = std.mem.split(u8, result, "\n");
                                var line = lines.next();
                                while (line != null) : (line = lines.next()) {
                                    num_lines += 1;
                                }
                            },
                            c.SDLK_UP => {
                                if (skip > 0) {
                                    skip -= 1;
                                }
                            },
                            c.SDLK_PAGEUP => {
                                if (skip < 10) {
                                    skip = 0;
                                } else {
                                    skip -= 10;
                                }
                            },
                            c.SDLK_DOWN => {
                                skip += 1;
                            },
                            c.SDLK_PAGEDOWN => {
                                skip += 10;
                            },
                            c.SDLK_HOME => {
                                skip = 0;
                            },
                            c.SDLK_END => {
                                if (num_lines > 10) {
                                    skip = num_lines - 10;
                                }
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

        _ = c.SDL_SetRenderDrawColor(renderer, 0, 0, 0, 100);
        //_ = c.SDL_SetRenderDrawBlendMode(renderer, c.SDL_BlendMode.SDL_BLENDMODE_BLEND);
        _ = c.SDL_RenderClear(renderer);

        // thanks to https://stackoverflow.com/questions/22886500/how-to-render-text-in-sdl2 for some actually useful code here
        const white: c.SDL_Color = c.SDL_Color{ .r = 255, .g = 255, .b = 255, .a = 255 };
        const black: c.SDL_Color = c.SDL_Color{ .r = 0, .g = 0, .b = 0, .a = 100 };
        // Shaded vs Solid gives a nicer output, with solid the output
        // was squiggly and not aligned with a baseline.
        const text = c.TTF_RenderUTF8_Shaded(font, &msg, white, black);
        const texture = c.SDL_CreateTextureFromSurface(renderer, text);
        c.SDL_FreeSurface(text);
        _ = c.SDL_RenderCopy(renderer, texture, null, &c.SDL_Rect{ .x = 0, .y = 0, .w = @intCast(c_int, msg.len) * glyph_width, .h = glyph_height });

        var i: c_int = 1;
        var lines = std.mem.split(u8, result, "\n");
        var line = lines.next();
        {
            var skipped: i32 = 0;
            while (skipped < skip and line != null) : (skipped += 1) {
                line = lines.next();
            }
        }

        try process.poll();
        std.debug.print("{s} {d} {d}\n", .{ process.is_running(), process.stdout_buf.items.len, process.stdout_buf.capacity });
        lines = std.mem.split(u8, process.stdout(), "\n");
        line = lines.next();
        while (line != null and i * glyph_height < window_height) {
            const line_c = try gpa.dupeZ(u8, line.?);
            // TODO: render tabs at correct width (or some width at least)
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
    else if (std.mem.startsWith(u8, cmd, "s "))
        &[_][]const u8{ "ag", cmd["s ".len..], "/home/luna/t/raylib/src", "/home/luna/k/the-thing/resources/ode-build/include" }
    else if (cmd.len > 0 and std.ascii.isDigit(cmd[0]))
        &[_][]const u8{ "qalc", "-terse", cmd }
    else
        &[_][]const u8{ "bash", "-c", cmd };
    for (argv) |arg| {
        std.debug.print("'{s}' ", .{arg});
    }
    const result = try std.ChildProcess.exec(.{ .allocator = allocator, .argv = argv, .max_output_bytes = 1024 * 1024 });
    std.debug.print("stdout: '{s}'\n", .{result.stdout[0..std.math.min(100, result.stdout.len)]});
    std.debug.print("stderr: '{s}'\n", .{result.stderr});

    if (result.stdout.len > 0) {
        allocator.free(result.stderr);
        return result.stdout;
    } else if (result.stderr.len > 0) {
        allocator.free(result.stdout);
        return result.stderr;
    } else {
        allocator.free(result.stdout);
        allocator.free(result.stderr);
        return std.fmt.allocPrint(allocator, "<no output>", .{});
    }
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
