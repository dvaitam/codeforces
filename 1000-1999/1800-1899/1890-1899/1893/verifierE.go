package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod int64 = 998244353

type edgeE struct {
	a int
	b int
	d int64
}

type testCaseE struct {
	n     int
	m     int
	edges []edgeE
}

func mul(x, y [4]int64) [4]int64 {
	return [4]int64{
		(x[0]*y[0] + x[1]*y[2]) % mod,
		(x[0]*y[1] + x[1]*y[3]) % mod,
		(x[2]*y[0] + x[3]*y[2]) % mod,
		(x[2]*y[1] + x[3]*y[3]) % mod,
	}
}

func powJacob(n int64) [4]int64 {
	res := [4]int64{1, 0, 0, 1}
	base := [4]int64{1, 2, 1, 0}
	for n > 0 {
		if n&1 == 1 {
			res = mul(res, base)
		}
		base = mul(base, base)
		n >>= 1
	}
	return res
}

func jacobsthal(n int64) int64 {
	if n == 0 {
		return 0
	}
	p := powJacob(n - 1)
	return p[0] % mod
}

func expectedE(tc testCaseE) string {
	n, m := tc.n, tc.m
	deg := make([]int, n+1)
	var total int64
	for _, e := range tc.edges {
		deg[e.a]++
		deg[e.b]++
		total += e.d + 1
	}
	if m != n {
		return "0"
	}
	for i := 1; i <= n; i++ {
		if deg[i] != 2 {
			return "0"
		}
	}
	ans := jacobsthal(total - 1)
	ans = ans * 12 % mod
	return fmt.Sprint(ans)
}

func genTestsE() []testCaseE {
	rand.Seed(5)
	tests := make([]testCaseE, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(5) + 1
		m := rand.Intn(n + 1)
		if rand.Intn(2) == 0 {
			m = n
		}
		edges := make([]edgeE, m)
		for i := 0; i < m; i++ {
			a := rand.Intn(n) + 1
			b := rand.Intn(n) + 1
			for b == a {
				b = rand.Intn(n) + 1
			}
			d := rand.Int63n(3)
			edges[i] = edgeE{a: a, b: b, d: d}
		}
		tests = append(tests, testCaseE{n: n, m: m, edges: edges})
	}
	return tests
}

func runCase(bin string, tc testCaseE) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.a, e.b, e.d))
	}
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expect := expectedE(tc)
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsE()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
