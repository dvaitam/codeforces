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

type testCase struct {
	a []int
	L int
	R int
}

func isMonotoneTriple(a, b, c int) bool {
	return (a <= b && b <= c) || (a >= b && b >= c)
}

func checkSubsequence(vals []int, idxs []int) bool {
	m := len(idxs)
	for i := 0; i < m; i++ {
		for j := i + 1; j < m; j++ {
			for k := j + 1; k < m; k++ {
				if isMonotoneTriple(vals[idxs[i]], vals[idxs[j]], vals[idxs[k]]) {
					return false
				}
			}
		}
	}
	return true
}

func solveCase(tc testCase) (int, []int) {
	vals := tc.a[tc.L : tc.R+1]
	n := len(vals)
	// search length 4 then 3
	for l := 4; l >= 3; l-- {
		idxs := make([]int, l)
		var dfs func(pos, last int) []int
		dfs = func(pos, last int) []int {
			if pos == l {
				if checkSubsequence(vals, idxs) {
					res := make([]int, l)
					for i, v := range idxs {
						res[i] = v + tc.L + 1
					}
					return res
				}
				return nil
			}
			for i := last; i < n; i++ {
				idxs[pos] = i
				if r := dfs(pos+1, i+1); r != nil {
					return r
				}
			}
			return nil
		}
		if res := dfs(0, 0); res != nil {
			return l, res
		}
	}
	return 0, nil
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 3 // 3..8
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(21) - 10
	}
	L := rng.Intn(n - 2)
	R := L + 2 + rng.Intn(n-L-2)
	return testCase{a, L, R}
}

func runCase(bin string, tc testCase) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", len(tc.a), 1))
	for i, v := range tc.a {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprintf("%d", v))
	}
	input.WriteString(fmt.Sprintf("\n%d %d\n", tc.L+1, tc.R+1))

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	k, idxs := solveCase(tc)
	outFields := strings.Fields(strings.TrimSpace(out.String()))
	if k == 0 {
		if len(outFields) != 1 || outFields[0] != "0" {
			return fmt.Errorf("expected 0 got %v", outFields)
		}
		return nil
	}
	if len(outFields) != k+1 {
		return fmt.Errorf("expected %d numbers got %d", k+1, len(outFields))
	}
	if outFields[0] != fmt.Sprintf("%d", k) {
		return fmt.Errorf("expected length %d got %s", k, outFields[0])
	}
	for i := 0; i < k; i++ {
		var v int
		if _, err := fmt.Sscan(outFields[i+1], &v); err != nil {
			return fmt.Errorf("bad index")
		}
		if v != idxs[i] {
			return fmt.Errorf("expected %v got %v", idxs, outFields[1:])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
