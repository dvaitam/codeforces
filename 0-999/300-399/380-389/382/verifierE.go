package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const MOD int64 = 1000000007

type Test struct {
	in  string
	out string
}

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func modInv(a int64) int64 { return modPow(a, MOD-2) }

func compute(n, K int) int64 {
	maxN := n
	fac := make([]int64, maxN+1)
	ifac := make([]int64, maxN+1)
	fac[0] = 1
	for i := 1; i <= maxN; i++ {
		fac[i] = fac[i-1] * int64(i) % MOD
	}
	ifac[maxN] = modInv(fac[maxN])
	for i := maxN; i > 0; i-- {
		ifac[i-1] = ifac[i] * int64(i) % MOD
	}
	C := func(n, k int) int64 {
		if k < 0 || k > n {
			return 0
		}
		return fac[n] * ifac[k] % MOD * ifac[n-k] % MOD
	}
	inv2 := modInv(2)
	maxM := (n + 1) / 2
	dp := make([][][]int64, n+1)
	for i := range dp {
		dp[i] = make([][]int64, maxM+2)
		for j := range dp[i] {
			dp[i][j] = make([]int64, maxM+2)
		}
	}
	dp[1][0][0] = 1
	for sz := 2; sz <= n; sz++ {
		i := sz - 1
		coeff1 := C(sz-1, i) * int64(i) % MOD
		for g0 := 0; g0 <= i/2; g0++ {
			for g1 := 0; g1 <= g0; g1++ {
				cnt := dp[i][g0][g1]
				if cnt == 0 {
					continue
				}
				parentG1 := g0
				mg0 := g0
				if 1+g1 > mg0 {
					mg0 = 1 + g1
				}
				dp[sz][mg0][parentG1] = (dp[sz][mg0][parentG1] + cnt*coeff1) % MOD
			}
		}
		for i = 1; i <= sz-2; i++ {
			j := sz - 1 - i
			comb := C(sz-1, i)
			var labelCoeff int64
			if i != j {
				labelCoeff = comb * int64(i) % MOD * int64(j) % MOD
			} else {
				labelCoeff = comb * int64(i) % MOD * int64(j) % MOD * inv2 % MOD
			}
			for g0_1 := 0; g0_1 <= i/2; g0_1++ {
				for g1_1 := 0; g1_1 <= g0_1; g1_1++ {
					cnt1 := dp[i][g0_1][g1_1]
					if cnt1 == 0 {
						continue
					}
					for g0_2 := 0; g0_2 <= j/2; g0_2++ {
						for g1_2 := 0; g1_2 <= g0_2; g1_2++ {
							cnt2 := dp[j][g0_2][g1_2]
							if cnt2 == 0 {
								continue
							}
							pairCount := cnt1 * cnt2 % MOD
							ways := pairCount * labelCoeff % MOD
							parentG1 := g0_1 + g0_2
							match1a := 1 + g1_1 + g0_2
							match1b := 1 + g0_1 + g1_2
							mg0 := parentG1
							if match1a > mg0 {
								mg0 = match1a
							}
							if match1b > mg0 {
								mg0 = match1b
							}
							dp[sz][mg0][parentG1] = (dp[sz][mg0][parentG1] + ways) % MOD
						}
					}
				}
			}
		}
	}
	var ans int64
	if K <= maxM {
		for g1 := 0; g1 <= K; g1++ {
			ans = (ans + dp[n][K][g1]) % MOD
		}
	}
	return ans
}

func genCase(r *rand.Rand) Test {
	n := r.Intn(6) + 1
	K := r.Intn((n+1)/2 + 1)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, K)
	out := fmt.Sprintf("%d", compute(n, K))
	return Test{sb.String(), out}
}

func runCase(bin string, t Test) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(t.in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := strings.TrimSpace(t.out)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(5))
	for i := 0; i < 25; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
