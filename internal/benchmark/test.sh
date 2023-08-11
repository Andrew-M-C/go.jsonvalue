#!/bin/bash
go mod tidy
go test -bench=. -run=none -benchmem -benchtime=2s
go version