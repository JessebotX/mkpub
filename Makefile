GOEXE = go
GOFMTEXE = gofmt
GOVETEXE = staticcheck

BINSRCDIR = cmd
BINEXE = mkpub

all: build

build:
	$(GOEXE) build ./$(BINSRCDIR)/$(BINEXE)

fmt:
	$(GOFMTEXE) -s -w -l *.go
	$(GOFMTEXE) -s -w -l ./$(BINSRCDIR)/$(BINEXE)/*.go

vet:
	$(GOVETEXE) ./...

clean:
	rm -f $(BINEXE) $(BINEXE).exe
	rm -rf testdata/build
