package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `100
5 10 ())(( 4 3 5 2 2 4 5 5 4 4
4 4 ()(( 2 1 3 1
6 8 ))))() 1 1 2 4 2 3 6 4
6 7 )))()( 3 5 6 6 2 6 3
10 10 (())(())() 2 7 3 1 5 7 7 2 1 10
2 7 )) 1 1 2 1 1 1 1
5 7 ))(() 3 3 2 4 4 4 5
8 10 ())()))) 6 1 7 6 1 7 3 1 6 8
7 6 ))((()) 6 4 3 5 5 3
4 6 ())) 3 4 1 1 2 3
10 4 )()()(())( 8 3 2 6
5 10 ))((( 5 2 3 5 2 3 3 1 5 3
4 7 )))) 4 3 4 4 1 4 2
5 1 ))(() 5
6 9 )(()(( 1 1 6 5 2 4 5 1 1
9 2 ()(()(()( 5 9
9 1 )(((((()( 1
9 10 )()())()) 1 1 3 1 2 1 2 8 1 2
10 9 ))()())))) 5 4 6 7 2 3 9 1 7
3 10 (() 2 3 3 3 2 3 1 3 2 1
7 8 ))))((( 5 3 6 5 1 7 4 2
8 3 ()))()() 2 6 2
2 8 (( 2 1 1 1 2 1 1 2
3 1 ()) 3
6 2 ((((() 5 2
3 4 (() 3 3 2 2
7 10 )))()() 7 6 2 4 6 5 7 5 6 5
9 3 )((()))(( 5 4 3
10 6 ()))()())) 4 3 10 6 4 6
9 3 ))()()()( 8 4 4
3 4 )(( 2 1 1 3
2 5 (( 2 1 2 1 1
3 2 ))( 2 2
7 1 ()))))( 2
9 7 ()()()()) 2 9 6 6 8 5 5
3 6 ()) 1 3 2 3 3 3
4 3 ())( 1 2 3
8 9 )()))((( 3 7 6 2 5 5 7 1 3
2 8 )( 2 2 2 2 1 2 2 1
4 5 (((( 2 4 4 2 2
8 4 ((()))() 7 7 6 5
9 9 ))((((()( 9 8 4 4 9 4 1 9 8
3 10 )(( 2 1 3 1 1 2 3 1 3 1
9 9 ())))(((( 6 4 2 4 7 4 8 6 2
2 7 (( 1 2 2 2 1 2 1
9 5 ())))()() 6 6 8 9 6
8 4 (())()(( 4 1 3 7
10 3 (((())((() 8 6 7
2 1 () 1
8 8 (((())(( 6 2 6 2 1 5 5 8
6 8 ()())) 1 1 6 4 1 5 5 1
3 1 ((( 3
2 8 (( 2 1 2 2 2 2 1 2
6 10 )()()( 5 4 5 3 5 6 4 2 6 6
8 9 )()))()) 6 5 3 7 8 1 3 3 1
9 2 ()((()))) 1 2
5 7 ())() 1 2 2 1 2 2 5
2 2 () 1 1
10 9 ))()(((((( 4 5 2 2 8 7 4 10 2
9 6 ))()))((( 8 5 2 6 6 4
9 2 ))((()()( 7 8
3 5 ()( 3 1 3 1 3
9 9 ))())))() 1 3 2 1 3 5 9 9 1
9 1 (()))())( 5
4 2 )))) 1 2
4 8 )))( 4 3 4 4 2 3 3 3
2 8 )) 2 1 1 1 1 1 1 1
2 1 () 2
7 9 ()((()) 3 7 1 5 5 7 1 4 3
6 4 ))))(( 4 1 6 6
4 10 (()) 2 3 1 2 1 1 4 1 4 4
10 10 ))()(((()) 8 8 5 5 8 5 9 3 10 6
4 9 ())( 4 3 1 3 1 4 2 3 3
3 7 ()( 3 2 3 2 1 2 2
4 7 ))() 2 1 3 3 4 3 1
4 7 ))() 2 2 1 1 1 1 3
3 7 ))) 2 2 3 3 1 1 1
9 3 ((())(()( 3 3 5
5 3 (((() 1 4 4
4 1 )((( 2
6 3 ))((() 2 4 2
3 6 )() 2 1 3 2 3 3
9 10 ))))())() 3 5 3 6 5 4 3 1 1 7
4 2 ((() 3 3
3 8 ))( 3 3 1 2 2 3 3 2
7 10 ))))))( 1 5 5 2 6 3 6 6 2 3
8 5 )())()(( 7 7 5 1 4
8 7 (()(()() 7 6 1 5 6 6 7
9 6 ()()(()(( 3 2 4 6 7 5
5 1 ()))( 4
7 10 ()()()) 5 2 1 7 2 5 1 4 1 7
7 7 ()())() 2 5 6 4 5 6 6
6 7 )((()) 1 1 3 4 4 6 1
7 7 ()()((( 5 5 2 5 7 1 2
5 8 ))))( 2 5 3 3 5 2 5 5
7 10 ))()()) 6 2 4 7 7 2 2 7 1 4
3 9 )() 1 1 2 2 1 2 3 3 1
5 7 ))(() 3 3 3 1 3 1 2
5 5 )((() 2 2 3 5 2
4 9 ))() 4 3 1 1 4 4 1 2 3`

type testCase struct {
	input    string
	expected string
}

type Fenwick struct {
	n int
	t []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, t: make([]int, n+1)}
}

func (f *Fenwick) Add(idx, val int) {
	for i := idx + 1; i <= f.n; i += i & -i {
		f.t[i] += val
	}
}

func (f *Fenwick) Sum(idx int) int {
	if idx < 0 {
		return 0
	}
	res := 0
	for i := idx + 1; i > 0; i -= i & -i {
		res += f.t[i]
	}
	return res
}

func (f *Fenwick) RangeSum(l, r int) int {
	if r < l {
		return 0
	}
	return f.Sum(r) - f.Sum(l-1)
}

func (f *Fenwick) FindFirst() int {
	total := f.Sum(f.n - 1)
	if total == 0 {
		return -1
	}
	idx := 0
	bit := 1 << bits.Len(uint(f.n))
	sum := 0
	for bit > 0 {
		next := idx + bit
		if next <= f.n && sum+f.t[next] == 0 {
			idx = next
			sum += f.t[next]
		}
		bit >>= 1
	}
	return idx
}

func (f *Fenwick) FindLast() int {
	total := f.Sum(f.n - 1)
	if total == 0 {
		return -1
	}
	idx := 0
	bit := 1 << bits.Len(uint(f.n))
	sum := 0
	target := total - 1
	for bit > 0 {
		next := idx + bit
		if next <= f.n && sum+f.t[next] <= target {
			idx = next
			sum += f.t[next]
		}
		bit >>= 1
	}
	return idx
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type SegTree struct {
	n   int
	sum []int
	pre []int
}

func NewSegTree(arr []int) *SegTree {
	n := 1
	for n < len(arr) {
		n <<= 1
	}
	st := &SegTree{n: n, sum: make([]int, 2*n), pre: make([]int, 2*n)}
	for i := 0; i < len(arr); i++ {
		st.sum[n+i] = arr[i]
		if arr[i] < 0 {
			st.pre[n+i] = arr[i]
		} else {
			st.pre[n+i] = 0
		}
	}
	for i := n - 1; i > 0; i-- {
		st.pull(i)
	}
	return st
}

func (st *SegTree) pull(v int) {
	l, r := v<<1, v<<1|1
	st.sum[v] = st.sum[l] + st.sum[r]
	st.pre[v] = min(st.pre[l], st.sum[l]+st.pre[r])
}

func (st *SegTree) Update(idx, val int) {
	v := st.n + idx
	st.sum[v] = val
	if val < 0 {
		st.pre[v] = val
	} else {
		st.pre[v] = 0
	}
	for v >>= 1; v > 0; v >>= 1 {
		st.pull(v)
	}
}

func (st *SegTree) query(v, l, r, L, R int) (int, int) {
	if L <= l && r <= R {
		return st.sum[v], st.pre[v]
	}
	m := (l + r) >> 1
	if R <= m {
		return st.query(v<<1, l, m, L, R)
	}
	if L > m {
		return st.query(v<<1|1, m+1, r, L, R)
	}
	sumL, preL := st.query(v<<1, l, m, L, R)
	sumR, preR := st.query(v<<1|1, m+1, r, L, R)
	return sumL + sumR, min(preL, sumL+preR)
}

func (st *SegTree) Prefix(idx int) (int, int) {
	if idx < 0 {
		return 0, 0
	}
	if idx >= st.n {
		idx = st.n - 1
	}
	return st.query(1, 0, st.n-1, 0, idx)
}

func (st *SegTree) RangeSum(l, r int) int {
	if r < l {
		return 0
	}
	s, _ := st.query(1, 0, st.n-1, l, r)
	return s
}

func solve(n int, s []byte, queries []int) string {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		if s[i] == '(' {
			arr[i] = 1
		} else {
			arr[i] = -1
		}
	}
	st := NewSegTree(arr)

	pairOpen := make([]bool, max(0, n-1))
	pairClose := make([]bool, max(0, n-1))
	fenOpen := NewFenwick(max(0, n-1))
	fenClose := NewFenwick(max(0, n-1))
	for i := 0; i < n-1; i++ {
		if s[i] == '(' && s[i+1] == '(' {
			pairOpen[i] = true
			fenOpen.Add(i, 1)
		}
		if s[i] == ')' && s[i+1] == ')' {
			pairClose[i] = true
			fenClose.Add(i, 1)
		}
	}

	check := func() bool {
		if n%2 == 1 || s[0] != '(' || s[n-1] != ')' {
			return false
		}
		total, minPref := st.Prefix(n - 1)
		if total == 0 && minPref >= 0 {
			return true
		}
		first := fenOpen.FindFirst()
		last := fenClose.FindLast()
		if first == -1 || last == -1 || first >= last {
			return false
		}
		_, pref := st.Prefix(first)
		if pref < 0 {
			return false
		}
		if st.RangeSum(last+2, n-1) > 0 {
			return false
		}
		return true
	}

	var out strings.Builder
	for _, pos := range queries {
		pos--
		if s[pos] == '(' {
			s[pos] = ')'
			st.Update(pos, -1)
		} else {
			s[pos] = '('
			st.Update(pos, 1)
		}
		for i := pos - 1; i <= pos; i++ {
			if i >= 0 && i < n-1 {
				newOpen := s[i] == '(' && s[i+1] == '('
				if pairOpen[i] != newOpen {
					if pairOpen[i] {
						fenOpen.Add(i, -1)
					} else {
						fenOpen.Add(i, 1)
					}
					pairOpen[i] = newOpen
				}
				newClose := s[i] == ')' && s[i+1] == ')'
				if pairClose[i] != newClose {
					if pairClose[i] {
						fenClose.Add(i, -1)
					} else {
						fenClose.Add(i, 1)
					}
					pairClose[i] = newClose
				}
			}
		}
		if check() {
			out.WriteString("YES\n")
		} else {
			out.WriteString("NO\n")
		}
	}
	return strings.TrimSpace(out.String())
}

func loadCases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	pos++
	cases := make([]testCase, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n", caseIdx+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseIdx+1, err)
		}
		pos++
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing q", caseIdx+1)
		}
		q, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad q: %w", caseIdx+1, err)
		}
		pos++
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing string", caseIdx+1)
		}
		str := fields[pos]
		pos++
		if len(str) != n {
			return nil, fmt.Errorf("case %d: string length mismatch", caseIdx+1)
		}
		queries := make([]int, q)
		for i := 0; i < q; i++ {
			if pos >= len(fields) {
				return nil, fmt.Errorf("case %d: missing query", caseIdx+1)
			}
			val, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad query: %w", caseIdx+1, err)
			}
			queries[i] = val
			pos++
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, q)
		sb.WriteString(str)
		sb.WriteByte('\n')
		for i, v := range queries {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: solve(n, []byte(str), queries),
		})
	}
	return cases, nil
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
		fmt.Println("usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
