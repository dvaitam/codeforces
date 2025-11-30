package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Embedded testcases (same format as original file).
const embeddedTestcases = `100
10
3 0
3 5
3 -5
3 -2
3 -5
1 -4 0
2 1
3 3
1 4 -2
1 -2 1

5
1 1 -3
1 -3 4
3 2
1 -3 -5
1 -2 -2

3
1 -1 0
1 3 5
3 -2

3
3 -2
3 -1
1 0 1

3
1 -1 -4
2 1
3 4

1
3 5

6
1 -1 0
2 1
3 0
1 2 2
3 -3
1 -1 -5

6
3 -5
3 1
3 1
3 -5
3 -5
3 -3

10
1 -4 -2
2 1
3 0
3 -1
3 -4
3 0
3 -5
3 -4
1 0 3
3 0

3
3 -1
3 3
1 -1 5

6
3 -3
1 5 -3
3 -1
2 1
3 -5
1 4 3

7
1 -2 4
2 1
3 5
3 -3
1 5 -5
2 1
1 -3 4

3
3 1
1 -3 1
2 1

1
3 -1

3
3 4
1 3 2
2 1

8
3 -1
3 1
1 -4 1
3 -3
3 2
2 1
1 2 -1
3 3

9
3 -4
3 4
3 -5
3 0
3 5
3 2
3 -1
3 5
1 4 -5

8
3 -1
3 5
3 2
3 3
3 5
3 0
3 5
3 1

6
1 2 0
2 1
3 -3
1 0 2
2 1
3 5

7
1 4 4
3 5
2 1
3 3
3 -1
3 -5
1 -3 4

8
3 5
1 -2 5
1 5 -5
2 1
1 -5 -3
1 0 -3
2 2
3 -5

7
3 0
3 5
3 -4
3 -2
1 0 -5
2 1
3 1

2
3 3
3 -5

9
3 -1
1 -1 3
3 0
3 -1
2 1
3 1
3 5
3 0
3 -3

3
3 1
3 2
1 -3 4

2
3 5
1 1 -4

6
3 4
3 -3
3 5
3 1
3 1
1 2 -1

8
3 1
3 -3
3 4
3 -1
3 -1
3 -5
3 -1
3 -1

3
3 -5
1 5 4
2 1

5
1 -3 1
1 2 3
3 -1
1 2 -5
1 2 -1

3
3 -1
3 2
3 2

9
3 4
3 -4
1 -3 -1
2 1
3 -1
3 3
1 2 0
2 1
1 -5 -1

9
3 5
1 5 3
1 -5 1
2 2
2 1
3 5
3 -2
3 2
1 -1 -5

7
3 -1
3 2
3 -3
1 2 3
2 1
3 -3
3 4

9
1 -4 -2
2 1
1 5 -5
2 1
1 1 2
1 -4 -4
1 4 5
1 -3 -1
3 2

3
1 4 -3
2 1
1 4 -2

1
3 -4

7
3 -4
3 -5
3 4
1 1 4
1 -5 1
1 0 5
3 2

8
3 5
3 -5
1 -1 -3
3 5
3 3
2 1
3 -4
3 -2

5
3 3
1 -2 0
2 1
1 -3 -5
3 0

2
3 5
1 -4 -4

9
3 2
1 1 2
2 1
1 3 -5
3 1
1 -3 1
3 -2
3 -4
1 5 2

4
1 4 1
2 1
3 0
3 0

7
3 -3
3 0
3 -1
3 1
3 3
1 4 4
1 -3 1

2
1 -5 -5
1 -2 -2

1
3 2

6
1 1 2
2 1
3 2
3 -4
1 -3 5
1 -4 0

7
3 -3
3 -1
1 -1 -3
2 1
1 2 -5
2 1
3 -5

1
3 2

8
1 -4 5
2 1
3 5
1 -3 -5
1 -5 4
1 5 -4
3 1
3 0

5
1 2 0
2 1
3 2
1 -2 -3
3 -2

1
3 4

7
3 -1
1 2 -5
2 1
1 -5 -5
3 3
3 0
1 -4 -2

8
1 2 -5
3 -2
3 -5
2 1
1 -5 -1
2 1
3 5
1 -5 -2

3
3 3
3 1
1 3 -2

2
3 -4
1 3 -3

10
1 -2 4
1 2 3
2 2
2 1
3 5
1 -5 -1
3 2
1 3 -4
1 -1 2
1 -1 1

6
1 -3 3
1 2 -2
3 -4
2 1
2 2
3 -2

1
3 -2

5
3 -4
1 1 0
2 1
3 5
1 4 -5

8
1 5 -5
3 -5
3 1
2 1
3 3
1 0 -4
2 1
1 -1 2

6
3 -5
1 3 3
1 -5 -4
1 -3 0
2 2
3 -5

4
3 1
1 0 -2
2 1
1 -3 4

10
1 -4 -1
3 4
1 5 -3
2 2
2 1
3 0
3 -2
3 4
3 -3
3 -5

3
3 2
3 4
1 -4 -3

7
3 -5
1 -5 -5
3 -2
2 1
3 2
1 -2 -1
1 -1 3

9
3 -2
3 0
1 5 -5
3 2
2 1
1 -3 -2
2 1
1 -4 -3
2 1

5
1 0 2
3 2
3 4
1 -2 1
3 -5

2
3 4
3 4

3
3 -5
3 5
1 4 -1

9
3 3
3 3
3 1
3 4
3 3
3 -5
3 -3
3 0
3 -5

9
1 -1 -1
1 -5 -4
2 2
2 1
3 -2
1 -5 -4
3 -1
2 1
3 3

2
1 -2 0
3 1

7
1 5 1
2 1
3 3
1 5 4
2 1
3 -4
3 3

6
1 -4 -2
2 1
3 4
1 -3 -4
3 1
1 -2 5

2
3 4
3 5

9
3 1
3 -2
3 2
3 4
1 -1 1
2 1
3 0
1 1 -5
1 -4 -1

10
3 4
3 -3
3 -2
3 4
1 5 0
3 -4
1 2 -3
3 -4
2 1
2 2

4
3 4
3 -4
1 -2 5
1 -4 -4

6
1 0 4
1 4 2
3 0
1 3 4
3 -4
2 1

2
1 -2 0
2 1

8
3 -4
1 1 3
2 1
1 -4 4
2 1
3 -5
3 1
3 -4

8
3 -1
1 -4 2
2 1
3 1
1 -1 2
3 -4
3 -4
2 1

6
1 3 1
1 -2 1
1 -1 5
3 -4
1 -1 -4
3 2

7
1 3 1
1 -4 2
2 1
1 -5 -3
2 2
1 -1 3
2 2

9
1 -5 -2
1 2 -5
2 1
2 2
3 -5
1 -3 1
3 1
3 5
3 4

10
1 -2 -5
2 1
3 1
1 -3 0
1 -1 -2
2 1
2 2
3 1
1 -1 -3
2 1

5
1 5 5
1 -3 2
1 -4 5
2 1
2 2

4
3 -5
3 2
3 4
3 -3

8
1 -4 0
3 -5
1 -3 1
1 -1 -3
3 0
2 1
1 -1 3
3 3

2
3 -4
1 5 1

6
3 0
3 3
3 2
1 4 -3
2 1
3 -5

7
1 -2 2
3 -5
1 -5 -2
2 1
3 4
1 -1 -5
3 -2

1
3 1

6
1 3 -1
2 1
1 5 -4
2 1
3 4
3 1

7
3 -3
3 5
1 5 2
1 -2 3
3 1
1 3 -1
3 0

7
1 2 0
1 2 -5
3 4
1 4 2
1 -5 -4
1 0 -3
2 3

5
3 2
1 -4 -1
2 1
3 -1
3 -3

4
3 -1
3 2
1 4 -3
3 -5

6
3 -4
1 2 -5
3 4
2 1
1 4 1
1 1 3`

type Query struct {
	t   int
	a   int64
	b   int64
	idx int
	q   int64
}

type SegmentLine struct {
	a, b int64
	l, r int
}

type Line struct {
	m int64
	c int64
}
type LCNode struct {
	ln          Line
	left, right *LCNode
}

var (
	seg     [][]int
	lines   []SegmentLine
	queries []Query
	answers []string
	xs      []int64
	m       int
	n       int
)

func eval(ln Line, x int64) int64 {
	return ln.m*x + ln.c
}

func insert(node *LCNode, l, r int, ln Line) *LCNode {
	if node == nil {
		return &LCNode{ln: ln}
	}
	newNode := &LCNode{ln: node.ln, left: node.left, right: node.right}
	node = newNode
	mid := (l + r) / 2
	midX := xs[mid]
	if eval(ln, midX) > eval(node.ln, midX) {
		node.ln, ln = ln, node.ln
	}
	if l == r {
		return node
	}
	if eval(ln, xs[l]) > eval(node.ln, xs[l]) {
		node.left = insert(node.left, l, mid, ln)
	} else if eval(ln, xs[r]) > eval(node.ln, xs[r]) {
		node.right = insert(node.right, mid+1, r, ln)
	}
	return node
}

func query(node *LCNode, l, r int, x int64) int64 {
	if node == nil {
		return math.MinInt64
	}
	res := eval(node.ln, x)
	if l == r {
		return res
	}
	mid := (l + r) / 2
	if x <= xs[mid] {
		if v := query(node.left, l, mid, x); v > res {
			res = v
		}
	} else {
		if v := query(node.right, mid+1, r, x); v > res {
			res = v
		}
	}
	return res
}

func addSeg(node, l, r, ql, qr, idx int) {
	if ql <= l && r <= qr {
		seg[node] = append(seg[node], idx)
		return
	}
	mid := (l + r) / 2
	if ql <= mid {
		addSeg(node*2, l, mid, ql, qr, idx)
	}
	if qr > mid {
		addSeg(node*2+1, mid+1, r, ql, qr, idx)
	}
}

func dfs(node, l, r int, root *LCNode) {
	for _, idx := range seg[node] {
		ln := lines[idx]
		root = insert(root, 0, m-1, Line{m: ln.a, c: ln.b})
	}
	if l == r {
		q := queries[l]
		if q.t == 3 {
			if root == nil {
				answers[l] = "EMPTY SET"
			} else {
				val := query(root, 0, m-1, q.q)
				answers[l] = fmt.Sprintf("%d", val)
			}
		}
		return
	}
	mid := (l + r) / 2
	dfs(node*2, l, mid, root)
	dfs(node*2+1, mid+1, r, root)
}

func solveCase(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return ""
	}
	if n <= 0 {
		return ""
	}

	queries = make([]Query, n+1)
	lines = make([]SegmentLine, 0)
	addMap := make(map[int]int)
	xsList := make([]int64, 0)

	for i := 1; i <= n; i++ {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var a, b int64
			fmt.Fscan(reader, &a, &b)
			queries[i] = Query{t: t}
			lineIdx := len(lines)
			lines = append(lines, SegmentLine{a: a, b: b, l: i, r: n})
			addMap[i] = lineIdx
			queries[i].idx = lineIdx
		} else if t == 2 {
			var idx int
			fmt.Fscan(reader, &idx)
			queries[i] = Query{t: t, idx: idx}
			lineIdx := addMap[idx]
			lines[lineIdx].r = i - 1
		} else {
			var qv int64
			fmt.Fscan(reader, &qv)
			queries[i] = Query{t: t, q: qv}
			xsList = append(xsList, qv)
		}
	}

	if len(xsList) == 0 {
		xsList = append(xsList, 0)
	}
	sort.Slice(xsList, func(i, j int) bool { return xsList[i] < xsList[j] })
	xs = make([]int64, 0, len(xsList))
	for i, v := range xsList {
		if i == 0 || v != xsList[i-1] {
			xs = append(xs, v)
		}
	}
	m = len(xs)

	seg = make([][]int, 4*n+5)
	for idx, ln := range lines {
		addSeg(1, 1, n, ln.l, ln.r, idx)
	}

	answers = make([]string, n+1)
	dfs(1, 1, n, nil)

	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if queries[i].t == 3 {
			sb.WriteString(answers[i])
			sb.WriteByte('\n')
		}
	}
	return strings.TrimSpace(sb.String())
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scan := bufio.NewScanner(strings.NewReader(embeddedTestcases))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t := 0
	fmt.Sscan(scan.Text(), &t)

	for caseIdx := 1; caseIdx <= t; caseIdx++ {
		if !scan.Scan() {
			fmt.Printf("missing q for case %d\n", caseIdx)
			os.Exit(1)
		}
		q, _ := strconv.Atoi(scan.Text())
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", q))
		for i := 0; i < q; i++ {
			if !scan.Scan() {
				fmt.Printf("bad test file at case %d\n", caseIdx)
				os.Exit(1)
			}
			tok := scan.Text()
			tt, _ := strconv.Atoi(tok)
			sb.WriteString(tok)
			switch tt {
			case 1:
				var a, b int
				if !scan.Scan() {
					fmt.Printf("bad test file at case %d\n", caseIdx)
					os.Exit(1)
				}
				fmt.Sscan(scan.Text(), &a)
				if !scan.Scan() {
					fmt.Printf("bad test file at case %d\n", caseIdx)
					os.Exit(1)
				}
				fmt.Sscan(scan.Text(), &b)
				sb.WriteString(fmt.Sprintf(" %d %d", a, b))
			case 2:
				var idx int
				if !scan.Scan() {
					fmt.Printf("bad test file at case %d\n", caseIdx)
					os.Exit(1)
				}
				fmt.Sscan(scan.Text(), &idx)
				sb.WriteString(fmt.Sprintf(" %d", idx))
			default:
				var v int
				if !scan.Scan() {
					fmt.Printf("bad test file at case %d\n", caseIdx)
					os.Exit(1)
				}
				fmt.Sscan(scan.Text(), &v)
				sb.WriteString(fmt.Sprintf(" %d", v))
			}
			sb.WriteByte('\n')
		}
		input := sb.String()
		expect := solveCase(input)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", caseIdx, err, got)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed:\nexpected:\n%s\n got:\n%s\n", caseIdx, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
