#ifndef METAL_H
#define METAL_H

#include "material.h"

vec3 reflect(const vec3 &v, const vec3 &n) { return v - 2 * dot(v, n) * n; }

class metal : public material {
public:
  metal(const vec3 &a) : albedo(a) {}

  virtual bool scatter(const ray &r_in, const hit_record &rec,
                       vec3 &attenuation, ray &scattered) const {
    vec3 reflected = reflect(unit_vector(r_in.direction()), rec.normal);
    scattered = ray(rec.p, reflected);
    attenuation = albedo;
    return (dot(scattered.direction(), rec.normal) > 0);
  }

  vec3 albedo;
};

#endif
