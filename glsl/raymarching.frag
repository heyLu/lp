/*
 * Ray marching in Fragmentarium.
 *
 * Based on the `trace` function from [1], runnable in
 * Fragmentarium [2]. Inspired by a talk about signed
 * distance functions [3], Fractal Lab [4] and lots of
 * other things [citation needed].
 *
 * [1]: http://blog.hvidtfeldts.net/index.php/2011/06/distance-estimated-3d-fractals-part-i/
 * [2]: http://syntopia.github.io/Fragmentarium/
 * [3]: https://www.youtube.com/watch?v=s8nFqwOho-s
 * [4]: http://sub.blue/fractal-lab
 */

// 3D.frag includes camera support and only requires
// us to define a `color` function.  (See below.)
#include "3D.frag"

#group DistranceEstimator
uniform float MinimumDistance; slider[0.0,0.01,10.0]
uniform int MaximumRaySteps; slider[1,10,100]

// Defined later, must be forward declared for use in `trace`.
float DistanceEstimator(vec3 pos);

// Adapted with minimal changes from [1].
float trace(vec3 from, vec3 direction) {
	float totalDistance = 0.0;
	int steps;
	for (steps=0; steps < MaximumRaySteps; steps++) {
		vec3 p = from + totalDistance * direction;
		float distance = DistanceEstimator(p);
		totalDistance += distance;
		if (distance < MinimumDistance) break;
	}
	return 1.0-float(steps)/float(MaximumRaySteps);
}

// This is where the actual calculation takes place!
float DistanceEstimator(vec3 pos) {
  return length(pos) - 1.0;
}

// Simple gray-scalish rendering.
vec3 color(vec3 pos, vec3 direction) {
  float dist = trace(pos, direction);
  return direction + vec3(dist, dist, dist);
}