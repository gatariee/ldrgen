BINARY_NAME = ldrgen
BIN_FOLDER = ./bin

.PHONY: all windows linux clean

all:
	@echo [*] No default target specified, available targets: `x64`, `x86`, `win64`, `win32`.
	@echo [*] Alternatively, you can run `make windows` or `make linux` to compile for both architectures

ldrgen: windows linux
	@echo [*] Done compiling for all platforms, binary is located at: $(BIN_FOLDER)

windows: win64 win32
	@echo [*] Done compiling for Windows x64 and x86, binary is located at: $(BIN_FOLDER)


linux: x64 x86
	@echo [*] Done compiling for Linux x64 and x86, binary is located at: $(BIN_FOLDER)

x64:
	@echo [*] Building $(BINARY_NAME)-linux-x64...
	@set GOOS=linux&& set GOARCH=amd64&& go build -o $(BIN_FOLDER)/$(BINARY_NAME)-linux-x64

x86:
	@echo [*] Building $(BINARY_NAME)-linux-x86...
	@set GOOS=linux&& set GOARCH=386&& go build -o $(BIN_FOLDER)/$(BINARY_NAME)-linux-x86

win64:
	@echo [*] Building $(BINARY_NAME)-win64.exe...
	@set GOOS=windows&& set GOARCH=amd64&& go build -o $(BIN_FOLDER)\$(BINARY_NAME)-win64.exe

win32:
	@echo [*] Building $(BINARY_NAME)-win32.exe...
	@set GOOS=windows&& set GOARCH=386&& go build -o $(BIN_FOLDER)\$(BINARY_NAME)-win32.exe

clean:
	@if exist "$(BIN_FOLDER)\$(BINARY_NAME)-linux-x64" del /f "$(BIN_FOLDER)\$(BINARY_NAME)-linux-x64"
	@if exist "$(BIN_FOLDER)\$(BINARY_NAME)-linux-x86" del /f "$(BIN_FOLDER)\$(BINARY_NAME)-linux-x86"
	@if exist "$(BIN_FOLDER)\$(BINARY_NAME)-win64.exe" del /f "$(BIN_FOLDER)\$(BINARY_NAME)-win64.exe"
	@if exist "$(BIN_FOLDER)\$(BINARY_NAME)-win32.exe" del /f "$(BIN_FOLDER)\$(BINARY_NAME)-win32.exe"
	@echo [*] Done.