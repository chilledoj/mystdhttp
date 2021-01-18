.PHONY: build

build:
	go build -o ./bin/mystdhttp ./cmd/.

.PHONY: run

run: build
	./bin/mystdhttp
