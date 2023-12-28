X64_GCC := gcc
X86_GCC := mingw32-gcc

CFLAGS := -Wall

SOURCES := $(wildcard *.c)

OUTPUT := implant.exe

all:
	@echo "Usage: make [x64|x86|clean]"

x64:
	$(X64_GCC) $(CFLAGS) $(SOURCES) -o $(OUTPUT)

x86:
	$(X86_GCC) $(CFLAGS) $(SOURCES) -o $(OUTPUT)

clean:
	rm -f $(OUTPUT)
