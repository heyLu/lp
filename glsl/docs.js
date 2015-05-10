/*document.head.innerHTML = "";
document.body.innerHTML = "";
document.body.style = "margin: 0";*/

var styleEl = document.createElement("style");
styleEl.textContent = `
#toggle {
position: absolute;
bottom: 0.5ex;
right: 1ex;
z-index: 3;
font-size: 20px;
font-weight: bold;
color: white;
text-decoration: none;
}

#docs {
  position: absolute:
  top: 100vh;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 2;
}

#docs:target {
  position: absolute;
  top: 0;
}

#docs:target div {
  display: flex;
  justify-content: center;
}

#docs:target div pre {
  margin: 0;
  margin-top: 5vh;
  padding: 1em 2em;
  box-sizing: border-box;

  max-width: 1000px;
  height: 90vh;

  font-size: larger;
  /*font-weight: bold;*/

  background-color: rgba(255, 255, 255, 0.8);

  white-space: pre-wrap;

  overflow-y: scroll;
}
`
document.head.appendChild(styleEl);

var toggleEl = document.createElement("a");
toggleEl.id = "toggle";
toggleEl.href = "#docs";
toggleEl.textContent = "?";
function handleToggle(ev) {
  if (location.hash == "#docs") {
    location.hash = "";
  } else {
    location.hash = "#docs";
  }
  ev.preventDefault();
}
toggleEl.onclick = handleToggle;
document.body.appendChild(toggleEl);

var docsEl = document.createElement("pre");
docsEl.textContent = `# shaders!

Use the sliders on the left to change values used in the shader.

## Keyboard shortcuts

- \`Ctrl-Enter\` reruns the shader
- \`Ctrl-e\` toggles the editor
- \`Space\` toggles animation
- \`Ctrl-h\` toggles help

## Default uniforms

- \`uniform vec2 iResolution\`: resolution of the canvas
- \`uniform vec2 iMouse\`: mouse position (updated on mousemove)
- \`uniform float iGlobalTime\`: animation time in seconds

## Special comments

Special comments can be used to create sliders and to include other shaders, making it possible to create libraries of reusable functions.

### Sliders (\`//#slider[...]\`)

\`uniform {float,vec2,vec3} name; //#slider[...]\`

  Creates a slider called \`name\`.

  The value format is as follows:

    - Numbers are specified as \`(min,default,max)\`.
      (Omit the parens for \`float\` values, though.)
    - For vectors, the values must be separated by commas.

  For example, the following creates a slider whose \`.x\` component ranges from 0.0 to 2.0 and defaults to 1.0 and whose \`.y\` component ranges from -3.0 to 3.0, defaulting to 0.0:

    uniform vec2 example; //#slider[(0.0,1.0,2.0),(-3.0,0.0,3.0)]

### Includes (\`//#include "<name>"\`)

For example, \`//#include "includes/iq-primitives.frag"\` includes the primitives from iq's distance functions page at http://www.iquilezles.org/www/articles/distfunctions/distfunctions.htm.

## Resources

...

## Fin

That's it, now go and write some shaders!  The source code (*very* dirty) is at https://github.com/heyLu/lp/tree/master/glsl if you want to look at it or fix things.
`

var docsContainerEl = document.createElement("div");
docsContainerEl.id = "docs";
var innerContainerEl = document.createElement("div");
innerContainerEl.appendChild(docsEl);
docsContainerEl.appendChild(innerContainerEl);
document.body.appendChild(docsContainerEl);

window.addEventListener("keydown", function(ev) {
  if (ev.ctrlKey && ev.keyCode == 72) { // Ctrl-h
    handleToggle(ev);
  }
});