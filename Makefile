build:
	CGO_ENABLED=0 go build -ldflags "-s -w" -o bin/blog

docker: FORCE
	rm -rf docker/build
	mkdir -p docker/build
	cp bin/blog docker/build
	# cp config.json docker/build
	cp sharpdown docker/build
	cp -r web docker/build
	cd docker && docker buildx build -t cracktc/blog .

FORCE:
