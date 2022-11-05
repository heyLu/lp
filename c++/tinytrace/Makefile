all: tinytrace example.png

CFLAGS = -O3 -march=native -g
TRACY_FLAGS = -DTRACY_ENABLE -DTRACY_ON_DEMAND 

tinytrace: *.cpp *.h tracy.o
	clang++ $(CFLAGS) $(TRACY_FLAGS) tinytrace.cpp ray.cpp tracy.o -o tinytrace 

tracy.o: tracy/public/TracyClient.cpp
	clang++ $(CFLAGS) $(TRACY_FLAGS) -c tracy/public/TracyClient.cpp -o tracy.o

example.png: example.ppm
	convert example.ppm example.png
