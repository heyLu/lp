// Playing around with SDL + TTF.
//
// Resources:
// - http://wiki.libsdl.org/CategoryAPI
// - https://www.libsdl.org/projects/SDL_ttf/docs/SDL_ttf.html

#include <sys/param.h>

#include <SDL.h>
#include <SDL_ttf.h>

char *font_file = "./FantasqueSansMono-Regular.ttf";

int main(int argc, char *argv[]) {
	printf("hello, world!\n");

	if (SDL_Init(SDL_INIT_EVERYTHING) != 0) {
		printf("shit: %s\n", SDL_GetError());
		return 1;
	}

	if (TTF_Init() != 0) {
		printf("ttf init: %s\n", TTF_GetError());
		return 1;
	}

	if (argc > 1) {
		font_file = argv[1];
	}
	printf("using font '%s'\n", font_file);

	TTF_Font *font = TTF_OpenFont(font_file, 20);
	if (font == NULL) {
		printf("open font: %s\n", TTF_GetError());
		return 1;
	}

	SDL_Window *window = SDL_CreateWindow("hello fonts", SDL_WINDOWPOS_CENTERED, SDL_WINDOWPOS_CENTERED, 600, 100, SDL_WINDOW_BORDERLESS);
	if (window == NULL) {
		printf("create window: %s\n", SDL_GetError());
		return 1;
	}

	SDL_Renderer *renderer = SDL_CreateRenderer(window, -1, SDL_RENDERER_ACCELERATED);
	if (renderer == NULL) {
		printf("create renderer: %s\n", SDL_GetError());
		return 1;
	}

	SDL_Surface *surface = SDL_GetWindowSurface(window);

	char msg[100] = "howdy there, enby! 🐘                                          ";

	// monospace -> fixed width (duh)
	int advance = 0;
	for (int i = 0; i < strlen(msg); i++) {
		TTF_GlyphMetrics(font, msg[i], NULL, NULL, NULL, NULL, &advance);
		printf("advance '%c': %d\n", msg[i], advance);
	}

	SDL_StartTextInput();

	int pos = 0;
	int max_chars = MIN(surface->w / advance, sizeof(msg));

	SDL_bool quit = SDL_FALSE;

	SDL_Event event;
	while (!quit) {
		while(SDL_PollEvent(&event)) {
			if (event.type == SDL_QUIT || (event.type == SDL_KEYDOWN && event.key.keysym.sym == SDLK_ESCAPE)) {
				printf("quit received\n");
				quit = SDL_TRUE;
				break;
			}

			if (event.type == SDL_KEYDOWN && event.key.keysym.sym == SDLK_BACKSPACE) {
				if (pos == 0) {
					pos = max_chars-1;
				} else {
					pos = (pos - 1) % (max_chars - 1);
				}
				msg[pos] = '_';
				printf("back to %d\n", pos);
				continue;
			}

			if (event.type == SDL_TEXTINPUT) {
				if (strlen(event.text.text) > 0) {
					printf("key: %s\n", event.text.text);

					msg[pos] = event.text.text[0];
					pos = (pos + 1) % (max_chars - 1);
				}
			}
		}

		// thanks to https://stackoverflow.com/questions/22886500/how-to-render-text-in-sdl2 for some actually useful code here
		SDL_Color white = {255, 255, 255, 255};
		SDL_Color black = {0, 0, 0};
		SDL_Surface* text = TTF_RenderUTF8_Shaded(font, msg, white, black);
		SDL_BlitSurface(text, NULL, surface, NULL);

		SDL_UpdateWindowSurface(window);
		SDL_Delay(16);
	}

	SDL_DestroyWindow(window);

	SDL_Quit();

	printf("done here.\n");
}
