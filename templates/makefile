X64_COMPILER := x86_64-w64-mingw32-gcc
X86_COMPILER := i686-w64-mingw32-gcc
STRIP := strip

CFLAGS := -Wall
LDFLAGS := -lkernel32

BIN_DIR := bin

SOURCES := $(wildcard *.c)
OUTPUT := implant

.PHONY: all
all:
	@echo "Usage: make [x64 | x86]"

.PHONY: both
both: x64 x86

.PHONY: x64
x64: $(BIN_DIR)/$(OUTPUT)_x64.exe

$(BIN_DIR)/$(OUTPUT)_x64.exe: $(SOURCES)
	@mkdir -p $(BIN_DIR)
	@$(X64_COMPILER) $(CFLAGS) $^ -o $@ $(LDFLAGS) -m64
	@$(STRIP) $@
	@echo [+] $@ -> OK

.PHONY: x86
x86: $(BIN_DIR)/$(OUTPUT)_x86.exe

$(BIN_DIR)/$(OUTPUT)_x86.exe: $(SOURCES)
	@mkdir -p $(BIN_DIR)
	@$(X86_COMPILER) $(CFLAGS) $^ -o $@ $(LDFLAGS) -m32
	@$(STRIP) $@
	@echo [+] $@ -> OK

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
	@echo "Clean -> OK"
