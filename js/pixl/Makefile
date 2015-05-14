run:
	node server.js &> .server.log

fetch_live:
	curl pixl.papill0n.org:8001/world > data/live.json

fetch_deps:
	curl -L -o public/glmatrix.js https://raw.github.com/toji/gl-matrix/master/dist/gl-matrix.js
	curl -L -o public/three.js https://raw.github.com/mrdoob/three.js/master/build/three.min.js
	curl -L -o public/three.flycontrols.js https://raw.github.com/mrdoob/three.js/master/examples/js/controls/FlyControls.js
