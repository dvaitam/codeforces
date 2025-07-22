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

const mod int64 = 1000000007

type mat [2][2]int64

func mul(a, b mat) mat {
	return mat{
		{(a[0][0]*b[0][0] + a[0][1]*b[1][0]) % mod, (a[0][0]*b[0][1] + a[0][1]*b[1][1]) % mod},
		{(a[1][0]*b[0][0] + a[1][1]*b[1][0]) % mod, (a[1][0]*b[0][1] + a[1][1]*b[1][1]) % mod},
	}
}

func powMat(m mat, e int64) mat {
	res := mat{{1, 0}, {0, 1}}
	for e > 0 {
		if e&1 == 1 {
			res = mul(res, m)
		}
		m = mul(m, m)
		e >>= 1
	}
	return res
}

func expected(n int64) int64 {
	base := mat{{0, 3}, {1, 2}}
	mn := powMat(base, n)
	return mn[0][0] % mod
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Int63n(1000000) + 1
	return fmt.Sprintf("%d\n", n), expected(n)
}

func runCase(exe, input string, expectedAns int64) error {
	cmd := exec.Command(exe)
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
	if got != expectedAns {
		return fmt.Errorf("expected %d got %d", expectedAns, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
