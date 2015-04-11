/// `inc` is a function that increments it's argument, except if the
/// argument is 42.
///
/// # Arguments
///
/// * `x` - The number to increment
///
/// # Examples
///
/// ```rust
/// inc(3)  == 4
/// inc(42) == 42
/// inc(99) == 100
/// ```
fn inc(x: i32) -> i32 {
    if x == 42 {
        return 42
    }

    x + 1
}

struct Point {
    x: i32,
    y: i32
}

struct Meters(i32);

fn main() {
    let x = 5; // x: i32

    println!("Hello, World!");
    if x == 42 { // 42 is important
        println!("It's THE ANSWER!");
    } else {
        println!("Meh, it's just a number, {}.", x);
    }

    // let's call some functions
    let n = 10;
    println!("inc({}) = {}", n, inc(n));
    println!("inc(42) = {}", inc(42));
    println!("inc(43) = {}", inc(43));

    // tuples
    let ns = (3, 4);
    println!("x, y = {}, {}", ns.0, ns.1);
    let (x, y) = ns;
    println!("x, y = {}, {}", x, y);

    if ns == (4, 5) {
        println!("All the rules are different!");
    } else {
        println!("The universe is basically ok.");
    }

    // structs
    let p = Point{x: 1, y: 2};
    println!("Meet you at ({}, {}).", p.x, p.y);

    let Meters(l) = Meters(3);
    println!("It's {}m until there, too long for me. Bye.", l)
}
