#[deriving(Show)]
pub enum List<A> {
    Nil,
    Cons(A, Box<List<A>>)
}

pub fn cons<A>(x: A, xs: List<A>) -> List<A> {
    match xs {
        Nil => Cons(x, box Nil),
        l => Cons(x, box l)
    }
}

pub fn first<A>(xs: List<A>) -> Option<A> {
    match xs {
        Nil => None,
        Cons(x, _) => Some(x)
    }
}

pub fn main() {
	let nil: List<int> = Nil;
	let ns:  List<int> = cons(1i, cons(2, cons(3, Nil)));

	println!("nil = {}, ns = {}", nil, ns);
	println!("first(nil) = {}, first(ns) = {}", first(nil), first(ns));
}
