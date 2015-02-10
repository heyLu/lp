require('./style.css');

class Greeter {
	greet(name = "World", suffix = "!") {
		console.log(`Hello, ${name}${suffix}`);
	}
}

new Greeter().greet("Alice");
