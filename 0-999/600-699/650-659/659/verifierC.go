package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solve(input string) string {
	var n int
	var m int64
	fmt.Sscan(strings.SplitN(input, "\n", 2)[0], &n, &m)
	parts := strings.Fields(strings.SplitN(input, "\n", 2)[1])
	have := make(map[int64]bool, n)
	for i := 0; i < n; i++ {
		var x int64
		fmt.Sscan(parts[i], &x)
		have[x] = true
	}
	var res []int64
	var cost int64
	for t := int64(1); cost+t <= m; t++ {
		if !have[t] {
			res = append(res, t)
			cost += t
		}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(res))
	for i, v := range res {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", v)
	}
	if len(res) > 0 {
		b.WriteByte('\n')
	}
	return b.String()
}

func generateTests() []string {
	rand.Seed(44)
	tests := make([]string, 100)
	for t := 0; t < 100; t++ {
		n := rand.Intn(20) + 1
		m := rand.Int63n(1000) + int64(rand.Intn(10))
		have := make(map[int64]bool)
		vals := make([]int64, 0, n)
		for len(vals) < n {
			x := rand.Int63n(1000) + 1
			if !have[x] {
				have[x] = true
				vals = append(vals, x)
			}
		}
		var b strings.Builder
		fmt.Fprintf(&b, "%d %d\n", n, m)
		for i, v := range vals {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
		tests[t] = b.String()
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		expect := strings.TrimSpace(solve(t))
		got, err := runBinary(bin, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if expect != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, t, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
