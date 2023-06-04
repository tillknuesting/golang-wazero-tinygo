# golang-wazero-tinygo

This is a proof of concept of a key-value store implemented in Go and compiled into WebAssembly using TinyGo.
The state is stored in a global variable in the wasm module, which is persistent as long as the host keeps the instance alive.
Strings are shared between the host and the wasm module using a shared memory buffer.

A set and a get function are exported and can be called from the host.

The set function takes keys and values as JSON strings and stores the value under the key in a global map.

The get function takes a key and returns the value as JSON.

To try it out, run go run host.go and use the provided precompiled wasm file.

### Versions used:
tinygo version 0.28.0-dev-e2e6570 darwin/amd64 (using go version go1.20.2 and LLVM version 15.0.0)

(build dev branch of tinygo yourself to get the latest version which supports reflection)

### build wasm file:

```bash
cd module; tinygo build -o store.wasm -scheduler=none --no-debug -target wasi store.go; cd ..
```

### run host that calls functions from wasm module:

```bash
go run host.go
```

Example output:

```bash
go run host.go
Setting the value:  {"key1": "value1", "key2": "value2"}
Returned value : success

Getting the value for key:  key1
Returned value : {"key1":"value1"}

Updating the value:  {"key1": "value1updated"}
Returned value : success

Getting the value for key:  key1
Returned value : {"key1":"value1updated"}

```
