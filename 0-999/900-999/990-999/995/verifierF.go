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

const MOD int64 = 1000000007

func modPow(a, b int64) int64 {
	res := int64(1)
	a %= MOD
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func solveF(n int, D int64, parent []int) int64 {
	children := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		p := parent[i]
		children[p] = append(children[p], i)
	}
	m := n
	dp := make([][]int64, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = make([]int64, m+1)
	}
	order := make([]int, 0, n)
	stack := []int{1}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, c := range children[v] {
			stack = append(stack, c)
		}
	}
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		for x := 1; x <= m; x++ {
			prod := int64(1)
			for _, c := range children[v] {
				prod = prod * dp[c][x] % MOD
			}
			dp[v][x] = (dp[v][x-1] + prod) % MOD
		}
	}
	f := make([]int64, m+1)
	for i := 0; i <= m; i++ {
		f[i] = dp[1][i]
	}
	if D <= int64(m) {
		return f[D]
	}
	fact := make([]int64, m+1)
	invFact := make([]int64, m+1)
	fact[0] = 1
	for i := 1; i <= m; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[m] = modPow(fact[m], MOD-2)
	for i := m - 1; i >= 0; i-- {
		invFact[i] = invFact[i+1] * int64(i+1) % MOD
	}
	pre := make([]int64, m+1)
	suf := make([]int64, m+2)
	pre[0] = 1
	for i := 1; i <= m; i++ {
		val := (D - int64(i-1)) % MOD
		if val < 0 {
			val += MOD
		}
		pre[i] = pre[i-1] * val % MOD
	}
	suf[m+1] = 1
	for i := m; i >= 0; i-- {
		val := (D - int64(i)) % MOD
		if val < 0 {
			val += MOD
		}
		suf[i] = suf[i+1] * val % MOD
	}
	ans := int64(0)
	for i := 0; i <= m; i++ {
		num := pre[i] * suf[i+1] % MOD
		term := f[i] * num % MOD
		term = term * invFact[i] % MOD
		term = term * invFact[m-i] % MOD
		if (m-i)%2 == 1 {
			term = (MOD - term) % MOD
		}
		ans = (ans + term) % MOD
	}
	return ans
}

type caseF struct {
	n      int
	D      int64
	parent []int
}

func genCaseF(rng *rand.Rand) caseF {
	n := rng.Intn(8) + 1
	parent := make([]int, n+1)
	for i := 2; i <= n; i++ {
		parent[i] = rng.Intn(i-1) + 1
	}
	D := int64(rng.Intn(15) + 1)
	return caseF{n: n, D: D, parent: parent}
}

func runCaseF(bin string, tc caseF) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.D)
	for i := 2; i <= tc.n; i++ {
		fmt.Fprintf(&sb, "%d\n", tc.parent[i])
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
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
	exp := solveF(tc.n, tc.D, tc.parent)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseF(rng)
		if err := runCaseF(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
