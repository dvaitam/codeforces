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

const MOD int64 = 998244353

func mobius(n int) []int {
	mu := make([]int, n+1)
	mu[1] = 1
	primes := []int{}
	isComp := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		if !isComp[i] {
			primes = append(primes, i)
			mu[i] = -1
		}
		for _, p := range primes {
			if i*p > n {
				break
			}
			isComp[i*p] = true
			if i%p == 0 {
				mu[i*p] = 0
				break
			} else {
				mu[i*p] = -mu[i]
			}
		}
	}
	return mu
}

func countForD(d, n, m int, L, R []int) int64 {
	limit := m / d
	base := 0
	diffs := make([]int, n)
	for i := 0; i < n; i++ {
		li := (L[i] + d - 1) / d
		ri := R[i] / d
		if li > ri {
			return 0
		}
		base += li
		diffs[i] = ri - li
	}
	limit -= base
	if limit < 0 {
		return 0
	}
	dp := make([]int64, limit+1)
	dp[0] = 1
	for _, r := range diffs {
		prefix := int64(0)
		ndp := make([]int64, limit+1)
		for s := 0; s <= limit; s++ {
			prefix += dp[s]
			if prefix >= MOD {
				prefix -= MOD
			}
			if s-r-1 >= 0 {
				prefix -= dp[s-r-1]
				if prefix < 0 {
					prefix += MOD
				}
			}
			ndp[s] = prefix
		}
		dp = ndp
	}
	var total int64
	for _, v := range dp {
		total += v
		if total >= MOD {
			total -= MOD
		}
	}
	return total
}

func solve(n, m int, L, R []int) string {
	mu := mobius(m)
	var ans int64
	for d := 1; d <= m; d++ {
		if mu[d] == 0 {
			continue
		}
		val := countForD(d, n, m, L, R)
		if val == 0 {
			continue
		}
		if mu[d] == 1 {
			ans += val
		} else {
			ans -= val
		}
		ans %= MOD
	}
	if ans < 0 {
		ans += MOD
	}
	return fmt.Sprintf("%d", ans)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 2
	m := rng.Intn(40) + 10
	L := make([]int, n)
	R := make([]int, n)
	for i := 0; i < n; i++ {
		l := rng.Intn(m) + 1
		r := rng.Intn(m-l+1) + l
		L[i] = l
		R[i] = r
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", L[i], R[i]))
	}
	return sb.String(), solve(n, m, L, R)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, want := generateCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
