run:
	node server.js &> .server.log

fetch_live:
	curl pixl.papill0n.org:8001/world > data/live.json

fetch_deps:
	curl -o public/glmatrix.js https://raw.github.com/toji/gl-matrix/master/dist/gl-matrix.js
	curl -o public/three.js https://rawgithub.com/mrdoob/three.js/master/build/three.min.js
