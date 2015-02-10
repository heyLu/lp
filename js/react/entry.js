import React from "react";

require('./style.css');

class Greeter {
	greet(name = "World", suffix = "!") {
		console.log(`Hello, ${name}${suffix}`);
	}
}

new Greeter().greet("Alice");

var container = document.createElement("div");
document.body.appendChild(container);

React.render(<h1>???</h1>, container);
