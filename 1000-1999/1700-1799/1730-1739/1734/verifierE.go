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

func genTests() []string {
	rand.Seed(5)
	tests := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(8) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(rand.Intn(10)))
		}
		sb.WriteByte('\n')
		tests = append(tests, sb.String())
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	tests := genTests()
	for i, input := range tests {
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := validateOutput(input, got); err != nil {
			fmt.Printf("Test %d failed\nInput:\n%sError: %v\nOutput:\n%s\n", i+1, input, err, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func runExe(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// validateOutput checks that candidate matrix satisfies constraints:
// - entries are in [0, n)
// - diagonal equals b mod n
// - anti-rectangle property: for all r1<r2, c1<c2: a[r1,c1]+a[r2,c2] != a[r1,c2]+a[r2,c1] (mod n)
func validateOutput(input, output string) error {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) < 2 {
		return fmt.Errorf("malformed input")
	}
	n, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("bad n: %v", err)
	}
	bfields := strings.Fields(lines[1])
	if len(bfields) != n {
		return fmt.Errorf("expected %d b's, got %d", n, len(bfields))
	}
	b := make([]int, n)
	for i := 0; i < n; i++ {
		bi, _ := strconv.Atoi(bfields[i])
		b[i] = ((bi % n) + n) % n
	}

	toks := strings.Fields(strings.TrimSpace(output))
	if len(toks) != n*n {
		return fmt.Errorf("expected %d numbers, got %d", n*n, len(toks))
	}
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, n)
	}
	k := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v, err := strconv.Atoi(toks[k])
			k++
			if err != nil {
				return fmt.Errorf("non-integer token: %v", err)
			}
			if v < 0 || v >= n {
				return fmt.Errorf("value out of range at (%d,%d): %d", i+1, j+1, v)
			}
			a[i][j] = v
		}
	}
	// diagonal check modulo n
	for i := 0; i < n; i++ {
		if a[i][i]%n != b[i] {
			return fmt.Errorf("diagonal mismatch at %d: have %d want %d (mod %d)", i+1, a[i][i], b[i], n)
		}
	}
	// anti-rectangle property
	for r1 := 0; r1 < n; r1++ {
		for r2 := r1 + 1; r2 < n; r2++ {
			for c1 := 0; c1 < n; c1++ {
				for c2 := c1 + 1; c2 < n; c2++ {
					lhs := (a[r1][c1] + a[r2][c2]) % n
					rhs := (a[r1][c2] + a[r2][c1]) % n
					if lhs == rhs {
						return fmt.Errorf("rectangle constraint violated at r1=%d r2=%d c1=%d c2=%d", r1+1, r2+1, c1+1, c2+1)
					}
				}
			}
		}
	}
	return nil
}
