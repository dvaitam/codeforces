package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveA(n int64) int64 {
	k := int64(1)
	for n > k {
		n -= k
		k++
	}
	return n
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(1)
	tests := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		100, 500, 1000, 99999999999999, 100000000000000}
	for len(tests) < 100 {
		tests = append(tests, rand.Int63n(100000000000000)+1)
	}
	for idx, n := range tests {
		inp := fmt.Sprintf("%d\n", n)
		expected := fmt.Sprintf("%d", solveA(n))
		got, err := runBinary(bin, inp)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed: n=%d expected %s got %s\n", idx+1, n, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
