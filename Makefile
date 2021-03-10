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

WASM_DIR=./wasm/

.PHONY: all test build build-wasm coverage clean resources mac-app

all: test build

build:
		$(GOBUILD) -tags prod -o $(BINARY) -v

build-wasm:
		rsync -tv $(shell go env GOROOT)/misc/wasm/wasm_exec.js $(WASM_DIR)
		GOOS=js GOARCH=wasm $(GOBUILD) -tags prod -o $(WASM_DIR)$(BINARY).wasm -v
		gzip --force --keep --best $(WASM_DIR)$(BINARY).wasm

test:
		$(GOTEST) -race -v $(TESTS)

coverage:
		$(GOTEST) -coverprofile=$(COVERAGE_FILE) $(TESTS)
		$(GOTOOL) cover -html=$(COVERAGE_FILE)

clean:
		$(GOCLEAN)
		rm -f $(BINARY) $(COVERAGE_FILE)

mac-app: build
		rm -rf Cavern.app
		mkdir -p "Cavern.app/Contents/"{MacOS,Resources}
		cp Info.plist Cavern.app/Contents/Info.plist
		mkdir icon.iconset

		sips -z 16 16 icon64.png --out icon.iconset/icon_16x16.png
		sips -z 32 32 icon64.png --out icon.iconset/icon_16x16@2x.png
		cp icon.iconset/icon_16x16@2x.png icon.iconset/icon_32x32.png
		cp icon64.png icon.iconset/icon_32x32@2x.png
		sips -z 128 128 icon128.png --out icon.iconset/icon_128x128.png
		sips -z 256 256 icon256.png --out icon.iconset/icon_128x128@2x.png
		cp icon.iconset/icon_128x128@2x.png icon.iconset/icon_256x256.png
		sips -z 512 512 icon.png --out icon.iconset/icon_256x256@2x.png
		cp icon.iconset/icon_256x256@2x.png icon.iconset/icon_512x512.png
		sips -z 1024 1024 icon.png --out icon.iconset/icon_512x512@2x.png

		# Create .icns file
		iconutil -c icns icon.iconset

		cp icon.iconset/icon_512x512@2x.png Cavern.app/Contents/Resources/Cavern.png
		mv icon.icns Cavern.app/Contents/Resources/

		# Cleanup
		rm -R icon.iconset

		# Copy executable
		cp cavern Cavern.app/Contents/MacOS/Cavern
