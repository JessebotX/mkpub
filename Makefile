GOEXE = go
GOFMTEXE = gofmt
GOFMTFLAGS = -s -w -l
GOFMTINPUT = *.go cmd/mkpub/*.go config/*.go
GOLINTEXE = staticcheck
GOLINTFLAGS =
GOLINTINPUT = ./...

BINSRCDIR = ./cmd
BINEXE = mkpub

all: build

build:
	$(GOEXE) mod download
	$(GOEXE) build $(BINSRCDIR)/$(BINEXE)

fmt: format
format:
	$(GOFMTEXE) $(GOFMTFLAGS) $(GOFMTINPUT)

vet: lint
check: lint
lint:
	$(GOLINTEXE) $(GOLINTFLAGS) $(GOLINTINPUT)

clean: clean-bin clean-test

clean-bin:
	rm -f $(BINEXE) $(BINEXE).exe

clean-test:
	rm -rf testdata/build

install:
	$(GOEXE) install ./cmd/mkpub
