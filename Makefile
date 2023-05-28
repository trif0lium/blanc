.PHONY: *

init:
	cd init && \
		go build --ldflags '-s -w -extldflags "-lm -lstdc++ -static"' -o bin/init main.go
