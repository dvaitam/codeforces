package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Query struct {
	l, r int
	k    int64
}

type Case struct {
	n       int
	a       []int64
	queries []Query
}

func runProg(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(1807))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(6) + 1
		q := rng.Intn(5) + 1
		arr := make([]int64, n)
		for j := range arr {
			arr[j] = int64(rng.Intn(10) + 1)
		}
		queries := make([]Query, q)
		for j := range queries {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			k := int64(rng.Intn(10) + 1)
			queries[j] = Query{l, r, k}
		}
		cases[i] = Case{n: n, a: arr, queries: queries}
	}
	return cases
}

func expected(a []int64, qs []Query) []string {
	prefix := make([]int64, len(a)+1)
	for i := 1; i <= len(a); i++ {
		prefix[i] = prefix[i-1] + a[i-1]
	}
	total := prefix[len(a)]
	ans := make([]string, len(qs))
	for i, qu := range qs {
		newSum := total - (prefix[qu.r] - prefix[qu.l-1]) + int64(qu.r-qu.l+1)*qu.k
		if newSum%2 == 1 {
			ans[i] = "YES"
		} else {
			ans[i] = "NO"
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for idx, c := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("1\n%d %d\n", c.n, len(c.queries)))
		for i, v := range c.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for _, qu := range c.queries {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", qu.l, qu.r, qu.k))
		}
		expList := expected(c.a, c.queries)
		got, err := runProg(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		outputs := strings.Fields(got)
		if len(outputs) != len(expList) {
			fmt.Printf("case %d failed: expected %d lines got %d\n", idx+1, len(expList), len(outputs))
			os.Exit(1)
		}
		for i, exp := range expList {
			if strings.ToUpper(outputs[i]) != exp {
				fmt.Printf("case %d failed on query %d: expected %s got %s\n", idx+1, i+1, exp, outputs[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
