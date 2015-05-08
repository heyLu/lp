// # Rendering Signed Distance Functions with WebGL
//
// Yes, you can do this in Shadertoy, but this is how you can do it from scratch.
// (It also loads a lot faster.)
//
// Use it by opening the Scratchpad in Firefox and evaluate the code with `Ctrl-R`
// to rerender.  Including via `<script>` should also work.
//
// # Resources
//
// - [Raw WebGL](http://nickdesaulniers.github.io/RawWebGL) ([video](https://www.youtube.com/watch?v=H4c8t6myAWU))
// - [WebGL Fundamentals](http://webglfundamentals.org/webgl/lessons/webgl-fundamentals.html)

try {
  function compileShader(gl, type, shaderSrc) {
    var shader = gl.createShader(type);
    gl.shaderSource(shader, shaderSrc);
    gl.compileShader(shader);

    if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS)) {
      throw new Error(gl.getShaderInfoLog(shader));
    }

    return shader;
  }
  
  function linkShaders(gl, vertexShader, fragmentShader) {
    var program = gl.createProgram();
    gl.attachShader(program, vertexShader);
    gl.attachShader(program, fragmentShader);
    gl.linkProgram(program);
    
    if (!gl.getProgramParameter(program, gl.LINK_STATUS)) {
      throw new Error(gl.getProgramInfoLog(program));
    }
    
    return program;
  }
  
  function initBuffer(gl, data, elemPerVertex, attribute) {
    var buffer = gl.createBuffer();
    if (!buffer) {
      throw new Error('failed to create buffer');
    }
    
    gl.bindBuffer(gl.ARRAY_BUFFER, buffer);
    gl.bufferData(gl.ARRAY_BUFFER, data, gl.STATIC_DRAW);
    gl.vertexAttribPointer(attribute, elemPerVertex, gl.FLOAT, false, 0, 0);
    gl.enableVertexAttribArray(attribute);
  }

  var vertexShaderSrc = `
attribute vec4 aPosition;

void main() {
  gl_Position = aPosition;
}
`

  var fragmentShaderSrc = `
precision highp float;

uniform vec2 iResolution;
uniform vec3 iMouse;

const int MaximumRaySteps = 150;
const float MinimumDistance = 0.0000001;

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
}

float DistanceEstimator(vec3 pos) {
  return length(pos) - 1.0;
}

/*uniform int MaxIterations; //#slider[1,50,200]
const float bailout = 4.0;
const float power = 8.0;
const float phaseX = 0.0;
const float phaseY = 0.0;

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
}*/

mat3 setCamera( in vec3 ro, in vec3 ta, float cr ) {
	vec3 cw = normalize(ta-ro);
	vec3 cp = vec3(sin(cr), cos(cr),0.0);
	vec3 cu = normalize( cross(cw,cp) );
	vec3 cv = normalize( cross(cu,cw) );
  return mat3( cu, cv, cw );
}

uniform vec3 origin; //#slider[(-10.0,1.0,10.0),(-10.0,2.0,10.0),(-10.0,-1.0,10.0)]
uniform vec3 angle; //#slider[(-3.0,0.0,3.0),(-3.0,0.0,3.0),(-3.0,0.0,3.0)]
uniform vec3 color; //#slider[(0.0, 1.0, 1.0),(0.0,0.0,1.0),(0.0,0.0,1.0)]
uniform float colorMix; //#slider[0.0,0.9,1.0]

void main() {
  vec2 q = gl_FragCoord.xy / iResolution.xy;
  vec2 p = -1.0 + 2.0*q;
  p.x *= iResolution.x / iResolution.y;
  vec2 mo = iMouse.xy/iResolution.xy;

  float time = 15.0 + 0.0; // iGlobalTime

  // camera	
  vec3 ro = origin; //vec3( -0.5+3.2*cos(0.1*time + 6.0*mo.x), 1.0 + 2.0*mo.y, 0.5 + 3.2*sin(0.1*time + 6.0*mo.x) );
  vec3 ta = angle; //vec3( -0.5, -0.4, 0.5 );

  // camera-to-world transformation
  mat3 ca = setCamera( ro, ta, 0.0 );

  // ray direction
  vec3 rd = ca * normalize( vec3(p.xy, 2.5) );

  // render	
  float dist = trace(ro, rd);
  vec3 col = vec3(dist, dist, dist);

  col = mix(color, col, colorMix);
  //col = pow( col, vec3(0.4545));

  gl_FragColor = vec4( col, 1.0 );
}
`
  
  document.body.style = "margin: 0; overflow: hidden;";
  document.body.innerHTML = "";
  
  var styleEl = document.createElement("style");
  styleEl.textContent = `

#sidebar {
  position: absolute;
  top: 0;
  right: -250px;

  padding: 1ex;

  font-family: monospace;
  font-weight: bold;

  background-color: rgba(255, 255, 255, 0.5);
}

#sidebar:hover {
  transition: right 0.1s;
  right: 0;
}
  `
  document.head.appendChild(styleEl);

  var canvas = document.createElement("canvas");
  var w = canvas.width = window.innerWidth;
  var h = canvas.height = window.innerHeight;
  document.body.appendChild(canvas);

  var gl = canvas.getContext("webgl");
  if (!gl) { alert("i think your browser does not support webgl"); }

  gl.clearColor(0.0, 0.0, 0.0, 1.0);
  gl.clear(gl.COLOR_BUFFER_BIT);
  
  var vertexShader = compileShader(gl, gl.VERTEX_SHADER, vertexShaderSrc);
  var fragmentShader = compileShader(gl, gl.FRAGMENT_SHADER, fragmentShaderSrc);
  
  var program = linkShaders(gl, vertexShader, fragmentShader);
  gl.useProgram(program);
  
  var aPosition = gl.getAttribLocation(program, 'aPosition');
  var iResolution = gl.getUniformLocation(program, 'iResolution');
  var iMouse = gl.getUniformLocation(program, 'iMouse');
  
  gl.vertexAttrib2f(aPosition, 0.0, 0.0);
  gl.uniform2f(iResolution, canvas.width, canvas.height);
  gl.uniform3f(iMouse, 0.0, 0.0, 0.0);
  
  // two triangles
  var positions = new Float32Array([
    -1.0, 1.0,
    -1.0, -1.0,
    1.0, 1.0,
    
    1.0, 1.0,
    -1.0, -1.0,
    1.0, -1.0
  ]);
  initBuffer(gl, positions, 2, aPosition);
  
  function render() {
   gl.drawArrays(gl.TRIANGLES, 0, 6);
  }
  
  document.body.addEventListener("mousemove", function(ev) {
    gl.uniform3f(iMouse, ev.mouseX, ev.mouseY, 0.0);
  });
  
  function draw() {
    requestAnimationFrame(draw);
    render();
  }
  
  window.onresize = function(ev) {
    w = canvas.width = window.innerWidth;
    h = canvas.height = window.innerHeight;
    gl.viewport(0, 0, w, h);
    gl.uniform2f(iResolution, w, h);
    render();
  };

  var sidebarEl = document.createElement("div");
  sidebarEl.id = "sidebar";
  document.body.appendChild(sidebarEl);
  
  var sliders = findSliders(fragmentShaderSrc);
  initSliders(gl, program, sliders, function(ev) {
    render();
  });
  
  addSliders(sidebarEl, sliders);
  
  render();
} catch (e) {
  window.error = e;
  var msg = document.createElement("pre");
  msg.style = "color: red; position: absolute; right: 0; bottom: 0";
  msg.textContent = e;
  document.body.appendChild(msg);
}
