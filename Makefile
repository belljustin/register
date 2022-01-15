BIN = ./bin

clean:
	rm -rf ${BIN}/*

buildplugins: clean
	go build -buildmode=plugin -o ${BIN}/currencyapi.so ./drivers/currencyapi/plugin
	go build -buildmode=plugin -o ${BIN}/freeforex.so ./drivers/freeforex/plugin

runcompile:
	go run cmd/compiletime/main.go CAD USD

runruntime: buildplugins
	go run cmd/runtime/main.go ${BIN}