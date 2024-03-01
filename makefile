BINARY_NAME := ldrgen
SRC_FOLDER := src
BIN_FOLDER := bin

GO_BUILD_CMD := GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o

.PHONY: all windows linux clean win64 win32 x64 x86

all:
	@echo "[*] No default target specified, available targets: 'win64', 'win32', 'x64', 'x86'."
	@echo "[*] Alternatively, you can run 'make windows' or 'make linux' to compile for both platforms."

windows: GOOS=windows
windows: win64 win32
	@echo "[*] Done compiling for Windows x64 and x86, binaries are located at: $(BIN_FOLDER)"

linux: GOOS=linux
linux: x64 x86
	@echo "[*] Done compiling for Linux x64 and x86, binaries are located at: $(BIN_FOLDER)"

win64: GOARCH=amd64
win64: binary
	@echo "[*] Done compiling for Windows x64, binary is located at: $(BIN_FOLDER)"

win32: GOARCH=386
win32: binary
	@echo "[*] Done compiling for Windows x86, binary is located at: $(BIN_FOLDER)"

x64: GOARCH=amd64
x64: binary
	@echo "[*] Done compiling for Linux x64, binary is located at: $(BIN_FOLDER)"

x86: GOARCH=386
x86: binary
	@echo "[*] Done compiling for Linux x86, binary is located at: $(BIN_FOLDER)"

binary:
	@echo "[*] Compiling for $(GOOS) $(GOARCH)"
	@mkdir -p $(BIN_FOLDER)
	@$(GO_BUILD_CMD) $(BIN_FOLDER)/$(BINARY_NAME)_$(GOOS)_$(GOARCH) $(SRC_FOLDER)/main.go

clean:
	@echo "[*] Cleaning up"
	@rm -rf $(BIN_FOLDER)
	@echo "[*] Done cleaning up"
