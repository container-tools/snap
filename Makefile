
build: vfs
	go build -o snap ./cmd/

vfs:
	GO111MODULE=on go run github.com/rakyll/statik -src=deploy -f

.PHONY: build vfs