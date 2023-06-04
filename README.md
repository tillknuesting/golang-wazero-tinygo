# golang-wazero-tinygo

cd testdata; tinygo build -o add.wasm -target=wasi add.go; cd ..
cd testdata; tinygo build -o add.wasm -scheduler=none --no-debug -target wasi add.go; cd ..
