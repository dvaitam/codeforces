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

const MOD int = 1000000007

func modPow(a, e int) int {
	res := 1
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func modInv(a int) int { return modPow(a, MOD-2) }

func solveF(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n, a, b int
	fmt.Fscan(in, &n, &a, &b)
	g := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &g[i])
	}
	L := make([]int, n)
	R := make([]int, n)
	maxN := n
	for i := 0; i < n; i++ {
		var s int
		var str string
		fmt.Fscan(in, &s, &str)
		R[i] = s
		cnt := 0
		for j := 0; j < len(str); j++ {
			if str[j] == '0' {
				cnt++
			}
		}
		L[i] = cnt
		if s > maxN {
			maxN = s
		}
	}
	fact := make([]int, maxN+1)
	invFact := make([]int, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * i % MOD
	}
	invFact[maxN] = modInv(fact[maxN])
	for i := maxN - 1; i >= 0; i-- {
		invFact[i] = invFact[i+1] * (i + 1) % MOD
	}
	comb := func(n, k int) int {
		if n < 0 || k < 0 || k > n {
			return 0
		}
		return fact[n] * invFact[k] % MOD * invFact[n-k] % MOD
	}
	cntTop := 0
	for i := 0; i < n; i++ {
		rank := 1
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			if L[j] > R[i] {
				rank++
			}
		}
		if rank <= a {
			cntTop++
		}
	}
	ans := comb(cntTop, b)
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

func genTest() string {
	r := rand.New(rand.NewSource(rand.Int63()))
	n := r.Intn(4) + 2
	a := r.Intn(n) + 1
	b := r.Intn(a) + 1
	mat := make([][]byte, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]byte, n)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if r.Intn(2) == 0 {
				mat[i][j] = '1'
				mat[j][i] = '0'
			} else {
				mat[i][j] = '0'
				mat[j][i] = '1'
			}
		}
		mat[i][i] = '0'
	}
	lines := []string{fmt.Sprintf("%d %d %d", n, a, b)}
	for i := 0; i < n; i++ {
		lines = append(lines, string(mat[i]))
	}
	for i := 0; i < n; i++ {
		s := r.Intn(3) + 1
		var sb strings.Builder
		for j := 0; j < s; j++ {
			if r.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		lines = append(lines, fmt.Sprintf("%d %s", s, sb.String()))
	}
	return strings.Join(lines, "\n") + "\n"
}

func generateTests() []string {
	rand.Seed(6)
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		tests[i] = genTest()
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
			fmt.Printf("test %d failed. expected %s got %s\ninput:\n%s", i+1, expected, got, t)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
