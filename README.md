# ldrgen

ldrgen is a golang cli tool to generate shellcode loaders

## Usage
```
go build .
./ldrgen -h
```

## Why?
when ur doing a box and your beacons die when dropping / running from disk 

## How?
1. generate your shellcode

![generate shellcode](./assets/3e27d7894ec76ece20e41fd290df7ded.png)

2. `./ldr -b <bin_path> -o <out_file> -ldr <loader_type> -enc <encryption_type> -key <encryption_key>`

![generate loader](./assets/beb0f93fce10788ff4fafa558c7bec54.png)

3. loader source code will be generated in `out_file`

![loader source code](./assets/d76dc3645cf50997bf17ba2c28ed3565.png)

4. compile & run the loader

![run](./assets/bad05d44ec8a4ad5b361d0e5eb3bf2a3.png)

5. profit?

![profit](./assets/c2f1fd7a899c87ffd61303b6d46a6e2b.png)
