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

const MOD int64 = 1000000007

var spf []int

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

func comb(n, k int, fact, invFact []int64) int64 {
	if n < k || k < 0 {
		return 0
	}
	return fact[n] * invFact[k] % MOD * invFact[n-k] % MOD
}

func precompute(limit int) ([]int64, []int64) {
	fact := make([]int64, limit+1)
	invFact := make([]int64, limit+1)
	fact[0] = 1
	for i := 1; i <= limit; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[limit] = modPow(fact[limit], MOD-2)
	for i := limit; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}
	return fact, invFact
}

func initSPF(max int) {
	spf = make([]int, max+1)
	for i := 2; i*i <= max; i++ {
		if spf[i] == 0 {
			for j := i * i; j <= max; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
	for i := 2; i <= max; i++ {
		if spf[i] == 0 {
			spf[i] = i
		}
	}
}

func solveQuery(r, n int, fact, invFact []int64) int64 {
	ans := int64(1)
	m := n
	for m > 1 {
		p := spf[m]
		cnt := 0
		for m%p == 0 {
			m /= p
			cnt++
		}
		val := comb(r+cnt+1, r+1, fact, invFact) - comb(r+cnt-1, r+1, fact, invFact)
		if val < 0 {
			val += MOD
		}
		ans = ans * val % MOD
	}
	return ans % MOD
}

func generateCase(rng *rand.Rand, fact, invFact []int64) (string, string) {
	q := rng.Intn(5) + 1
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		r := rng.Intn(5)
		n := rng.Intn(1000) + 1
		queries[i] = [2]int{r, n}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", q)
	for _, qq := range queries {
		fmt.Fprintf(&sb, "%d %d\n", qq[0], qq[1])
	}
	var out strings.Builder
	for _, qq := range queries {
		ans := solveQuery(qq[0], qq[1], fact, invFact)
		fmt.Fprintf(&out, "%d\n", ans)
	}
	return sb.String(), strings.TrimSpace(out.String())
}

func runCase(exe, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	initSPF(1000000)
	fact, invFact := precompute(10000)
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng, fact, invFact)
		got, err := runCase(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
