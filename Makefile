all: build

build:
	go build cmd/rentals-cli/rentals-cli.go
	mv rentals-cli ./scripts

clean:
	rm ./scripts/rentals-cli
