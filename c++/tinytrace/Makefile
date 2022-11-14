all: tinytrace example.png

CC = clang++
CFLAGS = -Wall -O3 -march=native -g
LDFLAGS = -lpthread $(shell pkg-config --libs luajit sdl2)
TRACY_FLAGS = -DTRACY_ENABLE -DTRACY_ON_DEMAND 

tinytrace: *.cpp *.h stb_image.h earthmap.jpg tracy.o
	$(CC) $(CFLAGS) $(LDFLAGS) $(TRACY_FLAGS) tinytrace.cpp ray.cpp lua.cpp tracy.o -o tinytrace

stb_image.h:
	curl -Lsfo stb_image.h https://raw.githubusercontent.com/nothings/stb/master/stb_image.h

earthmap.jpg:
	curl -Lsfo earthmap.jpg https://raw.githubusercontent.com/RayTracing/raytracing.github.io/master/images/earthmap.jpg

tracy.o: tracy/public/TracyClient.cpp
	$(CC) $(CFLAGS) $(TRACY_FLAGS) -c tracy/public/TracyClient.cpp -o tracy.o

tinytrace-0.1.0-x86_64.AppImage: tinytrace AppImageBuilder.yml Dockerfile
	podman build -t tinytrace-appimage:local -f Dockerfile .
	podman run -it --rm -v $(PWD):/app tinytrace-appimage:local cp /tmp/tinytrace-0.1.0-x86_64.AppImage /app

example.png: example.ppm
	convert example.ppm example.png
