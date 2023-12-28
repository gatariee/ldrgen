X64_COMPILER := x86_64-w64-mingw32-gcc # apt install mingw-w64
X86_COMPILER := i686-w64-mingw32-gcc

CFLAGS := -Wall -MMD -MP
LDFLAGS := -lkernel32

SOURCES := $(wildcard *.c)
OBJECTS := $(SOURCES:.c=.o)
DEPS := $(OBJECTS:.o=.d)
OUTPUT := implant.exe

.PHONY: all
all:
	@echo "Select target: 'make x64' for 64-bit or 'make x86' for 32-bit"

ifneq ($(SOURCES),)
.PHONY: x64
x64: $(OUTPUT)
	@echo "Building 64-bit PE32 for Windows..."

$(OUTPUT): $(OBJECTS)
	$(X64_COMPILER) $(LDFLAGS) $(OBJECTS) -o $(OUTPUT) -m64

.PHONY: x86
x86: $(OUTPUT)
	@echo "Built 32-bit version."

$(OUTPUT): $(OBJECTS)
	$(X86_COMPILER) $(LDFLAGS) $(OBJECTS) -o $(OUTPUT) -m32

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

-include $(DEPS)

.PHONY: clean
clean:
	rm -f $(OUTPUT) $(OBJECTS) $(DEPS)
	@echo "Cleaned."

else
$(error No source files found in the current directory.)
endif
