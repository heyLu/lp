#include <chrono>
#include <float.h>
#include <fstream>
#include <iostream>

#include "camera.h"
#include "hitable_list.h"
#include "lambertian.h"
#include "metal.h"
#include "sphere.h"

vec3 color(const ray &r, hitable *world, int depth) {
  hit_record rec;
  if (world->hit(r, 0.001, MAXFLOAT, rec)) {
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

int main(int argc, char **argv) {
  std::string out_name = "/dev/stdout";
  bool write_partial = false;
  if (argc > 1) {
    write_partial = true;
    out_name = argv[1];
  }

  int nx = 200;
  int ny = 100;
  int ns = 100;
  hitable *list[4];
  list[0] =
      new sphere(vec3(0, 0, -1), 0.5, new lambertian(vec3(0.8, 0.3, 0.3)));
  list[1] =
      new sphere(vec3(0, -100.5, -1), 100, new lambertian(vec3(0.8, 0.8, 0.0)));
  list[2] = new sphere(vec3(1, 0, -1), 0.5, new metal(vec3(0.8, 0.6, 0.2)));
  list[3] = new sphere(vec3(-1, 0, -1), 0.5, new metal(vec3(0.8, 0.8, 0.8)));
  hitable *world = new hitable_list(list, 4);
  camera cam;

  int c = 0;
  vec3 *image[nx * ny];
  for (int j = ny - 1; j >= 0; j--) {
    for (int i = 0; i < nx; i++) {
      image[c] = new vec3(254, 254, 254);
      c++;
    }
  }

  c = 0;
  auto start = std::chrono::high_resolution_clock::now();
  for (int j = ny - 1; j >= 0; j--) {
    for (int i = 0; i < nx; i++) {
      vec3 col(0, 0, 0);

      // anti-aliasing (sample #ns rays)
      for (int s = 0; s < ns; s++) {
        float u = float(i + drand48()) / float(nx);
        float v = float(j + drand48()) / float(ny);
        ray r = cam.get_ray(u, v);
        vec3 p = r.point_at_parameter(2.0);
        col += color(r, world, 0);
      }
      col /= float(ns);

      // gamme correct?
      col = vec3(sqrt(col.r()), sqrt(col.g()), sqrt(col.b()));

      image[c] = new vec3(col.r(), col.g(), col.b());

      if ((c) % int(nx * ny / 100.0) == 0) {
        std::cerr << "."; // progress dots ✨

        if (write_partial) {
          // write updated image
          std::ofstream out;
          out.open(out_name);
          out << "P3\n" << nx << " " << ny << "\n255\n";
          int cc = 0;
          for (int y = ny - 1; y >= 0; y--) {
            for (int x = 0; x < nx; x++) {
              vec3 col = *image[cc];
              int ir = int(255.99 * col.r());
              int ig = int(255.99 * col.g());
              int ib = int(255.99 * col.b());
              out << ir << " " << ig << " " << ib << "\n";
              cc++;
            }
          }
          out.flush();
          out.close();
        }
      }

      c++;
    }
  }
  std::cerr << "\n";
  auto finish = std::chrono::high_resolution_clock::now();
  std::cerr << "took "
            << std::chrono::duration_cast<std::chrono::milliseconds>(finish -
                                                                     start)
                   .count()
            << "ms"
            << "\n";

  std::ofstream out;
  out.open(out_name);
  out << "P3\n" << nx << " " << ny << "\n255\n";
  c = 0;
  for (int y = ny - 1; y >= 0; y--) {
    for (int x = 0; x < nx; x++) {
      vec3 col = *image[c];
      int ir = int(255.99 * col.r());
      int ig = int(255.99 * col.g());
      int ib = int(255.99 * col.b());
      out << ir << " " << ig << " " << ib << "\n";

      c++;
    }
  }
  out.flush();
  out.close();
}

inline std::istream &operator>>(std::istream &is, vec3 &t) {
  is >> t.e[0] >> t.e[1] >> t.e[2];
  return is;
}

inline std::ostream &operator<<(std::ostream &os, const vec3 &t) {
  os << t.e[0] << " " << t.e[1] << " " << t.e[2];
  return os;
}

inline void vec3::make_unit_vector() {
  float k = 1.0 / sqrt(e[0] * e[0] + e[1] * e[1] + e[2] * e[2]);
  e[0] *= k;
  e[1] *= k;
  e[2] *= k;
}
