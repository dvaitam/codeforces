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

const mod int64 = 998244353

func modPow(a, e, m int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % m
		}
		a = a * a % m
		e >>= 1
	}
	return res
}

func computeAnswers(D int64, queries [][2]int64) []int64 {
	primes := make([]int64, 0)
	exps := make([]int, 0)
	n := D
	for p := int64(2); p*p <= n; p++ {
		if n%p == 0 {
			cnt := 0
			for n%p == 0 {
				n /= p
				cnt++
			}
			primes = append(primes, p)
			exps = append(exps, cnt)
		}
	}
	if n > 1 {
		primes = append(primes, n)
		exps = append(exps, 1)
	}
	maxExp := 0
	for _, e := range exps {
		maxExp += e
	}
	fact := make([]int64, maxExp+1)
	invFact := make([]int64, maxExp+1)
	fact[0] = 1
	for i := 1; i <= maxExp; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[maxExp] = modPow(fact[maxExp], mod-2, mod)
	for i := maxExp; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}

	ans := make([]int64, len(queries))
	for idx, qu := range queries {
		v := qu[0]
		u := qu[1]
		diffV := make([]int, len(primes))
		diffU := make([]int, len(primes))
		sumV, sumU := 0, 0
		tempV := v
		tempU := u
		for i, p := range primes {
			cv, cu := 0, 0
			for tempV%p == 0 {
				tempV /= p
				cv++
			}
			for tempU%p == 0 {
				tempU /= p
				cu++
			}
			g := cv
			if cu < g {
				g = cu
			}
			diffV[i] = cv - g
			diffU[i] = cu - g
			sumV += diffV[i]
			sumU += diffU[i]
		}
		res := fact[sumV]
		for _, x := range diffV {
			res = res * invFact[x] % mod
		}
		res2 := fact[sumU]
		for _, x := range diffU {
			res2 = res2 * invFact[x] % mod
		}
		ans[idx] = res * res2 % mod
	}
	return ans
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randomDivisors(D int64, rng *rand.Rand) int64 {
	divisors := []int64{}
	for i := int64(1); i*i <= D; i++ {
		if D%i == 0 {
			divisors = append(divisors, i)
			if i*i != D {
				divisors = append(divisors, D/i)
			}
		}
	}
	return divisors[rng.Intn(len(divisors))]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for caseNum := 0; caseNum < 100; caseNum++ {
		D := int64(rng.Intn(1000)+1) * int64(rng.Intn(5)+1)
		q := rng.Intn(3) + 1
		queries := make([][2]int64, q)
		for i := 0; i < q; i++ {
			v := randomDivisors(D, rng)
			u := randomDivisors(D, rng)
			queries[i] = [2]int64{v, u}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", D))
		sb.WriteString(fmt.Sprintf("%d\n", q))
		for _, qu := range queries {
			sb.WriteString(fmt.Sprintf("%d %d\n", qu[0], qu[1]))
		}
		input := sb.String()
		answers := computeAnswers(D, queries)
		var expSB strings.Builder
		for i, v := range answers {
			if i > 0 {
				expSB.WriteByte('\n')
			}
			expSB.WriteString(fmt.Sprintf("%d", v))
		}
		exp := expSB.String()
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", caseNum+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
