package main

import (
	"encoding/json"
	"log"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"time"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
	jsoniter "github.com/json-iterator/go"
)

// brew install graphviz
// go run .
// go tool pprof -http=:6060 ./jsonvalue-unmarshal.profile
// go tool trace jsonvalue-unmarshal-trace.profile

const (
	iteration = 200000
)

var (
	unmarshalText = []byte(`{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!","object":{"int":123456,"float":123.456789,"string":"Hello, world!"},"array":[{"int":123456,"float":123.456789,"string":"Hello, world!"},{"int":123456,"float":123.456789,"string":"Hello, world!"}]}}},"array":[{"int":123456,"float":123.456789,"string":"Hello, world!"},{"int":123456,"float":123.456789,"string":"Hello, world!"}]}`)
	printf        = log.Printf
)

func jsonvalueUnmarshalTraceTest() {
	ft, err := os.OpenFile("jsonvalue-unmarshal-trace.profile", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer ft.Close()
	trace.Start(ft)
	defer trace.Stop()

	_, err = jsonvalue.Unmarshal(unmarshalText)
	if err != nil {
		printf("unmarshal error: %v", err)
		return
	}

	printf("jsonvalue unmarshal trace done")
}

func jsonvalueUnmarshalTest() {
	f, err := os.OpenFile("jsonvalue-unmarshal.profile", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < iteration; i++ {
		_, err := jsonvalue.Unmarshal(unmarshalText)
		if err != nil {
			printf("unmarshal error: %v", err)
			return
		}
	}

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

	defer f.Close()
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
}

func mapInterfaceUnmarshalTest() {
	f, err := os.OpenFile("mapinterface-unmarshal.profile", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
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

	defer f.Close()
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
}

type object struct {
	Int    int       `json:"int"`
	Float  float64   `json:"float"`
	String string    `json:"string"`
	Object *object   `json:"object,omitempty"`
	Array  []*object `json:"array,omitempty"`
}

func structUnmarshalTest() {
	f, err := os.OpenFile("struct-unmarshal.profile", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < iteration; i++ {
		o := object{}
		err := json.Unmarshal(unmarshalText, &o)
		if err != nil {
			printf("unmarshal error: %v", err)
			return
		}
	}

	printf("struct unmarshal done")
}

func structMarshalTest() {
	o := object{}
	json.Unmarshal(unmarshalText, &o)

	f, err := os.OpenFile("struct-marshal.profile", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < iteration; i++ {
		_, err := json.Marshal(&o)
		if err != nil {
			printf("marshal error: %v", err)
			return
		}
	}

	printf("struct marshal done")
}

func jsoniterGetTest() {
	f, err := os.OpenFile("jsoniter-get.profile", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < iteration; i++ {
		any := jsoniter.Get(unmarshalText)
		any.Get("object", "object", "object", "array", 1)
	}

	printf("jsoniter get done")
}

func jsoniterUnmarshalTest() {
	j := jsoniter.ConfigCompatibleWithStandardLibrary
	f, err := os.OpenFile("jsoniter-unmarshal.profile", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < iteration; i++ {
		o := object{}
		err := j.Unmarshal(unmarshalText, &o)
		if err != nil {
			printf("unmarshal error: %v", err)
			return
		}
	}

	printf("jsoniter unmarshal done")
}

func jsoniterMarshalTest() {
	j := jsoniter.ConfigCompatibleWithStandardLibrary
	o := object{}
	j.Unmarshal(unmarshalText, &o)

	f, err := os.OpenFile("jsoniter-marshal.profile", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < iteration; i++ {
		_, err := json.Marshal(&o)
		if err != nil {
			printf("marshal error: %v", err)
			return
		}
	}

	printf("jsoniter marshal done")
}

func main() {
	run := func(f func()) {
		start := time.Now().Local()
		printf("start: %v", start)
		f()
		printf("done, elapsed %v", time.Since(start))
	}

	run(jsonvalueUnmarshalTraceTest)

	run(jsonvalueUnmarshalTest)
	run(jsonvalueMarshalTest)

	run(mapInterfaceUnmarshalTest)
	run(mapInterfaceMarshalTest)

	run(structUnmarshalTest)
	run(structMarshalTest)

	run(jsoniterGetTest)
	run(jsoniterUnmarshalTest)
	run(jsoniterMarshalTest)
}
