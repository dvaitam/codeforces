package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const MOD int64 = 998244353

type TestCaseC struct {
	n int
	k int64
	b []int64
}

func genCaseC(rng *rand.Rand) TestCaseC {
	n := rng.Intn(6) + 1
	k := int64(rng.Intn(5) + 1)
	b := make([]int64, n)
	for i := range b {
		b[i] = rng.Int63n(20)
	}
	return TestCaseC{n: n, k: k, b: b}
}

func applyN(v []int64) []int64 {
	n := len(v)
	res := make([]int64, n)
	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = (prefix[i-1] + v[i-1]) % MOD
	}
	for i := 1; i <= n; i++ {
		lb := i & -i
		res[i-1] = (prefix[i-1] - prefix[i-lb]) % MOD
		if res[i-1] < 0 {
			res[i-1] += MOD
		}
	}
	return res
}

func powMod(a, e int64) int64 {
	r := int64(1)
	for e > 0 {
		if e&1 == 1 {
			r = r * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return r
}

func modInv(a int64) int64 { return powMod(a, MOD-2) }

func solveCaseC(n int, k int64, b []int64) []int64 {
	r := bits.Len(uint(n))
	powers := make([][]int64, r)
	powers[0] = append([]int64(nil), b...)
	for i := 1; i < r; i++ {
		powers[i] = applyN(powers[i-1])
	}
	coeff := make([]int64, r)
	coeff[0] = 1
	comb := int64(1)
	for i := 1; i < r; i++ {
		comb = comb * (k + int64(i-1)) % MOD
		comb = comb * modInv(int64(i)) % MOD
		val := comb
		if i%2 == 1 {
			val = (MOD - val) % MOD
		}
		coeff[i] = val
	}
	ans := make([]int64, n)
	for i := 0; i < r; i++ {
		c := coeff[i]
		if c == 0 {
			continue
		}
		for j := 0; j < n; j++ {
			ans[j] = (ans[j] + c*powers[i][j]) % MOD
		}
	}
	return ans
}

func runCaseC(bin string, tc TestCaseC, expect []int64) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != tc.n {
		return fmt.Errorf("expected %d numbers got %d", tc.n, len(fields))
	}
	for i, f := range fields {
		var v int64
		fmt.Sscan(f, &v)
		if v%MOD != expect[i] {
			return fmt.Errorf("index %d expected %d got %d", i, expect[i], v%MOD)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseC(rng)
		exp := solveCaseC(tc.n, tc.k, tc.b)
		if err := runCaseC(bin, tc, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: n=%d k=%d b=%v\n", i+1, err, tc.n, tc.k, tc.b)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
