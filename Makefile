# Copyright (c) Mainflux
# SPDX-License-Identifier: Apache-2.0

PROGRAM = mproxy
SOURCES = $(wildcard *.go) cmd/main.go

all: $(PROGRAM)

.PHONY: all clean $(PROGRAM)

$(PROGRAM): $(SOURCES)
	go build -mod=vendor -ldflags "-s -w" -o ./build/mproxy cmd/main.go

clean:
	rm -rf $(PROGRAM)
