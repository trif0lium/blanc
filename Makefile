.PHONY: *

init:
	cd init && \
		go build --ldflags '-s -w -extldflags "-lm -lstdc++ -static"' -o /var/lib/blanc/init main.go
