package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ACNode struct {
	next [26]int
	fail int
}

type Query struct {
	node  int
	index int
}

type Edge struct {
	to   int
	char byte
}

type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int, n+2)}
}

func (b *BIT) Add(i, v int) {
	for ; i <= b.n; i += i & -i {
		b.tree[i] += v
	}
}

func (b *BIT) Sum(i int) int {
	s := 0
	for ; i > 0; i -= i & -i {
		s += b.tree[i]
	}
	return s
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// Embedded copy of testcasesG.txt so the verifier is self-contained.
const testcasesRaw = `1
2 0 a
1
1 cbc

3
2 0 b
2 0 c
2 2 c
3
1 c
3 cbc
3 abb

1
2 0 b
2
1 c
1 a

2
2 0 a
1 c
2
1 bbb
1 c

2
1 c
2 1 c
1
2 aa

2
1 a
1 a
3
2 ccc
2 bc
2 ac

4
1 a
2 1 a
2 0 a
2 3 b
1
4 bc

4
1 a
2 1 b
1 a
1 c
1
3 bab

4
2 0 b
2 1 b
2 2 b
2 1 c
2
4 cb
3 ca

4
2 0 c
2 1 c
1 b
1 c
2
1 a
1 cbb

4
1 a
1 b
1 b
2 0 b
1
1 ba

2
2 0 a
1 c
1
1 bbc

1
1 c
3
1 aac
1 bbc
1 aab

4
2 0 c
2 1 c
1 b
2 2 b
1
1 bbb

4
2 0 b
1 c
1 a
1 c
3
1 abb
1 c
2 bba

2
1 b
1 b
3
2 aaa
2 aac
1 b

2
2 0 a
1 b
3
1 ac
1 bab
2 bab

2
1 c
2 1 c
1
1 c

2
2 0 a
1 a
2
1 c
1 bca

1
1 b
1
1 b

1
1 a
3
1 a
1 c
1 cbc

3
2 0 a
1 c
1 b
3
1 cc
3 cbc
1 cca

3
2 0 a
1 b
2 1 c
1
2 a

1
2 0 b
1
1 ac

3
1 a
1 b
1 a
3
2 bcc
3 bb
1 aa

4
1 a
1 a
1 a
1 c
2
4 cc
4 c

4
1 b
1 b
2 1 b
1 b
2
3 a
3 a

4
1 c
1 a
2 1 a
1 c
1
4 bc

4
2 0 b
1 a
1 c
1 c
3
1 aac
1 b
3 c

2
1 a
1 c
1
2 bca

2
2 0 c
2 0 b
3
2 b
1 c
1 cb

3
1 a
2 0 b
2 0 c
2
1 b
3 cb

1
1 a
2
1 a
1 ab

1
2 0 c
1
1 bc

3
2 0 b
1 c
1 c
3
1 bcc
2 a
2 a

3
1 a
2 0 b
2 2 b
1
3 aca

1
2 0 b
1
1 abb

3
1 c
1 b
1 c
3
1 bba
3 b
3 ccb

3
2 0 b
2 0 c
2 2 a
2
1 cbc
3 a

1
1 b
2
1 c
1 aa

3
2 0 b
2 0 c
1 b
2
2 bba
1 c

4
1 a
2 0 a
2 0 c
1 a
2
3 aa
3 a

2
1 a
1 c
2
2 bcc
2 bc

4
1 a
1 a
1 b
1 c
2
4 cbc
2 aca

4
1 b
1 c
1 b
2 3 c
2
4 bb
3 abc

2
2 0 c
2 0 a
2
1 cb
2 ab

4
2 0 a
1 b
2 2 a
1 c
3
3 aa
4 abc
1 b

2
2 0 b
2 0 c
2
2 ab
2 a

2
1 a
1 b
3
1 bbb
2 a
2 a

2
1 a
2 1 a
2
2 ab
2 c

1
1 a
2
1 bca
1 aac

2
1 c
2 1 c
1
2 a

2
2 0 c
2 0 c
1
1 ab

2
1 c
2 1 a
2
1 c
1 cbb

3
2 0 b
1 a
2 0 c
2
2 c
2 cb

4
2 0 c
2 1 a
2 2 a
2 1 c
3
1 a
2 bbc
1 c

3
1 c
1 b
2 0 b
2
1 ca
3 ab

3
1 c
2 0 c
1 a
2
3 a
1 ac

3
2 0 c
2 1 c
1 a
1
3 cac

2
2 0 b
1 c
2
2 ca
1 a

4
1 a
1 c
1 b
2 3 b
2
2 cc
4 cb

4
1 a
1 b
2 0 c
2 2 c
3
2 c
1 a
2 cc

2
1 b
1 a
1
2 b

4
1 a
2 0 a
1 b
2 1 a
2
4 ca
3 caa

4
2 0 a
2 1 c
1 c
2 1 c
2
3 bb
4 bbb

2
1 c
2 1 a
1
2 caa

1
1 a
3
1 b
1 a
1 b

2
1 a
1 c
1
1 cbc

4
1 a
2 0 b
1 c
2 2 a
3
4 cb
1 bba
1 b

1
1 c
1
1 ac

2
2 0 b
1 a
2
1 cb
1 bc

2
1 b
2 0 a
1
1 aab

1
2 0 a
2
1 baa
1 cc

3
1 a
1 a
2 0 c
2
3 bb
1 aab

2
1 a
2 0 b
2
2 aba
1 ca

2
1 a
1 a
2
1 cc
1 bc

1
1 b
1
1 cbb

2
1 b
2 0 b
1
1 aa

1
1 a
1
1 c

2
2 0 a
1 c
3
2 a
2 b
1 aca

4
2 0 c
2 0 a
2 0 b
1 c
1
1 aa

4
1 b
1 b
2 2 c
1 a
1
1 a

2
2 0 c
2 0 b
1
2 abb

3
1 c
2 1 b
1 c
2
3 bab
2 a

3
2 0 a
2 1 c
1 c
3
2 a
2 a
3 ab

2
2 0 b
2 0 a
3
1 a
2 aab
2 ac

4
1 a
1 a
1 b
2 2 c
1
2 bac

4
2 0 b
1 a
1 a
1 b
2
3 cbc
2 caa

4
2 0 b
2 0 c
2 1 c
1 b
1
4 b

3
2 0 b
1 b
1 b
2
3 ba
1 bc

1
1 a
3
1 bbc
1 c
1 bca

2
1 b
1 c
1
2 cc

3
1 c
2 1 a
1 c
1
2 c

1
2 0 a
2
1 bac
1 ba

1
1 c
1
1 ac

1
1 c
1
1 c

2
2 0 c
1 b
3
1 cac
2 c
2 c

3
1 b
2 0 a
1 c
3
2 ca
2 ab
2 b

1
2 0 b
1
1 a

4
2 0 b
2 0 c
1 c
1 b
2
3 ab
4 bab`

func parseTestcases() []string {
	raw := strings.TrimSpace(testcasesRaw)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, "\n\n")
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		res = append(res, strings.TrimSpace(p))
	}
	return res
}

// solve implements the logic from 1207G.go for a single test case input.
func solve(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return "", err
	}

	children := make([][]Edge, n+1)
	for i := 1; i <= n; i++ {
		var t int
		if _, err := fmt.Fscan(reader, &t); err != nil {
			return "", err
		}
		if t == 1 {
			var c string
			if _, err := fmt.Fscan(reader, &c); err != nil {
				return "", err
			}
			children[0] = append(children[0], Edge{to: i, char: c[0]})
		} else {
			var j int
			var c string
			if _, err := fmt.Fscan(reader, &j, &c); err != nil {
				return "", err
			}
			children[j] = append(children[j], Edge{to: i, char: c[0]})
		}
	}

	var m int
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return "", err
	}

	nodes := make([]ACNode, 1)
	insert := func(s string) int {
		v := 0
		for i := 0; i < len(s); i++ {
			c := int(s[i] - 'a')
			if nodes[v].next[c] == 0 {
				nodes = append(nodes, ACNode{})
				nodes[v].next[c] = len(nodes) - 1
			}
			v = nodes[v].next[c]
		}
		return v
	}

	queries := make([][]Query, n+1)
	for qi := 0; qi < m; qi++ {
		var idx int
		var t string
		if _, err := fmt.Fscan(reader, &idx, &t); err != nil {
			return "", err
		}
		node := insert(t)
		queries[idx] = append(queries[idx], Query{node: node, index: qi})
	}

	queue := make([]int, 0)
	for c := 0; c < 26; c++ {
		v := nodes[0].next[c]
		if v != 0 {
			queue = append(queue, v)
		}
	}
	for i := 0; i < len(queue); i++ {
		v := queue[i]
		for c := 0; c < 26; c++ {
			u := nodes[v].next[c]
			if u != 0 {
				nodes[u].fail = nodes[nodes[v].fail].next[c]
				queue = append(queue, u)
			} else {
				nodes[v].next[c] = nodes[nodes[v].fail].next[c]
			}
		}
	}

	size := len(nodes)
	childrenFail := make([][]int, size)
	for v := 1; v < size; v++ {
		p := nodes[v].fail
		childrenFail[p] = append(childrenFail[p], v)
	}

	tin := make([]int, size)
	tout := make([]int, size)
	timer := 0
	type Frame struct{ node, idx int }
	stack := []Frame{{0, 0}}
	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		if top.idx == 0 {
			timer++
			tin[top.node] = timer
		}
		if top.idx < len(childrenFail[top.node]) {
			v := childrenFail[top.node][top.idx]
			top.idx++
			stack = append(stack, Frame{v, 0})
		} else {
			tout[top.node] = timer
			stack = stack[:len(stack)-1]
		}
	}

	bit := NewBIT(timer + 2)
	ans := make([]int, m)

	type SFrame struct{ song, state, idx int }
	sstack := []SFrame{{0, 0, 0}}
	for len(sstack) > 0 {
		fr := &sstack[len(sstack)-1]
		if fr.idx == len(children[fr.song]) {
			if fr.song != 0 {
				bit.Add(tin[fr.state], -1)
			}
			sstack = sstack[:len(sstack)-1]
			continue
		}
		e := children[fr.song][fr.idx]
		fr.idx++
		ns := nodes[fr.state].next[int(e.char-'a')]
		bit.Add(tin[ns], 1)
		for _, q := range queries[e.to] {
			res := bit.Sum(tout[q.node]) - bit.Sum(tin[q.node]-1)
			ans[q.index] = res
		}
		sstack = append(sstack, SFrame{e.to, ns, 0})
	}

	var out strings.Builder
	for i := 0; i < m; i++ {
		fmt.Fprintln(&out, ans[i])
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	cand, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	cases := parseTestcases()
	for i, c := range cases {
		expect, err := solve(c)
		if err != nil {
			fmt.Printf("solver failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runCandidate(cand, c+"\n")
		if err != nil {
			fmt.Printf("candidate runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\nGot:\n%s\n", i+1, c, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
