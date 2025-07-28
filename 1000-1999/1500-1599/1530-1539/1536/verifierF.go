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

const modF int64 = 1_000_000_007

func modPowF(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % modF
		}
		a = a * a % modF
		b >>= 1
	}
	return res
}

func solveF(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(in, &n)
	fac := make([]int64, n+1)
	invFac := make([]int64, n+1)
	fac[0] = 1
	for i := 1; i <= n; i++ {
		fac[i] = fac[i-1] * int64(i) % modF
	}
	invFac[n] = modPowF(fac[n], modF-2)
	for i := n; i >= 1; i-- {
		invFac[i-1] = invFac[i] * int64(i) % modF
	}
	inv := make([]int64, n+1)
	inv[1] = 1
	for i := 2; i <= n; i++ {
		inv[i] = modF - int64(modF/int64(i))*inv[int(modF%int64(i))]%modF
	}
	comb := func(a, b int) int64 {
		if b < 0 || b > a {
			return 0
		}
		return fac[a] * invFac[b] % modF * invFac[a-b] % modF
	}
	ans := int64(0)
	for k := 2; k <= n; k += 2 {
		m := n - k
		var ways int64
		if m == 0 {
			if n%2 == 0 {
				ways = 2
			} else {
				continue
			}
		} else {
			ways = 2 * int64(n) % modF
			ways = ways * comb(k-1, m-1) % modF
			ways = ways * inv[m] % modF
		}
		ans = (ans + ways*fac[k]) % modF
	}
	return fmt.Sprintf("%d", ans)
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []string {
	r := rand.New(rand.NewSource(6))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := r.Intn(50) + 1
		tests[i] = fmt.Sprintf("%d\n", n)
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		expected := solveF(t)
		got, err := runBinary(bin, t)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("test %d failed. input: %sexpected %s got %s\n", i+1, t, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
