# ldrgen

**ldrgen** allows operators to generate shellcode loaders from a set of pre-defined templates and profiles with the aim of hopefully streamlining development and reducing the time it takes to get a new loader out the door.

## Getting Started
There are available binaries on the [releases](https://github.com/gatariee/ldrgen/releases) page, or you can build from source for the latest version.

### Usage
```bash
ldrgen profile --profile [path_to_profile] --shellcode [path_to_shellcode]
```
![hmm](https://i.gyazo.com/9b20c0e2699542268ed51aecba6eed2d.png)

### How does it work?
* Source files (`*.c`) or loaders are saved as templates in the `templates` directory, placeholders are also supported to be templated during generation.
* The template configuration is saved as [`config.yaml`](./templates/config.yaml) in the `templates` directory, this contains information about the loader, what source files it requires and whether any substitutions are required.
* Profiles are saved in the `profiles` directory, these can be used and customized to your personal preference to utilize the templates you've created.

#### Source
These are your loaders, most of the time these are the artifacts that get caught by AV. Put some effort into making some of these in your own time, and you'll have a nice collection of loaders to use!

[CreateThread_Xor.c](./templates/Source/CreateThread_Xor.c) is a good example of a simple loader that uses `CreateThread` to execute shellcode. It also has a `xorShellcode` function that will decrypt the shellcode before execution.
```c
...

    xorShellcode( shellcode, shellcode_size, "${KEY}" );

...
```

#### Templates
The ${KEY} placeholder is explicitly defined in the `config.yaml` file, and is used to substitute the key used to decrypt the shellcode. This is a simple example of how to use placeholders in your templates.
```yaml
  - token: "createthread_xor"
    key_required: true
    enc_type: "xor"
    method: "VirtualAlloc, xorShellcode, memcpy, CreateThread, WaitForSingleObject"
    files:
      - sourcePath: "Source/CreateThread_Xor.c"
        outputPath: "Main.c"
        substitutions: 
          key: "${KEY}"

      - sourcePath: "Source/Xor.c"
        outputPath: "Xor.c"
      
      - sourcePath: "Include/Xor.h"
        outputPath: "Xor.h"

      - sourcePath: "makefile"
        outputPath: "makefile"
```

#### Profiles
These are malleable configurations that can be used to generate a loader from a template.
```json
{
    "name": "Default Loader, but with XOR (CreateThread)",
    "author": "@gatari",
    "description": "A simple shellcode loader that uses CreateThread with XOR encrypted shellcode.",
    "template": {
        "path": "/opt/tools/ldrgen/templates/config.yaml",
        "token": "CreateThread_Xor",
        "enc_type": "xor",
        "substitutions": {
            "key": "as@&(!L@J#JKsn"
        }
    },
    "arch": "x64",
    "compile": {
        "automatic": true,
        "make": "make",
        "gcc": {
            "x64": "x86_64-w64-mingw32-gcc",
            "x86": "i686-w64-mingw32-gcc"
        },
        "strip": {
            "i686": "strip",
            "x86_64": "strip"
        }
    },
    "output_dir": "./createthreadxor"
}
```

### Releases
Standalone binaries are available, however `templates` and `profiles` will be ingested by `ldrgen`. The latest release should have these included.

### Building from Source
```bash
cd ldrgen
make build
```

### Design Principles
This tool is **not** designed to be turing complete by any means, it's designed to **streamline** the process of developing loaders- it is ultimately up to the operator to develop their own loaders and templates should they wish to do so.

`ldrgen` started as a personal project for encrypting beacon shellcode, then parsing it into a header file for use in a stageless loader. It was a pain to do this manually, so I decided to automate it- and it built from there.

### Contributing
If you'd like to contribute, feel free to open a PR or an issue. I'm always open to suggestions and improvements. However, before you make an issue, do remember that the objective of this tool is to streamline the process of developing loaders, not to be provide a cheat code for bypassing AV.