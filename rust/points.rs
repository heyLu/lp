// adopted from the lifetime guide: http://doc.rust-lang.org/guide-lifetimes.html

#[deriving(Show)]
struct Point {
    x: f64,
    y: f64
}

fn main() {
    let on_the_stack : Point      = Point{x: 3.4, y: 5.7};
    let on_the_heap  : Box<Point> = box Point{x: 9.6, y: 2.1};

    let dist = compute_distance(&on_the_stack, &*on_the_heap);

    // mismatched types: &Point vs Box<Point>
    // e.g., there is no automatic "promoting" of boxed to reference pointers
    //let dist2 = compute_distance(&on_the_stack, on_the_heap);
    
    println!("dist({}, {}) = {}", on_the_stack, on_the_heap, dist)
}

fn compute_distance(p1: &Point, p2: &Point) -> f64 {
    let dx = p1.x - p2.x;
    let dy = p1.y - p2.y;
    (dx*dx + dy*dy).sqrt()
}
