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

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	mat := make([][]int, n)
	cur := 0
	for i := 0; i < n; i++ {
		mat[i] = make([]int, n)
		for j := 0; j < n; j++ {
			mat[i][j] = rng.Intn(2)
			if i == j && mat[i][j] == 1 {
				cur ^= 1
			}
		}
	}
	q := rng.Intn(50) + 1
	var in strings.Builder
	fmt.Fprintf(&in, "%d\n", n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				in.WriteByte(' ')
			}
			fmt.Fprintf(&in, "%d", mat[i][j])
		}
		in.WriteByte('\n')
	}
	fmt.Fprintf(&in, "%d\n", q)
	var out strings.Builder
	for t := 0; t < q; t++ {
		typ := rng.Intn(3) + 1
		if typ == 3 {
			fmt.Fprintln(&in, 3)
			out.WriteByte(byte('0' + cur))
		} else {
			idx := rng.Intn(n) + 1
			fmt.Fprintf(&in, "%d %d\n", typ, idx)
			cur ^= 1
		}
	}
	out.WriteByte('\n')
	return in.String(), out.String()
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	got := strings.TrimSpace(buf.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, buf.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
