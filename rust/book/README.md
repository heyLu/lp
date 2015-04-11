# A journey through the Rust book

[This book.](http://doc.rust-lang.org/1.0.0-beta/book/)

## Compiling statically linked binaries

The following assumes you have a `main.rs` and want to get a statically
linked version of it in `main`.  I'm not sure how this would work with
multiple files or `cargo`.

```
# collect `.rlib`s that need to be linked in later
$ rlibs=`rustc -Z print-link-args main.rs | sed -e 's/ /\n/g' | sed -e 's/"//g' | grep rlib`
$ echo rlibs
/home/lu/.local/lib/rustlib/x86_64-unknown-linux-gnu/lib/libstd-4e7c5e5c.rlib
/home/lu/.local/lib/rustlib/x86_64-unknown-linux-gnu/lib/libcollections-4e7c5e5c.rlib
/home/lu/.local/lib/rustlib/x86_64-unknown-linux-gnu/lib/libunicode-4e7c5e5c.rlib
/home/lu/.local/lib/rustlib/x86_64-unknown-linux-gnu/lib/librand-4e7c5e5c.rlib
/home/lu/.local/lib/rustlib/x86_64-unknown-linux-gnu/lib/liballoc-4e7c5e5c.rlib
/home/lu/.local/lib/rustlib/x86_64-unknown-linux-gnu/lib/liblibc-4e7c5e5c.rlib
/home/lu/.local/lib/rustlib/x86_64-unknown-linux-gnu/lib/libcore-4e7c5e5c.rlib

# build `main.o`, i.e. compile `main.rs`, but do not link
$ rustc --emit obj main.rs

# link it statically (you need `librt.a` on your system, or in `./deps`)
#   (also, the warnings seem to be expected)
$ gcc -static -static-libgcc -o main main.o `tr $'\n' ' '<<<$rlibs` -lpthread -lm -ldl -L./deps -lrt
/home/lu/.local/lib/rustlib/x86_64-unknown-linux-gnu/lib/libstd-4e7c5e5c.rlib(std-4e7c5e5c.o): In function `dynamic_lib::DynamicLibrary::open::h232b8007ebe62ccakTe':
std.0.rs:(.text._ZN11dynamic_lib14DynamicLibrary4open20h232b8007ebe62ccakTeE+0x103): warning: Using 'dlopen' in statically linked applications requires at runtime the shared libraries from the glibc version used for linking
/home/lu/.local/lib/rustlib/x86_64-unknown-linux-gnu/lib/libstd-4e7c5e5c.rlib(std-4e7c5e5c.o): In function `env::home_dir::h89e202bbf866699a0cf':
std.0.rs:(.text._ZN3env8home_dir20h89e202bbf866699a0cfE+0x14d): warning: Using 'getpwuid_r' in statically linked applications requires at runtime the shared libraries from the glibc version used for linking
/home/lu/.local/lib/rustlib/x86_64-unknown-linux-gnu/lib/libstd-4e7c5e5c.rlib(std-4e7c5e5c.o): In function `net::lookup_host::h6fcc17ec54074bdeCik':
std.0.rs:(.text._ZN3net11lookup_host20h6fcc17ec54074bdeCikE+0x1ce): warning: Using 'getaddrinfo' in statically linked applications requires at runtime the shared libraries from the glibc version used for linking

# and now you can run it!
$ file main
main: ELF 64-bit LSB executable, x86-64, version 1 (GNU/Linux), statically linked, for GNU/Linux 2.6.32, BuildID[sha1]=02abd2412f0570e81b4ffd8184edc883fb264f92, not stripped
$ ./main
Hello, World!

# it also works on other systems
$ docker run -it --rm -v $PWD/main:/usr/local/bin/hello busybox /usr/local/bin/hello
Hello, World!
```

Thanks to Kai Noda for [figuring this out](https://mail.mozilla.org/pipermail/rust-dev/2014-November/011365.html).
