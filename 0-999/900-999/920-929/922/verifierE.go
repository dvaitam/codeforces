package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseE struct {
	n    int
	W    int64
	B    int64
	X    int64
	c    []int
	cost []int64
}

func expectedE(tc testCaseE) int {
	n := tc.n
	W := tc.W
	B := tc.B
	X := tc.X
	c := tc.c
	cost := tc.cost
	sum := 0
	for _, v := range c {
		sum += v
	}
	const neg = int64(-1 << 60)
	dp := make([]int64, sum+1)
	for i := 1; i <= sum; i++ {
		dp[i] = neg
	}
	dp[0] = W
	maxBirds := 0
	for i := 0; i < n; i++ {
		if i > 0 {
			for j := 0; j <= maxBirds; j++ {
				if dp[j] < 0 {
					continue
				}
				cap := W + int64(j)*B
				val := dp[j] + X
				if val > cap {
					val = cap
				}
				dp[j] = val
			}
		}
		newdp := make([]int64, sum+1)
		for k := 0; k <= sum; k++ {
			newdp[k] = neg
		}
		for j := 0; j <= maxBirds; j++ {
			if dp[j] < 0 {
				continue
			}
			maxT := c[i]
			if maxT+j > sum {
				maxT = sum - j
			}
			for t := 0; t <= maxT; t++ {
				req := int64(t) * cost[i]
				if req > dp[j] {
					break
				}
				val := dp[j] - req
				if val > newdp[j+t] {
					newdp[j+t] = val
				}
			}
		}
		maxBirds += c[i]
		dp = newdp
	}
	ans := 0
	for j := maxBirds; j >= 0; j-- {
		if dp[j] >= 0 {
			ans = j
			break
		}
	}
	return ans
}

func genTestsE() []testCaseE {
	rand.Seed(5)
	tests := make([]testCaseE, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(5) + 1
		W := int64(rand.Intn(20))
		B := int64(rand.Intn(10))
		X := int64(rand.Intn(10))
		c := make([]int, n)
		cost := make([]int64, n)
		for i := 0; i < n; i++ {
			c[i] = rand.Intn(5)
			cost[i] = int64(rand.Intn(10))
		}
		tests = append(tests, testCaseE{n: n, W: W, B: B, X: X, c: c, cost: cost})
	}
	return tests
}

func runCase(bin string, tc testCaseE) error {
	var input strings.Builder
	fmt.Fprintf(&input, "%d %d %d %d\n", tc.n, tc.W, tc.B, tc.X)
	for i, v := range tc.c {
		if i > 0 {
			input.WriteByte(' ')
		}
		fmt.Fprint(&input, v)
	}
	input.WriteByte('\n')
	for i, v := range tc.cost {
		if i > 0 {
			input.WriteByte(' ')
		}
		fmt.Fprint(&input, v)
	}
	input.WriteByte('\n')

	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotStr := strings.TrimSpace(out.String())
	expect := fmt.Sprint(expectedE(tc))
	if gotStr != expect {
		return fmt.Errorf("expected %s got %s", expect, gotStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsE()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
