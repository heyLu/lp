/*
 * playing with rust, it's quite interesting so far.
 *
 * http://static.rust-lang.org/doc/0.5/tutorial.html
 */

/** a pointy structure */
struct Point { x: float, y: float }
/** shapey things */
enum Shape {
	Circle(Point, float),
	Rectangle(Point, Point)
}

impl Shape {
	/** calculate the area of a shape */
	fn area(&self) -> float {
		area(*self)
	}
}

impl Shape {
	fn weird(&self) -> int {
		42
	}
}

fn area(shape: Shape) -> float {
	match shape {
		Circle(_, size) => float::consts::pi * size * size,
		// Rectangle(p1, p2) => (p2.x - p1.x) * (p2.y - p1.y)
		Rectangle(Point {x: x1, y: y1}, Point {x: x2, y:y2}) =>
			(x2 - x1) * (y2 - y1)
	}
}

fn main() {
	io::println("Hello, Rust!?");
	
	let p1 = Point {x: 10f, y: 5f};
	let p2 = Point {x: 50f, y: 10f};
	let r1 = Rectangle(p1, p2);
	io::println(fmt!("area: %f", area(r1)));
	
	let x = @10;
	let y = ~10;
	let z: ~int = y; // y is 'moved' (i.e. can't be used afterwards)
	let z_copy: ~int = copy z;
	io::println(fmt!("%? %? %?", *x, *z, *z_copy));
	io::println(fmt!("%? %? %?", **&x, **&z, **&z_copy));
	//assert x != z; // yields a type error (i.e. == is a sane operator)
	assert z == z_copy;
	
	let countables = ~[1, 2, 3];
	let not_countables = ~[42];
	let stuff = countables + not_countables;
	io::println(fmt!("%? %?", stuff, countables + not_countables));
	assert stuff == countables + not_countables;
	
	let mut cs = ~[1, 2, 3, 4];
	cs += [3];
	io::println(fmt!("%?", [cs[0], cs[2]]));
	
	for cs.each |c| {
		io::println(fmt!("%?", *c));
	}
	
	fn each(v: &[int], op: fn(v: &int)) {
		let mut n = 0;
		while n < v.len() {
			op(&v[n]);
			n += 1;
		}
	}
	
	do each(cs) |c| { // must really by of type fn(v:&int) -> ()
		*c + 3;
	}
	
	fn faked_each(v: &[int], op: fn(v: &int) -> bool) {
		if v.len() < 3 {
			return;
		}
		
		if !op(&v[0]) { return; }
		if !op(&v[1]) { return; }
		if !op(&v[2]) { return; }
	}
	
	for faked_each([1, 2, 3, 4, 5, 6]) |&n| {
		io::println(fmt!("%d", n));
		if n % 2 == 0 {
			break;
		}
	}

	fn first_even(v: &[int]) -> int {
		for faked_each(v) |&n| {
			if n % 2 == 0 {
				return n;
			}
		}
		return -1;
	}
	let x = first_even([1, 2, 3, 4, 5, 6, 7]); io::println(fmt!("%?", x));
	
	io::println(fmt!("%? %?", r1.area(), r1.weird()));
	
	fn map<T, U>(vector: &[T], function: fn(v: T) -> U) -> ~[U] {
		let mut accum = ~[];
		for vector.each |&v| {
			accum.push(function(v));
		}
		return accum;
	}
	
	let cs2 = do map(cs) |c| {
		c + 1
	};
	io::println(fmt!("%?", cs2));
}
