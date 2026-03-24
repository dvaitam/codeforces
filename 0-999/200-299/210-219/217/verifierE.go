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
const testcasesRaw = `TA
5
3
2 2
1 1
1 4

C
2
2
1 1
1 2

T
3
2
1 1
2 2

TATGCC
2
1
4 5

TATTC
5
3
4 5
4 7
3 6

AAA
1
3
3 3
4 4
2 3

CGTC
3
2
4 4
5 5

A
1
1
1 1

CTG
1
0

ATT
2
0

TCG
2
0

CGTCT
8
3
5 5
3 3
3 3

TT
2
2
1 2
2 3

A
1
0

CCAGCT
1
0

TGCGC
4
1
4 5

C
1
3
1 1
2 2
2 2

T
3
2
1 1
1 2

TA
7
2
1 2
1 4

G
2
1
1 1

GTGCTT
7
3
4 4
7 7
4 4

CTACG
1
3
4 5
2 2
4 5

TAGGC
2
0

TG
4
3
1 1
1 3
5 5

TCATA
3
0

AAG
2
3
1 2
4 4
4 5

A
3
2
1 1
2 2

TG
2
3
1 1
2 2
2 4

GTT
3
0

AT
1
0

GACG
6
2
3 3
2 2

T
5
3
1 1
2 2
1 2

TCG
3
0

CT
4
2
1 2
1 1

CTTTCT
5
0

A
1
2
1 1
1 2

C
3
3
1 1
2 2
2 2

CTG
1
3
3 3
3 4
4 5

CGT
11
3
2 2
2 4
1 4

TAT
4
3
1 3
4 6
2 4

T
2
1
1 1

TACCA
2
1
1 3

T
1
0

CCGCA
5
2
1 5
3 9

TG
2
0

CT
4
3
2 2
2 2
4 4

T
1
0

AC
1
0

CA
1
3
1 2
4 4
2 3

TTGA
4
1
3 4

TACTAA
3
2
5 6
6 6

AGCACA
2
2
5 5
6 7

AACTA
1
0

CCTC
1
2
3 4
6 6

GG
2
0

TCGTC
7
1
1 4

TGG
5
2
3 3
3 4

A
5
3
1 1
1 2
1 2

C
1
2
1 1
1 2

CCAC
1
0

TCA
2
1
3 3

CTATTT
1
1
2 3

C
1
0

GACCGG
1
3
5 6
6 7
7 7

CT
2
1
2 2

CTGTA
3
1
2 4

CAAAC
5
2
3 5
8 8

GCCGGT
1
2
1 4
3 10

GCATTA
4
0

AC
3
1
2 2

AGGCCA
5
1
2 6

GT
5
2
1 2
3 4

ACT
1
1
2 3

C
1
1
1 1

TAACC
5
2
5 5
3 6

GT
1
0

GAA
4
2
3 3
1 2

AGATCC
1
1
2 2

A
2
2
1 1
1 2

CAAACG
6
1
1 5

TTCAA
2
0

T
5
3
1 1
1 2
4 4

CGAC
3
2
4 4
2 4

TTC
3
3
3 3
3 4
6 6

C
2
2
1 1
2 2

A
1
0

CCCTGG
6
1
6 6

TATA
2
0

G
1
0

GA
3
1
2 2

CAT
4
2
1 2
5 5

GAG
4
1
1 2

CAATA
8
1
2 5

TGTGTG
3
0

AGGT
3
1
2 3

T
1
0

CGCTG
6
3
4 4
4 4
7 7

TCCT
3
2
2 3
5 6

GGTA
10
3
3 3
5 5
1 4

CAA
3
1
2 2`

type operation struct {
	l int
	r int
}

type testCase struct {
	s   string
	k   int
	ops []operation
}

// ---------- correct solver (rope / persistent approach from CF-accepted solution) ----------

const (
	typOrig   = 0
	typConcat = 1
	typView   = 2
)

type Node struct {
	typ uint8
	len int64
	a   int64
	b   int64
	l   *Node
	r   *Node
}

var pool []Node
var poolPtr int

func alloc() *Node {
	n := &pool[poolPtr]
	poolPtr++
	return n
}

func makeOrig(start, step, cnt int64) *Node {
	if cnt <= 0 {
		return nil
	}
	if cnt == 1 {
		step = 1
	}
	n := alloc()
	n.typ = typOrig
	n.len = cnt
	n.a = start
	n.b = step
	n.l = nil
	n.r = nil
	return n
}

func locate(node *Node, pos int64) int64 {
	for {
		switch node.typ {
		case typOrig:
			return node.a + node.b*(pos-1)
		case typView:
			pos = node.a + node.b*(pos-1)
			node = node.l
		default:
			if pos <= node.l.len {
				node = node.l
			} else {
				pos -= node.l.len
				node = node.r
			}
		}
	}
}

func makeView(child *Node, offset, step, cnt int64) *Node {
	if child == nil || cnt <= 0 {
		return nil
	}
	if cnt == 1 {
		return makeOrig(locate(child, offset), 1, 1)
	}
	for {
		if offset == 1 && step == 1 && cnt == child.len {
			return child
		}
		switch child.typ {
		case typOrig:
			return makeOrig(child.a+child.b*(offset-1), child.b*step, cnt)
		case typView:
			offset = child.a + child.b*(offset-1)
			step = child.b * step
			child = child.l
		default:
			lenL := child.l.len
			last := offset + step*(cnt-1)
			if last <= lenL {
				child = child.l
			} else if offset > lenL {
				offset -= lenL
				child = child.r
			} else {
				n := alloc()
				n.typ = typView
				n.len = cnt
				n.a = offset
				n.b = step
				n.l = child
				n.r = nil
				return n
			}
		}
	}
}

func makeConcat(a, b *Node) *Node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.typ == typOrig && b.typ == typOrig && a.b == b.b && a.a+a.b*a.len == b.a {
		return makeOrig(a.a, a.b, a.len+b.len)
	}
	if a.typ == typView && b.typ == typView && a.l == b.l && a.b == b.b && a.a+a.b*a.len == b.a {
		return makeView(a.l, a.a, a.b, a.len+b.len)
	}
	n := alloc()
	n.typ = typConcat
	n.len = a.len + b.len
	n.a = 0
	n.b = 0
	n.l = a
	n.r = b
	return n
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

type Task struct {
	node *Node
	a    int64
	d    int64
	m    int64
}

func materialize(root *Node, src []byte) []byte {
	out := make([]byte, int(root.len))
	idx := 0
	stack := make([]Task, 0, 1024)
	n := root
	var a, d, m int64 = 1, 1, root.len
	for {
		for {
			if n == nil || m <= 0 {
				break
			}
			switch n.typ {
			case typView:
				a = n.a + n.b*(a-1)
				if m == 1 {
					d = 1
				} else {
					d = n.b * d
				}
				n = n.l
			case typConcat:
				lenL := n.l.len
				if a > lenL {
					a -= lenL
					n = n.r
					continue
				}
				t := (lenL-a)/d + 1
				if t >= m {
					n = n.l
					continue
				}
				stack = append(stack, Task{node: n.r, a: a + d*t - lenL, d: d, m: m - t})
				m = t
				n = n.l
			default:
				start := n.a + n.b*(a-1)
				if m == 1 {
					out[idx] = src[int(start-1)]
					idx++
				} else {
					step := n.b * d
					if step == 1 {
						s := int(start - 1)
						e := s + int(m)
						copy(out[idx:idx+int(m)], src[s:e])
						idx += int(m)
					} else {
						pos := start - 1
						for i := int64(0); i < m; i++ {
							out[idx] = src[int(pos)]
							idx++
							pos += step
						}
					}
				}
				break
			}
			if n == nil || n.typ == typOrig {
				break
			}
		}
		if len(stack) == 0 {
			break
		}
		t := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		n = t.node
		a = t.a
		d = t.d
		m = t.m
	}
	return out
}

func solveCorrect(tc testCase) string {
	src := []byte(tc.s)
	k := int64(tc.k)
	nm := len(tc.ops)
	pool = make([]Node, 30*nm+100)
	poolPtr = 0
	curLen := int64(len(src))
	root := makeOrig(1, 1, min64(curLen, k))
	for i := 0; i < nm; i++ {
		l := int64(tc.ops[i].l)
		r := int64(tc.ops[i].r)
		m := r - l + 1
		curLen += m
		newLen := min64(curLen, k)
		if r >= newLen {
			root = makeView(root, 1, 1, newLen)
			continue
		}
		a := makeView(root, 1, 1, r)
		b := makeView(root, l, 1, m)
		copyLen := min64(m, newLen-r)
		evenCnt := m / 2
		var d *Node
		if copyLen <= evenCnt {
			d = makeView(b, 2, 2, copyLen)
		} else {
			d = makeConcat(makeView(b, 2, 2, evenCnt), makeView(b, 1, 2, copyLen-evenCnt))
		}
		tailLen := newLen - r - copyLen
		e := makeView(root, r+1, 1, tailLen)
		root = makeConcat(makeConcat(a, d), e)
	}
	out := materialize(root, src)
	return string(out)
}

// ---------- end correct solver ----------

func parseCase(lines []string, start int) (testCase, int, error) {
	if start >= len(lines) {
		return testCase{}, start, fmt.Errorf("no data")
	}
	s := strings.TrimSpace(lines[start])
	start++
	if start >= len(lines) {
		return testCase{}, start, fmt.Errorf("missing k")
	}
	k, err := strconv.Atoi(strings.TrimSpace(lines[start]))
	if err != nil {
		return testCase{}, start, fmt.Errorf("parse k: %w", err)
	}
	start++
	if start >= len(lines) {
		return testCase{}, start, fmt.Errorf("missing n")
	}
	n, err := strconv.Atoi(strings.TrimSpace(lines[start]))
	if err != nil {
		return testCase{}, start, fmt.Errorf("parse n: %w", err)
	}
	start++
	ops := make([]operation, n)
	for i := 0; i < n; i++ {
		if start >= len(lines) {
			return testCase{}, start, fmt.Errorf("missing op %d", i+1)
		}
		fields := strings.Fields(lines[start])
		if len(fields) != 2 {
			return testCase{}, start, fmt.Errorf("op %d: wrong field count", i+1)
		}
		l, err := strconv.Atoi(fields[0])
		if err != nil {
			return testCase{}, start, fmt.Errorf("op %d: parse l: %w", i+1, err)
		}
		r, err := strconv.Atoi(fields[1])
		if err != nil {
			return testCase{}, start, fmt.Errorf("op %d: parse r: %w", i+1, err)
		}
		ops[i] = operation{l: l, r: r}
		start++
	}
	return testCase{s: s, k: k, ops: ops}, start, nil
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for i := 0; i < len(lines); {
		for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
			i++
		}
		if i >= len(lines) {
			break
		}
		tc, next, err := parseCase(lines, i)
		if err != nil {
			return nil, fmt.Errorf("case starting at line %d: %w", i+1, err)
		}
		cases = append(cases, tc)
		i = next
	}
	return cases, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(tc.s)
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", tc.k))
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.ops)))
	for _, op := range tc.ops {
		fmt.Fprintf(&sb, "%d %d\n", op.l, op.r)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range testcases {
		expect := solveCorrect(tc)
		input := buildInput(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
