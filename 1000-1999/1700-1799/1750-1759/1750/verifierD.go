package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const mod int64 = 998244353

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func sieve(limit int) []int {
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}
	primes := make([]int, 0)
	for i := 2; i <= limit; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

func primeFactors(x int64, primes []int) []int64 {
	factors := make([]int64, 0)
	for _, p := range primes {
		if int64(p)*int64(p) > x {
			break
		}
		if x%int64(p) == 0 {
			factors = append(factors, int64(p))
			for x%int64(p) == 0 {
				x /= int64(p)
			}
		}
	}
	if x > 1 {
		factors = append(factors, x)
	}
	return factors
}

func countCoprime(limit int64, factors []int64) int64 {
	if limit <= 0 {
		return 0
	}
	m := len(factors)
	var bad int64
	for mask := 1; mask < (1 << m); mask++ {
		prod := int64(1)
		bits := 0
		for i := 0; i < m; i++ {
			if (mask>>i)&1 == 1 {
				prod *= factors[i]
				bits++
			}
		}
		if bits%2 == 1 {
			bad += limit / prod
		} else {
			bad -= limit / prod
		}
	}
	return limit - bad
}

func solveD(r *bufio.Reader) string {
	primes := sieve(31623)
	var T int
	if _, err := fmt.Fscan(r, &T); err != nil {
		return ""
	}
	var out strings.Builder
	for ; T > 0; T-- {
		var n int
		var m int64
		fmt.Fscan(r, &n, &m)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(r, &a[i])
		}
		if a[0] > m {
			out.WriteString("0\n")
			continue
		}
		ans := int64(1)
		ok := true
		for i := 1; i < n && ok; i++ {
			if a[i-1]%a[i] != 0 {
				ok = false
				break
			}
			x := a[i-1] / a[i]
			limit := m / a[i]
			pf := primeFactors(x, primes)
			cnt := countCoprime(limit, pf)
			ans = (ans * (cnt % mod)) % mod
		}
		if ok {
			out.WriteString(fmt.Sprintf("%d\n", ans%mod))
		} else {
			out.WriteString("0\n")
		}
	}
	return strings.TrimSpace(out.String())
}

func generateCaseD(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	m := rng.Intn(30) + 1
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", rng.Intn(30)+1)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseD(rng)
		expect := solveD(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
