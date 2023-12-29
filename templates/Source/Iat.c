#include "Hash.h"
/*
* Hash.c

 - Remember to include this for IAT resolution

 {
    uint64_t Hash( const char* str, uint64_t seed );
 }

*/

#include <stdint.h>
#include <stdio.h>
#include <windows.h>

FARPROC AddrFromHash( HMODULE hLib, uint64_t hashval, uint64_t seed ) {

    /*
    https://www.ired.team/offensive-security/defense-evasion/windows-api-hashing-in-malware
    */

    fprintf( stdout, "[-] Tasked to resolve hash: 0x%llx\n", hashval );

    FARPROC ret = NULL;
    IMAGE_DOS_HEADER *dosHeader = (IMAGE_DOS_HEADER *)hLib;
    IMAGE_NT_HEADERS *ntHeader = (IMAGE_NT_HEADERS *)( (uint64_t)hLib + dosHeader->e_lfanew );
    IMAGE_EXPORT_DIRECTORY *exportDir = (IMAGE_EXPORT_DIRECTORY *)( (uint64_t)hLib + ntHeader->OptionalHeader.DataDirectory[IMAGE_DIRECTORY_ENTRY_EXPORT].VirtualAddress );
    uint32_t *nameTable = (uint32_t *)( (uint64_t)hLib + exportDir->AddressOfNames );
    uint16_t *ordinalTable = (uint16_t *)( (uint64_t)hLib + exportDir->AddressOfNameOrdinals );
    uint32_t *functionTable = (uint32_t *)( (uint64_t)hLib + exportDir->AddressOfFunctions );

    for ( uint32_t i = 0; i < exportDir->NumberOfNames; i++ ) {
        char *name = (char *)( (uint64_t)hLib + nameTable[i] );
        uint64_t nameHash = Hash( name, seed );
        if ( nameHash == hashval ) {
            uint16_t ordinal = ordinalTable[i];
            uint32_t functionRVA = functionTable[ordinal];
            ret = (FARPROC)( (uint64_t)hLib + functionRVA );
            fprintf( stdout, "[+] Found: %s (0x%llx) at 0x%llx\n\n", name, hashval, ret );
            break;
        }
    }

    return ret;
};