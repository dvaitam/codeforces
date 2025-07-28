package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func modPow(a, e, mod int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func solveE(n int, mod int64) int64 {
	max := 3 * n
	fact := make([]int64, max+1)
	invFact := make([]int64, max+1)
	fact[0] = 1
	for i := 1; i <= max; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[max] = modPow(fact[max], mod-2, mod)
	for i := max; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
	comb := func(n, k int) int64 {
		if k < 0 || k > n {
			return 0
		}
		return fact[n] * invFact[k] % mod * invFact[n-k] % mod
	}
	factN := fact[n]
	fact2N := fact[2*n]
	fact3N := fact[3*n]
	pow3FactN := factN * factN % mod * factN % mod
	c2n := comb(2*n, n)
	countNoA := c2n * c2n % mod * pow3FactN % mod
	var s int64
	for k := 0; k <= n; k++ {
		term := comb(n, k) * comb(n, k) % mod * comb(2*n-k, n) % mod
		s = (s + term) % mod
	}
	intersection := s * pow3FactN % mod
	ans := (3*fact3N%mod - 2*fact2N%mod - 2*countNoA%mod + factN%mod + intersection%mod - 1) % mod
	if ans < 0 {
		ans += mod
	}
	return ans
}

var primes = []int64{998244353, 1000000007, 1000000009, 1000000033}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type test struct {
		n   int
		mod int64
	}

	var cases []test
	cases = append(cases, test{n: 1, mod: 998244353})
	cases = append(cases, test{n: 2, mod: 1000000007})
	cases = append(cases, test{n: 3, mod: 1000000009})

	for len(cases) < 100 {
		n := rng.Intn(10) + 1
		mod := primes[rng.Intn(len(primes))]
		cases = append(cases, test{n: n, mod: mod})
	}

	for i, tc := range cases {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.mod)
		expected := solveE(tc.n, tc.mod)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		fields := strings.Fields(strings.TrimSpace(out))
		if len(fields) != 1 {
			fmt.Fprintf(os.Stderr, "case %d: expected single integer got %q\n", i+1, out)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse integer\n", i+1)
			os.Exit(1)
		}
		if val%tc.mod != expected%tc.mod {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", i+1, expected%tc.mod, val%tc.mod)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
