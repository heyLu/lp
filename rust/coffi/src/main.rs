#![feature(libc)]
extern crate libc;

#[link(name = "m")]
extern {
    fn cos(d: f64) -> f64;
}

#[repr(C)]
struct PNGImage {
    opaque: *mut libc::c_void,
    version: libc::c_uint,
    width: libc::c_uint,
    height: libc::c_uint,
    format: libc::c_uint,
    flags: libc::c_uint,
    colormap_entries: libc::c_uint,
    warning_or_error:  libc::c_uint,
    message: [libc::c_char; 64],
}

impl PNGImage {
    fn new() -> PNGImage {
        let mut img: PNGImage = unsafe { std::mem::zeroed() };
        img.version = 1;
        return img
    }

    fn begin_read_from_file(&mut self, file_name: *const libc::c_char) -> u32 {
        unsafe { png_image_begin_read_from_file(self, file_name) as u32 }
    }
}

impl PNGImage {
    fn message(&self) -> String {
        String::from_utf8(self.message.iter().map(|&c| c as u8).collect()).unwrap()
    }
}

impl std::fmt::Display for PNGImage {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(f, "{}x{} {} {} {} {} {} {}", self.width, self.height, self.version,
               self.format, self.flags, self.colormap_entries, self.warning_or_error,
               self.message())
        }
}

#[link(name = "png")]
extern {
    fn png_image_begin_read_from_file(img: *mut PNGImage, file_name: *const libc::c_char) -> libc::c_int;
}

fn main() {
    let x = unsafe { cos(3.1415) };
    println!("cos(3.1415) = {}", x);
    println!("");

    let mut img = PNGImage::new();
    let file_name = std::env::args().nth(1).unwrap_or(String::from("mei.png"));
    let c_name = std::ffi::CString::new(file_name).unwrap();
    let res = img.begin_read_from_file(c_name.as_ptr());
    println!("read_from_file: {}", res);
    println!("{}", img);
}
