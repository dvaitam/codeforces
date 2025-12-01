package main

import (
	"bytes"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE.txt.
const testcasesEData = `
100
1
31 3
2
2
29 1
13 3
2
3
31 3
13 1
29 2
3
3
13 3
3 2
23 1
1
1
3 2
0
1
19 3
0
2
19 2
29 2
1
2
47 2
19 1
3
1
23 1
2
2
5 1
13 2
1
1
31 3
0
3
41 3
5 2
23 3
2
2
29 2
43 3
0
1
5 1
3
2
3 2
29 3
2
3
7 2
31 1
37 2
3
1
29 3
2
2
2 3
11 1
1
1
17 1
0
3
5 3
3 2
47 3
2
1
13 1
3
3
23 1
7 1
3 3
3
3
31 2
2 2
47 1
2
3
5 3
47 2
3 3
2
2
7 3
47 2
0
1
41 1
0
1
23 3
2
3
13 1
11 2
3 3
3
3
3 1
7 1
43 3
0
1
31 3
2
2
19 1
3 3
0
3
7 3
29 1
43 2
1
3
5 1
37 2
29 1
2
2
13 1
3 3
3
2
41 3
29 1
0
2
23 2
7 3
1
3
5 2
11 3
23 3
1
3
17 1
47 1
23 1
0
1
41 3
2
2
11 3
31 3
2
2
43 2
19 3
2
3
13 3
3 1
43 2
0
3
5 3
31 3
29 1
1
3
41 3
5 2
47 1
3
3
3 2
7 3
19 1
3
2
37 1
17 3
0
3
13 1
31 2
43 2
0
1
19 1
1
3
43 3
13 1
7 3
0
3
7 3
47 3
37 3
0
3
41 1
19 3
7 2
1
2
5 2
37 1
2
2
29 2
41 3
1
1
7 2
3
2
43 1
47 2
1
2
47 2
11 3
1
2
47 2
37 1
1
1
29 3
3
2
3 3
41 1
3
1
41 1
1
2
3 3
2 3
3
2
3 1
17 1
2
1
3 1
3
3
19 2
17 2
31 2
3
3
19 1
7 3
29 1
0
2
29 1
2 1
0
1
29 2
3
3
41 3
2 1
47 3
0
3
7 1
5 1
29 1
2
3
19 2
31 1
41 1
1
2
43 3
29 2
2
2
13 3
17 2
0
2
19 1
41 2
1
3
3 3
41 3
2 3
0
3
43 1
29 2
17 2
2
2
11 1
43 1
2
3
11 2
7 2
29 2
1
1
41 2
2
2
31 3
23 2
3
3
7 2
2 3
13 3
1
2
7 2
13 3
2
1
31 3
3
1
47 1
2
3
23 3
29 2
19 3
2
1
11 2
1
1
41 1
1
2
43 3
19 1
0
3
41 2
7 3
43 2
3
2
31 1
41 3
2
1
43 3
3
3
23 3
5 2
41 2
0
2
2 1
31 1
0
3
7 3
13 1
23 1
0
2
2 1
19 1
3
2
37 2
43 2
3
1
31 2
3
3
13 2
23 2
43 1
0
1
19 2
1
2
19 3
13 1
2
3
43 3
13 2
17 2
2
`

// ---------- reference solution (adapted from 396E.go) ----------

const maxP = 1000000

var spf []int

func buildSPF() {
	spf = make([]int, maxP+1)
	for i := 2; i <= maxP; i++ {
		if spf[i] == 0 {
			for j := i; j <= maxP; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
}

type testCase struct {
	primes []primeExp
	k      *big.Int
}

type primeExp struct {
	p int
	e *big.Int
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesEData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no data")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("parse T: %w", err)
	}
	pos := 1
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing m", i+1)
		}
		m, _ := strconv.Atoi(fields[pos])
		pos++
		prs := make([]primeExp, m)
		for j := 0; j < m; j++ {
			if pos+1 >= len(fields) {
				return nil, fmt.Errorf("case %d: missing prime/exp", i+1)
			}
			p, _ := strconv.Atoi(fields[pos])
			pos++
			expStr := fields[pos]
			pos++
			e := new(big.Int)
			if _, ok := e.SetString(expStr, 10); !ok {
				return nil, fmt.Errorf("case %d: bad exp %q", i+1, expStr)
			}
			prs[j] = primeExp{p: p, e: e}
		}
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing k", i+1)
		}
		kVal := new(big.Int)
		if _, ok := kVal.SetString(fields[pos], 10); !ok {
			return nil, fmt.Errorf("case %d: bad k %q", i+1, fields[pos])
		}
		pos++
		cases = append(cases, testCase{primes: prs, k: kVal})
	}
	return cases, nil
}

func solve(tc testCase) string {
	initE := make(map[int]*big.Int)
	primesList := make([]int, 0, len(tc.primes))
	for _, pe := range tc.primes {
		initE[pe.p] = new(big.Int).Set(pe.e)
		primesList = append(primesList, pe.p)
	}

	ans := make(map[int]*big.Int)
	seen := make(map[int]bool)

	curr := make([]int, 0, len(primesList))
	currE := make(map[int]*big.Int)
	for _, p := range primesList {
		curr = append(curr, p)
		currE[p] = new(big.Int).Set(initE[p])
		seen[p] = true
	}

	h := 0
	for len(curr) > 0 {
		if big.NewInt(int64(h)).Cmp(tc.k) > 0 {
			break
		}
		kh := new(big.Int).Sub(tc.k, big.NewInt(int64(h)))
		nextE := make(map[int]*big.Int)

		for _, p := range curr {
			e0 := currE[p]
			rem := new(big.Int).Sub(e0, kh)
			if rem.Sign() > 0 {
				if _, ok := ans[p]; !ok {
					ans[p] = new(big.Int)
				}
				ans[p].Add(ans[p], rem)
			}
			A := new(big.Int)
			if e0.Cmp(kh) <= 0 {
				A.Set(e0)
			} else {
				A.Set(kh)
			}
			if A.Sign() <= 0 {
				continue
			}
			x := p - 1
			for x > 1 {
				f := spf[x]
				cnt := 0
				for x%f == 0 {
					x /= f
					cnt++
				}
				add := new(big.Int).Mul(A, big.NewInt(int64(cnt)))
				if _, ok := nextE[f]; !ok {
					nextE[f] = new(big.Int)
				}
				nextE[f].Add(nextE[f], add)
			}
		}

		curr = curr[:0]
		currE = make(map[int]*big.Int)
		for p, e := range nextE {
			if !seen[p] {
				seen[p] = true
				curr = append(curr, p)
			} else {
				curr = append(curr, p)
			}
			currE[p] = e
		}
		h++
	}

	outPr := make([]int, 0, len(ans))
	for p := range ans {
		outPr = append(outPr, p)
	}
	sort.Ints(outPr)

	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(outPr)))
	sb.WriteByte('\n')
	for _, p := range outPr {
		sb.WriteString(strconv.Itoa(p))
		sb.WriteByte(' ')
		sb.WriteString(ans[p].String())
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.primes)))
	for _, pe := range tc.primes {
		sb.WriteString(fmt.Sprintf("%d %s\n", pe.p, pe.e.String()))
	}
	sb.WriteString(tc.k.String())
	sb.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	buildSPF()
	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\n\ngot:\n%s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
