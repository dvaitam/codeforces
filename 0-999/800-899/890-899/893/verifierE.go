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
const MAXN int = 1000020

var fac [MAXN + 1]int64
var ifac [MAXN + 1]int64
var primes []int

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func initFactorials() {
	fac[0] = 1
	for i := 1; i <= MAXN; i++ {
		fac[i] = fac[i-1] * int64(i) % MOD
	}
	ifac[MAXN] = modPow(fac[MAXN], MOD-2)
	for i := MAXN; i > 0; i-- {
		ifac[i-1] = ifac[i] * int64(i) % MOD
	}
}

func C(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fac[n] * ifac[k] % MOD * ifac[n-k] % MOD
}

func sievePrimes(limit int) {
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}
	for p := 2; p*p <= limit; p++ {
		if isPrime[p] {
			for j := p * p; j <= limit; j += p {
				isPrime[j] = false
			}
		}
	}
	for i := 2; i <= limit; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
}

type pair struct {
	p int
	e int
}

func factorize(x int) []pair {
	res := []pair{}
	tmp := x
	for _, p := range primes {
		if p*p > tmp {
			break
		}
		if tmp%p == 0 {
			cnt := 0
			for tmp%p == 0 {
				tmp /= p
				cnt++
			}
			res = append(res, pair{p, cnt})
		}
	}
	if tmp > 1 {
		res = append(res, pair{tmp, 1})
	}
	return res
}

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveQuery(x, y int) int64 {
	factors := factorize(x)
	ans := modPow(2, int64(y-1))
	for _, pe := range factors {
		ans = ans * C(pe.e+y-1, y-1) % MOD
	}
	return ans
}

func generateCase(r *rand.Rand) (string, string) {
	q := r.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", q))
	answers := make([]string, q)
	for i := 0; i < q; i++ {
		x := r.Intn(1000000) + 1
		y := r.Intn(1000) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
		answers[i] = fmt.Sprint(solveQuery(x, y))
	}
	return sb.String(), strings.Join(answers, "\n")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	sievePrimes(1000000)
	initFactorials()
	bin := os.Args[1]
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(r)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
