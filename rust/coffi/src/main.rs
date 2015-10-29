#[link(name = "m")]
extern {
    fn cos(d: f64) -> f64;
}

fn main() {
    let x = unsafe { cos(3.1415) };
    println!("cos(3.1415) = {}", x);
}
