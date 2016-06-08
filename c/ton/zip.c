// Utilities for getting entries from ZIP files.
// 
// Uses libzip, alternatives are minizip (from zlib) and zziplib.

#include <stdio.h>
#include <stdlib.h>

#include "zip.h"

void print_zip_err(char *prefix, zip_t *zip);

char *get_contents_zip(char *path, char *name) {
	zip_t *archive = zip_open(path, ZIP_RDONLY, NULL);
	if (archive == NULL) {
		print_zip_err("zip_open", archive);
		return NULL;
	}

	zip_stat_t stat;
	if (zip_stat(archive, name, 0, &stat) < 0) {
		print_zip_err("zip_stat", archive);
		goto close_archive;
	}

	zip_file_t *f = zip_fopen(archive, name, 0);
	if (f == NULL) {
		print_zip_err("zip_fopen", archive);
		goto close_archive;
	}

	char *buf = malloc(stat.size + 1);
	if (zip_fread(f, buf, stat.size) < 0) {
		print_zip_err("zip_fread", archive);
		goto free_buf;
	}
	buf[stat.size] = '\0';

	zip_fclose(f);
	zip_close(archive);

	return buf;

free_buf:
	free(buf);
close_file:
	zip_fclose(f);
close_archive:
	zip_close(archive);

	return NULL;
}

void print_zip_err(char *prefix, zip_t *zip) {
		zip_error_t *err = zip_get_error(zip);
		printf("%s: %s\n", prefix, zip_error_strerror(err));
		zip_error_fini(err);
}
