FROM docker.io/appimagecrafters/appimage-builder:latest

RUN apt-get update && apt-get install -y libsdl2-2.0-0
RUN apt-get install -y curl
RUN apt-get install -y libsdl2-dev
RUN apt-get install -y clang

COPY tracy/public /tmp/tracy/public
COPY Makefile *.cpp *.h /tmp
RUN cd /tmp \
 && touch empty.cpp && clang++ -c empty.cpp -o tracy.o \
 && make TRACY_FLAGS= tinytrace

RUN apt-get install -y adwaita-icon-theme

COPY AppImageBuilder.yml /tmp
RUN cd /tmp \
 && mkdir -p AppDir/usr/bin \
 && cp /tmp/tinytrace AppDir/usr/bin \
 && appimage-builder --skip-test
