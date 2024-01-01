def hash( value, seed=5 ):
    p = 31
    lp = 1000000007
    ret = seed % lp
    for char in str( value ):
        ret = ( ( ret << 2 ) ^ ( ord( char ) << 1 ) ) % lp
        ret = ( ret >> 1 ) ^ ( p * ret )

    return hex ( abs( ret ) )


apis = [
    "CreateProcessA",
    "VirtualAllocEx",
    "WriteProcessMemory",
    "QueueUserAPC",
    "ResumeThread",
    "WaitForSingleObject",
    "CloseHandle",
    "VirtualFreeEx",
    "1000",
]

for api in apis:
    print(f"{ api } -> { hash( api ) }")
