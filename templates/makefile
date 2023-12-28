X64_COMPILER := x86_64-w64-mingw32-gcc # apt install mingw-w64
X86_COMPILER := i686-w64-mingw32-gcc
CFLAGS := -Wall -MMD -MP
LDFLAGS := -lkernel32

SOURCES := $(wildcard *.c)
OBJECTS_X64 := $(SOURCES:%.c=%_x64.o)
OBJECTS_X86 := $(SOURCES:%.c=%_x86.o)
DEPS := $(OBJECTS_X64:.o=.d) $(OBJECTS_X86:.o=.d)
OUTPUT := implant

.PHONY: all
all:
	@echo "Usage: make [x64 | x86]"

.PHONY: x64
x64: $(OUTPUT)_x64.exe

$(OUTPUT)_x64.exe: $(OBJECTS_X64)
	$(X64_COMPILER) $(LDFLAGS) $^ -o $@ -m64
	@echo "[+] implant_x64.exe -> OK"

.PHONY: x86
x86: $(OUTPUT)_x86.exe

$(OUTPUT)_x86.exe: $(OBJECTS_X86)
	$(X86_COMPILER) $(LDFLAGS) $^ -o $@ -m32
	@echo "[+] implant_x86.exe -> OK"

%_x64.o: %.c
	$(X64_COMPILER) $(CFLAGS) -c $< -o $@ -m64

%_x86.o: %.c
	$(X86_COMPILER) $(CFLAGS) -c $< -o $@ -m32

-include $(DEPS)

.PHONY: clean
clean:
	rm -f $(OUTPUT)_x64.exe $(OUTPUT)_x86.exe $(OBJECTS_X64) $(OBJECTS_X86) $(DEPS)
	@echo "Cleaned."