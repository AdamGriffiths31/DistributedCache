build:
	go build -o bin/distributedcache

run: build
	./bin/distributedcache