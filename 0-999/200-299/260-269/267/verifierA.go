package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func ops(a, b int64) int64 {
	var res int64
	for a > 0 && b > 0 {
		if a >= b {
			res += a / b
			a %= b
		} else {
			res += b / a
			b %= a
		}
	}
	return res
}

func generate() (string, string) {
	const T = 150
	var in strings.Builder
	var out strings.Builder
	fmt.Fprintf(&in, "%d\n", T)
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < T; i++ {
		a := rng.Int63n(1_000_000_000) + 1
		b := rng.Int63n(1_000_000_000) + 1
		fmt.Fprintf(&in, "%d %d\n", a, b)
		fmt.Fprintf(&out, "%d\n", ops(a, b))
	}
	return in.String(), out.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	in, exp := generate()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("error running binary:", err)
		os.Exit(1)
	}
	got := strings.TrimSpace(buf.String())
	if got != strings.TrimSpace(exp) {
		fmt.Println("wrong answer")
		fmt.Println("expected:\n" + exp)
		fmt.Println("got:\n" + buf.String())
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
