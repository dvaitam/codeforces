package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Fenwick tree.
type bit struct {
	n    int
	data []int
}

func newBIT(n int) *bit {
	return &bit{n: n, data: make([]int, n+1)}
}

func (b *bit) add(i, v int) {
	for i <= b.n {
		b.data[i] += v
		i += i & -i
	}
}

func (b *bit) sum(i int) int {
	res := 0
	for i > 0 {
		res += b.data[i]
		i -= i & -i
	}
	return res
}

func (b *bit) findByOrder(k int) int {
	idx := 0
	bitMask := 1
	for bitMask<<1 <= b.n {
		bitMask <<= 1
	}
	for d := bitMask; d > 0; d >>= 1 {
		next := idx + d
		if next <= b.n && b.data[next] < k {
			idx = next
			k -= b.data[next]
		}
	}
	return idx + 1
}

type kv struct {
	val int
	pos int
}

// Embedded solver from 387E.go.
func solve(reader io.Reader) (string, error) {
	in := bufio.NewReader(reader)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return "", err
	}
	p := make([]int, n+1)
	pos := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if _, err := fmt.Fscan(in, &p[i]); err != nil {
			return "", err
		}
		pos[p[i]] = i
	}
	bvals := make([]int, k)
	isKeeper := make([]bool, n+1)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(in, &bvals[i]); err != nil {
			return "", err
		}
		isKeeper[bvals[i]] = true
	}
	bElems := make([]kv, 0, k)
	for _, v := range bvals {
		bElems = append(bElems, kv{val: v, pos: pos[v]})
	}
	sort.Slice(bElems, func(i, j int) bool { return bElems[i].val < bElems[j].val })

	remElems := make([]kv, 0, n-k)
	for v := 1; v <= n; v++ {
		if !isKeeper[v] {
			remElems = append(remElems, kv{val: v, pos: pos[v]})
		}
	}

	remBIT := newBIT(n)
	for i := 1; i <= n; i++ {
		remBIT.add(i, 1)
	}
	blockerBIT := newBIT(n)

	var ans int64
	idxB := 0
	for _, item := range remElems {
		x := item.val
		px := item.pos
		for idxB < len(bElems) && bElems[idxB].val < x {
			blockerBIT.add(bElems[idxB].pos, 1)
			idxB++
		}
		sumLeft := blockerBIT.sum(px - 1)
		leftBound := 0
		if sumLeft > 0 {
			leftBound = blockerBIT.findByOrder(sumLeft)
		}
		sumBefore := blockerBIT.sum(px)
		totalBl := blockerBIT.sum(n)
		rightBound := n + 1
		if sumBefore < totalBl {
			rightBound = blockerBIT.findByOrder(sumBefore + 1)
		}
		l := leftBound + 1
		r := rightBound - 1
		if l <= r {
			count := remBIT.sum(r) - remBIT.sum(leftBound)
			ans += int64(count)
		}
		remBIT.add(px, -1)
	}
	return fmt.Sprintf("%d\n", ans), nil
}

type testCase struct {
	n int
	k int
	a []int
	b []int
}

const testcasesRaw = `100
4 0 4 3 2 1
2 2 2 1 2 1
10 2 6 1 9 4 8 7 3 2 5 10 7 1
2 2 2 1 1 2
10 5 6 10 5 9 2 7 4 8 1 3 8 4 3 9 10
7 3 5 6 1 3 2 4 7 3 1 4
5 3 5 4 1 2 3 4 1 5
8 2 1 8 3 4 6 7 5 2 3 1
6 1 5 2 4 1 6 3 4
6 2 5 4 6 2 1 3 3 4
3 1 3 1 2 3
9 4 5 3 9 7 8 6 4 1 2 9 2 5 4
3 0 3 1 2
5 0 3 2 4 1 5
7 2 2 1 6 7 3 4 5 4 7
5 1 3 5 4 1 2 1
2 1 1 2 2
8 0 5 7 1 6 4 2 8 3
3 3 3 1 2 1 3 2
6 1 1 5 4 3 6 2 4
10 10 6 10 4 2 1 3 7 5 9 8 9 8 10 4 2 5 6 3 1 7
7 5 4 1 7 6 2 3 5 2 3 4 6 5
1 0 1
9 0 4 5 7 8 6 9 1 2 3
4 1 1 4 2 3 3
4 3 2 4 3 1 2 4 1
7 2 5 3 4 6 1 2 7 4 6
6 4 1 5 2 3 6 4 3 6 2 5
6 3 5 4 6 2 1 3 6 2 1
3 3 1 3 2 1 2
10 0 6 5 10 1 9 7 3 2 5 6 4 1
9 2 5 6 1 3 2 4 7 3 1 4
5 1 3 5 4 1 2 1
2 1 1 2 2
8 0 5 7 1 6 4 2 8 3
3 3 3 1 2 1 3 2
6 1 1 5 4 3 6 2 4
10 10 6 10 4 2 1 3 7 5 9 8 9 8 10 4 2 5 6 3 1 7
7 5 4 1 7 6 2 3 5 2 3 4 6 5
1 0 1
9 0 4 5 7 8 6 9 1 2 3
4 1 1 4 2 3 3
4 3 2 4 3 1 2 4 1
7 2 5 3 4 6 1 2 7 4 6
6 4 1 5 2 3 6 4 3 6 2 5
6 3 5 4 6 2 1 3 6 2 1
3 3 1 3 2 1 2
10 5 2 3 5 6 7 5 2 6 4 1 2 9 2 5 8
7 5 3 4 7 8 2 5 7 3 4 5 6 2
4 3 1 6 3 2
5 0 1 3 7 2 6 5 4
1 0 1
9 0 4 5 7 8 6 9 1 2 3
4 1 1 4 2 3 3
4 3 2 4 3 1 2 4 1
7 2 5 3 4 6 1 2 7 4 6
6 4 1 5 2 3 6 4 3 6 2 5
6 3 5 4 6 2 1 3 6 2 1
3 3 1 3 2 1 2
9 8 5 2 3 5 6 7 5 2 6 4 1 2 9 2 5 8
7 5 3 4 7 8 2 5 7 3 4 5 6 2
4 3 1 6 3 2
5 0 1 3 7 2 6 5 4
1 0 1
9 0 4 5 7 8 6 9 1 2 3
4 1 1 4 2 3 3
4 3 2 4 3 1 2 4 1
7 2 5 3 4 6 1 2 7 4 6
6 4 1 5 2 3 6 4 3 6 2 5
6 3 5 4 6 2 1 3 6 2 1
3 3 1 3 2 1 2
2 2 1 2 2 1
9 4 8 9 3 7 1 5 2 6 4 5 1 8 2
8 4 6 3 5 1 4 8 2 7 5 4 8 6
2 1 2 1 2
4 3 1 4 2 3 2 3 4
10 8 5 8 9 10 7 2 6 3 4 1 5 4 8 7 2 6 10 3
8 6 2 8 4 3 7 5 6 1 4 7 1 6 5 3
7 0 1 3 7 2 6 5 4
9 4 1 4 9 7 3 8 5 6 2 4 3 9 8
2 2 1 2 2 1`

func parseTestcases(raw string) ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(raw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("missing test count")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, fmt.Errorf("parse t: %w", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("case %d: missing n", i+1)
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			return nil, fmt.Errorf("case %d: missing k", i+1)
		}
		k, _ := strconv.Atoi(scan.Text())
		a := make([]int, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("case %d: missing a[%d]", i+1, j)
			}
			a[j], _ = strconv.Atoi(scan.Text())
		}
		b := make([]int, k)
		for j := 0; j < k; j++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("case %d: missing b[%d]", i+1, j)
			}
			b[j], _ = strconv.Atoi(scan.Text())
		}
		cases = append(cases, testCase{n: n, k: k, a: a, b: b})
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expected, err := solve(strings.NewReader(fmt.Sprintf("%d %d\n%s\n%s\n", tc.n, tc.k, joinInts(tc.a), joinInts(tc.b))))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}

		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
		input.WriteString(joinInts(tc.a))
		input.WriteByte('\n')
		if len(tc.b) > 0 {
			input.WriteString(joinInts(tc.b))
		}
		input.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n%s", idx+1, err, string(out))
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\n", idx+1, strings.TrimSpace(expected), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

func joinInts(vals []int) string {
	if len(vals) == 0 {
		return ""
	}
	var sb strings.Builder
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}
