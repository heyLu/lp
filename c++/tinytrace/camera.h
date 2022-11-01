#ifndef CAMERA_H
#define CAMERA_H

#include "ray.h"

vec3 random_in_unit_disk() {
  vec3 p;
  do {
    p = 2.0 * vec3(drand48(), drand48(), 0) - vec3(1, 1, 0);
  } while (dot(p, p) >= 1.0);
  return p;
}

class camera {
public:
  camera(vec3 look_from, vec3 look_at, vec3 up, float vfov, float aspect,
         float aperture, float focus_dist) { // vfov is top to bottom in degrees
    lens_radius = aperture / 2;
    float theta = vfov * M_PI / 180;
    float half_height = tan(theta / 2);
    float half_width = aspect * half_height;
    origin = look_from;
    w = unit_vector(look_from - look_at);
    u = unit_vector(cross(up, w));
    v = cross(w, u);
    lower_left_corner = origin - half_width * focus_dist * u -
                        half_height * focus_dist * v - focus_dist * w;
    horizontal = 2 * half_width * focus_dist * u;
    vertical = 2 * half_height * focus_dist * v;
  }

  ray get_ray(float s, float t) const {
    vec3 rd = lens_radius * random_in_unit_disk();
    vec3 offset = u * rd.x() + v * rd.y();
    return ray(origin + offset, lower_left_corner + s * horizontal +
                                    t * vertical - origin - offset);
  }

  vec3 origin;
  vec3 lower_left_corner;
  vec3 horizontal;
  vec3 vertical;
  vec3 u, v, w;
  float lens_radius;
};

#endif
