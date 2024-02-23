#include <stdint.h>
#include <stdio.h>
#include <windows.h>

uint64_t Hash( const char * str, uint64_t seed ) {
    const uint64_t p  = 31;
    const uint64_t lp = 1000000007;
    uint64_t ret      = seed % lp;

    while ( *str ) {
        ret = ( ( ret << 2 ) ^ ( (uint64_t)( *str ) << 1 ) ) % lp;
        ret = ( ret >> 1 ) ^ ( p * ret );
        str++;
    }
    return ret;
}