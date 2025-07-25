package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const MOD int64 = 1000000007

type Test struct {
	in  string
	out string
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func modPow(a, e int64) int64 {
	a %= MOD
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

func modInv(a int64) int64 {
	return modPow((a%MOD+MOD)%MOD, MOD-2)
}

func solveOne(n, k int64) int64 {
	g := gcd(n, k)
	n /= g
	k /= g
	pow2k := modPow(2, k)
	invPow2k := modInv(pow2k)
	Lterm := n % MOD * modInv((pow2k-1+MOD)%MOD) % MOD
	invK := modInv(k)
	invNK := modInv(n % MOD * k % MOD)
	B := (n%MOD*invK%MOD - invNK + MOD) % MOD
	inv2 := modInv(2)
	A := ((n+1)%MOD*inv2%MOD - ((k-1)%MOD)*B%MOD*inv2%MOD + MOD) % MOD
	Dnumer := (1 - ((k+1)%MOD)*invPow2k%MOD + MOD) % MOD
	Ddenom := (1 - invPow2k + MOD) % MOD
	D := Dnumer * modInv(Ddenom) % MOD
	return (Lterm + A + B*D%MOD) % MOD
}

func oracle(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return ""
	}
	var sb strings.Builder
	for ; T > 0; T-- {
		var n, k int64
		fmt.Fscan(reader, &n, &k)
		fmt.Fprintln(&sb, solveOne(n, k))
	}
	return strings.TrimSpace(sb.String())
}

func genCase(rng *rand.Rand) Test {
	T := rng.Intn(4) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", T)
	for i := 0; i < T; i++ {
		n := int64(rng.Intn(1000) + 1)
		k := int64(rng.Intn(int(n)) + 1)
		fmt.Fprintf(&sb, "%d %d\n", n, k)
	}
	input := sb.String()
	out := oracle(input)
	return Test{input, out}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(7))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		got, err := run(bin, tc.in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
		if got != tc.out {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.out, got, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
