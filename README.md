# ldrgen

ldrgen is a golang cli tool for rapid generation of shellcode loaders using pre-defined templates.

> ⚠️ this tool is meant to help with running beacon from disk, this does **not** help with evasion for memory scans or post-exploitation- OPSEC considerations are up to the discretion of the operator.

## Getting Started
There are available binaries on the [releases](https://github.com/gatariee/ldrgen/releases) page, or you can build from source for the latest version.

### Releases
Standalone binaries are available, however you will need a `templates` directory to be ingested by the generator. You can find the latest templates from source [here](./templates/) or it should be included as a zip file in the release.

Templates are expected to be in the same directory as where the binary is executed, but you can specify with the `--template` flag.

### Building from Source
```bash
cd ldrgen/src
make build

cd bin
ldrgen --help
```