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

const MOD int64 = 998244353
const MAXK = 2000

var stirling [MAXK + 1][MAXK + 1]int64

func init() {
	stirling[0][0] = 1
	for n := 1; n <= MAXK; n++ {
		for k := 1; k <= n; k++ {
			stirling[n][k] = (stirling[n-1][k-1] + int64(k)*stirling[n-1][k]) % MOD
		}
	}
}

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

func solveCase(n, m, k int64) int64 {
	a := (m + 1) / 2
	invm := modPow(m%MOD, MOD-2)
	p := a % MOD * invm % MOD
	mPow := modPow(m%MOD, n)
	maxI := k
	if n < maxI {
		maxI = n
	}
	fall := int64(1)
	powp := int64(1)
	ans := stirling[k][0] * fall % MOD * powp % MOD
	for i := int64(1); i <= maxI; i++ {
		fall = fall * ((n - i + 1) % MOD) % MOD
		powp = powp * p % MOD
		ans = (ans + stirling[k][i]*fall%MOD*powp) % MOD
	}
	ans = ans * mPow % MOD
	return ans
}

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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func verifyCase(bin string, n, m, k int64) error {
	input := fmt.Sprintf("1\n%d %d %d\n", n, m, k)
	expected := fmt.Sprint(solveCase(n, m, k))
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != expected {
		return fmt.Errorf("expected %s got %s", expected, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := int64(rng.Intn(1000) + 1)
		m := int64(rng.Intn(1000) + 1)
		k := int64(rng.Intn(10) + 1)
		if err := verifyCase(bin, n, m, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
