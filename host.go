// Package main of the host application
package main

import (
	"context"
	"fmt"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"log"
	"os"
)

func main() {
	// Choose the context to use for function calls.
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	runtime := wazero.NewRuntime(ctx)

	// This closes everything this Runtime created.
	defer runtime.Close(ctx)

	// Instantiate WASI
	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	// Load the WebAssembly module
	wasmPath := "./module/store.wasm"
	helloWasm, err := os.ReadFile(wasmPath)

	if err != nil {
		log.Panicln(err)
	}

	// Instantiate the guest Wasm into the same runtime.
	// It exports the `hello` function,
	// implemented in WebAssembly.
	mod, err := runtime.Instantiate(ctx, helloWasm)
	if err != nil {
		log.Panicln(err)
	}

	// Get the reference to the WebAssembly functions
	setFunction := mod.ExportedFunction("set")
	getFunction := mod.ExportedFunction("get")

	// state is stored for each wasm instance in a global variable that lives as long as the wasm instance

	// Call the WebAssembly functions
	stringData := `{"key1": "value1", "key2": "value2"}`
	log.Println("Setting the value: ", stringData)
	functionCallhelper(stringData, mod, err, ctx, setFunction)

	stringData2 := "key1"
	log.Println("Getting the value for key: ", stringData2)
	functionCallhelper(stringData2, mod, err, ctx, getFunction)

	stringData3 := `{"key1": "value1updated"}`
	log.Println("Updating the value: ", stringData3)
	functionCallhelper(stringData3, mod, err, ctx, setFunction)

	stringData4 := "key1"
	log.Println("Getting the value for key: ", stringData4)
	functionCallhelper(stringData4, mod, err, ctx, getFunction)
}

func functionCallhelper(stringData string, mod api.Module, err error, ctx context.Context, setFunction api.Function) {
	nameSize := uint64(len(stringData))

	// These function are exported by TinyGo
	malloc := mod.ExportedFunction("malloc")
	free := mod.ExportedFunction("free")

	// Allocate Memory
	results, err := malloc.Call(ctx, nameSize)
	if err != nil {
		log.Panicln(err)
	}
	namePosition := results[0]

	// This pointer is managed by TinyGo,
	// but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	defer free.Call(ctx, namePosition)

	// Copy string to memory
	if !mod.Memory().Write(uint32(namePosition), []byte(stringData)) {
		log.Panicf("out of range of memory size")
	}

	call(ctx, setFunction, namePosition, nameSize, mod)
}

func call(ctx context.Context, functionName api.Function, namePosition uint64, nameSize uint64, mod api.Module) {
	// the result type is []uint64
	result, err := functionName.Call(ctx, namePosition, nameSize)
	if err != nil {
		log.Panicln(err)
	}

	if result[0] == 0 {
		fmt.Println("Returned value is null")
		return
	}

	// Extract the position and size of the returned value
	valuePosition := uint32(result[0] >> 32)
	valueSize := uint32(result[0])

	// Read the value from the memory
	if bytes, ok := mod.Memory().Read(valuePosition, valueSize); !ok {
		log.Panicf("out of range of memory size")
	} else {
		fmt.Println("Returned value :", string(bytes))
	}
}
