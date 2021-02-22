BINDIR:=bin
BIN:=$(BINDIR)/autocmp

SRC=$(shell find . -name "*.go")

.PHONY: all
all: build

.PHONY: build
build: $(BIN)

$(BIN): $(BINDIR) $(SRC)
	go build -o $(BIN) .

$(BINDIR):
	mkdir -p $(BINDIR)

.PHONY: run
run: build
	./$(BIN)

.PHONY: test-unit
test-unit:
	go test -v -race ./...