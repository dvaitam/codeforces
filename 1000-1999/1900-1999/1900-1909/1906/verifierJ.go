package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod = 998244353

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierJ.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	inBytes, err := readAllStdin()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read stdin:", err)
		os.Exit(1)
	}
	n, a, err := parseInput(inBytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid input:", err)
		os.Exit(1)
	}
	want := solveExpected(n, a)
	candOut, err := runProgram(candidate, inBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	if err := compareAnswer(want, candOut); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "expected:")
		fmt.Fprintln(os.Stderr, want)
		fmt.Fprintln(os.Stderr, "candidate output:")
		fmt.Fprintln(os.Stderr, candOut)
		os.Exit(1)
	}

	fmt.Println("Accepted")
}

func readAllStdin() ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(os.Stdin)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func parseInput(input []byte) (int, []int, error) {
	fields := strings.Fields(string(input))
	if len(fields) < 1 {
		return 0, nil, fmt.Errorf("missing n")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil || n < 2 {
		return 0, nil, fmt.Errorf("invalid n")
	}
	if len(fields) != n+1 {
		return 0, nil, fmt.Errorf("expected %d permutation values, got %d", n, len(fields)-1)
	}
	a := make([]int, n+1)
	seen := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		v, err := strconv.Atoi(fields[i])
		if err != nil || v < 1 || v > n {
			return 0, nil, fmt.Errorf("invalid a[%d]", i)
		}
		if seen[v] {
			return 0, nil, fmt.Errorf("a is not a permutation")
		}
		seen[v] = true
		a[i] = v
	}
	if a[1] != 1 {
		return 0, nil, fmt.Errorf("a1 must be 1")
	}
	return n, a, nil
}

func add(x *int, y int) {
	*x += y
	if *x >= mod {
		*x -= mod
	}
}

func solveExpected(n int, a []int) int {
	f := make([]int, n+1)
	f[1] = 1
	for i := 1; i < n; i++ {
		p, s := 1, 0
		for j := i; j <= n; j++ {
			x := int((int64(f[j]) * int64(p)) % mod)
			f[j] = x
			p = int((int64(p) * 2) % mod)
			add(&f[j], s)
			if j < n && a[j+1] < a[j] {
				s = 0
			}
			add(&s, x)
		}
	}
	return f[n]
}

func runProgram(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	return out.String(), cmd.Run()
}

func compareAnswer(expected int, got string) error {
	fields := strings.Fields(got)
	if len(fields) == 0 {
		return fmt.Errorf("empty candidate output")
	}
	v, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("first output token is not an integer")
	}
	if v != expected {
		return fmt.Errorf("wrong answer: expected %d, got %d", expected, v)
	}
	return nil
}
