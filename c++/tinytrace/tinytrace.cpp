#include <chrono>
#include <float.h>
#include <iostream>

#include "camera.h"
#include "hitable_list.h"
#include "sphere.h"

vec3 random_in_unit_sphere() {
  vec3 p;
  do {
    p = 2.0 * vec3(drand48(), drand48(), drand48()) - vec3(1, 1, 1);
  } while (p.squared_length() >= 1.0);
  return p;
}

vec3 color(const ray &r, hitable *world) {
  hit_record rec;
  if (world->hit(r, 0.001, MAXFLOAT, rec)) {
    vec3 target = rec.p + rec.normal + random_in_unit_sphere();
    return 0.5 * color(ray(rec.p, target - rec.p), world);
  } else {
    vec3 unit_direction = unit_vector(r.direction());
    float t = 0.5 * (unit_direction.y() + 1.0);
    // lerp() for sky color
    return (1.0 - t) * vec3(1.0, 1.0, 1.0) + t * vec3(0.5, 0.7, 1.0);
  }
}

int main() {
  int nx = 200;
  int ny = 100;
  int ns = 100;
  std::cout << "P3\n" << nx << " " << ny << "\n255\n";
  hitable *list[2];
  list[0] = new sphere(vec3(0, 0, -1), 0.5);
  list[1] = new sphere(vec3(0, -100.5, -1), 100);
  hitable *world = new hitable_list(list, 2);
  camera cam;

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
        col += color(r, world);
      }
      col /= float(ns);

      // gamme correct?
      col = vec3(sqrt(col.r()), sqrt(col.g()), sqrt(col.b()));

      int ir = int(255.99 * col.r());
      int ig = int(255.99 * col.g());
      int ib = int(255.99 * col.b());
      std::cout << ir << " " << ig << " " << ib << "\n";
    }
  }
  auto finish = std::chrono::high_resolution_clock::now();
  std::cerr << "took "
            << std::chrono::duration_cast<std::chrono::milliseconds>(finish -
                                                                     start)
                   .count()
            << "ms"
            << "\n";
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
