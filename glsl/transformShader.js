var prelude = `precision highp float;

uniform float iGlobalTime;
uniform vec2 iResolution;
uniform vec3 iMouse;
uniform vec3 iDirection;

`;

function transformShader(shader, noPrelude, visited) {
  var visited = visited || {};
  var lines = shader.split("\n");
  lines = lines.map(function(line) {
    var match = line.match(/\/\/#include "(.*)"/);
    if (match) {
      var file = files.get(match[1]);
      if (!file.content) {
        throw new Error(`No such file: '${match[1]}'`);
      }
      if (file.name in visited) {
        throw new Error(`Include loop: '${match[1]}'`);
      }
      visited[file.name] = true;
      var expanded = transformShader(file.content, true, visited);
      return `//#includestart "${match[1]}" (start)
${expanded}
//#includeend "${match[1]}"`
    } else {
      return line;
    }
  });
  
  return (noPrelude ? "" : prelude) + lines.join("\n");
}