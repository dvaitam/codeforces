package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type test struct {
	n  int64
	q  int
	xs []int64
}

func genTests() []test {
	rand.Seed(2)
	tests := make([]test, 100)
	for i := range tests {
		n := rand.Int63n(1_000_000) + 1
		q := rand.Intn(5) + 1
		xs := make([]int64, q)
		for j := 0; j < q; j++ {
			xs[j] = rand.Int63n(n) + 1
		}
		tests[i] = test{n, q, xs}
	}
	return tests
}

func runBinary(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", path)
	} else {
		cmd = exec.CommandContext(ctx, path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func solve(n, x int64) int64 {
	if x%2 == 1 {
		return (x + 1) / 2
	}
	if n%2 == 0 {
		return n/2 + solve(n/2, x/2)
	}
	m := n / 2
	i := x / 2
	if i == 1 {
		return m + 1 + solve(m, m)
	}
	return m + 1 + solve(m, i-1)
}

func expected(t test) []int64 {
	res := make([]int64, t.q)
	for i, x := range t.xs {
		res[i] = solve(t.n, x)
	}
	return res
}

func buildInput(t test) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.q))
	for _, x := range t.xs {
		sb.WriteString(fmt.Sprintf("%d\n", x))
	}
	return sb.String()
}

func parseOutput(out string, q int) ([]int64, bool) {
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(out)))
	scanner.Split(bufio.ScanWords)
	vals := make([]int64, 0, q)
	for scanner.Scan() {
		v, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			return nil, false
		}
		vals = append(vals, v)
	}
	if len(vals) != q {
		return nil, false
	}
	return vals, true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		input := buildInput(t)
		out, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got, ok := parseOutput(out, t.q)
		if !ok {
			fmt.Printf("test %d: could not parse output\n", i+1)
			os.Exit(1)
		}
		exp := expected(t)
		match := true
		for j := 0; j < t.q; j++ {
			if got[j] != exp[j] {
				match = false
				break
			}
		}
		if !match {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%v\noutput:%v\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
