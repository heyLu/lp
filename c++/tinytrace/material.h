#ifndef MATERIAL_H
#define MATERIAL_H

#include "hitable.h"
#include "texture.h"

class material {
public:
  virtual bool scatter(const ray &r_in, const hit_record &rec,
                       vec3 &attenuation, ray &scattered) const = 0;
  virtual vec3 emitted(float u, float v, const vec3 &p) const {
    return vec3(0, 0, 0); // default emit nothing
  }
};

class diffuse_light : public material {
public:
  diffuse_light(texture *a) : emit(a) {}

  virtual bool scatter(const ray &r_in, const hit_record &rec,
                       vec3 &attentuation, ray &scattered) const {
    return false;
  }

  virtual vec3 emitted(float u, float v, const vec3 &p) const {
    return emit->value(u, v, p);
  }

  texture *emit;
};

#endif
