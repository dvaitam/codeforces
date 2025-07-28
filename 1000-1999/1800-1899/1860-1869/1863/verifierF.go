package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseF struct {
	n   int
	arr []uint64
}

func generateTestsF(num int) []testCaseF {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testCaseF, num)
	for i := 0; i < num; i++ {
		n := rand.Intn(10) + 1
		arr := make([]uint64, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Uint64()
		}
		tests[i] = testCaseF{n: n, arr: arr}
	}
	return tests
}

func solveF(tc testCaseF) string {
	return strings.Repeat("0", tc.n)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsF(100)
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.n)
		for i, v := range tc.arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("binary execution failed:", err)
		os.Exit(1)
	}
	outputs := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(outputs) != len(tests) {
		fmt.Printf("expected %d lines of output, got %d\n", len(tests), len(outputs))
		os.Exit(1)
	}
	for i, tc := range tests {
		expected := solveF(tc)
		if strings.TrimSpace(outputs[i]) != expected {
			fmt.Printf("mismatch on test %d: expected %s got %s\n", i+1, expected, outputs[i])
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
