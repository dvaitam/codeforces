package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded copy of testcasesG.txt so the verifier is self-contained.
const testcasesRaw = `3 3 c cab c 1 2 8 1 3 4 2 b
2 3 bb bbb 1 1 4 1 2 7 1 2 8
1 2 a 1 1 5 1 1 8
1 1 ba 1 1 10
1 2 a 2 ccb 2 c
3 1 bab cbb b 1 1 3
2 3 a bab 1 2 1 1 1 1 1 1 3
3 2 c c bca 1 2 1 1 3 7
1 1 bac 2 aaa
3 3 bcc ccb abc 1 3 5 1 2 1 1 2 5
2 3 bba acb 2 a 1 2 8 2 cc
2 1 c b 2 bc
3 2 b a bab 2 ab 1 2 8
2 2 ac a 1 2 1 2 b
2 1 ca ccb 1 2 2
2 2 cbc ccc 1 1 2 1 2 6
2 2 cba baa 2 cac 1 2 1
3 3 cc c cbb 1 3 8 2 abb 2 bc
3 3 bb a bb 1 1 2 1 3 10 2 c
3 1 ccb a bbc 1 1 9
3 2 ca b ba 2 ac 2 cc
2 3 bbc ca 1 1 3 1 1 4 2 cac
2 1 c a 2 a
2 2 aa cbc 1 2 2 2 cbc
2 2 cc a 2 ca 2 bbb
3 3 ab c bba 1 3 8 2 b 1 3 0
2 3 a acb 2 ac 2 a 1 1 8
2 1 cba cc 1 1 6
3 2 cb bbc c 1 2 4 1 2 1
2 2 bba caa 2 bba 1 1 5
1 2 cb 2 ca 2 b
1 3 bb 2 a 2 c 1 1 6
2 1 b bcb 1 2 4
1 3 bbc 2 cc 2 aba 1 1 7
3 3 ca cab cba 1 1 2 2 aa 2 b
2 1 bb bbc 2 abb
3 2 a bc cc 2 cb 2 abb
3 2 c b b 1 2 1 1 2 8
3 2 acc c bca 2 c 2 c
3 2 ac a bab 2 a 1 2 8
1 3 c 2 abc 1 1 4 1 1 7
2 2 ab c 1 1 9 1 1 1
1 3 bbc 2 a 1 1 1 1 1 8
1 3 cc 1 1 9 1 1 4 2 ac
1 3 ab 2 caa 2 a 1 1 8
1 2 aca 2 c 1 1 1
2 2 cc bb 2 aab 1 2 2
2 2 b ca 2 ba 2 cac
3 1 aaa abb c 2 cca
1 2 aaa 1 1 2 1 1 2
3 1 a caa a 1 1 6
3 3 ccc bcc c 2 bba 2 cab 1 1 3
2 1 ba bab 2 ab
3 2 cab cab acb 2 bcb 2 bc
3 1 b a bb 2 ba
1 1 c 2 c
3 3 ba ccb a 1 3 3 1 3 3 1 2 9
3 3 bac abc bbb 1 3 4 1 3 6 1 1 10
1 3 b 2 b 1 1 8 1 1 4
3 1 b bb bbc 2 cb
3 2 bac b abb 2 cb 1 1 0
2 3 a a 2 bc 2 b 1 1 3
3 2 cc b c 2 ab 2 b
1 2 b 1 1 0 1 1 8
3 2 bca a b 1 1 1 1 1 8
2 3 b a 2 aca 2 aa 2 aab
3 1 c bc aa 1 1 8
1 1 cc 1 1 3
2 1 a b 1 2 1
3 3 bcc cc ca 1 1 4 2 bc 1 2 4
3 2 a a aa 2 ab 1 3 5
2 3 ab cba 1 2 4 1 1 6 2 abc
3 2 bb bb bba 1 3 8 2 aba
2 3 ccb c 2 caa 1 1 5 1 1 3
1 1 cca 2 cbb
1 3 ca 2 aba 2 aa 2 c
3 3 c cb baa 2 cb 2 bbb 2 ca
3 3 a bc cba 1 2 9 1 3 10 2 ab
3 2 cca aab bc 2 ca 1 2 1
1 1 ccb 1 1 7`

const INF = int64(1) << 60

type IntHeap []int64

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int64))
}
func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// MultiSet supports insert, delete by value, and get max.
type MultiSet struct {
	maxH *IntHeap
	delH *IntHeap
}

func NewMultiSet() *MultiSet {
	mh := &IntHeap{}
	dh := &IntHeap{}
	heap.Init(mh)
	heap.Init(dh)
	return &MultiSet{mh, dh}
}

func (ms *MultiSet) Insert(v int64) {
	heap.Push(ms.maxH, v)
}

func (ms *MultiSet) Delete(v int64) {
	heap.Push(ms.delH, v)
}

func (ms *MultiSet) Top() int64 {
	for ms.maxH.Len() > 0 && ms.delH.Len() > 0 && (*ms.maxH)[0] == (*ms.delH)[0] {
		heap.Pop(ms.maxH)
		heap.Pop(ms.delH)
	}
	if ms.maxH.Len() == 0 {
		return -INF
	}
	return (*ms.maxH)[0]
}

// SegmentTree supports range max queries.
type SegmentTree struct {
	n int
	t []int64
}

func NewSegmentTree(sz int) *SegmentTree {
	n := 1
	for n < sz {
		n <<= 1
	}
	t := make([]int64, 2*n)
	for i := range t {
		t[i] = -INF
	}
	return &SegmentTree{n, t}
}

func (st *SegmentTree) Update(pos int, v int64) {
	i := pos + st.n
	st.t[i] = v
	for i >>= 1; i > 0; i >>= 1 {
		if st.t[2*i] > st.t[2*i+1] {
			st.t[i] = st.t[2*i]
		} else {
			st.t[i] = st.t[2*i+1]
		}
	}
}

// Query returns max in [l,r].
func (st *SegmentTree) Query(l, r int) int64 {
	res := -INF
	l += st.n
	r += st.n
	for l <= r {
		if (l & 1) == 1 {
			if st.t[l] > res {
				res = st.t[l]
			}
			l++
		}
		if (r & 1) == 0 {
			if st.t[r] > res {
				res = st.t[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func solve(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return "", err
	}
	names := make([]string, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(reader, &names[i]); err != nil {
			return "", err
		}
	}

	alpha := 26
	trie := [][]int{{}}
	trie[0] = make([]int, alpha)
	for i := range trie[0] {
		trie[0][i] = -1
	}
	patNode := make([]int, n)
	for i, s := range names {
		u := 0
		for _, ch := range s {
			c := int(ch - 'a')
			if trie[u][c] == -1 {
				trie = append(trie, make([]int, alpha))
				for j := range trie[len(trie)-1] {
					trie[len(trie)-1][j] = -1
				}
				trie[u][c] = len(trie) - 1
			}
			u = trie[u][c]
		}
		patNode[i] = u
	}
	sz := len(trie)

	fail := make([]int, sz)
	goTo := make([]int, sz*alpha)
	for i := 0; i < alpha; i++ {
		if trie[0][i] != -1 {
			fail[trie[0][i]] = 0
			goTo[i] = trie[0][i]
		} else {
			goTo[i] = 0
		}
	}
	q := make([]int, 0, sz)
	for i := 0; i < alpha; i++ {
		if trie[0][i] != -1 {
			q = append(q, trie[0][i])
		}
	}
	for qi := 0; qi < len(q); qi++ {
		u := q[qi]
		for c := 0; c < alpha; c++ {
			v := trie[u][c]
			if v != -1 {
				fail[v] = goTo[fail[u]*alpha+c]
				goTo[u*alpha+c] = v
				q = append(q, v)
			} else {
				goTo[u*alpha+c] = goTo[fail[u]*alpha+c]
			}
		}
	}

	adj := make([][]int, sz)
	for v := 1; v < sz; v++ {
		u := fail[v]
		adj[u] = append(adj[u], v)
	}

	size := make([]int, sz)
	heavy := make([]int, sz)
	for i := range heavy {
		heavy[i] = -1
	}
	order := make([]int, 0, sz)
	type stItem struct {
		u, idx int
	}
	stack := []stItem{{0, 0}}
	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		u := top.u
		if top.idx < len(adj[u]) {
			v := adj[u][top.idx]
			top.idx++
			stack = append(stack, stItem{v, 0})
		} else {
			order = append(order, u)
			stack = stack[:len(stack)-1]
		}
	}
	for _, u := range order {
		size[u] = 1
		maxSz := 0
		for _, v := range adj[u] {
			size[u] += size[v]
			if size[v] > maxSz {
				maxSz = size[v]
				heavy[u] = v
			}
		}
	}

	head := make([]int, sz)
	pos := make([]int, sz)
	curPos := 0
	type hp struct {
		u, h int
	}
	stk := []hp{{0, 0}}
	for len(stk) > 0 {
		x := stk[len(stk)-1]
		stk = stk[:len(stk)-1]
		u, h := x.u, x.h
		head[u] = h
		pos[u] = curPos
		curPos++
		for i := len(adj[u]) - 1; i >= 0; i-- {
			v := adj[u][i]
			if v != heavy[u] {
				stk = append(stk, hp{v, v})
			}
		}
		if heavy[u] != -1 {
			stk = append(stk, hp{heavy[u], h})
		}
	}

	counts := make([]int, sz)
	for _, u := range patNode {
		counts[u]++
	}
	msArr := make([]*MultiSet, sz)
	nodeVal := make([]int64, sz)
	for u := 0; u < sz; u++ {
		if counts[u] > 0 {
			ms := NewMultiSet()
			for i := 0; i < counts[u]; i++ {
				ms.Insert(0)
			}
			msArr[u] = ms
			nodeVal[u] = 0
		} else {
			nodeVal[u] = -INF
		}
	}

	st := NewSegmentTree(sz)
	for u := 0; u < sz; u++ {
		st.Update(pos[u], nodeVal[u])
	}

	patVal := make([]int64, n)

	var out bytes.Buffer
	writer := bufio.NewWriter(&out)
	for qi := 0; qi < m; qi++ {
		var tp int
		if _, err := fmt.Fscan(reader, &tp); err != nil {
			return "", err
		}
		if tp == 1 {
			var idx int
			var x int64
			if _, err := fmt.Fscan(reader, &idx, &x); err != nil {
				return "", err
			}
			idx--
			u := patNode[idx]
			old := patVal[idx]
			patVal[idx] = x
			ms := msArr[u]
			ms.Delete(old)
			ms.Insert(x)
			top := ms.Top()
			if top != nodeVal[u] {
				nodeVal[u] = top
				st.Update(pos[u], top)
			}
		} else {
			var qstr string
			if _, err := fmt.Fscan(reader, &qstr); err != nil {
				return "", err
			}
			cur := 0
			ans := -INF
			for _, ch := range qstr {
				c := int(ch - 'a')
				cur = goTo[cur*alpha+c]
				u := cur
				for u != 0 {
					h := head[u]
					res := st.Query(pos[h], pos[u])
					if res > ans {
						ans = res
					}
					u = fail[h]
				}
				if nodeVal[0] > ans {
					ans = nodeVal[0]
				}
			}
			if ans < 0 {
				fmt.Fprintln(writer, -1)
			} else {
				fmt.Fprintln(writer, ans)
			}
		}
	}
	writer.Flush()
	return strings.TrimSpace(out.String()), nil
}

func parseTestcases() ([]string, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var res []string
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(line, "\t") {
			return nil, fmt.Errorf("line %d contains tab characters", idx+1)
		}
		res = append(res, line)
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	return res, nil
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for i, line := range tests {
		input := line + "\n"
		expected, err := solve(input)
		if err != nil {
			fmt.Printf("test %d: failed to compute expected output: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
