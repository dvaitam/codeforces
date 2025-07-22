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

const MOD = 1051131

func modAdd(a, b int) int {
	c := a + b
	if c >= MOD {
		c -= MOD
	}
	return c
}

func modSub(a, b int) int {
	c := a - b
	if c < 0 {
		c += MOD
	}
	return c
}

func modMul(a, b int) int { return int((int64(a) * int64(b)) % MOD) }

func modPow(a int, e int64) int {
	res := 1
	base := a % MOD
	for e > 0 {
		if e&1 != 0 {
			res = modMul(res, base)
		}
		base = modMul(base, base)
		e >>= 1
	}
	return res
}

func modInv(a int) int {
	var egcd func(int, int) (int, int, int)
	egcd = func(a, b int) (int, int, int) {
		if b == 0 {
			return a, 1, 0
		}
		g, x1, y1 := egcd(b, a%b)
		return g, y1, x1 - (a/b)*y1
	}
	_, x, _ := egcd(a, MOD)
	x %= MOD
	if x < 0 {
		x += MOD
	}
	return x
}

func fwht(a []int, invert bool) {
	n := len(a)
	for step := 1; step < n; step <<= 1 {
		for i := 0; i < n; i += step << 1 {
			for j := 0; j < step; j++ {
				u := a[i+j]
				v := a[i+j+step]
				a[i+j] = modAdd(u, v)
				a[i+j+step] = modSub(u, v)
			}
		}
	}
	if invert {
		invN := modInv(n)
		for i := 0; i < n; i++ {
			a[i] = modMul(a[i], invN)
		}
	}
}

func expected(m int, t int64, s int, initial []int) int {
	N := 1 << m
	a0 := make([]int, N)
	copy(a0, initial)
	for i := s; i < N; i++ {
		a0[i] = int((101*int64(a0[i-s]) + 10007) % MOD)
	}
	fwht(a0, false)
	pow2m := N % MOD
	for i := 0; i < N; i++ {
		hi := pow2m + 1 - ((i * 2) % MOD)
		hi %= MOD
		if hi < 0 {
			hi += MOD
		}
		a0[i] = modMul(a0[i], modPow(hi, t))
	}
	fwht(a0, true)
	res := 0
	for i := 0; i < N; i++ {
		res ^= a0[i]
	}
	return res
}

func generateCase(rng *rand.Rand) (string, int) {
	m := rng.Intn(3) + 1
	N := 1 << m
	t := int64(rng.Intn(5) + 1)
	s := rng.Intn(min(N, 5)) + 1
	init := make([]int, s)
	for i := 0; i < s; i++ {
		init[i] = rng.Intn(1000) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", m, t, s))
	for i, v := range init {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	ans := expected(m, t, s, init)
	return sb.String(), ans
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Sscan(out, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %d\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
