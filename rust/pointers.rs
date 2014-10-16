extern crate debug;

mod wrapped {
    fn get_wrapped(v: &int) -> int {
        *v
    }

    #[deriving(Show)]
    struct Point {x: f64, y: f64}

    fn get_wrapped_struct(p: &Point) -> Point {
        *p
    }

    pub fn run() {
        println!("\nwrapped:");

        // mhh... i thought this would "move" the ownership, but somehow it doesn't. because it's immutable?
        // by the way, how is immutability defined? we can have mutable variables, but also mutable pointers,
        // what does that mean?
        println!("get_wrapped(&1) = {}", get_wrapped(&1));

        // even this does work, so what doesn't?
        let point = Point{x: 3.1, y: 4.7};
        println!("get_wrapped_struct(&Point{{x = 3.1, y = 4.7}}) = {}", get_wrapped_struct(&point));
    }
}

mod list_with_box {
    #[deriving(Clone, Show)]
    pub enum List<A> {
        Nil,
        Cons(A, Box<List<A>>)
    }

    // this could already be "ok", if we accept that the argument will be copied just to match on it. i think rust uses `memcpy` for that, so it should be quite fast, even if generally "bad".
    pub fn cons<A>(x: A, xs: List<A>) -> List<A> {
        match xs {
            Nil => Cons(x, box Nil),
            _   => Cons(x, box xs)
        }
    }

    // pub fn cons_with_ref<A>(x: A, xs: &List<A>) -> List<A> {
    //     match *xs {
    //         Nil => Cons(x, box Nil),
    //         _ => Cons(x, box *xs)
    //                       // ^    "cannot move out of dereference of `&`-pointer"
    //                       // same with `ref l => Cons(x, box *l)`
    //     }
    // }

    pub fn cons_with_box<A>(x: A, xs: Box<List<A>>) -> List<A> {
        match *xs {
            Nil => Cons(x, xs),
            _   => Cons(x, xs)
        }
    }

    // pub fn cons_with_referenced_box<A>(x: A, xs: &Box<List<A>>) -> List<A> {
    //     match **xs {
    //         Nil     => Cons(x, box Nil),
    //         ref l   => Cons(x, box *l) // "cannot move ...", again
    //     }
    // }

    pub fn run() {
        println!("\nlist_with_box:");

        let l: List<int> = cons(1, cons(2, Nil));
        // must be cloned before use as parameter to cons, which moves it
        let l2 = l.clone();

        println!("l = {}", l);
        println!("cons(3, l)  = {:?}", cons(3, l));
        println!("cons(4, l2) = {:?}", cons(4, l2));
        println!("");

        // println!("cons_with_ref(3, &l) = {:?}", cons_with_ref(3, &l));
        //                                                        // ^    "use of moved value: `l`"

        let lb: Box<List<int>> = box cons_with_box(1, box cons_with_box(2, box Nil));
        println!("lb = {:?}", lb);
        println!("cons_with_box(3, lb) = {:?}", cons_with_box(3, lb));
        // println!("cons_with_box(4, lb) = {:?}", cons_with_box(4, lb));
        //                                                       // ^     "use of moved value: `l`"
    }
}

mod list_with_ref {
    #[deriving(Show)]
    pub enum List<'a, A: 'a> {
        Nil,
        Cons(A, &'a List<'a A>) // are there alternatives to boxing? (i think not -- Rc, Arc, etc. *are* boxes. right?)
    }

    // we want to pass a reference to the list, pattern match on it and then return a new list that references the old one.
    pub fn cons<'a, A>(x: A, xs: &'a List<A>) -> List<'a, A> {
        match *xs {
            Nil => Cons(x, xs),
            ref xs  => Cons(x, xs)
        }
    }

    pub fn first<'a, A>(xs: &'a List<A>) -> Option<&'a A> {
        match *xs {
            Nil => None,
            // note the `ref x` here, the default would again be moving the value
            Cons(ref x, _) => Some(x)
        }
    }

    pub fn last<'a, A>(xs: &'a List<A>) -> Option<&'a A> {
        match *xs {
            Nil => None,
            Cons(ref x, &Nil) => Some(x),
            Cons(_, xs) => last(xs)
        }
    }

    // pub fn from_iterator<'a, I: Iterator<A>, A: 'a>(mut i: I) -> &'a List<'a A> {
    //     i.fold(&Nil, |xs, x| &cons(x, xs))
    // }

    pub fn run() {
        println!("\nlist_with_ref:");

        //let l: &List<int> = &Cons(1, &Cons(2, &Cons(3, &Nil)));
        let nil = &Nil;
        let l: &List<int> = &cons(1, nil);

        println!("l = {:?}", l);
        println!("cons(1, l) = {:?}", cons(1, l));

        let l2: &List<int> = &cons(2, l);
        println!("l2 = {:?}", l2);
        println!("first(l2) = {:?}", first(l2));
        println!("last(l2)  = {:?}", last(l2));
    }
}

mod list_with_ref_box {
    #[deriving(Show)]
    enum List<'a, A: 'a> {
        Nil,
        Cons(A, &'a Box<List<'a A>>)
    }

    pub fn cons<'a, A>(x: A, xs: &'a Box<List<'a A>>) -> List<'a A> {
        Cons(x, xs)
    }

    pub fn run() {
        println!("\nlist_with_ref_box:");

        let nil = &box Nil;
        let l: &Box<List<int>> = &box cons(1, nil);

        println!("l = {:?}", l);
        println!("cons(2, l) = {:?}", cons(2, l));
        println!("cons(3, l) = {:?}", cons(3, l));
    }
}

pub fn main() {
    wrapped::run();

    list_with_box::run();
    list_with_ref::run();
    list_with_ref_box::run();
}
