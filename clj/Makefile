all: paths

paths:
	mkdir -p data
	curl --progress-bar http://download.bbbike.org/osm/bbbike/Leipzig/Leipzig.osm.gz \
		--output data/leipzig.osm.gz
	gzip --decompress --stdout data/leipzig.osm.gz > data/leipzig.osm
	osmfilter --keep='highway=*' data/leipzig.osm > data/leipzig_highways.osm
