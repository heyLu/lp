#include <algorithm>
#include <chrono>
#include <float.h>
#include <fstream>
#include <iostream>
#include <mutex>
#include <sstream>
#include <thread>

#include "tracy/public/tracy/Tracy.hpp"

#include "camera.h"
#include "dielectric.h"
#include "distributor.h"
#include "hitable_list.h"
#include "lambertian.h"
#include "metal.h"
#include "sphere.h"

vec3 color(const ray &r, hitable *world, int depth) {
  ZoneScoped;
  TracyPlot("depth", int64_t(depth));

  hit_record rec;
  if (world->hit(r, 0.001, MAXFLOAT, rec)) {
    ZoneScopedN("hit");

    ray scattered;
    vec3 attenuation;
    if (depth < 50 && rec.mat_ptr->scatter(r, rec, attenuation, scattered)) {
      return attenuation * color(scattered, world, depth + 1);
    } else {
      return vec3(0, 0, 0);
    }
  } else {
    vec3 unit_direction = unit_vector(r.direction());
    float t = 0.5 * (unit_direction.y() + 1.0);
    // lerp() for sky color
    return (1.0 - t) * vec3(1.0, 1.0, 1.0) + t * vec3(0.5, 0.7, 1.0);
  }
}

void draw_image(std::string path, int nx, int ny, vec3 *image);
vec3 *read_image(std::string name, int &nx, int &ny);

hitable *random_scene(int n);

int main(int argc, char **argv) {
  // TracyAppInfo("tinytrace", 16);

  std::string out_name = "/dev/stdout";
  bool write_partial = false;
  if (argc > 1) {
    write_partial = true;
    out_name = argv[1];
  }

  bool continue_render = false;
  if (argc > 2 && std::string("continue").compare(argv[2]) == 0) {
    continue_render = true;
  }

  // int seed = time(0);
  // srand48(seed);
  // std::cerr << "s" << seed << " ";

  int nx = 2000;
  int ny = 1000;
  int ns = 100;

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
  hitable *world = random_scene(500);

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

  draw_image(out_name, nx, ny, image);
  // exit(0);

  auto d = new distributor(nx, ny);
  d->set_randomize(true);
  d->continue_from(image);
  std::mutex image_lock;

  // TODO: render small image first (10% per side if > 1000) -> blow up  to size
  // -> render full size

  int concurrency = std::thread::hardware_concurrency();
  std::cerr << concurrency;
  std::thread *threads = new std::thread[concurrency + 1];
  bool *done = new bool[concurrency + 1];
  int *counts = new int[concurrency + 1];

  auto start = std::chrono::high_resolution_clock::now();
  for (int t = 0; t < concurrency; t++) {
    done[t] = false;
    counts[t] = 0;

    vec3 look_from = vec3(11, 2, 3.5);
    vec3 look_at = vec3(5, 1, 1.5);
    float dist_to_focus = (look_from - look_at).length();
    float aperture = 0.05;
    camera cam(look_from, look_at, vec3(0, 1, 0), 20, float(nx) / float(ny),
               aperture, dist_to_focus);
    auto render = [t, done, counts, &d, cam, world, ns, image, &image_lock] {
      int c, i, j;
      while (d->next_pixel(c, i, j)) {
        ZoneScopedN("render");

        counts[t] += 1;

        vec3 col(0, 0, 0);

        // anti-aliasing (sample #ns rays)
        for (int s = 0; s < ns; s++) {
          float u = float(i + drand48()) / float(d->width());
          float v = float(j + drand48()) / float(d->height());
          ray r = cam.get_ray(u, v);
          vec3 p = r.point_at_parameter(2.0);
          col += color(r, world, 0);
        }
        col /= float(ns);

        // gamme correct?
        col = vec3(sqrt(col.r()), sqrt(col.g()), sqrt(col.b()));

        // image_lock.lock();
        image[c] = vec3(col.r(), col.g(), col.b());
        // image_lock.unlock();
      }

      std::cerr << "_";
      done[t] = true;
    };

    threads[t] = std::thread(render);
  }

  int i, j;
  while (!d->is_done()) {
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
}

hitable *random_scene(int n) {
  hitable **list = new hitable *[n + 1];
  list[0] =
      new sphere(vec3(0, -1000, 0), 1000, new lambertian(vec3(0.5, 0.5, 0.5)));
  int i = 1;
  for (int a = -11; a < 11; a++) {
    for (int b = -11; b < 11; b++) {
      float choose_mat = drand48();
      vec3 center(a + 0.9 * drand48(), 0.2, b + 0.9 * drand48());
      if ((center - vec3(4, 0.2, 0)).length() > 0.9) {
        if (choose_mat < 0.8) { // diffuse
          list[i++] = new sphere(
              center, 0.2,
              new lambertian(vec3(drand48() * drand48(), drand48() * drand48(),
                                  drand48() * drand48())));

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

  list[i++] = new sphere(vec3(0, 1, 0), 1.0, new dielectric(1.5));
  list[i++] =
      new sphere(vec3(-3, 1, 0), 1.0, new lambertian(vec3(0.4, 0.2, 0.1)));
  list[i++] =
      new sphere(vec3(4, 1, 0), 1.0, new metal(vec3(0.7, 0.6, 0.5), 0.0));

  return new hitable_list(list, i);
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
