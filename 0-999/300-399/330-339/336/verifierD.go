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

const modD = 1000000007

func modpow(a, e int64) int64 {
	res := int64(1)
	a %= modD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % modD
		}
		a = a * a % modD
		e >>= 1
	}
	return res
}

func expectedAnswerD(n, m, g int) int64 {
	if m == 0 {
		B := 0
		if n%2 == 0 {
			B = 1
		}
		if g == 1 {
			return int64(B)
		}
		return int64((1 - B + modD) % modD)
	}
	N := n + m
	fact := make([]int64, N+1)
	invfact := make([]int64, N+1)
	fact[0] = 1
	for i := 1; i <= N; i++ {
		fact[i] = fact[i-1] * int64(i) % modD
	}
	invfact[N] = modpow(fact[N], modD-2)
	for i := N; i >= 1; i-- {
		invfact[i-1] = invfact[i] * int64(i) % modD
	}
	comb := func(a, b int) int64 {
		if b < 0 || b > a {
			return 0
		}
		return fact[a] * invfact[b] % modD * invfact[a-b] % modD
	}
	B := make([]int64, n+1)
	if m == 1 {
		B[0] = 1
	} else {
		B[0] = 0
	}
	if n >= 1 {
		if m >= 2 {
			B[1] = 1
		} else {
			B[1] = 0
		}
	}
	for i := 2; i <= n; i++ {
		B[i] = (comb(i+m-2, i-1) + B[i-2]) % modD
	}
	total := comb(n+m, n)
	Bn := B[n]
	if g == 1 {
		return Bn
	}
	return (total - Bn + modD) % modD
}

func generateCaseD(rng *rand.Rand) (int, int, int) {
	n := rng.Intn(6)
	m := rng.Intn(6)
	g := rng.Intn(2)
	if n+m == 0 {
		n = 1
	}
	return n, m, g
}

func runCaseD(bin string, n, m, g int) error {
	input := fmt.Sprintf("%d %d %d\n", n, m, g)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(gotStr, 10, 64)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expected := expectedAnswerD(n, m, g)
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, g := generateCaseD(rng)
		if err := runCaseD(bin, n, m, g); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d %d %d\n", i+1, err, n, m, g)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
