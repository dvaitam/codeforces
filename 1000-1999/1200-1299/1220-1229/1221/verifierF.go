package main

import (
	"fmt"
	"os"
)

const (
	refSource  = "1000-1999/1200-1299/1220-1229/1221/1221F.go"
	refBinary  = "ref1221F.bin"
	totalTests = 60
)

type point struct {
	x int64
	y int64
	c int64
}

type testCase struct {
	points []point
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Problem F is interactive and cannot be automatically verified.")
		return
	}
}
