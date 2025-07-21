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

const mod = 1000000007

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

func expected(n int) int64 {
	N := 2 * n
	fac := make([]int64, N+1)
	invFac := make([]int64, N+1)
	fac[0] = 1
	for i := 1; i <= N; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
	}
	invFac[N] = modPow(fac[N], mod-2)
	for i := N; i > 0; i-- {
		invFac[i-1] = invFac[i] * int64(i) % mod
	}
	comb := fac[2*n-1] * invFac[n] % mod * invFac[n-1] % mod
	ans := (comb*2 - int64(n)) % mod
	if ans < 0 {
		ans += mod
	}
	return ans
}

func genCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(1000) + 1
	input := fmt.Sprintf("%d\n", n)
	return input, expected(n)
}

func runCase(bin, input string, exp int64) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
