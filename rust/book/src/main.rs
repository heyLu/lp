fn inc(x: i32) -> i32 {
    if x == 42 {
        return 42
    }

    x + 1
}

fn main() {
    let x = 5;

    println!("Hello, World!");
    if x == 42 {
        println!("It's THE ANSWER!");
    } else {
        println!("Meh, it's just a number, {}.", x);
    }

    let n = 10;
    println!("inc({}) = {}", n, inc(n));
    println!("inc(42) = {}", inc(42));
    println!("inc(43) = {}", inc(43));
}
