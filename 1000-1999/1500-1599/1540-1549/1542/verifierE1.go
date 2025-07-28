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
	input    string
	expected int64
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func solve(n int, mod int64) int64 {
	maxD := n * (n - 1) / 2
	offset := maxD
	size := 2*maxD + 1
	dp0 := make([]int64, size)
	dp1 := make([]int64, size)
	dp0[offset] = 1
	for m := n; m >= 1; m-- {
		next0 := make([]int64, size)
		next1 := make([]int64, size)
		for d := -maxD; d <= maxD; d++ {
			idx := d + offset
			if idx < 0 || idx >= size {
				continue
			}
			v0 := dp0[idx]
			if v0 != 0 {
				next0[idx] = (next0[idx] + v0*int64(m)) % mod
				for k := 1; k < m; k++ {
					delta := -k
					idx2 := d + delta + offset
					if idx2 >= 0 && idx2 < size {
						next1[idx2] = (next1[idx2] + v0*int64(m-k)) % mod
					}
				}
			}
			v1 := dp1[idx]
			if v1 != 0 {
				for delta := -m + 1; delta <= m-1; delta++ {
					idx2 := d + delta + offset
					if idx2 >= 0 && idx2 < size {
						cnt := m - absInt(delta)
						next1[idx2] = (next1[idx2] + v1*int64(cnt)) % mod
					}
				}
			}
		}
		dp0, dp1 = next0, next1
	}
	var ans int64
	for d := 1; d <= maxD; d++ {
		ans = (ans + dp1[d+offset]) % mod
	}
	return ans % mod
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	mod := int64(rng.Intn(1_000_000_000-1) + 2)
	input := fmt.Sprintf("%d %d\n", n, mod)
	return testCase{input: input, expected: solve(n, mod)}
}

func runCase(bin string, tc testCase) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != tc.expected {
		return fmt.Errorf("expected %d got %d", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []testCase{
		{input: "1 2\n", expected: solve(1, 2)},
		{input: "4 1000000007\n", expected: solve(4, 1000000007)},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for idx, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
