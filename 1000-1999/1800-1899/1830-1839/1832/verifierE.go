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

func solve(n, a1, x, y, m, k int) uint64 {
	a := int64(a1)
	xm := int64(x)
	ym := int64(y)
	mm := int64(m)
	b := make([]int64, k+1)
	var ans uint64
	for i := 1; i <= n; i++ {
		ai := a
		for t := k; t >= 1; t-- {
			val := (b[t] + b[t-1]) % mod
			if t == 1 {
				val = (val + ai) % mod
			}
			b[t] = val
		}
		b[0] = (b[0] + ai) % mod
		ans ^= uint64(b[k]) * uint64(i)
		a = (a*xm + ym) % mm
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	m := rng.Intn(50) + 2
	a1 := rng.Intn(m)
	x := rng.Intn(m)
	y := rng.Intn(m)
	k := rng.Intn(5) + 1
	input := fmt.Sprintf("%d %d %d %d %d %d\n", n, a1, x, y, m, k)
	exp := solve(n, a1, x, y, m, k)
	return input, fmt.Sprintf("%d", exp)
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
