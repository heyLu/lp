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

#[derive(Debug)]
enum Number {
    Integer(i32),
    Float(f32),
    NaN
}

fn print_num_type(num: Number) {
    println!("{}", match num {
        Number::Integer(_) => "integer",
        Number::Float(_) => "float",
        Number::NaN => "not a number"
    })
}

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
    println!("It's {}m until there, too long for me. Bye.", l);

    // enums
    let i: Number = Number::Integer(3);
    let f: Number = Number::Float(3.1415);
    let n: Number = Number::NaN;
    println!("Here are some numbers: {:?}, {:?} and {:?}", i, f, n);

    // match
    print_num_type(i);
    print_num_type(f);
    print_num_type(n);

    // loops
    for x in 0..10 {
        println!("{}", x);
    }

    // strings
    let string_slice: &str = "Hello, World!";
    let mut string: String = string_slice.to_string();
    println!("{} = {}", string_slice, string);

    string.push_str(" (again...)");
    println!("{}", string);

    // arrays
    let array: [i32; 3] = [1, 2, 3];

    println!("You have {} numbers:", array.len());
    for n in array.iter() {
        println!("{}", n);
    }

    // vectors
    let mut vec: Vec<i32> = vec![4, 5, 6];
    vec.push(7);

    println!("You have {} more numbers:", vec.len());
    for n in vec.iter() {
        println!("{}", n);
    }

    // slices
    let middle = &vec[1..3];
    println!("Only a few of those:");
    for n in middle.iter() {
        println!("{}", n);
    }
}
