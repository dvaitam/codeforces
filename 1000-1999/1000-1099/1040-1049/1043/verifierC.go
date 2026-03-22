package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func optimal(s string) string {
	k := strings.Count(s, "a")
	return strings.Repeat("a", k) + strings.Repeat("b", len(s)-k)
}

func simulate(s string, ops []int) string {
	b := []byte(s)
	for i, op := range ops {
		if op == 1 {
			// Reverse prefix of length i+1 (indices 0..i)
			for l, r := 0, i; l < r; l, r = l+1, r-1 {
				b[l], b[r] = b[r], b[l]
			}
		}
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	tests := []string{"a", "b", "ab", "ba", "aa", "bb", "abba", "aaaa", "bbbb", "abab", "baba"}
	for i := 0; i < 200; i++ {
		n := rng.Intn(20) + 1
		b := make([]byte, n)
		for j := range b {
			if rng.Intn(2) == 0 {
				b[j] = 'a'
			} else {
				b[j] = 'b'
			}
		}
		tests = append(tests, string(b))
	}

	for idx, s := range tests {
		input := s + "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s\ninput: %s\n", idx+1, err, errBuf.String(), s)
			os.Exit(1)
		}

		parts := strings.Fields(strings.TrimSpace(out.String()))
		n := len(s)
		if len(parts) != n {
			fmt.Fprintf(os.Stderr, "case %d: expected %d integers, got %d: %q\ninput: %s\n", idx+1, n, len(parts), out.String(), s)
			os.Exit(1)
		}
		ops := make([]int, n)
		for i, p := range parts {
			v, err := strconv.Atoi(p)
			if err != nil || (v != 0 && v != 1) {
				fmt.Fprintf(os.Stderr, "case %d: invalid value %q at position %d\ninput: %s\n", idx+1, p, i, s)
				os.Exit(1)
			}
			ops[i] = v
		}
		got := simulate(s, ops)
		want := optimal(s)
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d: resulting string %q != optimal %q\ninput: %s\nops: %v\n", idx+1, got, want, s, ops)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
