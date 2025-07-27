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

type query struct{ l, r int }

type testcase struct {
	n  int
	q  int
	s  string
	qs []query
}

func expected(tc testcase) []string {
	res := make([]string, tc.q)
	for idx, qr := range tc.qs {
		l := qr.l
		r := qr.r
		first := tc.s[l-1]
		last := tc.s[r-1]
		ok := false
		for i := 0; i < l-1; i++ {
			if tc.s[i] == first {
				ok = true
				break
			}
		}
		if !ok {
			for i := r; i < tc.n; i++ {
				if tc.s[i] == last {
					ok = true
					break
				}
			}
		}
		if ok {
			res[idx] = "YES"
		} else {
			res[idx] = "NO"
		}
	}
	return res
}

func runCase(bin string, tc testcase) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n%s\n", tc.n, tc.q, tc.s))
	for _, qr := range tc.qs {
		sb.WriteString(fmt.Sprintf("%d %d\n", qr.l, qr.r))
	}
	input := fmt.Sprintf("1\n%s", sb.String())

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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotLines := strings.Fields(out.String())
	expect := expected(tc)
	if len(gotLines) != len(expect) {
		return fmt.Errorf("expected %d lines got %d", len(expect), len(gotLines))
	}
	for i := range expect {
		if strings.ToUpper(gotLines[i]) != expect[i] {
			return fmt.Errorf("query %d expected %s got %s", i+1, expect[i], gotLines[i])
		}
	}
	return nil
}

func randomCase(rng *rand.Rand) testcase {
	n := rng.Intn(20) + 2
	q := rng.Intn(20) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	s := sb.String()
	qs := make([]query, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n-1) + 1
		r := rng.Intn(n-l) + l + 1
		qs[i] = query{l, r}
	}
	return testcase{n, q, s, qs}
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testcase, 0, 100)
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
