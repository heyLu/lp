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
  errorEl.style = "color: red; position: absolute; right: 0; bottom: 0; background-color: black; z-index: 4";
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
  height: 100%;
}

#editor textarea {
  display: block;
  width: 72ex;
  height: calc(100% - 2em + 1px); /* TODO: fix this */
  resize: horizontal;
  padding: 1ex;
  box-sizing: border-box;
}

#editor textarea, #editor input {
  border: none;
  background-color: rgba(255, 255, 255, 0.8);
}

#editor-changed {
  position: absolute;
  right: 0;
  top: 1.5em;
}
  `
  document.head.appendChild(styleEl);
  
  function TwoTriangles(canvas, fragmentShaderSrc, options) {
    var options = options || {};
    var tt = {};
    tt.canvas = canvas;
    tt.fragmentShaderSrc = transformShader(fragmentShaderSrc);
    
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
    
    tt.origin = gl.getUniformLocation(program, 'origin');

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
    
    tt.direction = [0, 0, 1];
    tt.right = [-1, 0, 0];
    tt.angle = {horizontal: 0.0, vertical: 0.0};
    
    var iDirection = gl.getUniformLocation(program, 'iDirection');
    gl.uniform3f(iDirection, tt.direction[0], tt.direction[1], tt.direction[2]);
    
    var mousePressed = false;
    tt.canvas.addEventListener("mousedown", function(ev) { if (ev.button == 0) { mousePressed = true; }});
    tt.canvas.addEventListener("mouseup", function(ev) { mousePressed = false; });
    
    tt.canvas.addEventListener("mousemove", function(ev) {
      gl.uniform3f(iMouse, ev.clientX, canvas.clientHeight - ev.clientY, 0.0);
      
      if (!mousePressed || !iDirection) {
        return;
      }
      
      tt.angle.horizontal += 0.0001 * (tt.w/2 - ev.clientX);
      tt.angle.vertical += 0.0001 * (tt.h/2 - ev.clientY);
      tt.direction = [Math.cos(tt.angle.vertical) * Math.sin(tt.angle.horizontal),
                      Math.sin(tt.angle.vertical),
                      Math.cos(tt.angle.vertical) * Math.cos(tt.angle.horizontal)];
      tt.right = [Math.sin(tt.angle.horizontal - Math.PI/2),
                  0,
                  Math.cos(tt.angle.horizontal - Math.PI/2)];
      
      gl.uniform3f(iDirection, tt.direction[0], tt.direction[1], tt.direction[2]);

      requestAnimationFrame(tt.render);
    });
    
    window.addEventListener('keydown', function(ev) {
      if (document.activeElement != document.body || ev.ctrlKey || ev.altKey) {
        return;
      }
      
      var origin = tt.gl.getUniform(tt.program, tt.origin);
      
      switch (ev.keyCode) {
        case 38: // Up
        case 87: // W
          origin = [origin[0]+tt.direction[0], origin[1]+tt.direction[1], origin[2]+tt.direction[2]];
          break;
        case 37: // Left
        case 65: // A
          origin = [origin[0]-tt.right[0], origin[1]-tt.right[1], origin[2]-tt.right[2]];
          break;
        case 40: // Down
        case 83: // S
          origin = [origin[0]-tt.direction[0], origin[1]-tt.direction[1], origin[2]-tt.direction[2]];
          break;
        case 39: // Right
        case 68: // D
          origin = [origin[0]+tt.right[0], origin[1]+tt.right[1], origin[2]+tt.right[2]];
          break;
        default:
          return
      }
      
      ev.preventDefault();
      
      gl.uniform3f(tt.origin, origin[0], origin[1], origin[2]);
      requestAnimationFrame(tt.render);
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
  
  var tt = TwoTriangles(canvas, files.open('default.frag').content);

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
  editor.el = document.createElement("textarea");
  editor.el.value = tt.fragmentShaderSrc;
  
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
      editor.container.style.display = "none";
    } else {
      editor.container.style.display = "inherit";
      editor.el.focus();
    }
    editor.visible = !editor.visible;
  }
  
  window.addEventListener('keydown', function(ev) {
    if (ev.ctrlKey && ev.keyCode == 69) { // Ctrl + e
      ev.preventDefault();
      editor.toggle();
    }
  });
  
  editor.ui = files.setupUI(editor.el);
  editor.container = document.createElement("div");
  editor.container.id = "editor";
  editor.visible = false;
  editor.container.style.display = "none";
  editor.container.appendChild(editor.ui);
  editor.container.appendChild(editor.el);
  
  document.body.appendChild(editor.container);
} catch (e) {
  displayError(e);
}
