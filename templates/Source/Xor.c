#include <string.h>
#include <windows.h>

void Xor( BYTE * shellcode, DWORD shellcodeSize, const char * key ) {
    size_t keyLen = strlen( key );
    BYTE   temp;

    for ( size_t i = 0; i < shellcodeSize; ++i ) {
        temp = key[i % keyLen];

        shellcode[i] ^= temp;

        temp ^= 0xFF;
        temp ^= 0xFF;

        shellcode[i] ^= ( temp ^ temp );
    }
}
