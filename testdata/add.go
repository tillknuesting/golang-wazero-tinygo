package main

import (
	"encoding/json"
	"fmt"
	"log"
	"unsafe"
)

func main() {}

type KeyValueStore struct {
	data map[string]string
}

var kvs = &KeyValueStore{
	data: make(map[string]string),
}

//export set
func set(valuePosition *uint32, length uint32) uint64 {
	// read the memory to get the parameter
	valueBytes := readBufferFromMemory(valuePosition, length)

	var returnValue []byte

	var data map[string]string
	err := json.Unmarshal(valueBytes, &data)
	if err != nil {
		returnValue = []byte(fmt.Sprintf("error: %s", err))
	} else {
		for key, value := range data {
			kvs.data[key] = value
		}
		returnValue = []byte(fmt.Sprintf("success"))
	}

	posSizePairValue := copyBufferToMemory(returnValue)

	// return the position and size
	return posSizePairValue
}

//export get
func get(valuePosition *uint32, length uint32) uint64 {
	// read the memory to get the parameter
	valueBytes := readBufferFromMemory(valuePosition, length)

	key := string(valueBytes)

	var returnValue []byte

	value, ok := kvs.data[key]
	if !ok {
		returnValue = []byte(fmt.Sprintf("key %s not found", key))
	} else {
		result := map[string]string{key: value}
		jsonData, err := json.Marshal(result)
		if err != nil {
			log.Panicln(err)
		}
		returnValue = jsonData
	}

	// copy the value to memory
	posSizePairValue := copyBufferToMemory(returnValue)

	// return the position and size
	return posSizePairValue
}

//export add
func add(valuePosition *uint32, length uint32) uint64 {
	// read the memory to get the parameter
	valueBytes := readBufferFromMemory(valuePosition, length)

	// copy the value to memory
	posSizePairValue := copyBufferToMemory([]byte("return message:" + string(valueBytes)))

	// return the position and size
	return posSizePairValue
}

// readBufferFromMemory returns a buffer from WebAssembly
func readBufferFromMemory(bufferPosition *uint32, length uint32) []byte {
	subjectBuffer := make([]byte, length)
	pointer := uintptr(unsafe.Pointer(bufferPosition))
	for i := 0; i < int(length); i++ {
		s := *(*int32)(unsafe.Pointer(pointer + uintptr(i)))
		subjectBuffer[i] = byte(s)
	}
	return subjectBuffer
}

// copyBufferToMemory returns a single value (a kind of pair with position and length)
func copyBufferToMemory(buffer []byte) uint64 {
	bufferPtr := &buffer[0]
	unsafePtr := uintptr(unsafe.Pointer(bufferPtr))

	ptr := uint32(unsafePtr)
	size := uint32(len(buffer))

	return (uint64(ptr) << uint64(32)) | uint64(size)
}
