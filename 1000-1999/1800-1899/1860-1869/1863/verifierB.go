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

type testCaseB struct {
	arr []int
}

func generateTestsB(num int) []testCaseB {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testCaseB, num)
	for i := 0; i < num; i++ {
		n := rand.Intn(50) + 1
		perm := rand.Perm(n)
		arr := make([]int, n)
		for j, v := range perm {
			arr[j] = v + 1
		}
		tests[i] = testCaseB{arr: arr}
	}
	return tests
}

func solveB(tc testCaseB) string {
	n := len(tc.arr)
	pos := make([]int, n+1)
	for i, x := range tc.arr {
		pos[x] = i + 1
	}
	segments := 1
	for v := 2; v <= n; v++ {
		if pos[v] < pos[v-1] {
			segments++
		}
	}
	return fmt.Sprint(segments - 1)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsB(100)
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, len(tc.arr))
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
		expected := solveB(tc)
		if strings.TrimSpace(outputs[i]) != expected {
			fmt.Printf("mismatch on test %d: expected %s got %s\n", i+1, expected, outputs[i])
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
