/*
 * A "Mandelbulb", i.e. a mandelbrot fractal in 3d.
 *
 * Adapted for Fragmentarium from [1].
 *
 * [1]: http://2008.sub.blue/blog/2009/12/13/mandelbulb.html
 */
#include "3D.frag"

#group Mandelbulb
uniform float MinimumDistance; slider[0.0,0.01,10.0]
uniform int MaximumRaySteps; slider[1,10,200]
uniform int MaxIterations; slider[1,2,30]
uniform float bailout; slider[0.5,4.0,12.0]
uniform float power; slider[-20.0,8.0,20.0]
uniform float phaseX; slider[-2.0,0.0,2.0]
uniform float phaseY; slider[-2.0,0.0,2.0]

// Defined later, must be forward declared for use in `trace`.
float DistanceEstimator(vec3 pos);

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

float DistanceEstimator(vec3 z0) {
	vec3 c = z0;
	vec3 z = z0;
	float pd = power - 1.0; // power for derivative
	
	// Convert z to polar coordinates
	float r = length(z);
	float th = atan(z.y, z.x);
	float ph = asin(z.z / r);
	
	vec3 dz;
	float ph_dz = 0.0;
	float th_dz = 0.0;
	float r_dz	= 1.0;
	float powR, powRsin;
	
	// Iterate to compute the distance estimator.
	for (int n = 0; n < MaxIterations; n++) {
		// Calculate derivative of
		powR = power * pow(r, pd);
		powRsin = powR * r_dz * sin(ph_dz + pd*ph);
		dz.x = powRsin * cos(th_dz + pd*th) + 1.0;
		dz.y = powRsin * sin(th_dz + pd*th);
		dz.z = powR * r_dz * cos(ph_dz + pd*ph);
		
		// polar coordinates of derivative dz
		r_dz  = length(dz);
		th_dz = atan(dz.y, dz.x);
		ph_dz = acos(dz.z / r_dz);
		
		// z iteration
		powR = pow(r, power);
		powRsin = sin(power*ph);
		z.x = powR * powRsin * cos(power*th);
		z.y = powR * powRsin * sin(power*th);
		z.z = powR * cos(power*ph);
		z += c;
		
		r  = length(z);
		if (r > bailout) break;
		
		th = atan(z.y, z.x) + phaseX;
		ph = acos(z.z / r) + phaseY;
		
	}
	
	// Return the distance estimation value which determines the next raytracing
	// step size, or if whether we are within the threshold of the surface.
	return 0.5 * r * log(r)/r_dz;
}

vec3 color(vec3 pos, vec3 direction) {
  float dist = trace(pos, direction);
  return vec3(dist, dist, dist);
}

#preset Default
Eye = 0.0,0.0,2.0
Target 0.0,0.0,1.0
MinimumDistance=0.00001
MaximumRaySteps=100
MaxIterations=4