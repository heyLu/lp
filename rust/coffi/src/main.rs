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
    message: [libc::c_char; 64],
}

impl png_image {
    fn new() -> png_image {
        let mut img: png_image = unsafe { std::mem::zeroed() };
        img.version = 1;
        return img
    }

    fn begin_read_from_file(&mut self, file_name: *const libc::c_char) -> u32 {
        unsafe { png_image_begin_read_from_file(self, file_name) as u32 }
    }
}

impl std::fmt::Display for png_image {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        fn get_message(msg: [libc::c_char; 64]) -> String {
            let mut vec = Vec::new();
            for i in 0..64 {
                vec.push(msg[i] as u8);
            }
            String::from_utf8(vec).unwrap()
        }

        write!(f, "{}x{} {} {} {} {} {} {}", self.width, self.height, self.version,
               self.format, self.flags, self.colormap_entries, self.warning_or_error,
               get_message(self.message))
    }
}

#[link(name = "png")]
extern {
    fn png_image_begin_read_from_file(img: *mut png_image, file_name: *const libc::c_char) -> libc::c_int;
}

fn main() {
    let x = unsafe { cos(3.1415) };
    println!("cos(3.1415) = {}", x);
    println!("");

    let mut img = png_image::new();
    let file_name = std::env::args().nth(1).unwrap_or(String::from("mei.png"));
    let c_name = std::ffi::CString::new(file_name).unwrap();
    let res = img.begin_read_from_file(c_name.as_ptr());
    println!("read_from_file: {}", res);
    println!("{}", img);
}
