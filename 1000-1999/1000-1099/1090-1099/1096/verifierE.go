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

const modE = 998244353
const M = 5220
const N = 120

var comb [M][N]int

func initComb() {
	for i := 0; i < M; i++ {
		comb[i][0] = 1
		for j := 1; j < N && j <= i; j++ {
			comb[i][j] = comb[i-1][j] + comb[i-1][j-1]
			if comb[i][j] >= modE {
				comb[i][j] -= modE
			}
		}
	}
}

func modPowE(a, b int) int {
	res := 1
	for b > 0 {
		if b&1 == 1 {
			res = int((int64(res) * int64(a)) % modE)
		}
		a = int((int64(a) * int64(a)) % modE)
		b >>= 1
	}
	return res
}

func solveE(p, sum, lower int) int {
	if lower == 0 {
		return modPowE(p, modE-2)
	}
	totalWays := comb[sum-lower+p-1][p-1]
	invTotal := modPowE(totalWays, modE-2)
	ans := 0
	for ties := 1; ties <= p; ties++ {
		save := modPowE(ties, modE-2)
		for top := lower; ties*top <= sum; top++ {
			rem := sum - ties*top
			t := p - ties
			if t == 0 {
				if rem == 0 {
					ans += save
					if ans >= modE {
						ans -= modE
					}
				}
				continue
			}
			res := 0
			for bad := 0; bad <= t && bad*top <= rem; bad++ {
				nsum := rem - bad*top
				add := int((int64(comb[nsum+t-1][t-1]) * int64(comb[t][bad])) % modE)
				if bad&1 == 1 {
					res = (res - add + modE) % modE
				} else {
					res = (res + add) % modE
				}
			}
			res = int((int64(res) * int64(comb[p-1][ties-1])) % modE)
			ans = int((int64(ans) + int64(res)*int64(save)) % modE)
		}
	}
	ans = int((int64(ans) * int64(invTotal)) % modE)
	return ans
}

type testCaseE struct {
	input    string
	expected int
}

func generateCaseE(rng *rand.Rand) testCaseE {
	p := rng.Intn(5) + 1
	lower := rng.Intn(5)
	sum := lower + rng.Intn(10)
	input := fmt.Sprintf("%d %d %d\n", p, sum, lower)
	return testCaseE{input: input, expected: solveE(p, sum, lower)}
}

func runCaseE(bin string, tc testCaseE) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	initComb()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	cases := []testCaseE{generateCaseE(rng)}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCaseE(rng))
	}
	for i, tc := range cases {
		if err := runCaseE(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
