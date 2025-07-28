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

func modPow(a, e int64) int64 {
	a %= mod
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

func expectedE(n int) int64 {
	if n%2 == 1 {
		return 0
	}
	fac := make([]int64, n+1)
	ifac := make([]int64, n+1)
	fac[0] = 1
	for i := 1; i <= n; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
	}
	ifac[n] = modPow(fac[n], mod-2)
	for i := n; i >= 1; i-- {
		ifac[i-1] = ifac[i] * int64(i) % mod
	}
	ans := int64(0)
	for s := 1; s <= n-1; s++ {
		comb := fac[n-2]
		comb = comb * ifac[s-1] % mod
		comb = comb * ifac[n-1-s] % mod
		term := comb
		term = term * modPow(int64(s), int64(s-1)) % mod
		term = term * modPow(int64(n-s), int64(n-s-1)) % mod
		if s%2 == 1 {
			ans = (ans - term) % mod
		} else {
			ans = (ans + term) % mod
		}
	}
	if ans < 0 {
		ans += mod
	}
	return ans
}

func genTestE(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	exp := expectedE(n)
	input := fmt.Sprintf("%d\n", n)
	return input, fmt.Sprintf("%d", exp)
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input, expect := genTestE(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n%s", t+1, err, got)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", t+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}
