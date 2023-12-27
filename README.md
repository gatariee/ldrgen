# ldrgen

ldrgen is a golang cli tool to generate shellcode loaders

## Usage
```
go build .
./ldrgen --help
```

## Why?
When you're doing a box and your beacons die when dropping / running from disk, and packers don't do a good enough job.

## Loader Tokens
Currently available injection tokens:
- [Inline.c](./templates/Source/Inline.c)
    - `VirtualAlloc` to allocate RWX memory, `memcpy` shellcode, execute inline with `( (void ( * )())exec )();`
- [Inline_Xor.c](./templates/Source/Inline_Xor.c)
    - `VirtualAlloc` to allocate RWX memory, `memcpy` shellcode, encrypt shellcode via [Xor.c](./templates/Source/Xor.c), execute inline with `( (void ( * )())exec )();`
- [CreateRemoteThread.c](./templates/Source/CreateRemoteThread.c)
    - `OpenProcess`, `VirtualAllocEx` with *PAGE_EXECUTE_READWRITE*, `WriteProcessMemory` and `CreateRemoteThread` to execute shellcode
- [CreateRemoteThreadRX.c](./templates/Source/CreateRemoteThreadRX.c)
    - `OpenProcess`, `VirtualAllocEx` with *PAGE_EXECUTE_READ*, `WriteProcessMemory`, `VirtualProtectEx` with *PAGE_EXECUTE_READ* and `CreateRemoteThread` to execute shellcode

## Shellcode Templates
These are globally accessible variables via loader templates, and will be replaced with the appropriate values when generating the loader source code.
1. [Shellcode.c](./templates/Source/Shellcode.c)
    - `${SHELLCODE}` -> shellcode in the format: { 0x00, 0x00, 0x00, 0x00, ... }
    - `${SHELLCODE_SIZE}` -> shellcode length in bytes
2. [Shellcode.h](./templates/Include/Shellcode.h)
    - `extern unsigned char shellcode[];`
    - `extern size_t shellcode_size;`

## Example (Sliver Session)
1. Generate shellcode
- `-l` -> disable symbol obfuscation (enable this if injection is flagged)
- `--format shellcode` -> output as .bin file
- `-G` -> optional: disable [shigata-ga-nai](https://unprotect.it/technique/shikata-ga-nai-sgn/#:~:text=Shikata%20Ga%20Nai%20(SGN)%20is,a%20self%2Ddecoding%20obfuscated%20shellcode.) (requires RWX memory, check loader support)

![generate shellcode](./assets/3e27d7894ec76ece20e41fd290df7ded.png)

2. `./ldr -b <bin_path> -o <out_file> -ldr <loader_type> -enc <encryption_type> -key <encryption_key>`

![generate loader](./assets/beb0f93fce10788ff4fafa558c7bec54.png)

3. loader source code will be generated in `out_file`

![loader source code](./assets/d76dc3645cf50997bf17ba2c28ed3565.png)

4. compile & run the loader

![run](./assets/bad05d44ec8a4ad5b361d0e5eb3bf2a3.png)

5. profit?

![profit](./assets/c2f1fd7a899c87ffd61303b6d46a6e2b.png)
