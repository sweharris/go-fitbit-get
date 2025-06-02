SRC=$(wildcard *.go)

fitbit_get: $(SRC)
	go build
