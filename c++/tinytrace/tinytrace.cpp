#include <algorithm>
#include <chrono>
#include <float.h>
#include <fstream>
#include <getopt.h>
#include <iostream>
#include <mutex>
#include <sstream>
#include <thread>

#include <SDL2/SDL.h>

#include "tracy/public/tracy/Tracy.hpp"
// #define STB_IMAGE_IMPLEMENTATION
#include "stb_image.h"

#include "bvh_node.h"
#include "camera.h"
#include "dielectric.h"
#include "distributor.h"
#include "hitable_list.h"
#include "image_texture.h"
#include "lambertian.h"
#include "metal.h"
#include "rectangle.h"
#include "sphere.h"
#include "texture.h"

vec3 color(const ray &r, hitable *world, int depth) {
  ZoneScoped;
  TracyPlot("depth", int64_t(depth));

  hit_record rec;
  if (world->hit(r, 0.001, MAXFLOAT, rec)) {
    ZoneScopedN("hit");

    ray scattered;
    vec3 attenuation;
    vec3 emitted = rec.mat_ptr->emitted(rec.u, rec.v, rec.p);
    if (depth < 50 && rec.mat_ptr->scatter(r, rec, attenuation, scattered)) {
      return emitted + attenuation * color(scattered, world, depth + 1);
    } else {
      return emitted;
    }
  } else {
    vec3 unit_direction = unit_vector(r.direction());
    float t = 0.5 * (unit_direction.y() + 1.0);
    // lerp() for sky color
    return (1.0 - t) * vec3(1.0, 1.0, 1.0) + t * vec3(0.5, 0.7, 1.0);
    // return vec3(0, 0, 0);
  }
}

void draw_image(std::string path, int nx, int ny, vec3 *image);
vec3 *read_image(std::string name, int &nx, int &ny);

hitable **random_scene(int &n);

int main(int argc, char **argv) {
  // TracyAppInfo("tinytrace", 16);

  std::string out_name = "/dev/stdout";
  bool write_partial = false;
  int nx = 800;
  int ny = 400;
  int ns = 100;
  int concurrency = std::thread::hardware_concurrency();
  bool continue_render = false;
  long *seed = NULL;
  bool verbose = false;

  while (true) {
    char opt_char = 0;
    int option_index = 0;
    long opt_seed = 0;
    static struct option long_options[] = {
        {"width", required_argument, 0, 'w'},
        {"height", required_argument, 0, 'h'},
        {"samples", required_argument, 0, 's'},
        {"seed", required_argument, 0, 'r'},
        {"jobs", required_argument, 0, 'j'},
        {"continue", no_argument, 0, 'c'},
        {"verbose", no_argument, 0, 'v'},
        {"help", no_argument, 0, 0},
        {0, 0, 0, 0},
    };

    opt_char =
        getopt_long(argc, argv, "w:h:s:j:cvh", long_options, &option_index);
    if (opt_char == -1) {
      break;
    }
    switch (opt_char) {
    case 'w':
      nx = std::stoi(optarg);
      break;
    case 'h':
      ny = std::stoi(optarg);
      break;
    case 's':
      ns = std::stoi(optarg);
    case 'r':
      opt_seed = std::stol(optarg);
      seed = &opt_seed;
      break;
    case 'j':
      concurrency = std::stoi(optarg);
      break;
    case 'c':
      continue_render = true;
      break;
    case 'v':
      verbose = true;
      break;
    default:
      std::cerr << "Usage: " << argv[0] << "[flags] [<filename>]"
                << "\n";
      std::cerr << R"(
  -w N, --width=N   Set width of output image to N
  -h N, --height=N  Set height of output image to N
  -s N, --samples=N Sample each ray per pixel N times (anti-aliasing)
  -r N, --seed=N    Use seed N
  -j N, --jobs=N    Render with N threads concurrently
  -c  , --continue  Continue rendering image contained in 'filename' (if it exists and is specified)
  -v  , --verbose   Enable verbose mode
        --help      This very message!
)";
      exit(1);
    }
  }

  // nx = 200;
  // ny = 100;
  // ns = 100;
  // continue_render = false;

  if (optind < argc) {
    write_partial = true;
    out_name = argv[optind];
  }

  if (seed != NULL) {
    srand48(*seed);
  }

  // hitable *list[5];
  // list[0] =
  //     new sphere(vec3(0, 0, -1), 0.5, new lambertian(vec3(0.1, 0.2, 0.5)));
  // list[1] =
  //     new sphere(vec3(0, -100.5, -1), 100, new lambertian(vec3(0.8, 0.8,
  //     0.0)));
  // list[2] =
  //     new sphere(vec3(1, 0, -1), 0.5, new metal(vec3(0.8, 0.6, 0.2), 0.0));
  // list[3] = new sphere(vec3(-1, 0, -1), 0.5, new dielectric(1.5));
  // list[4] = new sphere(vec3(-1, 0, -1), -0.45, new dielectric(1.5));
  // hitable *world = new hitable_list(list, 5);
  int size = 500;
  hitable **objects = random_scene(size);
  hitable *world = new bvh_node(objects, size, 0, 0);
  world = new hitable_list(objects, size);

  vec3 *image = new vec3[nx * ny + 1];
  for (int j = ny - 1, c = 0; j >= 0; j--) {
    for (int i = 0; i < nx; i++) {
      image[c] = vec3(255, 255, 255);
      c++;
    }
  }

  if (continue_render) {
    vec3 *old_image = read_image(out_name, nx, ny);
    if (old_image != NULL) {
      image = old_image;
    }
  }

  SDL_Init(SDL_INIT_VIDEO);
  atexit(SDL_Quit);

  int bit_depth = 16;
  auto window = SDL_CreateWindow("tinytrace!", 0, 0, nx, ny, SDL_WINDOW_VULKAN);

  auto screen = SDL_CreateRGBSurface(0, nx, ny, 32, 0, 0, 0, 0);
  auto renderer = SDL_CreateRenderer(window, -1, SDL_RENDERER_ACCELERATED);

  auto texture = SDL_CreateTexture(renderer, SDL_PIXELFORMAT_RGBA8888,
                                   SDL_TEXTUREACCESS_STATIC, nx, ny);
  Uint32 *pixels = new Uint32[nx * ny];
  memset(pixels, 255, nx * ny * sizeof(Uint32)); // init to white

  SDL_UpdateTexture(texture, NULL, pixels, nx * sizeof(Uint32));

  SDL_RenderClear(renderer);
  SDL_RenderCopy(renderer, texture, NULL, NULL);
  SDL_RenderPresent(renderer);

  draw_image(out_name, nx, ny, image);
  // exit(0);

  auto d = new distributor(nx, ny);
  d->set_randomize(true);
  d->continue_from(image);
  std::mutex image_lock;

  // TODO: render small image first (10% per side if > 1000) -> blow up  to size
  // -> render full size

  bool quit = false;

  std::cerr << concurrency;
  std::thread *threads = new std::thread[concurrency + 1];
  bool *done = new bool[concurrency + 1];
  int *counts = new int[concurrency + 1];

  vec3 look_from = vec3(11, 2, 3.5);
  vec3 look_at = vec3(5, 1, 1.5);
  float dist_to_focus = (look_from - look_at).length();
  float aperture = 0.05;
  camera cam(look_from, look_at, vec3(0, 1, 0), 20, float(nx) / float(ny),
             aperture, dist_to_focus);

  int samples = 1;

  auto start = std::chrono::high_resolution_clock::now();
  for (int t = 0; t < concurrency; t++) {
    done[t] = false;
    counts[t] = 0;

    auto render = [t, &samples, done, counts, &d, &cam, world, ns, image,
                   &image_lock, pixels, &quit] {
      for (int local_samples = samples; local_samples < ns;
           local_samples += 2) {
        if (quit) {
          break;
        }

        int c, i, j;
        while (d->next_pixel(c, i, j)) {
          ZoneScopedN("render");
          if (quit) {
            break;
          }

          counts[t] += 1;

          vec3 col(0, 0, 0);

          // anti-aliasing (sample #ns rays)
          for (int s = 0; s < local_samples; s++) {
            float u = float(i + drand48()) / float(d->width());
            float v = float(j + drand48()) / float(d->height());
            ray r = cam.get_ray(u, v);
            vec3 p = r.point_at_parameter(2.0);
            col += color(r, world, 0);
          }
          col /= float(local_samples);

          // gamme correct?
          col = vec3(sqrt(col.r()), sqrt(col.g()), sqrt(col.b()));

          // image_lock.lock();
          pixels[c] =
              (std::max(0, std::min(int(254.99 * col.r()), 255)) << 24) +
              (std::max(0, std::min(int(254.99 * col.g()), 255)) << 16) +
              (std::max(0, std::min(int(254.99 * col.b()), 255)) << 8);
          image[c] = vec3(col.r(), col.g(), col.b());
          // image_lock.unlock();
        }

        if (t == 0) {
          d->reset();
        }
      }

      std::cerr << "_";
      done[t] = true;
    };

    threads[t] = std::thread(render);
  }

  std::thread check([d, concurrency, done, write_partial, threads, out_name, nx,
                     ny, image, start, &quit] {
    int i, j;
    while (!quit && !d->is_done()) {
      bool all_done = true;
      for (int k = 0; k < concurrency; k++) {
        all_done = all_done && done[k];
      }
      if (all_done) {
        break;
      }

      if ((d->count()) % int(nx * ny / 100.0) == 0) {
        std::cerr << "."; // progress dots ✨

        if (write_partial) {
          // write updated image
          draw_image(out_name, nx, ny, image);
        }
      }
    }

    std::cerr << "!";
    for (int t = 0; t < concurrency; t++) {
      // std::cerr << " " << counts[t];
      threads[t].join();
    }

    std::cerr << "\n";
    auto finish = std::chrono::high_resolution_clock::now();
    std::cerr << "took "
              << std::chrono::duration_cast<std::chrono::milliseconds>(finish -
                                                                       start)
                     .count()
              << "ms"
              << "\n";

    draw_image(out_name, nx, ny, image);
  });

  bool fullscreen = false;
  while (!quit) {
    SDL_UpdateTexture(texture, NULL, pixels, nx * sizeof(Uint32));

    SDL_Event event;
    SDL_WaitEventTimeout(&event, 16);

    switch (event.type) {
    case SDL_KEYDOWN:
      switch (event.key.keysym.sym) {
      case SDLK_s:
        // memset(pixels, 255, nx * ny * sizeof(Uint32)); // init to white

        look_from = vec3(look_from.x() + 1, look_from.y(), look_from.z());
        cam = camera(look_from, look_at, vec3(0, 1, 0), 20,
                     float(nx) / float(ny), aperture, dist_to_focus);

        d->reset();
        samples = 1;

        break;

      case SDLK_F11:
        fullscreen = !fullscreen;
        SDL_SetWindowFullscreen(window, fullscreen ? SDL_WINDOW_FULLSCREEN : 0);
        break;
      case SDLK_ESCAPE:
        quit = true;
        break;
      }
      break;
    case SDL_QUIT:
      quit = true;
      break;
    }

    SDL_RenderClear(renderer);
    SDL_RenderCopy(renderer, texture, NULL, NULL);
    SDL_RenderPresent(renderer);
  }

  check.join();

  SDL_DestroyTexture(texture);
  SDL_DestroyRenderer(renderer);
  SDL_DestroyWindow(window);
  SDL_Quit();
}

hitable **random_scene(int &n) {
  hitable **list = new hitable *[n + 1];
  texture *base_texture =
      new checker_texture(new constant_texture(vec3(0.2, 0.3, 0.1)),
                          new constant_texture(vec3(0.9, 0.9, 0.9)));
  // texture *base_texture = new sine_texture();
  list[0] = new sphere(vec3(0, -1000, 0), 1000, new lambertian(base_texture));
  int i = 1;
  for (int a = -11; a < 11; a++) {
    for (int b = -11; b < 11; b++) {
      float choose_mat = drand48();
      vec3 center(a + 0.9 * drand48(), 0.2, b + 0.9 * drand48());
      if ((center - vec3(4, 0.2, 0)).length() > 0.9) {
        if (choose_mat < 0.8) { // diffuse
          list[i++] =
              new sphere(center, 0.2,
                         new lambertian(new constant_texture(
                             vec3(drand48() * drand48(), drand48() * drand48(),
                                  drand48() * drand48()))));

        } else if (choose_mat < 0.95) { // metal
          list[i++] = new sphere(
              center, 0.2,
              new metal(vec3(0.5 * (1 + drand48()), 0.5 * (1 + drand48()),
                             0.5 * (1 + drand48())),
                        0.5 * drand48()));
        } else { // glass
          list[i++] = new sphere(center, 0.2, new dielectric(1.5));
        }
      }
    }
  }

  // lights
  list[i++] = new xy_rect(
      3, 5, 1, 3, -2, new diffuse_light(new constant_texture(vec3(4, 4, 4))));
  list[i++] =
      new sphere(vec3(0, 0, 3), 1.0,
                 new diffuse_light(new constant_texture(vec3(1, 1, 1))));

  // big spheres
  list[i++] = new sphere(vec3(0, 1, 0), 1.0, new dielectric(1.5));
  int nx, ny, nn;
  texture *tex;
#ifdef STB_IMAGE_IMPLEMENTATION
  unsigned char *tex_data = stbi_load("earthmap.jpg", &nx, &ny, &nn, 0);
  tex = new image_texture(tex_data, nx, ny)
#else
  tex = new constant_texture(vec3(0.4, 0.2, 0.1));
#endif
      list[i++] = new sphere(vec3(-3, 1, 0), 1.0, new lambertian(tex));
  list[i++] =
      new sphere(vec3(4, 1, 0), 1.0, new metal(vec3(0.7, 0.6, 0.5), 0.0));

  n = i;
  return list;
}

void draw_image(std::string path, int nx, int ny, vec3 *image) {
  ZoneScoped;

  std::ofstream out;
  out.open(path + ".tmp", std::ios_base::trunc);
  out << "P3\n" << nx << " " << ny << "\n255\n";
  int c = 0;
  for (int y = ny - 1; y >= 0; y--) {
    for (int x = 0; x < nx; x++) {
      vec3 col = image[c];
      int ir = std::max(0, std::min(int(254.99 * col.r()), 255));
      int ig = std::max(0, std::min(int(254.99 * col.g()), 255));
      int ib = std::max(0, std::min(int(254.99 * col.b()), 255));
      out << ir << " " << ig << " " << ib << "\n";
      c++;
    }
  }
  out.flush();
  out.close();

  // use `rename` to get atomic appearance of the file to not confuse image
  // viewers
  rename((path + ".tmp").c_str(), path.c_str());

  // write a single space at the end so the image viewer notices
  out.open(path, std::fstream::out | std::fstream::app);
  out.write(" ", 1);
  out.flush();
  out.close();
}

vec3 *read_image(std::string name, int &nx, int &ny) {
  std::ifstream in;
  in.open(name, std::ios::in);

  if (!in.is_open()) {
    std::cerr << "could not read old file";
    return NULL;
  }

  // header
  std::string line;
  getline(in, line);
  if (line.compare("P3") != 0) {
    throw "'P3' not found, not a PPM?";
  }

  // resolution
  getline(in, line);
  std::istringstream line_stream(line);
  std::string part;
  getline(line_stream, part, ' ');
  nx = std::stoi(part);
  getline(line_stream, part, ' ');
  ny = std::stoi(part);

  std::cerr << "r" << nx << "x" << ny << ":";

  // max colors
  getline(in, line);
  int max_color = std::stoi(line);

  // pixels!
  vec3 *image = new vec3[nx * ny + 1];

  int c = 0;
  while (getline(in, line)) {
    if (line.compare(" ") == 0) {
      continue;
    }

    int r, g, b;

    std::istringstream line_stream(line);
    std::string part;
    getline(line_stream, part, ' ');
    r = std::max(0, std::min(std::stoi(part), max_color));
    getline(line_stream, part, ' ');
    g = std::max(0, std::min(std::stoi(part), max_color));
    getline(line_stream, part, ' ');
    b = std::max(0, std::min(std::stoi(part), max_color));

    image[c] =
        vec3(double(r) / double(max_color), double(g) / double(max_color),
             double(b) / double(max_color));
    // std::cerr << image[c].r() << " " << image[c].g() << " " << image[c].b()
    //           << "\t" << r << " " << g << " " << b << "\n";

    c++;
  }

  if (c != nx * ny) {
    std::cerr << "not enough pixels, expected " << nx * ny << " but only got "
              << c;
    exit(1);
  }

  in.close();

  return image;
}
