var geometry = {};
geometry.cube = function() {
	return [
		// front
		0, 0, 0,
		1, 0, 0,
		0, 1, 0,
		0, 1, 0,
		1, 0, 0,
		1, 1, 0,

		// back
		1, 0, 1,
		0, 0, 1,
		1, 1, 1,
		1, 1, 1,
		0, 0, 1,
		0, 1, 1,

		// left
		0, 0, 1,
		0, 0, 0,
		0, 1, 1,
		0, 1, 1,
		0, 0, 0,
		0, 1, 0,

		// right
		1, 0, 0,
		1, 0, 1,
		1, 1, 0,
		1, 1, 0,
		1, 0, 1,
		1, 1, 1,

		// top
		0, 1, 0,
		1, 1, 0,
		0, 1, 1,
		0, 1, 1,
		1, 1, 0,
		1, 1, 1,

		// bottom
		0, 0, 1,
		1, 0, 1,
		0, 0, 0,
		0, 0, 0,
		1, 0, 1,
		1, 0, 0
	];
};

geometry.cube.normals = function() {
	var front = [0.0, 0.0, 1.0],
	    back = [0.0, 0.0, -1.0],
	    left = [-1.0, 0.0, 0.0],
	    right = [1.0, 0.0, 0.0],
	    top = [0.0, 1.0, 0.0],
	    bottom = [0.0, -1.0, 0.0];

	return [].concat(
		front, front, front, front, front, front,
		back, back, back, back, back, back,
		left, left, left, left, left, left,
		right, right, right, right, right, right,
		top, top, top, top, top, top,
		bottom, bottom, bottom, bottom, bottom, bottom);
}
