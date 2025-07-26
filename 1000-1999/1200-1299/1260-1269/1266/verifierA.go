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

func compute(s string) string {
	sum := 0
	zeros := 0
	even := 0
	for i := 0; i < len(s); i++ {
		d := int(s[i] - '0')
		sum += d
		if d == 0 {
			zeros++
			even++
		} else if d%2 == 0 {
			even++
		}
	}
	if zeros > 0 && sum%3 == 0 && (zeros > 1 || even > 1) {
		return "red"
	}
	return "cyan"
}

func randomTests(n int) []string {
	res := make([]string, n)
	for i := 0; i < n; i++ {
		length := rand.Intn(99) + 2 // 2..100 digits
		b := make([]byte, length)
		for j := 0; j < length; j++ {
			b[j] = byte('0' + rand.Intn(10))
		}
		res[i] = string(b)
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run verifierA.go /path/to/binary\n")
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	tests := randomTests(100)

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, t := range tests {
		fmt.Fprintln(&input, t)
	}

	cmd := exec.Command(os.Args[1])
	cmd.Stdin = &input
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "execution failed: %v\n", err)
		os.Exit(1)
	}

	outputs := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(outputs) != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d lines of output, got %d\n", len(tests), len(outputs))
		os.Exit(1)
	}

	for i, in := range tests {
		want := compute(in)
		got := strings.TrimSpace(outputs[i])
		if want != got {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput: %s\nexpected: %s\ngot: %s\n", i+1, in, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
