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

type pair struct{ a, b int }

type test struct {
	input          string
	expectedCount  int
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solve(input string) int {
	var n int
	fmt.Sscan(strings.TrimSpace(input), &n)
	used := make([]bool, n+1)
	count := 0
	var p int
	for i := 3; i <= n/2; i += 2 {
		if !used[i] {
			p = i
			for j := i * 3; j <= n; j += i {
				if !used[j] {
					if p != 0 {
						count++
						used[p], used[j] = true, true
						p = 0
					} else {
						p = j
					}
				}
			}
			if p != 0 {
				pairv := i * 2
				if pairv <= n && !used[p] && !used[pairv] {
					count++
					used[p], used[pairv] = true, true
				}
			}
		}
	}
	p = 0
	for i := 2; i <= n; i += 2 {
		if !used[i] {
			if p != 0 {
				count++
				used[p], used[i] = true, true
				p = 0
			} else {
				p = i
			}
		}
	}
	return count
}

func generateTests() []test {
	rand.Seed(451)
	var tests []test
	fixed := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 20, 30, 50}
	for _, n := range fixed {
		inp := fmt.Sprintf("%d\n", n)
		tests = append(tests, test{inp, solve(inp)})
	}
	for len(tests) < 100 {
		n := rand.Intn(100) + 1
		inp := fmt.Sprintf("%d\n", n)
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func validateOutput(output string, n int, expectedCount int) error {
	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}
	m, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("cannot parse number of pairs: %v", err)
	}
	if m != expectedCount {
		return fmt.Errorf("expected %d pairs, got %d", expectedCount, m)
	}
	if len(lines) != m+1 {
		return fmt.Errorf("expected %d pair lines, got %d", m, len(lines)-1)
	}
	used := make(map[int]bool)
	for i := 1; i <= m; i++ {
		parts := strings.Fields(strings.TrimSpace(lines[i]))
		if len(parts) != 2 {
			return fmt.Errorf("pair line %d: expected 2 numbers, got %d", i, len(parts))
		}
		a, err1 := strconv.Atoi(parts[0])
		b, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("pair line %d: cannot parse numbers", i)
		}
		if a < 1 || a > n || b < 1 || b > n {
			return fmt.Errorf("pair line %d: values %d, %d out of range [1, %d]", i, a, b, n)
		}
		if a == b {
			return fmt.Errorf("pair line %d: duplicate value %d", i, a)
		}
		if gcd(a, b) <= 1 {
			return fmt.Errorf("pair line %d: gcd(%d, %d) = %d, must be > 1", i, a, b, gcd(a, b))
		}
		if used[a] {
			return fmt.Errorf("pair line %d: value %d already used", i, a)
		}
		if used[b] {
			return fmt.Errorf("pair line %d: value %d already used", i, b)
		}
		used[a] = true
		used[b] = true
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		var n int
		fmt.Sscan(strings.TrimSpace(t.input), &n)
		if verr := validateOutput(got, n, t.expectedCount); verr != nil {
			fmt.Printf("Wrong answer on test %d\nInput: %sError: %v\nGot:\n%s\n", i+1, t.input, verr, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
