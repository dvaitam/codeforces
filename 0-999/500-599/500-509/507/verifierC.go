package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type TestC struct {
	h int
	n int64
}

func generateTests() []TestC {
	rand.Seed(3)
	tests := make([]TestC, 100)
	for i := range tests {
		h := rand.Intn(25) + 1
		maxN := int64(1) << h
		n := rand.Int63n(maxN) + 1
		tests[i] = TestC{h, n}
	}
	return tests
}

func expected(t TestC) int64 {
	h := t.h
	n := t.n
	var ans int64
	dir := 0
	for i := h; i > 0; i-- {
		half := int64(1) << (i - 1)
		if (n <= half && dir == 0) || (n > half && dir == 1) {
			ans += 1
			dir ^= 1
		} else {
			ans += int64(1) << i
		}
		if n > half {
			n -= half
		}
	}
	return ans
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n", t.h, t.n)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.TrimSpace(got), 10, 64)
		if err != nil {
			fmt.Printf("test %d: invalid output\n", i+1)
			os.Exit(1)
		}
		exp := expected(t)
		if val != exp {
			fmt.Printf("test %d: expected %d got %d\n", i+1, exp, val)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
