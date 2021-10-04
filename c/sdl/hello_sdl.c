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

	char *msg = "howdy there, enby!🐘";

	// thanks to https://stackoverflow.com/questions/22886500/how-to-render-text-in-sdl2 for some actually useful code here
	SDL_Color white = {255, 255, 255};
	SDL_Color black = {0, 0, 0};
	SDL_Surface* surfaceMessage = TTF_RenderUTF8_Blended(font, msg, white); 
	SDL_Texture* message = SDL_CreateTextureFromSurface(renderer, surfaceMessage);

	SDL_Rect text_rect;
	text_rect.x = 0;
	text_rect.y = 0;
	text_rect.w = surfaceMessage->w;
	text_rect.h = surfaceMessage->h;

	int advance = 0;
	for (int i = 0; i < strlen(msg); i++) {
		TTF_GlyphMetrics(font, msg[i], NULL, NULL, NULL, NULL, &advance);
		printf("advance '%c': %d\n", msg[i], advance);
	}

	SDL_Event event;
	while (1) {
		SDL_PollEvent(&event);
		if (event.type == SDL_QUIT || (event.type == SDL_KEYDOWN && event.key.keysym.sym == SDLK_ESCAPE)) {
			printf("quit received\n");
			break;
		}

		SDL_RenderClear(renderer);
		SDL_RenderCopy(renderer, message, NULL, &text_rect);
		SDL_RenderPresent(renderer);

		SDL_Delay(16);
	}

	SDL_FreeSurface(surfaceMessage);
	SDL_DestroyTexture(message);
	SDL_DestroyWindow(window);

	SDL_Quit();

	printf("done here.\n");
}
