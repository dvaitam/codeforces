package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type node struct {
	val int
	cnt int
}

type testCase struct {
	n   int
	q   int
	arr []int
	l   []int
	r   []int
}

// Embedded testcases from testcasesD.txt.
const testcaseData = `
4 10 2 3 4 1 1 4 3 3 2 4 4 4 4 4 2 4 2 4 4 4 1 2 1 3
1 5 1 1 1 1 1 1 1 1 1 1 1
4 5 4 3 4 4 3 4 2 3 1 3 2 4 3 3
4 10 3 3 1 1 4 4 1 3 1 4 2 2 3 4 4 4 1 1 4 4 3 3 1 3
1 2 1 1 1 1 1
10 5 3 1 6 6 6 3 7 7 8 9 7 7 10 10 7 8 5 8 5 9
5 9 3 1 4 5 3 1 4 5 5 1 3 4 5 3 5 3 5 4 4 5 5 1 3
5 8 3 5 5 3 2 3 3 3 4 5 5 3 4 1 1 5 5 3 5 2 4
4 6 2 4 1 1 3 4 2 3 2 2 3 3 4 4 2 2
1 9 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
3 4 1 2 3 3 3 3 3 1 3 2 3
9 6 4 2 5 2 4 1 1 9 4 7 9 1 1 8 8 3 7 5 6 1 9
9 7 1 2 6 3 5 9 8 1 6 4 5 2 3 3 4 5 6 1 8 7 7 5 6
5 10 5 5 4 1 4 3 3 1 2 1 1 1 1 4 4 1 5 5 5 3 3 3 3 3 4
7 10 3 3 3 2 3 4 1 2 6 1 6 6 7 7 7 5 5 1 3 4 7 6 6 5 6 1 3
8 6 7 7 8 1 4 4 5 2 7 7 7 7 1 6 6 8 5 5 8 8
9 7 2 6 9 2 1 8 3 4 7 1 9 2 3 7 7 1 6 2 2 2 9 5 9
5 2 1 5 5 5 2 1 5 1 5
1 9 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
10 10 9 8 3 7 3 3 2 8 8 9 8 10 3 5 5 6 3 8 4 10 9 10 7 9 4 6 1 5 8 9
4 3 3 2 3 4 2 3 4 4 4 4
8 2 7 1 8 4 4 2 4 5 4 5 5 6
3 10 3 3 1 2 2 1 2 1 2 1 3 1 1 1 2 2 2 2 3 3 3 1 1
6 6 4 4 4 1 2 6 5 6 4 4 5 6 1 3 1 6 4 4
8 9 5 2 6 6 8 5 5 2 6 8 2 7 8 8 1 5 3 8 3 4 6 8 8 8 2 6
3 6 3 3 3 3 3 3 3 3 3 2 3 2 2 3 3
2 3 2 2 1 2 2 2 1 1
1 8 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
10 7 7 3 10 10 3 7 4 9 9 3 10 10 4 10 5 7 5 5 8 9 7 9 9 10
8 9 5 8 1 4 1 2 4 8 3 7 8 8 4 8 4 4 8 8 5 6 3 6 2 6 1 1
6 10 2 5 1 4 5 1 3 5 3 5 6 6 1 5 1 6 1 6 3 4 1 2 4 6 2 5
6 2 1 4 1 2 6 2 4 5 4 6
2 9 2 1 2 2 1 2 2 2 2 2 2 2 2 2 2 2 2 2 1 1
4 5 3 2 4 2 2 2 1 2 4 4 1 2 1 2
8 8 3 1 1 7 8 6 7 1 1 4 7 7 7 8 1 4 4 4 7 8 4 5 6 8
2 6 1 1 2 2 2 2 2 2 2 2 2 2 2 2
1 7 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
8 6 1 7 5 7 2 5 3 7 2 6 7 8 7 7 7 8 3 5 7 8
4 8 4 3 3 2 4 4 1 2 2 4 1 4 1 1 3 3 1 1 4 4
6 6 1 1 2 4 1 3 5 6 1 4 1 6 2 3 6 6 2 3
10 1 2 4 5 4 4 9 9 7 9 10 6 10
4 8 2 1 1 1 1 1 2 3 1 1 4 4 2 4 1 4 3 4 4 4
5 10 4 4 2 1 5 1 4 3 5 1 3 3 3 4 4 5 5 4 4 5 5 2 4 1 3
4 7 4 1 3 2 3 3 2 4 2 4 4 4 3 3 4 4 4 4
10 1 5 1 3 2 1 3 5 9 9 1 8 8
4 4 3 4 4 1 3 4 2 4 3 3 1 4
5 7 4 1 2 2 4 3 4 3 3 4 5 4 5 5 5 3 4 3 3
8 6 6 5 4 1 1 3 3 1 3 3 1 4 8 8 6 8 1 8 3 4
1 5 1 1 1 1 1 1 1 1 1 1 1
6 1 1 4 1 6 6 2 5 5
3 7 3 2 3 1 2 1 1 1 1 2 2 2 3 3 3 2 2
5 3 1 2 5 1 4 3 4 4 5 3 4
5 9 2 5 3 2 5 1 4 4 4 4 5 1 3 5 5 4 4 3 5 3 3 4 4
6 1 5 4 6 4 2 3 3 4
7 7 3 1 3 2 1 3 3 5 6 5 6 1 2 4 6 4 4 4 5 2 7
2 9 1 1 1 2 1 2 2 2 2 2 2 2 1 1 2 2 1 1 1 2
7 1 2 4 1 2 7 6 2 3 4
3 3 1 3 1 3 3 2 2 2 3
3 1 2 1 1 1 1
5 3 4 3 5 1 2 1 3 2 5 2 2
6 9 4 1 6 3 4 1 6 6 5 6 5 6 4 6 5 6 4 4 5 6 4 4 5 6
3 5 1 2 2 1 1 1 3 1 2 3 3 1 3
1 3 1 1 1 1 1 1 1
6 3 5 5 2 3 4 5 5 6 6 6 5 6
10 6 6 5 5 5 3 2 10 9 4 6 4 6 7 9 8 8 8 9 9 9 10 10
10 1 2 7 7 9 10 5 1 4 7 7 4 4
6 9 2 1 5 4 1 4 6 6 6 6 1 3 3 5 4 5 3 3 5 6 2 4 2 2
6 2 5 1 2 1 2 3 4 5 2 2
1 9 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
6 7 1 4 2 4 3 2 4 4 5 6 5 6 4 5 1 6 6 6 1 4
6 2 1 3 4 4 6 1 6 6 6 6
10 4 6 4 7 2 1 9 1 9 9 4 10 10 3 6 8 10 7 9
6 8 6 5 2 2 5 3 3 4 5 6 5 6 3 4 4 4 6 6 3 4 4 4
3 2 2 1 3 2 2 2 2
1 6 1 1 1 1 1 1 1 1 1 1 1 1 1
5 1 3 1 2 2 3 4 4
10 3 4 8 4 4 6 10 4 3 9 7 10 10 10 10 5 8
6 1 1 5 4 4 1 5 2 4
9 7 8 4 7 9 8 6 7 8 3 2 8 6 7 6 6 6 7 5 8 9 9 8 9
8 1 6 4 2 5 5 6 7 5 1 7
9 6 7 8 6 3 4 5 4 7 2 6 7 6 9 3 7 7 8 4 7 9 9
2 1 2 1 1 1
2 2 1 2 2 2 2 2
8 5 6 4 5 2 4 2 8 6 7 7 5 7 3 5 4 5 4 6
6 6 1 1 6 2 1 2 4 5 1 3 2 3 1 6 4 5 5 5
6 8 6 4 5 5 6 5 4 5 3 5 5 6 5 6 2 2 1 2 4 6 1 2
3 9 3 1 3 3 3 2 2 3 3 1 1 1 1 3 3 2 3 3 3 1 1
9 6 3 5 3 2 2 3 2 1 2 9 9 8 9 5 6 1 2 5 8 3 3
1 6 1 1 1 1 1 1 1 1 1 1 1 1 1
1 6 1 1 1 1 1 1 1 1 1 1 1 1 1
7 1 5 3 2 2 4 5 5 3 7
7 4 5 5 4 6 3 7 1 2 5 5 7 1 3 5 7
3 8 3 1 2 3 3 1 2 3 3 1 1 2 2 2 2 3 3 1 1
1 8 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
2 4 1 1 1 2 1 1 1 2 1 2
2 6 2 1 2 2 2 2 2 2 1 2 2 2 2 2
5 4 3 1 3 4 4 4 4 3 5 1 2 5 5
10 7 4 7 3 6 2 3 9 8 5 7 7 9 7 10 1 6 7 9 7 9 9 9 3 6
4 10 1 4 1 4 4 4 4 4 2 3 4 4 2 2 2 4 1 2 4 4 1 1 2 3
9 3 7 7 3 3 4 7 8 6 8 9 9 5 9 1 6
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	res := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		pos := 0
		nextInt := func() (int, error) {
			if pos >= len(fields) {
				return 0, fmt.Errorf("case %d unexpected end of data", i+1)
			}
			v, err := strconv.Atoi(fields[pos])
			pos++
			return v, err
		}
		n, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		q, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d bad q: %v", i+1, err)
		}
		arr := make([]int, n+1)
		for j := 1; j <= n; j++ {
			v, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d bad a[%d]: %v", i+1, j, err)
			}
			arr[j] = v
		}
		l := make([]int, q)
		r := make([]int, q)
		for j := 0; j < q; j++ {
			x, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d bad l[%d]: %v", i+1, j, err)
			}
			y, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d bad r[%d]: %v", i+1, j, err)
			}
			l[j], r[j] = x, y
		}
		if pos != len(fields) {
			return nil, fmt.Errorf("case %d extra tokens", i+1)
		}
		res = append(res, testCase{n: n, q: q, arr: arr, l: l, r: r})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	return res, nil
}

// solve mirrors 1514D.go.
func solve(tc testCase) string {
	n := tc.n
	q := tc.q
	a := tc.arr
	pos := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		pos[a[i]] = append(pos[a[i]], i)
	}
	tree := make([]node, 4*(n+1))
	var build func(int, int, int)
	build = func(idx, l, r int) {
		if l == r {
			tree[idx] = node{a[l], 1}
			return
		}
		mid := (l + r) / 2
		build(idx*2, l, mid)
		build(idx*2+1, mid+1, r)
		tree[idx] = merge(tree[idx*2], tree[idx*2+1])
	}
	build(1, 1, n)

	var query func(int, int, int, int, int) node
	query = func(idx, l, r, L, R int) node {
		if L <= l && r <= R {
			return tree[idx]
		}
		mid := (l + r) / 2
		if R <= mid {
			return query(idx*2, l, mid, L, R)
		}
		if L > mid {
			return query(idx*2+1, mid+1, r, L, R)
		}
		left := query(idx*2, l, mid, L, R)
		right := query(idx*2+1, mid+1, r, L, R)
		return merge(left, right)
	}

	outputs := make([]string, 0, q)
	for i := 0; i < q; i++ {
		l, r := tc.l[i], tc.r[i]
		cand := query(1, 1, n, l, r).val
		freq := 0
		if cand != 0 && cand < len(pos) {
			arr := pos[cand]
			left := sort.SearchInts(arr, l)
			right := sort.SearchInts(arr, r+1)
			freq = right - left
		}
		length := r - l + 1
		ans := 1
		if tmp := 2*freq - length; tmp > 1 {
			ans = tmp
		}
		outputs = append(outputs, strconv.Itoa(ans))
	}
	return strings.Join(outputs, "\n")
}

func merge(left, right node) node {
	if left.val == right.val {
		return node{left.val, left.cnt + right.cnt}
	}
	if left.cnt > right.cnt {
		return node{left.val, left.cnt - right.cnt}
	}
	if right.cnt > left.cnt {
		return node{right.val, right.cnt - left.cnt}
	}
	return node{0, 0}
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
	for i := 1; i <= tc.n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(tc.arr[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < tc.q; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.l[i], tc.r[i]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
