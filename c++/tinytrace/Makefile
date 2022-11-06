all: tinytrace example.png

CFLAGS = -O3 -march=native -g $(shell pkg-config --libs sdl2)
TRACY_FLAGS = -DTRACY_ENABLE -DTRACY_ON_DEMAND 

tinytrace: *.cpp *.h stb_image.h earthmap.jpg tracy.o
	clang++ $(CFLAGS) $(TRACY_FLAGS) tinytrace.cpp ray.cpp tracy.o -o tinytrace 

stb_image.h:
	curl -Lsfo stb_image.h https://raw.githubusercontent.com/nothings/stb/master/stb_image.h

earthmap.jpg:
	curl -Lsfo earthmap.jpg https://raw.githubusercontent.com/RayTracing/raytracing.github.io/master/images/earthmap.jpg

tracy.o: tracy/public/TracyClient.cpp
	clang++ $(CFLAGS) $(TRACY_FLAGS) -c tracy/public/TracyClient.cpp -o tracy.o

example.png: example.ppm
	convert example.ppm example.png
