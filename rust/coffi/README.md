# coffi - Some FFI experiments in Rust

Playing with the FFI in Rust.  Surprisingly easy, but there are likely
dragons elsewhere.

## Quickstart

Display some information about a PNG file, most notably it's width and
height.

```
cargo run <path-to-png>
```

Should build with Rust 1.6.0-nightly, which is necessary because it
uses the `feature(libc)` flag.

If you use emacs, `flycheck-mode` can check the code on the *fly*, which
is awesome!

## Why "coffi"?

Not sure, only the "ffi" part makes sense now.  Ah, no, the "co" comes
from `cos`, which was the first function we wanted to wrap.  That turned
out so easy that we tried `libpng` next.
