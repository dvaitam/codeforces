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

const mod = 1000000007

func runCandidate(bin, input string) (string, error) {
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

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func solveCase(n int, s int64, f []int64) string {
	fplus := make([]int64, n)
	for i := 0; i < n; i++ {
		fplus[i] = f[i] + 1
	}
	maxR := n
	fact := make([]int64, maxR+1)
	invFact := make([]int64, maxR+1)
	fact[0] = 1
	for i := 1; i <= maxR; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[maxR] = modPow(fact[maxR], mod-2)
	for i := maxR; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
	N := 1 << n
	sumF := make([]int64, N)
	pc := make([]int, N)
	for mask := 1; mask < N; mask++ {
		lb := mask & -mask
		i := bits.TrailingZeros(uint(lb))
		prev := mask ^ lb
		sumF[mask] = sumF[prev] + fplus[i]
		pc[mask] = pc[prev] + 1
	}
	r := n - 1
	var ans int64
	for mask := 0; mask < N; mask++ {
		k := pc[mask]
		total := sumF[mask]
		a := s - total + int64(n-1)
		if a < int64(r) {
			continue
		}
		var num int64 = 1
		for i := 0; i < r; i++ {
			t := (a - int64(i)) % mod
			if t < 0 {
				t += mod
			}
			num = num * t % mod
		}
		comb := num * invFact[r] % mod
		if k&1 == 1 {
			ans = (ans - comb + mod) % mod
		} else {
			ans = (ans + comb) % mod
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	s := int64(rng.Intn(50))
	f := make([]int64, n)
	for i := 0; i < n; i++ {
		f[i] = int64(rng.Intn(20))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, s))
	for i, v := range f {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	expect := solveCase(n, s, f)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
