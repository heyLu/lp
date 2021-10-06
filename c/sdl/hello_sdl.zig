// Playing around with SDL + TTF.
//
// Based on `hello_sdl.c`, zig + SDL code from https://github.com/andrewrk/sdl-zig-demo.
//
// Resources:
// - http://wiki.libsdl.org/CategoryAPI
// - https://www.libsdl.org/projects/SDL_ttf/docs/SDL_ttf.html
const c = @cImport({
    @cInclude("SDL2/SDL.h");
    @cInclude("SDL2/SDL_ttf.h");
});
const assert = @import("std").debug.assert;
const process = @import("std").process;

pub fn main() !void {
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

    var font_file = "./FantasqueSansMono-Regular.ttf";
    const font = c.TTF_OpenFont(font_file, 20) orelse {
        c.SDL_Log("Unable to load font: %s", c.TTF_GetError());
        return error.TTFInitializationFailed;
    };
    defer c.TTF_CloseFont(font);

    const msg = "howdy there, enby! 🐘                                          ";
    c.SDL_Log(msg);

    const screen = c.SDL_CreateWindow("hello fonts", c.SDL_WINDOWPOS_CENTERED, c.SDL_WINDOWPOS_CENTERED, 600, 100, c.SDL_WINDOW_BORDERLESS) orelse {
        c.SDL_Log("Unable to create window: %s", c.SDL_GetError());
        return error.SDLInitializationFailed;
    };
    defer c.SDL_DestroyWindow(screen);

    const surface = c.SDL_GetWindowSurface(screen);

    var quit = false;
    while (!quit) {
        var event: c.SDL_Event = undefined;
        while (c.SDL_PollEvent(&event) != 0) {
            switch (event.@"type") {
                c.SDL_QUIT => {
                    quit = true;
                },
                c.SDL_KEYDOWN => {
                    switch (event.key.keysym.sym) {
                        c.SDLK_ESCAPE => {
                            quit = true;
                        },
                        else => {},
                    }
                },
                else => {},
            }
        }

        // thanks to https://stackoverflow.com/questions/22886500/how-to-render-text-in-sdl2 for some actually useful code here
        const white: c.SDL_Color = c.SDL_Color{ .r = 255, .g = 255, .b = 255, .a = 255 };
        const black: c.SDL_Color = c.SDL_Color{ .r = 0, .g = 0, .b = 0, .a = 255 };
        const text = c.TTF_RenderUTF8_Shaded(font, msg, white, black);
        _ = c.SDL_BlitSurface(text, null, surface, null);

        _ = c.SDL_UpdateWindowSurface(screen);

        c.SDL_Delay(16);
    }
}
