document.body.innerHTML = "";
var styleEl = document.createElement("style");
styleEl.textContent = `
.container .slider {
  margin-left: 1em;
}
`;
document.head.appendChild(styleEl);

function makeSlider(name, range, onChange) {
  var sliderContainerEl = document.createElement("div");
  sliderContainerEl.className = "slider";

  if (name) {
    var sliderLabelEl = document.createElement("label");
    sliderLabelEl.textContent = name;
    sliderContainerEl.appendChild(sliderLabelEl);
  }

  var sliderEl = document.createElement("input");
  sliderEl.type = "range";
  sliderEl.min = range[0];
  sliderEl.max = range[2];
  sliderEl.step = Math.min(0.01, (range[2] - range[0]) / 100);
  sliderEl.value = range[1];
  sliderEl.style.verticalAlign = "middle";
  sliderEl.addEventListener("input", onChange);
  sliderContainerEl.appendChild(sliderEl);

  var sliderValueEl = document.createElement("span");
  sliderValueEl.textContent = sliderEl.value;
  sliderContainerEl.appendChild(sliderValueEl);

  sliderEl.addEventListener("input", () => sliderValueEl.textContent = sliderEl.value);
  
  return sliderContainerEl;
}

function makeMultiSlider(name, ranges, onChange) {
  var containerEl = document.createElement("div");
  containerEl.className = "container";
  document.body.appendChild(containerEl);
      
  var labelEl = document.createElement("label");
  labelEl.textContent = name;
  containerEl.appendChild(labelEl);
    
  var names = [".x", ".y", ".z", ".y"];
  ranges.forEach((range, i) => {
    containerEl.appendChild(makeSlider(names[i], range));
  });
  
  return containerEl;
}

function addSliders(parent, sliders) {
  sliders.forEach((slider) => {
    switch (slider.type) {
      case "float":
        parent.appendChild(makeSlider(slider.name, slider.range));
        break;
      case "vec2":
      case "vec3":
        parent.appendChild(makeMultiSlider(slider.name, slider.range));
        break;
      default:
        parent.innerHTML += "unknown slider type " + slider.type;
    }
  });
}

var sliders = [{"name":" iPosition","type":"vec3","range":[[0.01,0.52,1.03],[0.04,0.55,1.06],[0.07,0.58,1.09]]},{"name":" iResolution","type":"vec2","range":[[0.01,0.52,1.03],[0.04,0.55,1.06]]},{"name":"fancyness","type":"float","range":[0.01,0.52,1.03]}];

addSliders(document.body, sliders);