extern mod xlib;

use XlibT = xlib::xlib;
use Xlib = xlib::xlib::bindgen;

extern mod xlib_helpers;
use xh = xlib_helpers::xlib_helpers;

extern mod keylog_xlib {
	fn print_next_event(display: *XlibT::Display);
}

fn main() {
	let display = Xlib::XOpenDisplay(ptr::null());
	
	loop { keylog_xlib::print_next_event(display); }

	/*do xevents(KeyPressMask | FocusChangeMask) |&ev| {
		match ev.type {
			KeyPress => (),
			_ => ()
		}
	}*/
}
