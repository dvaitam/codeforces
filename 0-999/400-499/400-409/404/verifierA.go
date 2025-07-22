package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseA struct {
	n    int
	grid []string
}

func generateTestsA() []testCaseA {
	r := rand.New(rand.NewSource(1))
	tests := make([]testCaseA, 0, 100)
	for len(tests) < 100 {
		n := 3 + 2*r.Intn(5) // odd size 3..11
		diag := byte('a' + r.Intn(26))
		other := byte('a' + r.Intn(26))
		for other == diag {
			other = byte('a' + r.Intn(26))
		}
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			row := make([]byte, n)
			for j := 0; j < n; j++ {
				if i == j || i+j == n-1 {
					row[j] = diag
				} else {
					row[j] = other
				}
			}
			grid[i] = string(row)
		}
		// randomly corrupt one cell to create invalid cases
		if r.Intn(2) == 0 {
			i := r.Intn(n)
			j := r.Intn(n)
			row := []byte(grid[i])
			row[j] = byte('a' + r.Intn(26))
			grid[i] = string(row)
		}
		tests = append(tests, testCaseA{n: n, grid: grid})
	}
	return tests
}

func expectedA(t testCaseA) string {
	diag := t.grid[0][0]
	other := t.grid[0][1]
	if diag == other {
		return "NO"
	}
	for i := 0; i < t.n; i++ {
		for j := 0; j < t.n; j++ {
			c := t.grid[i][j]
			if i == j || i+j == t.n-1 {
				if c != diag {
					return "NO"
				}
			} else {
				if c != other {
					return "NO"
				}
			}
		}
	}
	return "YES"
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("run failed: %v\n%s", err, errb.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsA()
	for i, t := range tests {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%d\n", t.n))
		for _, row := range t.grid {
			b.WriteString(row)
			b.WriteByte('\n')
		}
		out, err := runBinary(bin, b.String())
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expected := expectedA(t)
		if strings.TrimSpace(out) != expected {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
