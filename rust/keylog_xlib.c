#include <time.h>
#include <ctype.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <X11/Xlib.h>
#include <X11/Xutil.h>
#define IsTTYFunctionKey(keysym) (((keysym >= 0xff08) && (keysym <= 0xff1b)) || keysym == 0xffff)

/* internal state */
Window focused_window = None;

char *rfc3339_format_time(struct tm *time) {
	char *rfc3339 = malloc(sizeof(char) * 100);
	strftime(rfc3339, 100, "%Y-%m-%d %H:%M:%S", time);
	int tzd_sec = abs(timezone) / 60 / 60;
	int tzd_mon = abs(timezone) % 60;
	sprintf(rfc3339 + strlen(rfc3339), "%+03d:%02d", tzd_sec, tzd_mon);
	return rfc3339;
}

char *rfc3339_now_str() {
	time_t now_t = time(NULL);
	struct tm *now = gmtime(&now_t);
	return rfc3339_format_time(now);
}

Window wparent(Display *display, Window window) {
	Window parent;
	Window root;
	Window *children;
	unsigned int nchildren;
	XQueryTree(display, window, &root, &parent, &children, &nchildren);
	return parent;
}

XTextProperty *get_window_name(Display *display, Window window) {
	XTextProperty *wm_name_prop = malloc(sizeof(XTextProperty));
	Atom net_wm_name = XInternAtom(display, "_NET_WM_NAME", False);
	XGetTextProperty(display, window, wm_name_prop, net_wm_name);
	if (wm_name_prop->value == NULL) {
		XGetTextProperty(display, wparent(display, window), wm_name_prop, net_wm_name);
	}
	return wm_name_prop;
}

void print_next_event(Display *display) {
	XEvent ev;
	KeyCode kc= -1;
	KeySym ksym;
	char *ksymname = NULL;
	char *kname = malloc(sizeof(char) * 2);
	
	//printf("hi there, nice to see you.\n");
	
	if (focused_window == None) {
		//printf("no focus, let's look for it.\n");
		int r = 0;
		XGetInputFocus(display, &focused_window, &r);
		//printf("got focus\n");
		if (focused_window != None) {
			XSelectInput(display, focused_window, KeyPressMask | FocusChangeMask | PropertyChangeMask);
			//printf("selected input");
		} else {
			
		}
		return;
	}

	XNextEvent(display, &ev);
	if(ev.xany.type == FocusOut) {
		focused_window = None;
	} else if (ev.xany.type == PropertyNotify) {
		Atom _NET_WM_NAME = XInternAtom(display, "_NET_WM_NAME", False);
		if (ev.xproperty.atom == _NET_WM_NAME) {
			XTextProperty *wm_name_prop = get_window_name(display, focused_window);
			printf("focus '%s' %s\n", wm_name_prop->value, rfc3339_now_str());
			free(wm_name_prop);
		}
	} else if (ev.xany.type == KeyPress) {
		/* XLookupString  handle keyboard input events in Latin-1 */
		XLookupString(&ev.xkey, kname, 2, &ksym, 0);

		/* Find out string representation */
		if(ksym == NoSymbol) {
			ksymname = "NoSymbol";
		} else {
			if (!(ksymname = XKeysymToString (ksym))) {
				ksymname = "(no name)";
			}
			kc = XKeysymToKeycode(display, ksym);
		}
	}
	
	if (ksymname != NULL) {
		const int non_printable = IsKeypadKey(ksym) || IsTTYFunctionKey(ksym) || IsFunctionKey(ksym) || IsCursorKey(ksym) || IsModifierKey(ksym);
		printf("key ");
		if (non_printable) {
			printf("%s", ksymname);
		} else {
			printf("'%s'", kname);
		}
		printf(" %s\n", rfc3339_now_str());
		fflush(stdout);
	}
}
