package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func fastPowMod(a, e, m int64) int64 {
	a %= m
	var res int64 = 1
	for e > 0 {
		if e&1 == 1 {
			res = (res * a) % m
		}
		a = (a * a) % m
		e >>= 1
	}
	return res
}

func fibPair(n, m int64) (int64, int64) {
	if n == 0 {
		return 0, 1
	}
	a, b := fibPair(n>>1, m)
	c := (a * ((2*b%m - a + m) % m)) % m
	d := (a*a%m + b*b%m) % m
	if n&1 == 0 {
		return c, d
	}
	return d, (c + d) % m
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
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

func solveE(r *bufio.Reader) string {
	var n int
	var x, y, p int64
	fmt.Fscan(r, &n, &x, &y, &p)
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &a[i])
	}
	if n == 1 {
		return fmt.Sprintf("%d", a[0]%p)
	}
	var midSumModP int64
	for i := 1; i+1 < n; i++ {
		midSumModP = (midSumModP + a[i]) % p
	}
	a0 := a[0]
	alast := a[n-1]
	mod2p := p * 2
	pow3x := fastPowMod(3, x, mod2p)
	num := ((2*midSumModP%mod2p)*pow3x%mod2p + ((a0+alast)%mod2p)*((pow3x+1)%mod2p)%mod2p) % mod2p
	Sx := (num / 2) % p
	Fx, Fx1 := fibPair(x, mod2p)
	L := a[n-2] % mod2p
	R := a[n-1] % mod2p
	Mx := (Fx*L + Fx1*R) % mod2p
	D := (a0%mod2p + Mx) % mod2p
	pow3y := fastPowMod(3, y, mod2p)
	twoSx := (2 * Sx) % mod2p
	term1 := (pow3y * twoSx) % mod2p
	term2 := (D * ((pow3y - 1 + mod2p) % mod2p)) % mod2p
	twoSy := (term1 - term2 + mod2p) % mod2p
	Sy := (twoSy / 2) % p
	return fmt.Sprintf("%d", Sy)
}

func generateCaseE(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	x := rng.Int63n(1000)
	y := rng.Int63n(1000)
	if x == 0 && y == 0 {
		x = 1
	}
	p := int64(rng.Intn(1000) + 2)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", n, x, y, p)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", rng.Intn(1000))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseE(rng)
		expect := solveE(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
