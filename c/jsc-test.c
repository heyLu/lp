#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>

#include <JavaScriptCore/JavaScript.h>

#define CONSOLE_LOG_BUF_SIZE 1000
char console_log_buf[CONSOLE_LOG_BUF_SIZE];

JSStringRef to_string(JSContextRef ctx, JSValueRef val);

JSValueRef function_console_log(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, const JSValueRef arguments[], JSValueRef* exception) {
	for (int i = 0; i < argumentCount; i++) {
		if (i > 0) {
			fprintf(stdout, " ");
		}

		JSStringRef str = to_string(ctx, arguments[i]);
		JSStringGetUTF8CString(str, console_log_buf, CONSOLE_LOG_BUF_SIZE);
		fprintf(stdout, "%s", console_log_buf);
	}
	fprintf(stdout, "\n");

	return JSValueMakeUndefined(ctx);
}

JSValueRef function_console_error(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, const JSValueRef arguments[], JSValueRef* exception) {
	for (int i = 0; i < argumentCount; i++) {
		if (i > 0) {
			fprintf(stderr, " ");
		}

		JSStringRef str = to_string(ctx, arguments[i]);
		JSStringGetUTF8CString(str, console_log_buf, CONSOLE_LOG_BUF_SIZE);
		fprintf(stderr, "%s", console_log_buf);
	}
	fprintf(stderr, "\n");

	return JSValueMakeUndefined(ctx);
}

JSValueRef function_import_script(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, const JSValueRef arguments[], JSValueRef* exception) {
	if (argumentCount == 1 && JSValueGetType(ctx, arguments[0]) == kJSTypeString) {
		JSStringRef path_str_ref = JSValueToStringCopy(ctx, arguments[0], NULL);
		char path[100];
		path[0] = '\0';
		JSStringGetUTF8CString(path_str_ref, path, 100);

		FILE *f = fopen(path, "r");
		if (f == NULL) {
			perror("fopen");
			goto err;
		}

		struct stat f_stat;
		if (fstat(fileno(f), &f_stat) < 0) {
			perror("fstat");
			goto err;
		}

		char *buf = malloc(f_stat.st_size * sizeof(char));
		fread(buf, sizeof(char), f_stat.st_size, f);
		if (ferror(f)) {
			perror("fread");
			free(buf);
			goto err;
		}

		JSStringRef script_ref = JSStringCreateWithUTF8CString(buf);
		free(buf);

		JSValueRef ex = NULL;
		JSEvaluateScript(ctx, script_ref, NULL, path_str_ref, 0, &ex);
		JSStringRelease(script_ref);

#ifdef DEBUG
		if (ex != NULL) {
			JSStringRef ex_str = to_string(ctx, ex);
			char ex_buf[1000];
			ex_buf[0] = '\0';
			JSStringGetUTF8CString(ex_str, ex_buf, 1000);
			printf("import: %s\n", ex_buf);
			JSStringRelease(ex_str);
		}
#endif
	}

	return JSValueMakeUndefined(ctx);

err:
	// TODO: Fill exception with error from errno
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

	register_global_function(ctx, "IMPORT_SCRIPT", function_import_script);

	register_global_function(ctx, "CONSOLE_LOG", function_console_log);
	register_global_function(ctx, "CONSOLE_ERROR", function_console_error);

	JSStringRef init_source = JSStringCreateWithUTF8CString("<init>");
	JSStringRef init_script = JSStringCreateWithUTF8CString("var console = {};"\
			"console.log = CONSOLE_LOG;"\
			"console.error = CONSOLE_ERROR;");
	JSEvaluateScript(ctx, init_script, NULL, init_source, 0, NULL);

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
