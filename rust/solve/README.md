# solve

**Warning**: *Don't use this, it's basically my "Hello World" in Rust
and looks and smells accordingly.  If you Know Things™ and want to tell
me how to do stuff better please do so, but don't expect to find pretty
things yet.  Maybe at some point in the future...*

That being said, this is me playing around with Rust and writing some
code related to constraint programming at the same time.

There are a few things here:

- an incomplete parser for CNFs in DIMAC format ([./src/cnf.rs](./src/cnf.rs))
- a naïve implementation of the DPLL algorithm
- a binary that uses these things to parse CNFs on stdin and solves them
    using the DPLL implementation

    you can try it yourself:

        $ cat examples/external-02-quinn.cnf | cargo run --release
