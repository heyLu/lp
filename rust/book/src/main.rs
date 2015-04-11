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
}
