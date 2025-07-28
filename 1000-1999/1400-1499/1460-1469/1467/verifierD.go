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

type queryD struct {
	idx int
	x   int64
}

type testCaseD struct {
	n   int
	k   int
	q   int
	arr []int64
	qs  []queryD
}

const MOD int64 = 1000000007

func solveD(tc testCaseD) []int64 {
	n, k := tc.n, tc.k
	dp := make([][]int64, k+1)
	for i := 0; i <= k; i++ {
		dp[i] = make([]int64, n+2)
	}
	for i := 1; i <= n; i++ {
		dp[0][i] = 1
	}
	for t := 1; t <= k; t++ {
		for i := 1; i <= n; i++ {
			var v int64
			if i > 1 {
				v += dp[t-1][i-1]
			}
			if i < n {
				v += dp[t-1][i+1]
			}
			if v >= MOD {
				v %= MOD
			}
			dp[t][i] = v
		}
	}
	coef := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		var c int64
		for t := 0; t <= k; t++ {
			c = (c + dp[t][i]*dp[k-t][i]) % MOD
		}
		coef[i] = c
	}
	a := append([]int64{0}, tc.arr...)
	total := int64(0)
	for i := 1; i <= n; i++ {
		total = (total + a[i]*coef[i]) % MOD
	}
	res := make([]int64, tc.q)
	for qi, qu := range tc.qs {
		total = (total - a[qu.idx]*coef[qu.idx]) % MOD
		if total < 0 {
			total += MOD
		}
		a[qu.idx] = qu.x
		total = (total + a[qu.idx]*coef[qu.idx]) % MOD
		res[qi] = total
	}
	return res
}

func buildInputD(tc testCaseD) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.k, tc.q))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for _, qu := range tc.qs {
		sb.WriteString(fmt.Sprintf("%d %d\n", qu.idx, qu.x))
	}
	return sb.String()
}

func runCaseD(bin string, tc testCaseD) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(buildInputD(tc))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	results := []int64{}
	for scanner.Scan() {
		var v int64
		if _, err := fmt.Sscan(scanner.Text(), &v); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		results = append(results, v)
	}
	expect := solveD(tc)
	if len(results) != len(expect) {
		return fmt.Errorf("expected %d lines got %d", len(expect), len(results))
	}
	for i := range expect {
		if results[i]%MOD != expect[i]%MOD {
			return fmt.Errorf("on query %d expected %d got %d", i+1, expect[i], results[i])
		}
	}
	return nil
}

func generateCasesD() []testCaseD {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCaseD, 0, 100)
	for len(cases) < 100 {
		n := rng.Intn(5) + 2
		k := rng.Intn(5) + 1
		q := rng.Intn(5) + 1
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			arr[i] = int64(rng.Intn(10) + 1)
		}
		qs := make([]queryD, q)
		for i := 0; i < q; i++ {
			qs[i] = queryD{idx: rng.Intn(n) + 1, x: int64(rng.Intn(10) + 1)}
		}
		cases = append(cases, testCaseD{n, k, q, arr, qs})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCasesD()
	for i, tc := range cases {
		if err := runCaseD(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
