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

const mod int64 = 998244353

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
	return strings.TrimSpace(out.String()), nil
}

func modPow(a, b int64) int64 {
	res := int64(1)
	a %= mod
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func solveF(n, m int64, k int) int64 {
	limit := k
	if n < int64(limit) {
		limit = int(n)
	}
	invM := modPow(m%mod, mod-2)
	powM := make([]int64, limit+1)
	powM[0] = 1
	for i := 1; i <= limit; i++ {
		powM[i] = powM[i-1] * invM % mod
	}
	fall := make([]int64, limit+1)
	fall[0] = 1
	for i := 1; i <= limit; i++ {
		fall[i] = fall[i-1] * ((n - int64(i) + 1) % mod) % mod
	}
	S := make([]int64, limit+1)
	S[0] = 1
	for i := 1; i <= k; i++ {
		upper := i
		if upper > limit {
			upper = limit
		}
		for j := upper; j >= 1; j-- {
			S[j] = (S[j-1] + S[j]*int64(j)) % mod
		}
		S[0] = 0
	}
	var ans int64
	for j := 0; j <= limit; j++ {
		ans = (ans + S[j]*fall[j]%mod*powM[j]) % mod
	}
	return ans
}

func generateCase(r *rand.Rand) (string, string) {
	n := int64(r.Intn(200) + 1)
	m := int64(r.Intn(200) + 1)
	k := r.Intn(10) + 1
	expect := fmt.Sprintf("%d", solveF(n, m, k))
	input := fmt.Sprintf("%d %d %d\n", n, m, k)
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
