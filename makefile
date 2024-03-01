BINARY_NAME := ldrgen
SRC_FOLDER := src
BIN_FOLDER := bin

.PHONY: all windows linux clean win64 win32 x64 x86

all:
	@echo "[*] Specify a target: windows, linux"

build: windows linux

windows: win64 win32

linux: x64 x86
	@echo "[*] Done compiling for Linux x64 and x86, binaries are located in: $(BIN_FOLDER)"

win64:
	@echo "[*] Compiling for Windows x64"
	@GOOS=windows GOARCH=amd64 go build -o $(BIN_FOLDER)/win/$(BINARY_NAME)_win64.exe $(SRC_FOLDER)/main.go
	@echo "[*] Done compiling for Windows x64, binary is located in: $(BIN_FOLDER)"

win32:
	@echo "[*] Compiling for Windows x86"
	@GOOS=windows GOARCH=386 go build -o $(BIN_FOLDER)/win/$(BINARY_NAME)_win32.exe $(SRC_FOLDER)/main.go
	@echo "[*] Done compiling for Windows x86, binary is located in: $(BIN_FOLDER)"

x64:
	@echo "[*] Compiling for Linux x64"
	@GOOS=linux GOARCH=amd64 go build -o $(BIN_FOLDER)/nix/$(BINARY_NAME)_x64 $(SRC_FOLDER)/main.go

x86:
	@echo "[*] Compiling for Linux x86"
	@GOOS=linux GOARCH=386 go build -o $(BIN_FOLDER)/nix/$(BINARY_NAME)_x86 $(SRC_FOLDER)/main.go


clean:
	@echo "[*] Cleaning up"
	@rm -rf $(BIN_FOLDER)/*
	@echo "[*] Done cleaning up"
