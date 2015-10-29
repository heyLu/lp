#![feature(libc)]
extern crate libc;

#[link(name = "m")]
extern {
    fn cos(d: f64) -> f64;
}

#[allow(non_camel_case_types)]
enum png_opaque {}

#[repr(C)]
struct png_image {
    opaque: *mut png_opaque,
    version: libc::uint32_t,
    width: libc::uint32_t,
    height: libc::uint32_t,
    format: libc::uint32_t,
    flags: libc::uint32_t,
    colormap_entries: libc::uint32_t,
    warning_or_error:  libc::uint32_t,
    message: [u8; 64],
}

#[link(name = "png")]
extern {
    fn png_image_begin_read_from_file(img: *mut png_image, file_name: *const u8) -> libc::c_int;
}

fn main() {
    let x = unsafe { cos(3.1415) };
    println!("cos(3.1415) = {}", x);

    let mut img: png_image;
    unsafe {
        img = std::mem::zeroed();
        let res = png_image_begin_read_from_file(&mut img, "mei.png".as_ptr());
        println!("read_from_file: {}", res);
        println!("{}x{}", img.width, img.height);
    }
}
