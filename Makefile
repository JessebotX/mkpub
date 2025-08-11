GOEXE = go
GOFMTEXE = gofmt
GOVETEXE = staticcheck

BINSRCDIR = cmd
BINEXE = mkpub

all: build

build:
	$(GOEXE) mod download
	$(GOEXE) build ./$(BINSRCDIR)/$(BINEXE)

fmt:
	$(GOFMTEXE) -s -w -l *.go
	$(GOFMTEXE) -s -w -l ./$(BINSRCDIR)/$(BINEXE)/*.go

vet:
	$(GOVETEXE) ./...

clean: clean-bin clean-test

clean-bin:
	rm -f $(BINEXE) $(BINEXE).exe

clean-test:
	rm -rf testdata/build

install:
	$(GOEXE) install ./cmd/mkpub
