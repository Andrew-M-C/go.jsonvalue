package main

import (
	"encoding/json"
	"log"
	"os"
	"runtime/pprof"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

// brew install graphviz
// go tool pprof -http=:6060 ./profile

const (
	iteration = 1000000
)

var (
	unmarshalText = []byte(`{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!"},"array":[{"int":123456,"float":123.456789,"string":"Hello, world!"},{"int":123456,"float":123.456789,"string":"Hello, world!"}]}}},"array":[{"int":123456,"float":123.456789,"string":"Hello, world!"},{"int":123456,"float":123.456789,"string":"Hello, world!"}]}`)
	printf        = log.Printf
)

func jsonvalueTest() {
	f, err := os.Create("jsonvalue.profile")
	if err != nil {
		log.Fatal(err)
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < iteration; i++ {
		_, err := jsonvalue.Unmarshal(unmarshalText)
		if err != nil {
			printf("unmarshal error: %v", err)
			return
		}
	}

	printf("jsonvalue done")
	return
}

func mapInterfaceTest() {
	f, err := os.Create("mapinterface.profile")
	if err != nil {
		log.Fatal(err)
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < iteration; i++ {
		m := map[string]interface{}{}
		err := json.Unmarshal(unmarshalText, &m)
		if err != nil {
			printf("unmarshal error: %v", err)
			return
		}
	}

	printf("mapinterface done")
	return
}

func main() {
	printf("start")
	jsonvalueTest()
	mapInterfaceTest()
}
