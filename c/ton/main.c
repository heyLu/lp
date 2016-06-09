#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>

#include <JavaScriptCore/JavaScript.h>

#include "zip.h"

#define CONSOLE_LOG_BUF_SIZE 1000
char console_log_buf[CONSOLE_LOG_BUF_SIZE];

JSStringRef to_string(JSContextRef ctx, JSValueRef val);
JSValueRef evaluate_script(JSContextRef ctx, char *script, char *source);

char *munge(char *s);

void bootstrap(JSContextRef ctx, char *deps_file_path, char *goog_base_path);

char* get_contents(char *path);

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
		JSStringRelease(path_str_ref);

		char full_path[150];
		snprintf(full_path, 150, "%s/%s", "out", path);
		char *buf = get_contents(full_path);
		if (buf == NULL) {
			goto err;
		}

		evaluate_script(ctx, buf, full_path);
		free(buf);
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

	evaluate_script(ctx, "var console = {};"\
			"console.log = CONSOLE_LOG;"\
			"console.error = CONSOLE_ERROR;", "<init>");

	bootstrap(ctx, "out/main.js", "out/goog/base.js");

	char *script;
	if (argc == 0) {
		script = "CONSOLE_LOG(\"Hello, World!\");";
	} else {
		script = argv[1];
	}
	JSValueRef res = evaluate_script(ctx, script, "<inline>");

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

JSValueRef evaluate_script(JSContextRef ctx, char *script, char *source) {
	JSStringRef script_ref = JSStringCreateWithUTF8CString(script);
	JSStringRef source_ref = NULL;
	if (source != NULL) {
		source_ref = JSStringCreateWithUTF8CString(source);
	}

	JSValueRef ex = NULL;
	JSValueRef val = JSEvaluateScript(ctx, script_ref, NULL, source_ref, 0, &ex);
	JSStringRelease(script_ref);
	if (source != NULL) {
		JSStringRelease(source_ref);
	}

#ifdef DEBUG
	if (ex != NULL) {
		JSStringRef ex_str = to_string(ctx, ex);
		char ex_buf[1000];
		ex_buf[0] = '\0';
		JSStringGetUTF8CString(ex_str, ex_buf, 1000);
		printf("eval %s: %s\n", source, ex_buf);
		JSStringRelease(ex_str);
	}
#endif

	return val;
}

char *munge(char *s) {
	int len = strlen(s);
	int new_len = 0;
	for (int i = 0; i < len; i++) {
		switch (s[i]) {
		case '!':
			new_len += 6; // _BANG_
			break;
		case '?':
			new_len += 7; // _QMARK_
			break;
		default:
			new_len += 1;
		}
	}

	char *ms = malloc((new_len+1) * sizeof(char));
	int j = 0;
	for (int i = 0; i < len; i++) {
		switch (s[i]) {
		case '-':
			ms[j++] = '_';
			break;
		case '!':
			ms[j++] = '_';
			ms[j++] = 'B';
			ms[j++] = 'A';
			ms[j++] = 'N';
			ms[j++] = 'G';
			ms[j++] = '_';
			break;
		case '?':
			ms[j++] = '_';
			ms[j++] = 'Q';
			ms[j++] = 'M';
			ms[j++] = 'A';
			ms[j++] = 'R';
			ms[j++] = 'K';
			ms[j++] = '_';
			break;

		default:
			ms[j++] = s[i];
		}
	}
	ms[new_len] = '\0';

	return ms;
}

void bootstrap(JSContextRef ctx, char *deps_file_path, char *goog_base_path) {
	char source[] = "<bootstrap>";

	// Setup CLOSURE_IMPORT_SCRIPT
	evaluate_script(ctx, "CLOSURE_IMPORT_SCRIPT = function(src) { IMPORT_SCRIPT('goog/' + src); return true; }", source);

	// Load goog base
	char *base_script_str = get_contents(goog_base_path);
	if (base_script_str == NULL) {
		fprintf(stderr, "The goog base JavaScript text could not be loaded");
		exit(1);
	}
	evaluate_script(ctx, base_script_str, "<bootstrap:base>");
	free(base_script_str);

	// Load the deps file
	char *deps_script_str = get_contents(deps_file_path);
	if (deps_script_str == NULL) {
		fprintf(stderr, "The goog base JavaScript text could not be loaded");
		exit(1);
	}
	evaluate_script(ctx, deps_script_str, "<bootstrap:deps>");
	free(deps_script_str);

	evaluate_script(ctx, "goog.isProvided_ = function(x) { return false; };", source);

	evaluate_script(ctx, "goog.require = function (name) { return CLOSURE_IMPORT_SCRIPT(goog.dependencies_.nameToPath[name]); };", source);

	evaluate_script(ctx, "goog.require('cljs.core');", source);

	// redef goog.require to track loaded libs
	evaluate_script(ctx, "cljs.core._STAR_loaded_libs_STAR_ = cljs.core.into.call(null, cljs.core.PersistentHashSet.EMPTY, [\"cljs.core\"]);\n"
			"goog.require = function (name, reload) {\n"
			"    if(!cljs.core.contains_QMARK_(cljs.core._STAR_loaded_libs_STAR_, name) || reload) {\n"
			"        var AMBLY_TMP = cljs.core.PersistentHashSet.EMPTY;\n"
			"        if (cljs.core._STAR_loaded_libs_STAR_) {\n"
			"            AMBLY_TMP = cljs.core._STAR_loaded_libs_STAR_;\n"
			"        }\n"
			"        cljs.core._STAR_loaded_libs_STAR_ = cljs.core.into.call(null, AMBLY_TMP, [name]);\n"
			"        CLOSURE_IMPORT_SCRIPT(goog.dependencies_.nameToPath[name]);\n"
			"    }\n"
			"};", source);
}

char *get_contents(char *path) {
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

	char *buf = malloc(f_stat.st_size + 1);
	memset(buf, 0, f_stat.st_size);
	fread(buf, f_stat.st_size, 1, f);
	buf[f_stat.st_size] = '\0';
	if (ferror(f)) {
		perror("fread");
		free(buf);
		goto err;
	}

	if (fclose(f) < 0) {
		perror("fclose");
		goto err;
	}

	return buf;

err:
	return NULL;
}
