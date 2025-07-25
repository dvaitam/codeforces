package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveC(n uint64) uint64 {
	l0, l1 := uint64(1), uint64(2)
	h := uint64(1)
	for l1 <= n {
		h++
		l0, l1 = l1, l0+l1
	}
	return h - 1
}

func genTestC() uint64 {
	return uint64(rand.Int63n(1e18-1)) + 2
}

func runBinary(path string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	path := os.Args[1]
	for i := 0; i < 100; i++ {
		n := genTestC()
		input := fmt.Sprintf("%d\n", n)
		expected := solveC(n)
		gotStr, err := runBinary(path, input)
		if err != nil {
			fmt.Printf("test %d: execution error: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		var got uint64
		_, err = fmt.Sscanf(gotStr, "%d", &got)
		if err != nil {
			fmt.Printf("test %d: parse output error: %v\ninput:%soutput:%s\n", i+1, err, input, gotStr)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed\ninput:%sexpected:%d\ngot:%d\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
