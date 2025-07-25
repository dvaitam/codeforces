package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	in  string
	out string
}

func solveD(a, b, k, t int) int64 {
	mod := int64(1e9 + 7)
	d := b - a
	N := 4*k*t + 1
	base := 2 * k * t
	dp := make([]int64, N)
	newdp := make([]int64, N)
	s0 := make([]int64, N)
	s1 := make([]int64, N)
	dp[base] = 1
	twok1 := int64(2*k + 1)
	for step := 1; step <= t; step++ {
		currLow := base - 2*k*step
		currHigh := base + 2*k*step
		s0[0] = dp[0]
		s1[0] = 0
		for i := 1; i < N; i++ {
			s0[i] = s0[i-1] + dp[i]
			if s0[i] >= mod {
				s0[i] -= mod
			}
			s1[i] = s1[i-1] + dp[i]*int64(i)%mod
			if s1[i] >= mod {
				s1[i] -= mod
			}
		}
		for i := currLow; i <= currHigh; i++ {
			newdp[i] = 0
		}
		for dIdx := currLow; dIdx <= currHigh; dIdx++ {
			L := dIdx - 2*k
			if L < 0 {
				L = 0
			}
			R := dIdx + 2*k
			if R >= N {
				R = N - 1
			}
			var sumPrev int64
			if L > 0 {
				sumPrev = s0[R] - s0[L-1]
			} else {
				sumPrev = s0[R]
			}
			if sumPrev < 0 {
				sumPrev += mod
			}
			Rleft := dIdx
			if Rleft > R {
				Rleft = R
			}
			Lleft := L
			var sumLeft, sumS1Left int64
			if Rleft >= Lleft {
				if Lleft > 0 {
					sumLeft = s0[Rleft] - s0[Lleft-1]
					sumS1Left = s1[Rleft] - s1[Lleft-1]
				} else {
					sumLeft = s0[Rleft]
					sumS1Left = s1[Rleft]
				}
				if sumLeft < 0 {
					sumLeft += mod
				}
				if sumS1Left < 0 {
					sumS1Left += mod
				}
			}
			Lright := dIdx + 1
			if Lright > R {
				newdp[dIdx] = twok1 * sumPrev % mod
			} else {
				sumPrevRight := s0[R] - s0[dIdx]
				sumS1Right := s1[R] - s1[dIdx]
				if sumPrevRight < 0 {
					sumPrevRight += mod
				}
				if sumS1Right < 0 {
					sumS1Right += mod
				}
				part1 := (sumLeft*int64(dIdx)%mod - sumS1Left) % mod
				if part1 < 0 {
					part1 += mod
				}
				part2 := (sumS1Right - sumPrevRight*int64(dIdx)%mod) % mod
				if part2 < 0 {
					part2 += mod
				}
				W1 := part1 + part2
				if W1 >= mod {
					W1 -= mod
				}
				val := twok1*sumPrev%mod - W1
				if val < 0 {
					val += mod
				}
				newdp[dIdx] = val
			}
		}
		dp, newdp = newdp, dp
	}
	s0[0] = dp[0]
	for i := 1; i < N; i++ {
		s0[i] = s0[i-1] + dp[i]
		if s0[i] >= mod {
			s0[i] -= mod
		}
	}
	start := base + d + 1
	if start < 0 {
		start = 0
	}
	if start > N-1 {
		return 0
	}
	res := s0[N-1]
	if start > 0 {
		res = (res - s0[start-1] + mod) % mod
	}
	return res
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(4))
	tests := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		a := rng.Intn(20) + 1
		b := rng.Intn(20) + 1
		k := rng.Intn(5) + 1
		t := rng.Intn(3) + 1
		expect := solveD(a, b, k, t)
		in := fmt.Sprintf("%d %d %d %d\n", a, b, k, t)
		out := fmt.Sprintf("%d\n", expect)
		tests[i] = testCase{in: in, out: out}
	}
	return tests
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := strings.TrimSpace(tc.out)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
