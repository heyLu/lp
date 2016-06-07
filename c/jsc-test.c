#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <JavaScriptCore/JavaScript.h>

#define CONSOLE_LOG_BUF_SIZE 1000
char console_log_buf[CONSOLE_LOG_BUF_SIZE];

JSStringRef to_string(JSContextRef ctx, JSValueRef val);

JSValueRef console_log(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, const JSValueRef arguments[], JSValueRef* exception) {
	for (int i = 0; i < argumentCount; i++) {
		if (i > 0) {
			putchar(' ');
		}

		JSStringRef str = to_string(ctx, arguments[i]);
		JSStringGetUTF8CString(str, console_log_buf, CONSOLE_LOG_BUF_SIZE);
		printf("%s", console_log_buf);
	}
	putchar('\n');

	return JSValueMakeUndefined(ctx);
}

void register_global_function(JSContextRef ctx, char *name, JSObjectCallAsFunctionCallback handler) {
	JSObjectRef global_obj = JSContextGetGlobalObject(ctx);

	JSStringRef fn_name = JSStringCreateWithUTF8CString(name);
	JSObjectRef fn_obj = JSObjectMakeFunctionWithCallback(ctx, fn_name, handler);

	JSObjectSetProperty(ctx, global_obj, fn_name, fn_obj, kJSPropertyAttributeNone, NULL);
}

int main(int argc, char **argv) {
	JSGlobalContextRef ctx = JSGlobalContextCreate(NULL);

	JSStringRef nameRef = JSStringCreateWithUTF8CString("jsc-test");
	JSGlobalContextSetName(ctx, nameRef);

	JSObjectRef global_obj = JSContextGetGlobalObject(ctx);

	register_global_function(ctx, "CONSOLE_LOG", console_log);

	JSStringRef source = JSStringCreateWithUTF8CString("<inline>");
	JSStringRef script;
	if (argc == 0) {
		script = JSStringCreateWithUTF8CString("CONSOLE_LOG(\"Hello, World!\");");
	} else {
		script = JSStringCreateWithUTF8CString(argv[1]);
	}
	JSValueRef res = JSEvaluateScript(ctx, script, global_obj, source, 0, NULL);

	char res_buf[1000];
	res_buf[0] = '\0';
	JSStringRef res_str = to_string(ctx, res);
	JSStringGetUTF8CString(res_str, res_buf, 1000);
	printf("%s\n", res_buf);
}

JSStringRef to_string(JSContextRef ctx, JSValueRef val) {
	if (JSValueIsUndefined(ctx, val)) {
		return JSStringCreateWithUTF8CString("undefined");
	} else if (JSValueIsNull(ctx, val)) {
		return JSStringCreateWithUTF8CString("null");
	} else {
		JSStringRef to_string_name = JSStringCreateWithUTF8CString("toString");
		JSObjectRef obj = JSValueToObject(ctx, val, NULL);
		JSValueRef to_string = JSObjectGetProperty(ctx, obj, to_string_name, NULL);
		JSObjectRef to_string_obj = JSValueToObject(ctx, to_string, NULL);
		JSValueRef obj_val = JSObjectCallAsFunction(ctx, to_string_obj, obj, 0, NULL, NULL);

		return JSValueToStringCopy(ctx, obj_val, NULL);
	}
}
