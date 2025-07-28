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

type testCaseA struct {
	n int
	a int
	q int
	s string
}

func generateTestsA(num int) []testCaseA {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testCaseA, num)
	for i := 0; i < num; i++ {
		n := rand.Intn(100) + 1
		a := rand.Intn(n + 1)
		q := rand.Intn(100) + 1
		sb := strings.Builder{}
		for j := 0; j < q; j++ {
			if rand.Intn(2) == 0 {
				sb.WriteByte('+')
			} else {
				sb.WriteByte('-')
			}
		}
		tests[i] = testCaseA{n: n, a: a, q: q, s: sb.String()}
	}
	return tests
}

func solveA(tc testCaseA) string {
	cur := tc.a
	maxOnline := tc.a
	plus := 0
	for _, ch := range tc.s {
		if ch == '+' {
			cur++
			plus++
		} else {
			cur--
		}
		if cur > maxOnline {
			maxOnline = cur
		}
	}
	if maxOnline >= tc.n {
		return "YES"
	}
	if tc.a+plus < tc.n {
		return "NO"
	}
	return "MAYBE"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsA(100)
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d %d\n%s\n", tc.n, tc.a, tc.q, tc.s)
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
		expected := solveA(tc)
		if strings.TrimSpace(outputs[i]) != expected {
			fmt.Printf("mismatch on test %d: expected %s got %s\n", i+1, expected, outputs[i])
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
