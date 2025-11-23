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

// nx returns the index of the last child of the node at index x
func nx(x int) int {
	if x <= 0 {
		return 0
	}
	// Count powers of 2 less than or equal to x
	c := 0
	for p := 1; p <= x; p *= 2 {
		c++
	}
	return x + c
}

func solveRef(input string) string {
	r := strings.NewReader(input)
	var n, m int
	if _, err := fmt.Fscan(r, &n, &m); err != nil {
		return ""
	}
	
type Op struct {
		l, r, x int
	}
	// Store ops by level
	ops := make([][]Op, n+1)
	
	var sb strings.Builder
	
	for k := 0; k < m; k++ {
		var type_ int
		fmt.Fscan(r, &type_)
		if type_ == 1 {
			var t, l, r_, x int
			fmt.Fscan(r, &t, &l, &r_, &x)
			ops[t] = append(ops[t], Op{l, r_, x})
		} else {
			var t, v int
			fmt.Fscan(r, &t, &v)
			
			curL, curR := v, v
			
			seen := make(map[int]struct{})
			
			for i := t; i <= n; i++ {
				// Check ops at this level
				for _, op := range ops[i] {
					// Check intersection
					if max(op.l, curL) <= min(op.r, curR) {
						seen[op.x] = struct{}{}
					}
				}
				
				// Move to next level
				// L -> first child of L
				// R -> last child of R
				if i < n {
					curL = nx(curL - 1) + 1
					curR = nx(curR)
				}
			}
			fmt.Fprintln(&sb, len(seen))
		}
	}
	return sb.String()
}

func min(a, b int) int { if a < b { return a }; return b }
func max(a, b int) int { if a > b { return a }; return b }

func genCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1 // 1 to 10 levels
	m := rng.Intn(20) + 1 // 1 to 20 ops
	
	// Precompute counts
	cnt := make([]int, n+1)
	cnt[1] = 1
	for i := 1; i < n; i++ {
		cnt[i+1] = nx(cnt[i])
	}
	
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	
	for i := 0; i < m; i++ {
		kind := rng.Intn(2) + 1
		t := rng.Intn(n) + 1
		if kind == 1 {
			// 1 t l r x
			c := cnt[t]
			l := rng.Intn(c) + 1
			r := rng.Intn(c) + 1
			if l > r {
				l, r = r, l
			}
			x := rng.Intn(20) + 1
			fmt.Fprintf(&sb, "1 %d %d %d %d\n", t, l, r, x)
		} else {
			// 2 t v
			c := cnt[t]
			v := rng.Intn(c) + 1
			fmt.Fprintf(&sb, "2 %d %d\n", t, v)
		}
	}
	return sb.String()
}

func runAndVerify(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execution failed: %v\nOutput: %s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("mismatch\nInput:\n%s\nExpected:\n%s\nGot:\n%s", input, exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go ./solution")
		os.Exit(1)
	}
	bin := os.Args[1]
	
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		expected := solveRef(input)
		if err := runAndVerify(bin, input, expected); err != nil {
			fmt.Printf("Test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}