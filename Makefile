build:
	go build -o bin/distributedcache

run: build
	./bin/distributedcache

runfollower: build
	./bin/distributedcache --listenaddr :4000 --leaderaddr :8080

test:
	go test -v ./...