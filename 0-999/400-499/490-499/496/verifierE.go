package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type part struct {
	a, b int64
	idx  int
}

type actor struct {
	c, d, k int64
	idx     int
}

type testCase struct {
	parts  []part
	actors []actor
}

// Embedded testcases (from testcasesE.txt) to keep the verifier self contained.
const rawTestcases = `1
3 3
1
5 8 3

3
2 4
8 12
6 8
3
4 4 3
9 10 2
9 14 1

4
6 6
9 10
2 6
4 4
3
3 8 1
10 13 3
8 11 1

4
5 6
1 3
5 10
6 11
3
4 7 1
1 6 1
3 7 2

2
1 1
9 10
2
2 3 1
9 10 2

4
1 3
5 7
10 10
4 8
1
6 8 2

1
10 15
3
6 6 2
4 6 3
5 6 1

2
10 11
9 9
3
10 14 3
10 10 3
3 5 2

4
5 8
3 6
7 9
2 4
2
5 6 2
8 8 2

3
5 10
6 7
7 8
2
3 7 2
10 13 3

2
1 3
3 3
4
7 11 1
1 4 2
10 15 1
7 7 3

2
7 11
8 11
4
7 10 1
5 6 1
6 6 3
10 13 1

4
4 9
3 5
1 4
6 9
3
9 10 2
10 10 2
3 5 1

3
1 3
10 13
4 4
2
10 13 1
5 10 3

2
3 6
10 15
4
6 11 2
10 15 1
1 6 3
7 12 2

3
5 6
5 8
5 8
2
8 9 3
3 6 2

1
7 8
3
7 7 2
5 6 1
4 6 1

2
6 7
10 14
2
1 1 1
5 9 2

3
3 3
5 9
6 6
4
5 7 3
10 14 3
5 5 3
9 9 1

3
7 11
8 10
10 11
1
7 10 3

1
10 12
3
3 3 1
3 6 2
1 3 2

4
9 14
2 2
10 14
4 5
3
9 12 1
1 1 1
6 7 2

1
7 8
3
2 2 3
3 5 1
3 4 2

4
4 5
1 5
5 8
4 8
1
10 12 3

3
5 9
4 7
5 5
4
6 8 1
8 10 1
8 10 1
2 4 1

1
5 5
3
1 3 1
3 5 2
3 4 2

2
1 1
8 12
2
7 8 1
9 13 1

3
5 6
8 12
2 2
3
6 9 2
1 1 2
5 8 3

2
7 8
7 8
1
6 7 2

2
7 8
4 4
2
4 6 2
7 10 3

2
9 9
6 10
3
9 11 1
1 1 2
7 10 2

4
1 4
5 7
4 4
4 4
3
4 6 3
4 6 3
2 6 1

2
1 2
2 3
4
2 2 1
2 2 2
1 2 3
4 6 1

4
1 1
5 9
7 7
5 10
4
8 10 1
7 8 3
10 13 2
2 3 1

2
4 6
3 5
4
3 4 1
2 2 2
3 3 3
2 3 2

1
7 8
3
7 10 1
2 5 1
5 7 1

4
1 2
6 7
3 3
9 11
1
4 6 3

2
9 12
6 8
1
7 10 3

2
3 4
3 6
3
2 6 1
4 4 1
3 6 1

2
2 3
6 7
2
6 7 2
2 4 3

1
2 2
4
3 3 1
5 5 3
2 2 2
2 3 2

1
3 3
4
1 2 3
3 5 1
7 9 3
1 3 3

3
8 9
1 2
1 2
1
6 6 1

2
2 4
2 3
3
3 5 1
2 6 2
1 2 1

3
2 2
4 5
3 3
4
3 5 1
6 8 3
6 8 3
5 7 3

2
1 3
1 2
4
5 9 2
4 7 3
2 6 3
4 5 1

2
3 5
4 4
4
3 4 2
6 8 1
1 1 1
6 6 1

4
2 6
3 5
2 6
5 10
2
9 9 3
8 10 3

4
1 2
1 2
2 3
7 11
1
2 6 2

3
6 7
4 6
3 5
3
3 6 1
6 7 3
5 9 2

3
1 4
7 7
6 10
3
3 3 3
9 10 2
9 13 2

3
3 5
4 6
2 4
3
8 11 1
5 8 2
1 1 2

4
1 1
2 3
2 2
2 2
3
2 2 2
1 1 3
6 7 1

4
6 7
1 2
1 3
1 2
3
7 11 1
1 2 3
2 3 1

3
1 2
3 5
1 3
2
3 6 1
5 6 2

2
3 5
1 1
3
4 6 1
5 6 2
7 9 2

3
5 7
3 5
2 4
2
3 5 1
8 11 2

1
1 2
4
3 4 2
9 10 2
2 4 2
3 3 3

2
2 3
7 10
2
4 5 1
6 8 2

2
2 3
7 8
2
1 2 2
2 4 1

2
6 8
5 6
2
8 11 1
5 6 2

1
3 4
4
3 6 3
8 8 1
8 10 3
4 7 1

1
2 6
1
1 3 2

2
3 3
1 1
2
4 5 1
3 4 2

2
3 5
1 3
1
4 6 2

3
2 4
4 6
2 2
1
7 10 2

2
2 2
1 1
1
1 3 3

3
1 2
1 2
1 2
2
2 2 1
8 11 2

4
2 2
2 2
3 5
2 3
4
9 11 3
7 7 1
7 10 2
6 9 1

1
3 3
3
4 5 2
9 11 2
8 8 1

4
1 3
1 3
2 2
1 3
2
3 5 2
3 5 2

2
6 8
3 5
2
5 9 2
6 7 3

2
1 1
1 1
3
3 3 2
7 7 2
2 4 3

1
2 2
3
5 5 3
3 3 3
1 1 3

2
6 8
3 5
3
4 6 2
5 6 1
3 5 2

1
5 7
2
3 6 3
1 2 3

1
1 2
2
3 3 1
3 4 1

1
1 1
2
1 2 1
1 2 1

2
5 5
6 6
2
7 9 3
5 6 1

2
1 1
2 3
2
3 5 2
2 2 2

1
2 4
1
3 3 2

1
1 1
3
1 2 3
3 5 3
7 10 2

1
2 3
1
6 7 3

1
1 2
2
4 7 2
8 11 2

1
1 2
1
5 6 2

1
1 2
2
2 2 1
5 7 1

1
1 1
1
3 5 1

1
1 3
2
5 7 1
5 8 2

2
7 8
3 4
3
6 7 1
7 10 3
8 11 2

2
5 6
6 7
1
5 8 1

2
1 3
4 4
1
3 5 3

4
1 1
1 1
1 1
1 1
3
2 2 2
1 2 3
8 8 1

1
3 5
1
2 2 1

1
2 2
1
6 7 2

1
2 3
1
4 7 3

1
2 4
1
3 5 2

1
3 4
1
4 6 1

1
5 6
1
8 11 2

1
1 2
1
5 6 2

1
2 2
1
6 7 1

1
2 2
1
8 10 3

1
1 1
1
4 5 3

2
1 2
1 1
1
6 7 2

1
1 2
1
7 10 2

2
1 1
1 1
1
2 3 3

1
2 3
1
4 4 3

2
3 3
1 2
1
1 3 2

2
3 3
3 5
1
3 4 1

2
2 4
1 3
1
2 3 2

1
1 2
1
3 3 1

1
1 3
1
2 3 2

1
2 3
1
1 2 1

1
2 3
1
4 6 2

1
1 3
1
5 7 2

1
1 2
1
2 2 2

1
1 2
1
3 5 2

1
1 2
1
6 9 1

1
1 2
1
6 7 1

1
2 2
1
6 8 1

1
1 2
1
8 10 3

1
1 1
1
1 1 1
`

// Segment tree for range sums and find-first queries.
type segTree struct {
	n    int
	data []int64
}

func newSegTree(size int) *segTree {
	n := 1
	for n < size {
		n <<= 1
	}
	return &segTree{n: n, data: make([]int64, 2*n)}
}

func (st *segTree) add(pos int, delta int64) {
	i := pos + st.n
	st.data[i] += delta
	for i >>= 1; i > 0; i >>= 1 {
		st.data[i] = st.data[i<<1] + st.data[i<<1|1]
	}
}

// find first position >= l with positive value, or -1 if none.
func (st *segTree) findFirst(l int) int {
	return st.findFirstRec(1, 0, st.n-1, l)
}

func (st *segTree) findFirstRec(node, nl, nr, ql int) int {
	if nr < ql || st.data[node] == 0 {
		return -1
	}
	if nl == nr {
		return nl
	}
	mid := (nl + nr) >> 1
	if ql <= mid {
		if res := st.findFirstRec(node<<1, nl, mid, ql); res != -1 {
			return res
		}
	}
	return st.findFirstRec(node<<1|1, mid+1, nr, ql)
}

func parseTestcases() ([]testCase, error) {
	scanner := bufio.NewScanner(strings.NewReader(rawTestcases))
	var cases []testCase
	for {
		var line string
		for scanner.Scan() {
			line = strings.TrimSpace(scanner.Text())
			if line != "" {
				break
			}
		}
		if line == "" {
			break
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("parse n: %w", err)
		}
		parts := make([]part, n)
		for i := 0; i < n; i++ {
			if !scanner.Scan() {
				return nil, fmt.Errorf("unexpected EOF reading parts")
			}
			fields := strings.Fields(scanner.Text())
			if len(fields) != 2 {
				return nil, fmt.Errorf("part %d malformed", i+1)
			}
			a, err1 := strconv.ParseInt(fields[0], 10, 64)
			b, err2 := strconv.ParseInt(fields[1], 10, 64)
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("part %d parse error", i+1)
			}
			parts[i] = part{a: a, b: b, idx: i}
		}
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected EOF reading m")
		}
		m, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil {
			return nil, fmt.Errorf("parse m: %w", err)
		}
		actors := make([]actor, m)
		for i := 0; i < m; i++ {
			if !scanner.Scan() {
				return nil, fmt.Errorf("unexpected EOF reading actor %d", i+1)
			}
			fields := strings.Fields(scanner.Text())
			if len(fields) != 3 {
				return nil, fmt.Errorf("actor %d malformed", i+1)
			}
			c, err1 := strconv.ParseInt(fields[0], 10, 64)
			d, err2 := strconv.ParseInt(fields[1], 10, 64)
			k, err3 := strconv.ParseInt(fields[2], 10, 64)
			if err1 != nil || err2 != nil || err3 != nil {
				return nil, fmt.Errorf("actor %d parse error", i+1)
			}
			actors[i] = actor{c: c, d: d, k: k, idx: i + 1}
		}
		cases = append(cases, testCase{parts: parts, actors: actors})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func uniqueInt64(a []int64) []int64 {
	if len(a) == 0 {
		return a
	}
	j := 1
	for i := 1; i < len(a); i++ {
		if a[i] != a[j-1] {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}

func solve(tc testCase) string {
	parts := make([]part, len(tc.parts))
	copy(parts, tc.parts)
	actors := make([]actor, len(tc.actors))
	copy(actors, tc.actors)

	diVals := make([]int64, len(actors))
	for i := 0; i < len(actors); i++ {
		diVals[i] = actors[i].d
	}
	sort.Slice(diVals, func(i, j int) bool { return diVals[i] < diVals[j] })
	diVals = uniqueInt64(diVals)

	sort.Slice(parts, func(i, j int) bool { return parts[i].a < parts[j].a })
	sort.Slice(actors, func(i, j int) bool { return actors[i].c < actors[j].c })

	st := newSegTree(len(diVals))
	type actorEntry struct {
		idx int
		k   int64
	}
	actorLists := make([][]*actorEntry, len(diVals))
	ans := make([]int, len(parts))

	ai := 0
	for _, p := range parts {
		for ai < len(actors) && actors[ai].c <= p.a {
			pos := sort.Search(len(diVals), func(i int) bool { return diVals[i] >= actors[ai].d })
			st.add(pos, actors[ai].k)
			entry := &actorEntry{idx: actors[ai].idx, k: actors[ai].k}
			actorLists[pos] = append(actorLists[pos], entry)
			ai++
		}
		j := sort.Search(len(diVals), func(i int) bool { return diVals[i] >= p.b })
		pos := st.findFirst(j)
		if pos == -1 {
			return "NO"
		}
		e := actorLists[pos][0]
		ans[p.idx] = e.idx
		st.add(pos, -1)
		e.k--
		if e.k == 0 {
			actorLists[pos] = actorLists[pos][1:]
		}
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return strings.TrimSpace(sb.String())
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tc.parts)))
	sb.WriteByte('\n')
	for _, p := range tc.parts {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.a, p.b))
	}
	sb.WriteString(strconv.Itoa(len(tc.actors)))
	sb.WriteByte('\n')
	for _, a := range tc.actors {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", a.c, a.d, a.k))
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
