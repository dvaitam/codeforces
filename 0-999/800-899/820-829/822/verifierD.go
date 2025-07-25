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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func sieve(n int) []int {
	spf := make([]int, n+1)
	primes := []int{}
	for i := 2; i <= n; i++ {
		if spf[i] == 0 {
			spf[i] = i
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p > spf[i] || i*p > n {
				break
			}
			spf[i*p] = p
		}
	}
	return spf
}

func precompute(spf []int) []int64 {
	n := len(spf) - 1
	f := make([]int64, n+1)
	f[1] = 0
	for i := 2; i <= n; i++ {
		p := spf[i]
		f[i] = f[i/p] + int64(i)*(int64(p)-1)/2
	}
	return f
}

func solveCase(t int64, l, r int) string {
	spf := sieve(r)
	f := precompute(spf)
	res := int64(0)
	pow := int64(1)
	for i := l; i <= r; i++ {
		res = (res + pow*(f[i]%MOD)) % MOD
		pow = pow * t % MOD
	}
	return fmt.Sprintf("%d", res%MOD)
}

func genCase(rng *rand.Rand) (string, string) {
	l := rng.Intn(50) + 2
	r := l + rng.Intn(50)
	tVal := rng.Int63n(1_000_000_000) + 1
	input := fmt.Sprintf("%d %d %d\n", tVal, l, r)
	expected := solveCase(tVal, l, r)
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
