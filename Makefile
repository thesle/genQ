APP_NAME := genQ
APP_ID := com.genq.GenQ
BUILD_DIR := build

.PHONY: run build-linux build-windows build-macos build-all flatpak clean

run:
	go run .

build-linux:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 .

build-windows:
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe .

build-macos-amd64:
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 .

build-macos-arm64:
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 .

build-macos: build-macos-amd64 build-macos-arm64

build-all: build-linux build-windows build-macos

flatpak: build-linux
	flatpak-builder --force-clean --user --install $(BUILD_DIR)/flatpak flatpak/$(APP_ID).yml

clean:
	rm -rf $(BUILD_DIR)
