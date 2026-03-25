package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

// solveC computes the maximum spiral sum using the same logic as the accepted solution.
// A spiral starts at center (cr, cc) and expands outward in rings of radius 2, 4, ...
// For odd-sized spirals (3x3, 7x7, ...) the center is a single cell.
// For even-expansion spirals (starting from 3x3 ring without center), they start at radius 1.
// The spiral is NOT a full square; it's a border-walking spiral that excludes one cell per ring.
func solveC(n, m int, A [][]int64) int64 {
	pref := make([][]int64, n+1)
	for i := 0; i <= n; i++ {
		pref[i] = make([]int64, m+1)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			pref[i+1][j+1] = pref[i][j+1] + pref[i+1][j] - pref[i][j] + A[i][j]
		}
	}

	getSum := func(r1, c1, r2, c2 int) int64 {
		return pref[r2+1][c2+1] - pref[r1][c2+1] - pref[r2+1][c1] + pref[r1][c1]
	}

	ans := int64(-1e18)

	// Type 1: spirals with a single-cell center (sizes 1, 5, 9, ... -> radius 0, 2, 4, ...)
	for cr := 0; cr < n; cr++ {
		for cc := 0; cc < m; cc++ {
			sum := A[cr][cc]
			rad := 2
			k := 5
			for cr-rad >= 0 && cr+rad < n && cc-rad >= 0 && cc+rad < m {
				r := cr - rad
				c := cc - rad
				border := getSum(r, c, r+k-1, c+k-1) - getSum(r+1, c+1, r+k-2, c+k-2)
				sum += border - A[r+1][c] + A[r+2][c+1]
				if sum > ans {
					ans = sum
				}
				rad += 2
				k += 4
			}
		}
	}

	// Type 2: spirals starting from a 3x3 ring (no single-cell center), sizes 3, 7, 11, ... -> radius 1, 3, 5, ...
	for cr := 0; cr < n; cr++ {
		for cc := 0; cc < m; cc++ {
			rad := 1
			k := 3
			if cr-rad >= 0 && cr+rad < n && cc-rad >= 0 && cc+rad < m {
				r := cr - rad
				c := cc - rad
				border := getSum(r, c, r+k-1, c+k-1) - getSum(r+1, c+1, r+k-2, c+k-2)
				sum := border - A[r+1][c]
				if sum > ans {
					ans = sum
				}

				rad += 2
				k += 4

				for cr-rad >= 0 && cr+rad < n && cc-rad >= 0 && cc+rad < m {
					r := cr - rad
					c := cc - rad
					border := getSum(r, c, r+k-1, c+k-1) - getSum(r+1, c+1, r+k-2, c+k-2)
					sum += border - A[r+1][c] + A[r+2][c+1]
					if sum > ans {
						ans = sum
					}
					rad += 2
					k += 4
				}
			}
		}
	}

	return ans
}

func runCase(bin string, n, m int, a [][]int64) error {
	var input bytes.Buffer
	fmt.Fprintf(&input, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", a[i][j])
		}
		input.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(&out, &got); err != nil {
		return fmt.Errorf("parse error: %v", err)
	}
	want := solveC(n, m, a)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	const tests = 100
	for t := 0; t < tests; t++ {
		n := rng.Intn(8) + 3
		m := rng.Intn(8) + 3
		a := make([][]int64, n)
		for i := 0; i < n; i++ {
			row := make([]int64, m)
			for j := 0; j < m; j++ {
				row[j] = int64(rng.Intn(2001) - 1000)
			}
			a[i] = row
		}
		if err := runCase(bin, n, m, a); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
