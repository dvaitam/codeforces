package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type node struct {
	x, t, s float64
	sz      int
}

type query struct {
	op   int
	a, b int
	c, d int
}

type testCase struct {
	n   int
	q   int
	arr []float64
	qs  []query
}

// Embedded testcases (from testcasesE.txt) to keep the verifier self contained.
const rawTestcases = `3 3 4 13 16 2 1 1 1 1 2 3 3 1 2 2 3 3
4 2 1 9 9 7 2 2 3 2 3 3
4 4 17 8 6 8 2 4 4 1 3 3 4 4 2 2 3 2 3 4
5 2 8 10 9 2 3 2 1 4 2 5 5
4 2 7 3 14 7 2 4 4 1 2 3 4 4
4 1 2 8 9 19 2 2 2
4 2 10 15 1 2 2 3 3 2 3 3
4 3 11 5 14 20 2 1 3 1 2 3 4 4 2 2 3
2 3 2 15 2 1 2 2 2 2 2 1 2
3 1 2 2 2 2 3 3
3 5 2 18 16 2 3 3 2 1 1 2 2 2 2 1 1 2 2 3
2 2 14 15 2 1 2 1 1 1 2 2
4 3 8 17 7 8 2 4 4 1 2 2 4 4 2 1 4
5 1 14 7 19 6 11 2 3 5
5 3 14 17 7 9 11 2 4 5 1 3 3 4 5 1 3 3 4 5
5 4 13 7 1 7 6 2 1 5 2 1 4 2 2 2 1 2 4 5 5
5 1 3 2 20 4 16 2 5 5
3 1 12 3 17 2 1 2
4 1 3 18 15 13 2 2 3
5 2 16 13 4 3 4 2 5 5 2 4 5
2 2 10 16 2 2 2 1 1 1 2 2
4 4 11 9 18 1 2 2 2 2 1 1 2 4 4 1 2 3 4 4
3 1 2 20 11 2 3 3
4 5 15 20 20 15 2 4 4 2 3 4 1 2 2 3 3 2 3 3 2 1 1
5 2 11 12 11 6 10 2 1 5 1 1 3 4 4
3 5 16 19 3 2 1 1 2 3 3 2 3 3 1 2 2 3 3 1 2 2 3 3
4 3 15 14 19 1 2 3 3 2 3 4 2 1 2
2 4 12 20 2 2 2 2 2 2 2 1 1 1 1 1 2 2
3 4 19 20 18 2 3 3 1 1 1 2 2 1 2 2 3 3 1 2 2 3 3
3 1 18 12 1 2 3 3
5 2 7 10 16 17 3 2 4 4 1 3 4 5 5
4 2 13 2 14 1 2 3 3 2 2 2
3 1 20 1 8 2 1 3
2 4 18 6 2 2 2 2 1 1 1 1 1 2 2 1 1 1 2 2
4 1 8 11 4 2 2 3 4
4 1 6 2 16 17 2 2 2
3 1 11 20 5 2 3 3
2 5 17 1 2 1 1 2 1 1 1 1 1 2 2 1 1 1 2 2 1 1 1 2 2
4 5 10 7 14 11 2 2 2 1 2 3 4 4 2 4 4 1 2 3 4 4 2 3 4
3 5 11 3 6 2 3 3 1 2 2 3 3 2 1 2 1 1 1 2 2 1 2 2 3 3
3 1 6 7 9 2 1 3
4 3 20 12 6 6 2 2 2 1 1 3 4 4 2 1 2
3 4 16 20 18 2 2 3 2 2 3 1 1 1 2 3 1 2 2 3 3
4 4 12 18 4 8 2 2 3 1 3 3 4 4 1 1 2 4 4 1 2 3 4 4
5 5 15 12 12 1 15 2 2 5 1 4 4 5 5 2 5 5 1 4 4 5 5 1 3 3 4 5
5 2 19 5 20 2 5 2 2 4 1 4 4 5 5
4 5 19 12 9 20 2 3 3 1 2 2 4 4 1 1 3 4 4 1 3 3 4 4 1 3 3 4 4
4 5 2 10 16 9 2 2 4 2 3 3 1 2 3 4 4 2 3 3 1 2 2 4 4
4 2 7 6 7 3 2 1 1 1 2 3 4 4
4 2 7 17 18 2 2 1 1 1 2 2 3 4
4 3 5 19 19 4 2 1 1 1 1 2 4 4 2 2 4
4 5 15 8 20 15 2 2 3 1 3 3 4 4 2 2 3 1 2 3 4 4 2 3 4
4 2 11 4 8 11 2 2 2 2 1 3
2 5 1 2 2 1 1 1 1 1 2 2 1 1 1 2 2 2 1 1 1 1 1 2 2
3 4 2 18 8 2 2 3 2 3 3 1 2 2 3 3 1 2 2 3 3
5 2 3 15 13 17 7 2 1 1 2 1 2
3 4 15 12 14 2 2 3 2 2 2 2 2 2 2 2 3
3 5 20 6 14 2 2 2 2 3 3 1 1 2 3 3 2 2 3 1 2 2 3 3
5 3 7 5 11 12 5 2 3 5 2 5 5 2 4 5
4 5 17 2 14 9 2 2 4 1 1 2 4 4 1 1 1 3 3 2 3 4 1 2 2 3 3
5 1 8 8 16 18 3 2 5 5
3 2 1 7 17 2 1 1 1 2 2 3 3
5 5 5 11 4 20 9 2 3 3 1 3 3 5 5 1 2 2 4 5 2 1 1 1 4 4 5 5
3 2 2 6 15 2 3 3 2 2 3
3 3 9 19 1 2 2 3 2 3 3 1 1 1 2 3
3 2 20 7 18 2 1 3 1 2 2 3 3
2 5 1 9 2 2 2 1 1 1 2 2 2 1 1 2 2 2 2 2 2
5 2 6 6 11 2 7 2 2 5 2 2 3
4 5 17 14 15 12 2 3 3 2 4 4 2 2 2 2 4 4 2 1 3
4 5 4 20 10 20 2 1 2 2 2 4 2 4 4 1 2 2 3 4 1 1 2 4 4
5 3 11 14 7 2 19 2 1 5 1 1 4 5 5 1 2 3 4 5
5 3 8 12 15 12 13 2 5 5 2 2 3 2 4 5
3 4 17 12 12 2 3 3 1 1 2 3 3 2 3 3 2 3 3
2 4 14 9 2 2 2 1 1 1 2 2 2 1 2 2 2 2
2 3 9 9 2 2 2 2 1 1 2 2 2
3 5 4 4 16 2 3 3 2 3 3 2 2 3 2 3 3 2 3 3
4 3 9 3 15 17 2 2 3 2 3 4 2 1 3
5 3 18 12 5 8 8 2 5 5 1 3 4 5 5 2 5 5
5 2 13 10 11 5 9 2 4 4 1 4 4 5 5
2 3 20 1 2 1 2 2 2 2 1 1 1 2 2
3 3 12 3 12 2 3 3 2 3 3 2 1 1
4 4 5 18 10 15 2 2 2 1 3 3 4 4 1 3 3 4 4 1 3 3 4 4
4 4 4 2 12 17 2 3 3 2 1 3 1 3 3 4 4 1 1 1 2 4
2 2 19 9 2 1 1 2 2 2
5 4 5 4 17 11 10 2 5 5 2 1 5 1 4 4 5 5 1 3 3 5 5
5 4 15 17 16 18 16 2 3 5 1 2 4 5 5 2 5 5 2 1 3
3 4 7 18 14 2 2 2 1 2 2 3 3 1 1 1 2 3 1 2 2 3 3
3 2 6 3 9 2 1 3 2 1 2
2 5 3 14 2 2 2 2 2 2 1 1 1 2 2 2 2 2 1 1 1 2 2
3 4 11 14 5 2 2 3 1 2 2 3 3 1 1 1 3 3 2 3 3
4 2 5 3 7 14 2 2 2 1 1 2 4 4
2 2 17 7 2 2 2 2 1 2
5 2 13 10 5 14 6 2 3 3 1 2 2 4 5
5 3 7 6 11 3 9 2 5 5 2 4 5 2 5 5
5 3 14 16 11 3 14 2 4 4 1 1 4 5 5 1 4 4 5 5
5 5 7 8 13 3 5 2 3 3 1 3 3 4 5 2 4 4 1 1 1 2 2 2 1 2
3 3 14 4 5 2 2 2 2 1 3 2 3 3
5 2 16 6 20 13 16 2 5 5 1 3 4 5 5
2 1 3 11 2 1 1
5 5 5 16 15 14 17 2 3 3 2 2 5 2 2 3 1 1 4 5 5 2 3 3`

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// build and solve the segment tree problem for one test case.
func solve(tc testCase) string {
	n := tc.n
	tree := make([]node, 4*n+5)

	var build func(i, l, r int, arr []float64)
	build = func(i, l, r int, arr []float64) {
		tree[i].t = 1
		tree[i].s = 0
		tree[i].sz = r - l
		if r-l == 1 {
			tree[i].x = arr[l]
			return
		}
		m := (l + r) >> 1
		build(i<<1, l, m, arr)
		build(i<<1|1, m, r, arr)
		tree[i].x = tree[i<<1].x + tree[i<<1|1].x
	}

	applyUpd := func(nd *node, tt, ts float64) {
		nd.x = tt*nd.x + ts*float64(nd.sz)
		nd.t *= tt
		nd.s = nd.s*tt + ts
	}

	var push func(i int)
	push = func(i int) {
		tt := tree[i].t
		ts := tree[i].s
		if tt != 1 || ts != 0 {
			applyUpd(&tree[i<<1], tt, ts)
			applyUpd(&tree[i<<1|1], tt, ts)
			tree[i].t = 1
			tree[i].s = 0
		}
	}

	var query func(i, l, r, ql, qr int) float64
	query = func(i, l, r, ql, qr int) float64 {
		if ql <= l && r <= qr {
			return tree[i].x
		}
		push(i)
		m := (l + r) >> 1
		var ans float64
		if ql < m && qr > l {
			ans += query(i<<1, l, m, ql, qr)
		}
		if qr > m && ql < r {
			ans += query(i<<1|1, m, r, ql, qr)
		}
		return ans
	}

	var modify func(i, l, r, ql, qr int, tt, ts float64)
	modify = func(i, l, r, ql, qr int, tt, ts float64) {
		if ql <= l && r <= qr {
			applyUpd(&tree[i], tt, ts)
			return
		}
		push(i)
		m := (l + r) >> 1
		if ql < m && qr > l {
			modify(i<<1, l, m, ql, qr, tt, ts)
		}
		if qr > m && ql < r {
			modify(i<<1|1, m, r, ql, qr, tt, ts)
		}
		tree[i].x = tree[i<<1].x + tree[i<<1|1].x
	}

	build(1, 0, n, tc.arr)

	var out []string
	for _, qu := range tc.qs {
		if qu.op == 1 {
			a := qu.a - 1
			b := qu.b
			c := qu.c - 1
			d := qu.d
			len1 := b - a
			len2 := d - c
			sum2 := query(1, 0, n, c, d)
			sum1 := query(1, 0, n, a, b)
			s1 := sum2 / float64(len2) / float64(len1)
			s2 := sum1 / float64(len2) / float64(len1)
			t1 := 1 - 1/float64(len1)
			t2 := 1 - 1/float64(len2)
			modify(1, 0, n, a, b, t1, s1)
			modify(1, 0, n, c, d, t2, s2)
		} else {
			a := qu.a - 1
			b := qu.b
			sum := query(1, 0, n, a, b)
			out = append(out, fmt.Sprintf("%.7f", sum))
		}
	}
	return strings.Join(out, "\n")
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(rawTestcases, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		toks := strings.Fields(line)
		if len(toks) < 2 {
			return nil, fmt.Errorf("line %d: too few tokens", idx+1)
		}
		n, err := strconv.Atoi(toks[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		q, err := strconv.Atoi(toks[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse q: %w", idx+1, err)
		}
		if len(toks) < 2+n {
			return nil, fmt.Errorf("line %d: not enough values for array", idx+1)
		}
		arr := make([]float64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(toks[2+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse arr[%d]: %w", idx+1, i, err)
			}
			arr[i] = float64(v)
		}
		pos := 2 + n
		var qs []query
		for qi := 0; qi < q; qi++ {
			if pos >= len(toks) {
				return nil, fmt.Errorf("line %d: missing query %d", idx+1, qi+1)
			}
			op, err := strconv.Atoi(toks[pos])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse op: %w", idx+1, err)
			}
			pos++
			if op == 1 {
				if pos+3 >= len(toks) {
					return nil, fmt.Errorf("line %d: query %d incomplete", idx+1, qi+1)
				}
				a, _ := strconv.Atoi(toks[pos])
				b, _ := strconv.Atoi(toks[pos+1])
				c, _ := strconv.Atoi(toks[pos+2])
				d, _ := strconv.Atoi(toks[pos+3])
				pos += 4
				qs = append(qs, query{op: op, a: a, b: b, c: c, d: d})
			} else {
				if pos+1 >= len(toks) {
					return nil, fmt.Errorf("line %d: query %d incomplete", idx+1, qi+1)
				}
				a, _ := strconv.Atoi(toks[pos])
				b, _ := strconv.Atoi(toks[pos+1])
				pos += 2
				qs = append(qs, query{op: op, a: a, b: b})
			}
		}
		cases = append(cases, testCase{n: n, q: q, arr: arr, qs: qs})
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(int64(v), 10))
	}
	sb.WriteByte('\n')
	for _, qu := range tc.qs {
		if qu.op == 1 {
			sb.WriteString(fmt.Sprintf("1 %d %d %d %d\n", qu.a, qu.b, qu.c, qu.d))
		} else {
			sb.WriteString(fmt.Sprintf("2 %d %d\n", qu.a, qu.b))
		}
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expected := solve(tc)
		input := buildInput(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}
