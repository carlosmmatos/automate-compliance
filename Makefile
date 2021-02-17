BINDIR:=bin
BIN:=$(BINDIR)/autocmp

.PHONY: all
all: build

.PHONY: build
build: $(BIN)

$(BIN): $(BINDIR) main.go 
	go build -o $(BIN) .

$(BINDIR):
	mkdir -p $(BINDIR)