package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type test struct {
	input    string
	expected string
}

func solve(a, b []int) int {
	maxA := a[0]
	minA := a[0]
	for _, v := range a {
		if v > maxA {
			maxA = v
		}
		if v < minA {
			minA = v
		}
	}
	minB := b[0]
	for _, v := range b {
		if v < minB {
			minB = v
		}
	}
	v := maxA
	if 2*minA > v {
		v = 2 * minA
	}
	if v < minB {
		return v
	}
	return -1
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(42))
	var tests []test
	for len(tests) < 100 {
		n := rng.Intn(10) + 1
		m := rng.Intn(10) + 1
		a := make([]int, n)
		b := make([]int, m)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(100) + 1
		}
		for i := 0; i < m; i++ {
			b[i] = rng.Intn(100) + 1
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for i, v := range b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		exp := fmt.Sprintf("%d", solve(a, b))
		tests = append(tests, test{sb.String(), exp})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
