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
	iteration = 100000
)

var (
	unmarshalText = []byte(`{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!"},"array":[{"int":123456,"float":123.456789,"string":"Hello, world!"},{"int":123456,"float":123.456789,"string":"Hello, world!"}]}}},"array":[{"int":123456,"float":123.456789,"string":"Hello, world!"},{"int":123456,"float":123.456789,"string":"Hello, world!"}]}`)
	printf        = log.Printf
)

func jsonvalueUnmarshalTest() {
	f, err := os.OpenFile("jsonvalue-unmarshal.profile", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
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

	printf("jsonvalue unmarshal done")
	return
}

func jsonvalueMarshalTest() {
	j, err := jsonvalue.Unmarshal(unmarshalText)
	if err != nil {
		printf("marshal error: %v", err)
		return
	}

	f, err := os.OpenFile("jsonvalue-marshal.profile", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < iteration; i++ {
		_, err = j.Marshal()
		if err != nil {
			printf("marshal error: %v", err)
			return
		}
	}

	printf("jsonvalue marshal done")
	return
}

func mapInterfaceUnmarshalTest() {
	f, err := os.OpenFile("mapinterface-unmarshal.profile", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
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

	printf("mapinterface unmarshal done")
	return
}

func mapInterfaceMarshalTest() {
	m := map[string]interface{}{}
	err := json.Unmarshal(unmarshalText, &m)
	if err != nil {
		printf("unmarshal error: %v", err)
		return
	}

	f, err := os.OpenFile("mapinterface-marshal.profile", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < iteration; i++ {
		_, err = json.Marshal(&m)
		if err != nil {
			printf("marshal error: %v", err)
			return
		}
	}

	printf("mapinterface marshal done")
	return
}

func main() {
	printf("start")
	jsonvalueUnmarshalTest()
	jsonvalueMarshalTest()
	mapInterfaceUnmarshalTest()
	mapInterfaceMarshalTest()
}
