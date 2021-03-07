CFLAGS = -Wall -g -lGL $(shell sdl2-config --cflags --libs) -ldl -lstdc++ -lm
TRACY_FLAGS = -DTRACY_ENABLE=1

main: main.c tracy.o load_shaders.o
	gcc $(CFLAGS) $(TRACY_FLAGS) load_shaders.o tracy.o main.c -o main

load_shaders.o: load_shaders.cpp
	g++ -c $(TRACY_FLAGS) -o load_shaders.o load_shaders.cpp

tracy.o: tracy/TracyClient.cpp
	g++ -c $(TRACY_FLAGS) -o tracy.o tracy/TracyClient.cpp
