X64_COMPILER := x86_64-w64-mingw32-gcc # apt install mingw-w64
X86_COMPILER := i686-w64-mingw32-gcc
CFLAGS := -Wall -MMD -MP
LDFLAGS := -lkernel32

BIN_DIR := bin
OBJ_DIR := obj
DEP_DIR := dep

$(shell mkdir -p $(BIN_DIR) $(OBJ_DIR) $(DEP_DIR))

SOURCES := $(wildcard *.c)
OBJECTS_X64 := $(SOURCES:%.c=$(OBJ_DIR)/%_x64.o)
OBJECTS_X86 := $(SOURCES:%.c=$(OBJ_DIR)/%_x86.o)
DEPS := $(OBJECTS_X64:$(OBJ_DIR)/%.o=$(DEP_DIR)/%.d) $(OBJECTS_X86:$(OBJ_DIR)/%.o=$(DEP_DIR)/%.d)
OUTPUT := implant

.PHONY: all
all:
	@echo "Usage: make [x64 | x86]"

.PHONY: both
both: x64 x86

.PHONY: x64
x64: $(BIN_DIR)/$(OUTPUT)_x64.exe

$(BIN_DIR)/$(OUTPUT)_x64.exe: $(OBJECTS_X64)
	$(X64_COMPILER) $(LDFLAGS) $^ -o $@ -m64
	@echo "[+] $(BIN_DIR)/implant_x64.exe -> OK"

.PHONY: x86
x86: $(BIN_DIR)/$(OUTPUT)_x86.exe

$(BIN_DIR)/$(OUTPUT)_x86.exe: $(OBJECTS_X86)
	$(X86_COMPILER) $(LDFLAGS) $^ -o $@ -m32
	@echo "[+] $(BIN_DIR)/implant_x86.exe -> OK"

$(OBJ_DIR)/%_x64.o: %.c
	$(X64_COMPILER) $(CFLAGS) -c $< -o $@ -m64

$(OBJ_DIR)/%_x86.o: %.c
	$(X86_COMPILER) $(CFLAGS) -c $< -o $@ -m32

-include $(DEPS)

.PHONY: clean
clean:
	rm -rf $(BIN_DIR) $(OBJ_DIR) $(DEP_DIR)
	@echo "Clean -> OK"
