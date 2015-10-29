#![feature(libc)]
extern crate libc;

#[link(name = "m")]
extern {
    fn cos(d: f64) -> f64;
}

#[repr(C)]
struct png_image {
    opaque: *mut libc::c_void,
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

    fn get_message(msg: [u8; 64]) -> String {
        let mut vec = Vec::new();
        for i in 0..64 {
            vec.push(msg[i]);
        }
        String::from_utf8(vec).unwrap()
    }

    fn print_img(img: &png_image) {
        println!("{}x{} {} {} {} {} {} {}", img.width, img.height, img.version, img.format, img.flags, img.colormap_entries, img.warning_or_error, get_message(img.message))
    }

    unsafe {
        let mut img: png_image = std::mem::zeroed();
        print_img(&img);
        let res = png_image_begin_read_from_file(&mut img, "mei.png\0".as_ptr());
        println!("read_from_file: {}", res);
        print_img(&img);
    }
}
