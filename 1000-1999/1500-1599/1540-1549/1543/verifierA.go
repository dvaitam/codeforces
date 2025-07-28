package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

func solve(a, b int64) (int64, int64) {
	if a == b {
		return 0, 0
	}
	if a < b {
		a, b = b, a
	}
	diff := a - b
	rem := a % diff
	moves := rem
	if diff-rem < moves {
		moves = diff - rem
	}
	return diff, moves
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(42)
	const T = 100
	tests := make([][2]int64, T)
	for i := 0; i < T; i++ {
		a := rand.Int63n(1e12)
		b := rand.Int63n(1e12)
		tests[i] = [2]int64{a, b}
	}
	var input bytes.Buffer
	fmt.Fprintln(&input, T)
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n", tc[0], tc[1])
	}
	expected := make([]string, T)
	for i, tc := range tests {
		d, m := solve(tc[0], tc[1])
		expected[i] = fmt.Sprintf("%d %d", d, m)
	}
	cmd := exec.Command(binary)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	idx := 0
	for idx < T {
		if !scanner.Scan() {
			fmt.Fprintln(os.Stderr, "insufficient output")
			os.Exit(1)
		}
		diffStr := scanner.Text()
		if !scanner.Scan() {
			fmt.Fprintln(os.Stderr, "insufficient output")
			os.Exit(1)
		}
		movesStr := scanner.Text()
		got := diffStr + " " + movesStr
		if got != expected[idx] {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %s got %s\n", idx+1, expected[idx], got)
			os.Exit(1)
		}
		idx++
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output after", T, "tests")
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
