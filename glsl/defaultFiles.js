files.current = "default.frag";
files.builtin = {
  "colors.frag": `uniform float blue; //#slider[0.0,0.0,1.0]

void main() {
  gl_FragColor = vec4(gl_FragCoord.xy / iResolution, blue, 1.0);
}`,
  "colors2.frag": `/*
 * A shader that maps the RGB scale onto the canvas.
 *
 * Blue controls the amount of blue to mix in, or the speed with which
 * to animate from/to blue when animating.
 */

uniform float blue; //#slider[0.0,0.0,1.0]

void main() {
  if (iGlobalTime > 0.0) {
    // if animating, control the speed of the color shift
    gl_FragColor = vec4(gl_FragCoord.xy / iResolution,
                        0.5 + sin(iGlobalTime * blue)*0.5, 1.0);
  } else {
    // else just control the blue component
    gl_FragColor = vec4(gl_FragCoord.xy / iResolution, blue, 1.0);
  }
}`,
  "default.frag": `//#include "includes/sphere-tracer.frag"
//#include "includes/default-main.frag"
//#include "includes/iq-primitives.frag"
//#include "includes/cupe-primitives.frag"

uniform vec3 offset; //#slider[(0.0,10.0,20.0),(0.0,10.0,20.0),(0.0,2.5,20.0)]

float DistanceEstimator(vec3 pos) {
  pMod1(pos.x, offset.x);
  pMod1(pos.y, offset.y);
  pMod1(pos.z, offset.z + sin(iGlobalTime));
  //return sphere(pos);
  //return min(sphere(vec3(pos.x, pos.y - 0.5, pos.z), 0.75),
  //           udBox(pos, vec3(1.0, 0.3, 1.0)));
  return min(max(-sphere(pos), udBox(pos, vec3(0.75))),
             sphere(pos, 0.05 + 0.25 * (1.0 + sin(iGlobalTime * 0.5)*0.5)));
}`,
  "mouse.frag": `void main() {
  gl_FragColor = vec4(iMouse.xy / iResolution, 0.0, 1.0);
}`,
  "includes/iq-primitives.frag": `float sphere(vec3 pos) {
  return length(pos) - 1.0;
}

float sphere(vec3 pos, float size) {
  return length(pos) - size;
}

float udBox( vec3 p, vec3 b ) {
  return length(max(abs(p)-b,0.0));
}
`,
  "includes/cupe-primitives.frag": `float pMod1(inout float p, float size) {
  float halfsize = size * 0.5;
  float c = floor((p + halfsize)/size);
  p = mod(p+halfsize, size)-halfsize;
  return c;
}`,
  "includes/sphere-tracer.frag": `const int MaximumRaySteps = 150;
const float MinimumDistance = 0.0001;

float DistanceEstimator(vec3 pos);

float trace(vec3 from, vec3 direction) {
	float totalDistance = 0.0;
  int stepsDone = 0;
	for (int steps = 0; steps < MaximumRaySteps; steps++) {
		vec3 p = from + totalDistance * direction;
		float distance = DistanceEstimator(p);
		totalDistance += distance;
    stepsDone = steps;
		if (distance < MinimumDistance) break;
	}
	return 1.0-float(stepsDone)/float(MaximumRaySteps);
}`,
  "includes/default-main.frag": `mat3 setCamera( in vec3 ro, in vec3 ta, float cr ) {
	vec3 cw = normalize(ta-ro);
	vec3 cp = vec3(sin(cr), cos(cr),0.0);
	vec3 cu = normalize( cross(cw,cp) );
	vec3 cv = normalize( cross(cu,cw) );
  return mat3( cu, cv, cw );
}

uniform vec3 origin; //#slider[(-10.0,0.41,10.0),(-10.0,2.03,10.0),(-10.0,-1.34,10.0)]
uniform float whee; //#slider[-3.1415,0.0,3.1415]

void main() {
  // screen coords (0,0)..(width,height) to "centered" & normalized coords (-1,-1)..(1,1)
  vec2 normalizedCoord = -1.0 + gl_FragCoord.xy / iResolution * 2.0;
  // aspect ratio (?)
  normalizedCoord.x *= iResolution.x / iResolution.y;

  mat3 ca = setCamera(origin, origin + iDirection, whee);
  vec3 direction = ca * normalize(vec3(normalizedCoord.xy, 2.5));

  // render
  float dist = trace(origin, direction);
  vec3 col = vec3(dist, dist, dist);

  gl_FragColor = vec4( col, 1.0 );
}`
};
