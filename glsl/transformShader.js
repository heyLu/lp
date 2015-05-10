var prelude = `precision highp float;

uniform float iGlobalTime;
uniform vec2 iResolution;
uniform vec3 iMouse;

`;

function transformShader(shader) {
  var lines = shader.split("\n");
  lines = lines.map(function(line) {
    var match = line.match(/\/\/#include "(.*)"/);
    if (match) {
      var file = files.get(match[1]);
      if (!file.content) {
        throw new Error(`No such file: '${match[1]}'`);
      }
      return `//#includestart "${match[1]}" (start)
${file.content}
//#includeend "${match[1]}"
`
    } else {
      return line;
    }
  });
  
  return prelude + lines.join("\n");
}