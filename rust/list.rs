#[deriving(Show)]
pub enum List<A> {
    Nil,
    Cons(A, Box<List<A>>)
}

// xs: doesn't work as a reference, probably as Box?
pub fn cons<A>(x: A, xs: List<A>) -> List<A> {
    match xs {
        Nil => Cons(x, box Nil),
        l => Cons(x, box l)
    }
}

pub fn first<A>(xs: &List<A>) -> Option<&A> {
    match xs {
        &Nil => None,
        &Cons(ref x, _) => Some(x)
    }
}

// nth
// map - does this make sense? (e.g. we probably don't want to copy everything.) should it be lazy?
//   fn map(..., &List<&A>) -> List<B> // or -> List<&B>?
// filter - we want references to the old values and keep the existing cons cells
// reduce

pub fn main() {
	let nil: &List<int> = &Nil;
	let ns:  &List<int> = &cons(1i, cons(2, cons(3, Nil)));

	println!("nil = {}, ns = {}", nil, ns);
	println!("first(nil) = {}, first(ns) = {}", first(nil), first(ns));
}
