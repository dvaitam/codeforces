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

// Embedded solver (same logic as the accepted solution cf_t24_1667_C.go)

func oracleSolve(input string) string {
	n, _ := strconv.Atoi(strings.TrimSpace(input))

	if n == 1 {
		return "1\n1 1"
	}

	k := (2*n + 1) / 3
	m := n - k

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", k)

	if m > 0 {
		for i := 1; i <= m/2; i++ {
			a := i
			b := m - i + 1
			fmt.Fprintf(&sb, "%d %d\n", a, b)
			fmt.Fprintf(&sb, "%d %d\n", b, a)
		}

		for i := 1; i <= (m-1)/2; i++ {
			a := m + i
			b := 2*m - i
			fmt.Fprintf(&sb, "%d %d\n", a, b)
			fmt.Fprintf(&sb, "%d %d\n", b, a)
		}

		var rem int
		if m%2 == 0 {
			rem = m + m/2
		} else {
			rem = (m + 1) / 2
		}
		fmt.Fprintf(&sb, "%d %d\n", rem, rem)
	}

	for x := 2 * m; x <= k; x++ {
		fmt.Fprintf(&sb, "%d %d\n", x, x)
	}

	return strings.TrimSpace(sb.String())
}

// Checker: verify that the placement is valid for half-queens on n x n board
func check(n int, output string) error {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}
	k, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("invalid k: %v", err)
	}
	expectedK := (2*n + 1) / 3
	if k != expectedK {
		return fmt.Errorf("expected k=%d, got k=%d", expectedK, k)
	}
	if len(lines) != k+1 {
		return fmt.Errorf("expected %d position lines, got %d", k, len(lines)-1)
	}

	// A half-queen on (a,b) attacks: same row, same column, and diagonal a-b=c-d
	coveredRows := make(map[int]bool)
	coveredCols := make(map[int]bool)
	coveredDiag := make(map[int]bool)

	for i := 1; i <= k; i++ {
		parts := strings.Fields(lines[i])
		if len(parts) != 2 {
			return fmt.Errorf("line %d: expected 2 numbers, got %d", i, len(parts))
		}
		a, err1 := strconv.Atoi(parts[0])
		b, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("line %d: invalid numbers", i)
		}
		if a < 1 || a > n || b < 1 || b > n {
			return fmt.Errorf("line %d: position (%d,%d) out of range [1,%d]", i, a, b, n)
		}
		coveredRows[a] = true
		coveredCols[b] = true
		coveredDiag[a-b] = true
	}

	// Check all cells are covered
	for r := 1; r <= n; r++ {
		if coveredRows[r] {
			continue
		}
		for c := 1; c <= n; c++ {
			if coveredCols[c] {
				continue
			}
			if coveredDiag[r-c] {
				continue
			}
			return fmt.Errorf("cell (%d,%d) not covered", r, c)
		}
	}
	return nil
}

func runBin(path, in string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	user := os.Args[1]

	// Test specific values and random ones
	tests := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < 90; i++ {
		tests = append(tests, rand.Intn(100)+1)
	}

	for i, n := range tests {
		input := fmt.Sprintf("%d\n", n)
		out, err := runBin(user, input)
		if err != nil {
			fmt.Printf("runtime error on test %d (n=%d): %v\n", i+1, n, err)
			os.Exit(1)
		}
		if err := check(n, out); err != nil {
			fmt.Printf("wrong answer on test %d (n=%d): %v\noutput:\n%s\n", i+1, n, err, out)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
