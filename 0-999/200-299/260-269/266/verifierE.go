package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE.txt so the verifier is self-contained.
const testcasesRaw = `4 5
3 12 15 4
1 4 17
? 1 2 4
? 3 3 0
1 3 8
3 4 20

6 2
19 10 12 16 7 5
3 3 17
? 1 3 4

5 9
6 13 13 19 9
2 3 9
1 1 1
3 5 17
? 3 3 5
4 4 20
? 3 3 2
5 5 20
? 3 3 0
? 3 5 4

4 2
10 5 9 14
3 3 9
? 3 3 2

5 6
4 20 13 19 2
2 5 9
4 4 10
? 3 3 3
3 4 18
2 5 6
1 1 5

10 3
19 1 17 15 18 7 10 1 3 16
7 8 15
8 9 15
7 10 7

7 4
15 6 1 1 8 8 7
? 7 7 3
? 2 4 0
? 3 7 0
? 4 4 3

7 2
13 6 18 5 10 9 15
? 6 7 3
? 6 7 2

7 8
2 8 20 6 1 12 19
3 3 5
? 4 7 5
4 5 0
2 2 19
? 1 4 3
? 5 5 1
5 6 17
? 4 7 0

2 1
19 3
2 2 1

6 2
16 0 9 11 2 2
? 4 5 1
? 3 6 1

8 7
3 2 3 19 11 16 13 13
? 8 8 5
? 4 6 3
3 5 5
3 5 15
? 5 5 5
5 5 17
8 8 16

4 7
9 11 7 5
? 1 1 4
? 4 4 4
? 4 4 1
3 4 4
2 4 19
3 3 18
? 1 1 3

3 6
20 11 10
1 3 0
? 1 2 0
3 3 18
? 1 1 1
? 3 3 4
? 3 3 5

7 10
7 15 7 9 11 7 10
? 5 6 3
? 7 7 2
5 5 8
5 6 17
? 3 6 3
? 2 6 0
6 7 19
? 1 7 2
6 7 19
? 2 3 0

5 3
1 5 12 18 19
? 5 5 0
? 2 2 1
3 3 3

5 1
20 9 19 20 5
? 4 4 5

2 9
11 0
? 1 2 1
2 2 16
1 1 8
? 2 2 2
1 2 0
1 2 4
1 1 7
1 2 17
2 2 10

3 9
18 6 3
2 3 20
? 2 2 5
3 3 15
? 2 2 1
? 1 3 5
2 3 11
? 1 3 0
3 3 1
1 2 19

3 5
3 16 16
? 1 1 4
1 3 1
3 3 6
3 3 9
? 3 3 0

3 7
15 2 16
? 2 3 1
? 2 3 1
2 3 8
2 3 10
? 1 2 2
? 1 2 2
? 3 3 4

6 2
5 19 2 6 8 6
? 4 4 5
4 6 6

3 2
7 16 3
2 3 18
? 1 3 0

3 1
5 6 8
2 3 19

6 3
5 6 1 1 5 1
? 3 6 0
? 2 5 3
? 4 6 2

6 7
12 1 1 5 4 11
2 4 18
? 5 6 3
6 6 7
3 4 17
1 1 6
6 6 17
1 3 12`

const mod = 1000000007

type node struct {
	t    [6]int64
	lazy int64
	set  bool
}

type testCase struct {
	n  int
	m  int
	a  []int64
	qs []string
}

func add(a, b int64) int64 {
	x := a + b
	if x >= mod {
		x -= mod
	}
	return x
}

func mul(a, b int64) int64 {
	return (a * b) % mod
}

func powMod(a, e int) int64 {
	res := int64(1)
	base := int64(a % mod)
	for e > 0 {
		if e&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		e >>= 1
	}
	return res
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for i := 0; i < len(lines); {
		for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
			i++
		}
		if i >= len(lines) {
			break
		}
		var n, m int
		if _, err := fmt.Sscan(lines[i], &n, &m); err != nil {
			return nil, fmt.Errorf("line %d: parse n m: %w", i+1, err)
		}
		i++
		if i >= len(lines) {
			return nil, fmt.Errorf("line %d: missing array", i+1)
		}
		arrFields := strings.Fields(lines[i])
		if len(arrFields) != n {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", i+1, n, len(arrFields))
		}
		a := make([]int64, n+1)
		for idx, v := range arrFields {
			val, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse value %d: %w", i+1, idx+1, err)
			}
			a[idx+1] = val
		}
		i++
		qs := make([]string, 0, m)
		for j := 0; j < m; j++ {
			if i >= len(lines) {
				return nil, fmt.Errorf("line %d: missing query %d", i+1, j+1)
			}
			line := strings.TrimSpace(lines[i])
			if line == "" {
				return nil, fmt.Errorf("line %d: empty query", i+1)
			}
			qs = append(qs, line)
			i++
		}
		cases = append(cases, testCase{n: n, m: m, a: a, qs: qs})
	}
	return cases, nil
}

func solveCase(tc testCase) ([]int64, error) {
	n := tc.n
	m := tc.m
	a := tc.a
	P := [6][]int64{}
	for j := 0; j <= 5; j++ {
		P[j] = make([]int64, n+1)
		for i := 1; i <= n; i++ {
			P[j][i] = add(P[j][i-1], mul(powMod(i, j), 1))
		}
	}
	var C [6][6]int64
	for i := 0; i <= 5; i++ {
		C[i][0] = 1
		for j := 1; j <= i; j++ {
			C[i][j] = add(C[i-1][j-1], C[i-1][j])
		}
	}

	st := make([]node, 4*(n+5))

	var build func(int, int, int)
	build = func(node, l, r int) {
		if l == r {
			ai := a[l] % mod
			p := int64(1)
			for j := 0; j <= 5; j++ {
				st[node].t[j] = ai * p % mod
				p = p * int64(l) % mod
			}
			return
		}
		mid := (l + r) >> 1
		lc, rc := node<<1, node<<1|1
		build(lc, l, mid)
		build(rc, mid+1, r)
		for j := 0; j <= 5; j++ {
			st[node].t[j] = add(st[lc].t[j], st[rc].t[j])
		}
	}

	var applySet func(int, int, int, int64)
	applySet = func(node, l, r int, x int64) {
		st[node].set = true
		st[node].lazy = x
		for j := 0; j <= 5; j++ {
			sum := (P[j][r] - P[j][l-1]) % mod
			if sum < 0 {
				sum += mod
			}
			st[node].t[j] = mul(x%mod, sum)
		}
	}

	var push func(int, int, int)
	push = func(node, l, r int) {
		if !st[node].set || l == r {
			return
		}
		mid := (l + r) >> 1
		lc, rc := node<<1, node<<1|1
		applySet(lc, l, mid, st[node].lazy)
		applySet(rc, mid+1, r, st[node].lazy)
		st[node].set = false
	}

	var update func(int, int, int, int, int, int64)
	update = func(node, l, r, ql, qr int, x int64) {
		if ql <= l && r <= qr {
			applySet(node, l, r, x)
			return
		}
		push(node, l, r)
		mid := (l + r) >> 1
		if ql <= mid {
			update(node<<1, l, mid, ql, qr, x)
		}
		if qr > mid {
			update(node<<1|1, mid+1, r, ql, qr, x)
		}
		for j := 0; j <= 5; j++ {
			st[node].t[j] = add(st[node<<1].t[j], st[node<<1|1].t[j])
		}
	}

	var query func(int, int, int, int, int, *[6]int64)
	query = func(node, l, r, ql, qr int, res *[6]int64) {
		if ql <= l && r <= qr {
			for j := 0; j <= 5; j++ {
				res[j] = add(res[j], st[node].t[j])
			}
			return
		}
		push(node, l, r)
		mid := (l + r) >> 1
		if ql <= mid {
			query(node<<1, l, mid, ql, qr, res)
		}
		if qr > mid {
			query(node<<1|1, mid+1, r, ql, qr, res)
		}
	}

	build(1, 1, n)

	results := make([]int64, 0, m)
	for idx, raw := range tc.qs {
		fields := strings.Fields(raw)
		if len(fields) == 0 {
			return nil, fmt.Errorf("query %d: empty", idx+1)
		}
		if fields[0] == "?" {
			if len(fields) != 4 {
				return nil, fmt.Errorf("query %d: bad format", idx+1)
			}
			l, _ := strconv.Atoi(fields[1])
			r, _ := strconv.Atoi(fields[2])
			k, _ := strconv.Atoi(fields[3])
			var res [6]int64
			query(1, 1, n, l, r, &res)
			pw := make([]int64, k+1)
			pw[0] = 1
			neg := (mod - int64(l-1)%mod) % mod
			for t := 1; t <= k; t++ {
				pw[t] = pw[t-1] * neg % mod
			}
			var ans int64
			for j := 0; j <= k; j++ {
				coef := C[k][j] * pw[k-j] % mod
				ans = (ans + res[j]*coef) % mod
			}
			results = append(results, ans)
		} else {
			if len(fields) != 3 {
				return nil, fmt.Errorf("query %d: bad update format", idx+1)
			}
			l, _ := strconv.Atoi(fields[0])
			r, _ := strconv.Atoi(fields[1])
			x, _ := strconv.ParseInt(fields[2], 10, 64)
			update(1, 1, n, l, r, x)
		}
	}
	return results, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i := 1; i <= tc.n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", tc.a[i]))
	}
	sb.WriteByte('\n')
	for _, q := range tc.qs {
		sb.WriteString(q)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range testcases {
		expectVals, err := solveCase(tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: solver error: %v\n", idx+1, err)
			os.Exit(1)
		}
		var expectBuf bytes.Buffer
		for _, v := range expectVals {
			fmt.Fprint(&expectBuf, v, '\n')
		}
		expect := strings.TrimSpace(expectBuf.String())
		input := buildInput(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected:\n%s\ngot:\n%s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
