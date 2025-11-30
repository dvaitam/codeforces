package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesD.txt so the verifier is self-contained.
const testcasesRaw = `13 98 54 6 34 66 63 52 39 62 46 75 28 65 2 5 7 13 13
20 33 69 91 78 19 40 13 94 10 88 43 61 72 13 46 56 41 79 82 27 5 16 19 17 19 2 19 1 3 13 13
20 64 43 32 94 42 91 9 25 73 29 31 19 70 58 12 11 41 66 63 14 3 18 19 4 14 18 18
20 71 76 37 57 12 77 50 41 74 31 38 24 25 24 5 79 85 34 61 9 1 5 9
2 11 90 5 2 2 1 1 2 2 2 2 2 2
11 79 15 63 76 81 43 25 32 3 94 35 1 4 9
6 43 55 8 13 19 90 2 1 5 6 6
1 16 2 1 1 1 1
4 5 78 3 25 2 1 4 2 4
2 87 3 5 2 2 2 2 1 1 2 2 2 2
2 65 60 1 1 2
7 34 46 94 61 73 22 90 2 7 7 7 7
6 44 68 33 16 77 57 2 1 4 6 6
19 66 40 84 46 50 85 33 20 72 89 2 59 95 11 43 95 6 70 36 2 8 15 12 16
12 76 82 80 17 92 40 50 96 54 84 11 1 5 4 9 3 6 4 11 7 12 11 12
2 52 90 5 2 2 1 2 1 2 1 2 2 2
2 64 42 3 2 2 2 2 1 1
1 52 4 1 1 1 1 1 1 1 1
20 84 26 39 36 89 24 13 61 51 81 11 3 36 58 15 33 18 84 67 84 3 4 8 9 9 2 3
7 88 34 72 41 47 73 6 5 6 7 6 7 6 7 3 7 2 3
13 76 38 2 18 20 35 43 44 48 92 12 44 100 5 1 1 5 7 3 12 5 10 7 11
5 38 15 62 94 31 1 3 3
17 94 10 39 52 43 39 54 14 13 72 62 61 44 44 16 62 15 4 14 14 10 15 5 7 13 13
3 11 26 96 2 1 2 1 1
13 72 67 38 58 63 75 92 87 28 55 11 48 29 3 10 11 7 8 6 7
3 90 4 68 4 3 3 1 2 2 3 1 3
2 28 80 2 1 1 2 2
12 70 20 14 77 63 19 73 52 82 88 55 67 4 11 12 8 11 11 11 9 10
1 44 3 1 1 1 1 1 1
19 38 92 91 61 9 11 67 6 9 29 17 6 39 2 98 58 43 21 20 4 12 18 17 19 2 4 17 19
3 96 55 97 2 2 3 2 3
20 76 30 3 85 1 95 24 39 65 73 33 43 9 64 34 39 99 53 50 50 1 6 16
5 31 37 94 43 8 1 4 5
5 63 78 92 11 87 2 3 4 1 5
15 50 59 7 13 61 100 20 3 5 77 80 17 81 42 14 5 11 13 4 10 13 14 2 2 10 15
15 79 81 44 84 16 88 92 80 38 17 50 38 96 88 16 5 14 14 1 13 7 14 6 9 8 13
3 6 6 63 3 1 3 3 3 1 1
17 90 68 54 65 40 15 19 55 73 55 11 14 54 9 13 54 100 2 1 15 14 17
1 64 3 1 1 1 1 1 1
1 45 3 1 1 1 1 1 1
7 1 27 85 87 94 16 96 1 3 5
1 78 2 1 1 1 1
16 45 91 34 17 4 27 47 43 61 38 38 71 82 42 24 76 1 4 12
19 40 21 49 19 17 29 41 66 32 31 97 24 38 48 54 85 6 17 77 1 13 13
3 17 54 39 5 2 2 3 3 2 3 1 1 2 3
17 8 49 53 2 54 94 42 57 27 48 38 61 12 24 14 36 15 5 5 17 15 16 6 12 14 15 8 15
11 67 19 46 60 81 82 12 62 97 27 38 1 8 11
1 28 3 1 1 1 1 1 1
16 98 30 70 98 52 36 81 3 16 35 86 6 1 33 51 68 5 13 16 4 15 9 14 10 16 7 16
3 5 10 34 3 3 3 1 3 1 1
3 54 38 37 5 1 3 3 3 3 3 2 3 3 3
10 57 48 73 81 18 21 16 90 16 49 4 10 10 3 7 6 9 7 8
16 63 89 65 41 64 84 8 57 39 19 96 64 7 80 28 4 3 16 16 1 3 3 13
13 1 47 6 15 80 1 35 82 90 38 94 30 19 5 5 8 2 8 8 13 6 12 3 8
14 83 88 56 19 58 91 19 68 41 17 27 24 57 45 4 7 14 7 10 13 13 8 9
19 91 7 50 5 30 82 11 24 47 8 96 82 87 23 30 79 39 79 12 5 10 15 14 17 2 18 18 19 19 19
16 33 91 61 28 44 35 6 6 7 21 45 1 38 84 1 18 1 14 16
8 78 51 72 29 59 25 44 78 1 2 8
11 42 69 59 42 33 4 67 6 25 48 11 2 9 10 4 7
9 87 94 95 39 40 67 50 33 62 3 4 4 5 9 2 2
15 64 93 57 7 53 64 59 57 16 11 11 31 13 98 20 4 15 15 8 9 14 15 9 15
13 6 24 32 63 29 17 36 46 41 56 14 72 37 5 9 10 12 13 13 13 9 13 8 12
9 35 30 3 16 79 92 13 23 94 4 4 5 5 5 9 9 1 2
13 83 35 16 95 73 46 30 87 92 91 70 85 37 2 12 12 2 10
10 87 42 30 48 81 62 37 75 22 18 1 9 10
12 75 82 4 17 51 20 23 66 10 18 98 27 4 10 12 4 7 12 12 4 10
12 78 76 17 81 64 14 79 4 68 77 46 63 4 5 5 4 12 11 11 11 12
16 70 41 91 11 34 18 78 52 91 25 41 38 50 8 27 5 3 8 13 15 15 9 14
6 40 3 46 74 70 8 2 3 3 4 6
2 4 31 1 1 1
11 9 8 45 85 55 18 28 58 56 19 46 3 3 8 7 10 1 7
9 69 69 95 88 91 60 98 6 73 1 7 8
6 1 65 18 80 85 66 2 1 3 2 3
8 3 22 96 88 72 22 92 11 4 2 6 8 8 1 5 6 8
13 4 81 5 64 12 46 38 86 20 59 31 65 46 2 12 13 6 10
16 51 2 40 68 37 71 61 5 99 69 74 71 34 88 5 59 4 4 10 12 15 2 2 9 9
9 88 88 75 90 100 38 88 98 27 5 9 9 7 8 4 4 6 7 9 9
6 20 43 96 2 75 7 5 2 4 3 5 6 6 3 6 4 6
14 22 1 19 73 6 57 17 44 2 93 62 86 85 100 3 12 14 4 5 9 12
9 23 68 22 9 85 82 21 75 15 5 9 9 7 8 5 7 1 7 5 7
18 68 71 41 44 25 91 56 19 1 97 66 20 85 90 73 50 47 60 1 18 18
20 98 30 3 47 68 21 88 25 81 46 81 90 64 3 94 95 32 74 31 36 2 14 14 19 20
8 96 92 58 66 91 13 25 22 4 2 5 7 8 5 8 6 8
11 12 40 4 64 2 98 33 26 98 51 50 4 11 11 1 10 8 10 10 10
19 91 36 42 4 51 61 67 18 6 11 73 45 47 1 9 25 92 15 86 5 16 16 11 11 11 17 5 17 9 15
5 77 19 52 40 66 1 2 3
5 62 91 83 92 98 1 5 5
18 89 96 90 81 50 24 45 75 11 11 72 23 34 26 34 42 90 91 3 9 17 15 16 15 16
2 81 75 2 1 2 1 1
15 79 31 59 67 21 91 43 84 18 61 99 72 8 70 11 5 6 6 14 14 2 8 10 12 10 13
11 49 66 47 82 16 18 41 3 24 94 17 1 6 10
7 6 53 83 8 90 40 50 1 5 7
6 46 10 53 7 57 46 5 5 6 6 6 5 6 4 4 1 4
9 25 50 9 46 13 16 4 45 3 2 7 9 1 6
15 71 96 89 64 61 11 7 69 52 34 4 83 67 13 11 3 6 7 8 8 3 11
10 5 1 49 44 21 71 90 19 21 23 2 4 8 6 6
16 81 87 51 6 29 31 82 37 43 22 31 46 29 21 54 60 3 5 11 1 6 1 13
6 20 3 4 42 66 1 1 1 1
19 79 19 100 20 87 49 4 54 56 73 88 43 92 32 18 47 66 28 69 4 3 7 14 18 12 13 14 17
8 61 49 29 51 31 83 62 51 5 2 8 5 7 6 8 1 8 4 6`

type query struct {
	l int
	r int
}

type testCase struct {
	h []int64
	q []query
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: too few fields", idx+1)
		}
		pos := 0
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		pos++
		if len(fields) < pos+n+1 {
			return nil, fmt.Errorf("line %d: not enough data for heights", idx+1)
		}
		h := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[pos+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse height %d: %w", idx+1, i+1, err)
			}
			h[i] = int64(v)
		}
		pos += n
		if pos >= len(fields) {
			return nil, fmt.Errorf("line %d: missing q", idx+1)
		}
		qn, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse q: %w", idx+1, err)
		}
		pos++
		if len(fields) != pos+2*qn {
			return nil, fmt.Errorf("line %d: field count mismatch", idx+1)
		}
		qs := make([]query, qn)
		for i := 0; i < qn; i++ {
			l, err := strconv.Atoi(fields[pos+2*i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse l%d: %w", idx+1, i+1, err)
			}
			r, err := strconv.Atoi(fields[pos+2*i+1])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse r%d: %w", idx+1, i+1, err)
			}
			qs[i] = query{l: l, r: r}
		}
		cases = append(cases, testCase{h: h, q: qs})
	}
	return cases, nil
}

// solveCase mirrors 232D.go and returns answers for all queries.
func solveCase(tc testCase) []int {
	n := len(tc.h)
	Dlen := n - 1
	if Dlen <= 0 {
		res := make([]int, len(tc.q))
		for i, qu := range tc.q {
			if qu.r-qu.l+1 == 1 {
				res[i] = n - 1
			} else {
				res[i] = 0
			}
		}
		return res
	}
	D := make([]int64, Dlen)
	F := make([]int64, Dlen)
	for i := 0; i < Dlen; i++ {
		D[i] = tc.h[i+1] - tc.h[i]
		F[i] = -D[i]
	}

	const mod1 = 1000000007
	const mod2 = 1000000009
	const base = 91138233

	pow1 := make([]int64, Dlen+5)
	pow2 := make([]int64, Dlen+5)
	pow1[0], pow2[0] = 1, 1
	for i := 1; i <= Dlen; i++ {
		pow1[i] = pow1[i-1] * base % mod1
		pow2[i] = pow2[i-1] * base % mod2
	}

	H1 := make([]int64, Dlen+1)
	H2 := make([]int64, Dlen+1)
	HF1 := make([]int64, Dlen+1)
	HF2 := make([]int64, Dlen+1)
	for i := 0; i < Dlen; i++ {
		v1 := (D[i]%mod1 + mod1) % mod1
		H1[i+1] = (H1[i]*base + v1) % mod1
		v2 := (D[i]%mod2 + mod2) % mod2
		H2[i+1] = (H2[i]*base + v2) % mod2
		f1 := (F[i]%mod1 + mod1) % mod1
		HF1[i+1] = (HF1[i]*base + f1) % mod1
		f2 := (F[i]%mod2 + mod2) % mod2
		HF2[i+1] = (HF2[i]*base + f2) % mod2
	}

	sa := buildSA(D)
	mst := newMergeTree(sa)

	getD := func(u, ln int) (int64, int64) {
		a := (H1[u+ln] - H1[u]*pow1[ln]%mod1 + mod1) % mod1
		b := (H2[u+ln] - H2[u]*pow2[ln]%mod2 + mod2) % mod2
		return a, b
	}
	getF := func(l, ln int) (int64, int64) {
		a := (HF1[l+ln] - HF1[l]*pow1[ln]%mod1 + mod1) % mod1
		b := (HF2[l+ln] - HF2[l]*pow2[ln]%mod2 + mod2) % mod2
		return a, b
	}

	isSufLess := func(u, l, m int) bool {
		maxLen := Dlen - u
		lim := m
		if maxLen < lim {
			lim = maxLen
		}
		lo, hi := 0, lim
		for lo < hi {
			mid := (lo + hi + 1) >> 1
			d1, d2 := getD(u, mid)
			f1, f2 := getF(l, mid)
			if d1 == f1 && d2 == f2 {
				lo = mid
			} else {
				hi = mid - 1
			}
		}
		lcp := lo
		if lcp == lim {
			return maxLen < m
		}
		dv := D[u+lcp]
		fv := F[l+lcp]
		return dv < fv
	}

	isSufLe := func(u, l, m int) bool {
		maxLen := Dlen - u
		lim := m
		if maxLen < lim {
			lim = maxLen
		}
		lo, hi := 0, lim
		for lo < hi {
			mid := (lo + hi + 1) >> 1
			d1, d2 := getD(u, mid)
			f1, f2 := getF(l, mid)
			if d1 == f1 && d2 == f2 {
				lo = mid
			} else {
				hi = mid - 1
			}
		}
		lcp := lo
		if lcp == lim {
			return maxLen <= m
		}
		dv := D[u+lcp]
		fv := F[l+lcp]
		return dv <= fv
	}

	results := make([]int, len(tc.q))
	for qi, qu := range tc.q {
		l := qu.l - 1
		r := qu.r - 1
		w := r - l + 1
		if w == 1 {
			results[qi] = n - 1
			continue
		}
		m := w - 1
		lo, hi := 0, len(sa)-1
		L := len(sa)
		for lo <= hi {
			mid := (lo + hi) >> 1
			if isSufLess(sa[mid], l, m) {
				lo = mid + 1
			} else {
				L = mid
				hi = mid - 1
			}
		}
		lo, hi = 0, len(sa)-1
		R := -1
		for lo <= hi {
			mid := (lo + hi) >> 1
			if isSufLe(sa[mid], l, m) {
				R = mid
				lo = mid + 1
			} else {
				hi = mid - 1
			}
		}
		if L > R {
			results[qi] = 0
			continue
		}
		max1 := l - m - 1
		cnt := 0
		if max1 >= 0 {
			cnt += mst.query(0, 0, Dlen-1, L, R, 0, max1)
		}
		min2 := r + 1
		if min2 <= Dlen-m {
			cnt += mst.query(0, 0, Dlen-1, L, R, min2, Dlen-m)
		}
		results[qi] = cnt
	}
	return results
}

// build suffix array for int64 slice D
func buildSA(D []int64) []int {
	n := len(D)
	if n == 0 {
		return []int{}
	}
	sa := make([]int, n)
	rnk := make([]int, n)
	tmp := make([]int, n)

	vals := append([]int64(nil), D...)
	m := len(vals)
	vs := make([]int64, m)
	copy(vs, vals)
	for i := 0; i < m; i++ {
		for j := i + 1; j < m; j++ {
			if vs[j] < vs[i] {
				vs[i], vs[j] = vs[j], vs[i]
			}
		}
	}
	mp := make(map[int64]int, m)
	uniq := make([]int64, 0, m)
	for i, v := range vs {
		if i == 0 || v != vs[i-1] {
			uniq = append(uniq, v)
		}
	}
	for i, v := range uniq {
		mp[v] = i
	}
	for i := 0; i < n; i++ {
		rnk[i] = mp[D[i]]
		sa[i] = i
	}

	k := 1
	tmpSA := make([]int, n)
	for k < n {
		maxv := n
		cnt := make([]int, maxv+1)
		for i := 0; i < n; i++ {
			key := 0
			if sa[i]+k < n {
				key = rnk[sa[i]+k] + 1
			}
			cnt[key]++
		}
		for i := 1; i <= maxv; i++ {
			cnt[i] += cnt[i-1]
		}
		for i := n - 1; i >= 0; i-- {
			key := 0
			if sa[i]+k < n {
				key = rnk[sa[i]+k] + 1
			}
			cnt[key]--
			tmpSA[cnt[key]] = sa[i]
		}
		cnt = make([]int, maxv+1)
		for i := 0; i < n; i++ {
			key := rnk[tmpSA[i]] + 1
			cnt[key]++
		}
		for i := 1; i <= maxv; i++ {
			cnt[i] += cnt[i-1]
		}
		for i := n - 1; i >= 0; i-- {
			key := rnk[tmpSA[i]] + 1
			cnt[key]--
			sa[cnt[key]] = tmpSA[i]
		}
		tmp[sa[0]] = 0
		p := 0
		for i := 1; i < n; i++ {
			a, b := sa[i-1], sa[i]
			if rnk[a] != rnk[b] ||
				(a+k < n && b+k < n && rnk[a+k] != rnk[b+k]) ||
				(a+k >= n && b+k < n) || (a+k < n && b+k >= n) {
				p++
			}
			tmp[b] = p
		}
		copy(rnk, tmp)
		if p == n-1 {
			break
		}
		k <<= 1
	}
	return sa
}

type mergeTree struct {
	sa []int
	t  [][]int
}

func newMergeTree(sa []int) *mergeTree {
	n := len(sa)
	if n == 0 {
		return &mergeTree{sa: sa, t: [][]int{}}
	}
	t := make([][]int, n*4)
	mt := &mergeTree{sa: sa, t: t}
	mt.build(0, 0, n-1)
	return mt
}

func (mt *mergeTree) build(node, l, r int) {
	if l == r {
		mt.t[node] = []int{mt.sa[l]}
		return
	}
	mid := (l + r) >> 1
	lc, rc := node*2+1, node*2+2
	mt.build(lc, l, mid)
	mt.build(rc, mid+1, r)
	a, b := mt.t[lc], mt.t[rc]
	merged := make([]int, len(a)+len(b))
	i, j := 0, 0
	for i < len(a) || j < len(b) {
		if j == len(b) || (i < len(a) && a[i] < b[j]) {
			merged[i+j] = a[i]
			i++
		} else {
			merged[i+j] = b[j]
			j++
		}
	}
	mt.t[node] = merged
}

func (mt *mergeTree) query(node, l, r, L, R, lo, hi int) int {
	if len(mt.t) == 0 || r < L || l > R {
		return 0
	}
	if L <= l && r <= R {
		arr := mt.t[node]
		cntL := lower(arr, lo)
		cntR := lower(arr, hi+1)
		return cntR - cntL
	}
	mid := (l + r) >> 1
	return mt.query(node*2+1, l, mid, L, R, lo, hi) + mt.query(node*2+2, mid+1, r, L, R, lo, hi)
}

func lower(arr []int, x int) int {
	lo, hi := 0, len(arr)
	for lo < hi {
		mid := (lo + hi) >> 1
		if arr[mid] < x {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.h)))
	for i, v := range tc.h {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.q)))
	for _, qu := range tc.q {
		fmt.Fprintf(&sb, "%d %d\n", qu.l, qu.r)
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range testcases {
		expectVals := solveCase(tc)
		expectParts := make([]string, len(expectVals))
		for i, v := range expectVals {
			expectParts[i] = strconv.Itoa(v)
		}
		expectStr := strings.Join(expectParts, "\n")
		input := buildInput(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expectStr {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expectStr, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
