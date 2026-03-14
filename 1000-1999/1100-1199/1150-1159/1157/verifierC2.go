package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const numTestsC2 = 100

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: verifierC2.go <binary>")
		os.Exit(1)
	}
	binPath, cleanup, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		os.Exit(1)
	}
	if cleanup != nil {
		defer cleanup()
	}
	r := rand.New(rand.NewSource(1))
	for t := 1; t <= numTestsC2; t++ {
		n := r.Intn(20) + 1
		a := make([]int, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			a[i] = r.Intn(100) + 1
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		optimalLen := solveC2Len(a)
		out, err := run(binPath, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		candK, candMoves, err := parseOutput(out)
		if err != nil {
			fmt.Printf("test %d: parse error: %v\noutput: %s\n", t, err, out)
			os.Exit(1)
		}
		if candK != len(candMoves) {
			fmt.Printf("test %d: declared length %d but moves string has length %d\n", t, candK, len(candMoves))
			os.Exit(1)
		}
		if err := validateMoves(a, candMoves); err != nil {
			fmt.Printf("test %d: invalid moves: %v\ninput: %soutput: %s\n", t, err, input, out)
			os.Exit(1)
		}
		if candK != optimalLen {
			fmt.Printf("test %d: suboptimal length %d, expected %d\ninput: %s", t, candK, optimalLen, input)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verify_binC2")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, string(out))
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return buf.String(), err
}

func parseOutput(s string) (int, string, error) {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	if len(lines) == 0 {
		return 0, "", fmt.Errorf("empty output")
	}
	var k int
	if _, err := fmt.Sscanf(lines[0], "%d", &k); err != nil {
		return 0, "", fmt.Errorf("cannot parse k: %v", err)
	}
	moves := ""
	if len(lines) > 1 {
		moves = strings.TrimSpace(lines[1])
	}
	return k, moves, nil
}

func validateMoves(a []int, moves string) error {
	n := len(a)
	l, r := 0, n-1
	last := 0
	for i, ch := range moves {
		var val int
		switch ch {
		case 'L':
			if l > r {
				return fmt.Errorf("move %d: no elements left", i)
			}
			val = a[l]
			l++
		case 'R':
			if l > r {
				return fmt.Errorf("move %d: no elements left", i)
			}
			val = a[r]
			r--
		default:
			return fmt.Errorf("move %d: invalid character '%c'", i, ch)
		}
		if val <= last {
			return fmt.Errorf("move %d: value %d not strictly greater than previous %d", i, val, last)
		}
		last = val
	}
	return nil
}

// solveC2Len returns the optimal (maximum) number of elements that can be taken.
// Uses DP/greedy: at each step pick the smaller available end if both are valid,
// trying both sides when equal and returning the better result.
func solveC2Len(a []int) int {
	n := len(a)
	// Use a recursive approach with memoization for correctness.
	type key struct{ l, r, last int }
	memo := make(map[key]int)
	var solve func(l, r, last int) int
	solve = func(l, r, last int) int {
		if l > r {
			return 0
		}
		k := key{l, r, last}
		if v, ok := memo[k]; ok {
			return v
		}
		best := 0
		if a[l] > last {
			v := 1 + solve(l+1, r, a[l])
			if v > best {
				best = v
			}
		}
		if a[r] > last {
			v := 1 + solve(l, r-1, a[r])
			if v > best {
				best = v
			}
		}
		memo[k] = best
		return best
	}
	return solve(0, n-1, 0)
}
