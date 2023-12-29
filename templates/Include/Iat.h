#ifndef IAT_H
#define IAT_H

#include <windows.h>
#include <stdint.h>

FARPROC AddrFromHash( HMODULE hLib, uint64_t hashval, uint64_t seed );

#endif 