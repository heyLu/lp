function findSlider(input) {
var type = "unknown";
var match = input.match(/^(uniform|attribute) (float|vec2|vec3) (.*);/);
if (!match) {
  return false;
}
type = match[2];
var name = match[3];

var numRegexp = /[+-]?(?:\d+(?:\.\d+)?)/;
var singleValueRegexp = new RegExp(`\\((${numRegexp.source}),\\s*(${numRegexp.source}),\\s*(${numRegexp.source})\\)`);
var vec2Regexp = new RegExp(`^${singleValueRegexp.source},\\s*${singleValueRegexp.source}$`)
var vec3Regexp = new RegExp(`^${singleValueRegexp.source},\\s*${singleValueRegexp.source},\\s*${singleValueRegexp.source}$`)

var match = input.match(/\/\/#slider\[(.*)\]/);
if (!match) {
  return false;
}

var rawVal = match[1];

if (type == "float") {
  var vals = rawVal.split(",");

  vals = vals.map(function(v) {
    return Number(v);
  });
} else if (type == "vec2" || type == "vec3") {
  var re = type == "vec2" ? vec2Regexp : vec3Regexp;
  var match = rawVal.match(re);
  if (!match) {
    throw new Error("invalid " + type + " slider");
  }

  var vals = [];
  for (var i = 0; i < Number(type[type.length-1]); i++) {
    vals.push([match[1+i*3+0], match[1+i*3+1], match[1+i*3+2]].map(Number));
  }
} else {
  throw new Error("unknown type");
}

return {name: name, type: type, range: vals};
}

function findSliders(input) {
var vars = [];
input.split("\n").forEach(function(line) {
  var slider = findSlider(line);
  if (slider) {
    vars.push(slider);
  }
});

return vars;
}
