GOCMD=env go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
GOGET=$(GOCMD) get

BINARY=cavern
TESTS=./...
COVERAGE_FILE=coverage.out

.PHONY: all test build build-wasm coverage clean resources

all: test build

build:
		$(GOBUILD) -o $(BINARY) -v

build-wasm:
		GOOS=js GOARCH=wasm $(GOBUILD) -o $(BINARY).wasm -v

test:
		$(GOTEST) -v $(TESTS)

coverage:
		$(GOTEST) -coverprofile=$(COVERAGE_FILE) $(TESTS)
		$(GOTOOL) cover -html=$(COVERAGE_FILE)

clean:
		$(GOCLEAN)
		rm -f $(BINARY) $(COVERAGE_FILE)

resources:
		$(GORUN) github.com/markbates/pkger/cmd/pkger -include /images -include /music -include /sounds
