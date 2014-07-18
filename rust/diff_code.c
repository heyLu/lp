#include <stdio.h>
#include <stdlib.h>

int* diff_encode(int* values, int length) {
	if (length < 1) {
		return NULL;
	}

	int* result = malloc(sizeof(int) * length);
	int last = result[0] = values[0];
	for (int i = 1; i < length; i++) {
		result[i] = values[i] - last;
		last = values[i];
	}

	return result;
}

int* diff_decode(int* encoded, int length) {
	if (length < 1) {
		return NULL;
	}

	int* result = malloc(sizeof(int) * length);
	int last = result[0] = encoded[0];
	for (int i = 1; i < length; i++) {
		result[i] = last + encoded[i];
		last = result[i];
	}

	return result;
}

int main() {
	int years[] = {1913, 2020, 1931, 1947, 1978, 1970, 2001, 2023, 1801, 2807};
	int length = 10;

	int* encoded = diff_encode(years, length);
	int* decoded = diff_decode(encoded, length);

	for (int i = 0; i < length; i++) {
		printf("%d -> %d -> %d\n", years[i], encoded[i], decoded[i]);
	}

	return 0;
}
