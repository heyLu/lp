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
  document.body.style = "margin: 0; overflow: hidden;";
  document.body.innerHTML = "";
  
  var errorEl = document.createElement("pre");
  errorEl.style = "color: red; position: absolute; right: 0; bottom: 0;";
  document.body.appendChild(errorEl);
  
  function displayError(e) {
    window.error = e;
    errorEl.textContent = e;
    console.error(e);
  }
  
  function clearError() {
    errorEl.textContent = "";
  }
  
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

  var fragmentShaderSrc = `precision highp float;

uniform float iGlobalTime;
uniform vec2 iResolution;
uniform vec3 iMouse;

const int MaximumRaySteps = 150;
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
}

float sphere(vec3 pos) {
  return length(pos) - 1.0;
}

float pMod1(inout float p, float size) {
  float halfsize = size * 0.5;
  float c = floor((p + halfsize)/size);
  p = mod(p+halfsize, size)-halfsize;
  return c;
}

uniform vec3 offset; //#slider[(0.0,10.0,20.0),(0.0,10.0,20.0),(0.0,2.5,20.0)]

float DistanceEstimator(vec3 pos) {
  pMod1(pos.x, offset.x);
  pMod1(pos.y, offset.y);
  pMod1(pos.z, offset.z + sin(iGlobalTime));
  return sphere(pos);
}

mat3 setCamera( in vec3 ro, in vec3 ta, float cr ) {
	vec3 cw = normalize(ta-ro);
	vec3 cp = vec3(sin(cr), cos(cr),0.0);
	vec3 cu = normalize( cross(cw,cp) );
	vec3 cv = normalize( cross(cu,cw) );
  return mat3( cu, cv, cw );
}

uniform vec3 origin; //#slider[(-10.0,0.41,10.0),(-10.0,2.03,10.0),(-10.0,-1.34,10.0)]
uniform vec3 angle; //#slider[(-3.0,0.31,3.0),(-3.0,1.77,3.0),(-3.0,-0.18,3.0)]
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
   
  var styleEl = document.createElement("style");
  styleEl.textContent = `

#sidebar {
  position: absolute;
  top: 0;
  right: -250px;
  z-index: 1;

  padding: 1ex;

  font-family: monospace;
  font-weight: bold;

  background-color: rgba(255, 255, 255, 0.5);
}

#sidebar:hover {
  transition: right 0.1s;
  right: 0;
}

#editor {
  position: absolute;
  top: 0;
  left: 0;

  border: none;
  background-color: rgba(255, 255, 255, 0.8);

  min-width: 72ex;
  height: 100vh;
}
  `
  document.head.appendChild(styleEl);
  
  function TwoTriangles(canvas, fragmentShaderSrc, options) {
    var options = options || {};
    var tt = {};
    tt.canvas = canvas;
    tt.fragmentShaderSrc = fragmentShaderSrc;
    
    tt.w = canvas.width = window.innerWidth;
    tt.h = canvas.height = window.innerHeight;

    var gl = tt.gl = canvas.getContext("webgl");
    if (!gl) { alert("i think your browser does not support webgl"); }

    gl.clearColor(0.0, 0.0, 0.0, 1.0);
    gl.clear(gl.COLOR_BUFFER_BIT);

    var vertexShader = compileShader(gl, gl.VERTEX_SHADER, vertexShaderSrc);
    var fragmentShader = compileShader(gl, gl.FRAGMENT_SHADER, tt.fragmentShaderSrc);

    var program = tt.program = linkShaders(gl, vertexShader, fragmentShader);
    gl.useProgram(program);

    var aPosition = gl.getAttribLocation(program, 'aPosition');
    var iResolution = gl.getUniformLocation(program, 'iResolution');
    var iGlobalTime = gl.getUniformLocation(program, 'iGlobalTime');
    var iMouse = gl.getUniformLocation(program, 'iMouse');

    gl.vertexAttrib2f(aPosition, 0.0, 0.0);
    gl.uniform2f(iResolution, canvas.width, canvas.height);
    gl.uniform1f(iGlobalTime, 0.0);
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

    tt.render = function() {
      gl.drawArrays(gl.TRIANGLES, 0, 6);
    };
    
    tt.playing = options.playing || false;
    tt.draw = function(timestamp) {
      if (tt.playing) {
        requestAnimationFrame(tt.draw);
      }
      gl.uniform1f(iGlobalTime, timestamp / 1000);
      tt.render();
    }
    
    tt.playingControls = document.createElement("div");
    var playPauseButton = document.createElement("button");
    playPauseButton.textContent = "Play";
    playPauseButton.onclick = function(ev) {
      tt.togglePlaying();
    }
    tt.playingControls.appendChild(playPauseButton);
    
    tt.togglePlaying = function() {
      if (!tt.playing) {
        tt.playing = true;
        tt.draw();
      } else {
        tt.playing = false;
      }
      playPauseButton.textContent = tt.playing ? "Pause" : "Play";
    }
    
    window.addEventListener('keydown', function(ev) {
      if (ev.keyCode == 32 && document.activeElement != editor.el) { // Space
        ev.preventDefault();
        tt.togglePlaying();
      }
    });
    
    tt.canvas.addEventListener("mousemove", function(ev) {
      gl.uniform3f(iMouse, ev.mouseX, ev.mouseY, 0.0);
    });
    
      
    window.onresize = function(ev) {
      tt.w = tt.canvas.width = window.innerWidth;
      tt.h = tt.canvas.height = window.innerHeight;
      gl.viewport(0, 0, tt.w, tt.h);
      gl.uniform2f(iResolution, tt.w, tt.h);
      tt.render();
    };
    
    tt.render();
    if (tt.playing) {
      tt.draw();
    }
    return tt;
  }
  
  var canvas = document.createElement("canvas");
  document.body.appendChild(canvas);
  
  var tt = TwoTriangles(canvas, fragmentShaderSrc);

  var sidebarEl = document.createElement("div");
  sidebarEl.id = "sidebar";
  document.body.appendChild(sidebarEl);
  
  var sliders = findSliders(tt.fragmentShaderSrc);
  initSliders(tt.gl, tt.program, sliders, function(ev) {
    requestAnimationFrame(tt.render);
  });
  addSliders(sidebarEl, sliders);
  
  sidebarEl.appendChild(tt.playingControls);
  
  tt.render();
  
  var editor = {};
  editor.visible = false;
  editor.el = document.createElement("textarea");
  editor.el.id = "editor";
  editor.el.style.display = "none";
  editor.el.value = tt.fragmentShaderSrc;
  document.body.appendChild(editor.el);
  
  editor.el.onkeydown = function(ev) {
    try {
      if (ev.ctrlKey && ev.keyCode == 13) {
        tt = TwoTriangles(canvas, editor.el.value, {playing: tt.playing});

        sidebarEl.innerHTML = "";
        sliders = findSliders(tt.fragmentShaderSrc);
        initSliders(tt.gl, tt.program, sliders, function(ev) {
          // rerender if not in animation mode (otherwise `tt.draw` will already do that)
          if (!tt.playing) {
            requestAnimationFrame(tt.render);
          }
        });
        addSliders(sidebarEl, sliders);
        
        sidebarEl.appendChild(tt.playingControls);

        tt.render();
        clearError();
      }
    } catch (e) {
      displayError(e);
    }
  }
  
  editor.toggle = function() {
    if (editor.visible) {
      editor.el.style.display = "none";
    } else {
      editor.el.style.display = "inherit";
      editor.el.focus();
    }
    editor.visible = !editor.visible;
  }
  
  window.addEventListener('keydown', function(ev) {
    if (ev.ctrlKey && ev.keyCode == 72) { // Ctrl + H
      ev.preventDefault();
      editor.toggle();
    }
  });
} catch (e) {
  displayError(e);
}
