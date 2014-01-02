var matrix = {};
matrix.multiply = function(a, dim_a, b, dim_b) {
	if (dim_a[1] != dim_b[0]) {
		throw new Error("cols a != rows b");
	} else {
		var c = [];
		var s = dim_a[1];
		var n = dim_a[0],
			 m = dim_b[1];
		for (var i = 0; i < n; i++) {
			for (var j = 0; j < m; j++) {
				var sum = 0;
				for (var k = 0; k < s; k++) {
					sum += a[i * s + k] * b[k * m + j];
				}
				c[i * m + j] = sum;
			}
		}
		return c;
	}
};

matrix.multiplyMany4x4 = function() {
	var m = arguments[0];
	for(var i = 1; i < arguments.length; i++) {
		m = matrix.multiply(m, [4, 4], arguments[i], [4, 4]);
	}
	return m;
}

matrix.identity = function(dim) {
	var m = [];
	for(var i = 0; i < dim; i++) {
		for(var j = 0; j < dim; j++) {
			m[i * dim + j] = i == j ? 1.0 : 0.0;
		}
	}
	return m;
}

matrix.transpose = function(m, dim) {
	var t = [];
	for(var i = 0; i < dim[0]; i++) {
		for(var j = 0; j < dim[1]; j++) {
			t[j * dim[0] + i] = m[i * dim[1] + j];
		}
	}
	return t;
}

matrix.translate = function(tx, ty, tz) {
	var tx = tx || 0.0, ty = ty || 0.0, tz = tz || 0.0;
	return [
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		tx, ty, tz, 1
	];
}

matrix.scale = function(sx, sy, sz) {
	var sx = sx || 1.0, sy = sy || 1.0, sz = sz || 1.0;
	return [
		sx, 0,  0,  0,
		0,  sy, 0,  0,
		0,  0,  sz, 0,
		0,  0,  0,  1
	];
}

matrix.rotateX = function(angle) {
	var s = Math.sin(angle);
	var c = Math.cos(angle);
	return [
		1,  0, 0, 0,
		0,  c, s, 0,
		0, -s, c, 0,
		0,  0, 0, 1
	];
}

matrix.rotateY = function(angle) {
	var s = Math.sin(angle);
	var c = Math.cos(angle);
	return [
		c, 0, -s, 0,
		0, 1,  0, 0,
		s, 0,  c, 0,
		0, 0,  0, 1
	];
}

matrix.rotateZ = function(angle) {
	var s = Math.sin(angle);
	var c = Math.cos(angle);
	return [
		 c, s, 0, 0,
		-s, c, 0, 0,
		 0, 0, 1, 0,
		 0, 0, 0, 1
	];
}

matrix.perspective = function(fov, aspect) {
	return [
		Math.PI / 3 / aspect, 0, 0, 0,
		0, Math.PI / 3, 0, 0,
		0, 0, 1, 1,
		0, 0, 0, 1
	];
}
