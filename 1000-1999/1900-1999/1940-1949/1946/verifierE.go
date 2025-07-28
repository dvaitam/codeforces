package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const MOD int64 = 1000000007
const MAXN = 200005

type test struct {
	input    string
	expected string
}

var fac [MAXN]int64
var ifac [MAXN]int64
var invNum [MAXN]int64

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func initComb() {
	fac[0] = 1
	for i := 1; i < MAXN; i++ {
		fac[i] = fac[i-1] * int64(i) % MOD
	}
	ifac[MAXN-1] = modPow(fac[MAXN-1], MOD-2)
	for i := MAXN - 1; i > 0; i-- {
		ifac[i-1] = ifac[i] * int64(i) % MOD
	}
	for i := 1; i < MAXN; i++ {
		invNum[i] = modPow(int64(i), MOD-2)
	}
}

func comb(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fac[n] * ifac[k] % MOD * ifac[n-k] % MOD
}

func prefixCount(L int, P []int) int64 {
	if L == 0 {
		if len(P) == 0 {
			return 1
		}
		return 0
	}
	if len(P) == 0 || P[0] != 1 || P[len(P)-1] > L {
		return 0
	}
	for i := 1; i < len(P); i++ {
		if P[i] <= P[i-1] {
			return 0
		}
	}
	ans := fac[L-1]
	for i := 1; i < len(P); i++ {
		ans = ans * invNum[P[i]-1] % MOD
	}
	return ans
}

func suffixCount(L int, S []int) int64 {
	if L == 0 {
		if len(S) == 0 {
			return 1
		}
		return 0
	}
	if len(S) == 0 || S[len(S)-1] != L {
		return 0
	}
	for i := 1; i < len(S); i++ {
		if S[i] <= S[i-1] {
			return 0
		}
	}
	ans := fac[L-1]
	for i := 0; i < len(S)-1; i++ {
		ans = ans * invNum[L-S[i]] % MOD
	}
	return ans
}

func solveCase(n int, P, S []int) int64 {
	if len(P) == 0 || len(S) == 0 {
		return 0
	}
	if P[0] != 1 || S[len(S)-1] != n {
		return 0
	}
	mp := make(map[int]struct{}, len(P))
	for _, v := range P {
		mp[v] = struct{}{}
	}
	inter := 0
	for _, v := range S {
		if _, ok := mp[v]; ok {
			inter++
		}
	}
	if inter != 1 {
		return 0
	}
	x := P[len(P)-1]
	if S[0] != x {
		return 0
	}
	for _, v := range P {
		if v > x {
			return 0
		}
	}
	for _, v := range S {
		if v < x {
			return 0
		}
	}
	L := x - 1
	R := n - x
	var Pleft []int
	for _, v := range P {
		if v < x {
			Pleft = append(Pleft, v)
		}
	}
	var Sright []int
	for _, v := range S {
		if v > x {
			Sright = append(Sright, v-x)
		}
	}
	A := prefixCount(L, Pleft)
	B := suffixCount(R, Sright)
	if A == 0 || B == 0 {
		return 0
	}
	ans := comb(n-1, L) * A % MOD
	ans = ans * B % MOD
	return ans
}

func solve(input string) string {
	reader := strings.NewReader(input)
	var t int
	fmt.Fscan(reader, &t)
	var out strings.Builder
	for ; t > 0; t-- {
		var n, m1, m2 int
		fmt.Fscan(reader, &n, &m1, &m2)
		P := make([]int, m1)
		for i := 0; i < m1; i++ {
			fmt.Fscan(reader, &P[i])
		}
		S := make([]int, m2)
		for i := 0; i < m2; i++ {
			fmt.Fscan(reader, &S[i])
		}
		ans := solveCase(n, P, S)
		out.WriteString(fmt.Sprintf("%d\n", ans))
	}
	return strings.TrimSpace(out.String())
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(50))
	var tests []test
	fixed := []string{
		"1\n1 0 0\n\n\n",
	}
	for _, f := range fixed {
		tests = append(tests, test{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rng.Intn(5) + 1
		m1 := rng.Intn(n + 1)
		m2 := rng.Intn(n + 1)
		var P []int
		for i := 0; i < m1; i++ {
			v := rng.Intn(n) + 1
			P = append(P, v)
		}
		var S []int
		for i := 0; i < m2; i++ {
			v := rng.Intn(n) + 1
			S = append(S, v)
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m1, m2))
		for i, v := range P {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		for i, v := range S {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	initComb()
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
