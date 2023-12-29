#include <string.h>
#include <windows.h>

void xorShellcode( unsigned char *shellcode, size_t shellcodeSize, const char *key ) {
    size_t keyLen = strlen( key );

    for ( size_t i = 0; i < shellcodeSize; ++i ) {
        shellcode[i] ^= key[i % keyLen];
    }
}
