package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded correct solver (accepted on CF)

type oSegTree struct {
	tree []int
	lazy []int
}

func oNewSegTree(n int) oSegTree {
	sz := n
	if sz < 1 {
		sz = 1
	}
	return oSegTree{
		tree: make([]int, 4*sz+4),
		lazy: make([]int, 4*sz+4),
	}
}

func (st *oSegTree) push(node int) {
	if st.lazy[node] != 0 {
		st.tree[2*node] += st.lazy[node]
		st.lazy[2*node] += st.lazy[node]
		st.tree[2*node+1] += st.lazy[node]
		st.lazy[2*node+1] += st.lazy[node]
		st.lazy[node] = 0
	}
}

func (st *oSegTree) add(node, l, r, ql, qr, val int) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.tree[node] += val
		st.lazy[node] += val
		return
	}
	st.push(node)
	mid := l + (r-l)/2
	st.add(2*node, l, mid, ql, qr, val)
	st.add(2*node+1, mid+1, r, ql, qr, val)
	if st.tree[2*node] > st.tree[2*node+1] {
		st.tree[node] = st.tree[2*node]
	} else {
		st.tree[node] = st.tree[2*node+1]
	}
}

func (st *oSegTree) query(node, l, r, ql, qr int) int {
	if ql > r || qr < l {
		return -1000000000
	}
	if ql <= l && r <= qr {
		return st.tree[node]
	}
	st.push(node)
	mid := l + (r-l)/2
	left := st.query(2*node, l, mid, ql, qr)
	right := st.query(2*node+1, mid+1, r, ql, qr)
	if left > right {
		return left
	}
	return right
}

func oracleSolve(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, 1024*1024*10), 1024*1024*10)

	scanString := func() string {
		scanner.Scan()
		return scanner.Text()
	}
	scanInt := func() int {
		scanner.Scan()
		res := 0
		s := scanner.Text()
		for _, c := range s {
			res = res*10 + int(c-'0')
		}
		return res
	}

	n := scanInt()
	q := scanInt()

	parent := make([]int, n+1)
	adj := make([][]int, n+1)
	initialChar := make([]byte, n+1)
	edgeChar := make([]byte, n+1)

	for i := 1; i <= n-1; i++ {
		p := scanInt()
		cStr := scanString()
		u := i + 1
		parent[u] = p
		adj[p] = append(adj[p], u)
		initialChar[u] = cStr[0]
		edgeChar[u] = '?'
	}

	sz := make([]int, n+1)
	depth := make([]int, n+1)
	heavyChild := make([]int, n+1)

	var dfsSz func(int)
	dfsSz = func(u int) {
		sz[u] = 1
		maxSub := 0
		for _, v := range adj[u] {
			depth[v] = depth[u] + 1
			dfsSz(v)
			sz[u] += sz[v]
			if sz[v] > maxSub {
				maxSub = sz[v]
				heavyChild[u] = v
			}
		}
	}
	dfsSz(1)

	hldHead := make([]int, n+1)
	hldPos := make([]int, n+1)
	hldNode := make([]int, n+1)
	hldPosIn := make([]int, n+1)
	hldPosOut := make([]int, n+1)
	hldTimer := 0

	var dfsHld func(int, int)
	dfsHld = func(u, h int) {
		hldHead[u] = h
		hldTimer++
		hldPos[u] = hldTimer
		hldNode[hldTimer] = u
		hldPosIn[u] = hldTimer
		if heavyChild[u] != 0 {
			dfsHld(heavyChild[u], h)
		}
		for _, v := range adj[u] {
			if v != heavyChild[u] {
				dfsHld(v, v)
			}
		}
		hldPosOut[u] = hldTimer
	}
	dfsHld(1, 1)

	leafIn := make([]int, n+1)
	leafOut := make([]int, n+1)
	leafTimer := 0
	possible := true
	D := -1

	var dfsLeaves func(int)
	dfsLeaves = func(u int) {
		isLeaf := true
		leafIn[u] = leafTimer
		for _, v := range adj[u] {
			isLeaf = false
			dfsLeaves(v)
		}
		if isLeaf {
			if D == -1 {
				D = depth[u]
			} else if D != depth[u] {
				possible = false
			}
			leafTimer++
		}
		leafOut[u] = leafTimer - 1
	}
	dfsLeaves(1)

	numLeaves := leafTimer
	if numLeaves < 1 {
		numLeaves = 1
	}
	var leafST [26]oSegTree
	for i := 0; i < 26; i++ {
		leafST[i] = oNewSegTree(numLeaves)
	}
	globalST := oNewSegTree(n)

	for i := 1; i <= n; i++ {
		v := hldNode[i]
		globalST.add(1, 1, n, i, i, depth[v]-D)
	}

	updateEdge := func(u int, newC byte) {
		oldC := edgeChar[u]
		if oldC == newC {
			return
		}
		if oldC == '?' {
			globalST.add(1, 1, n, hldPosIn[u], hldPosOut[u], -1)
		} else {
			c := int(oldC - 'a')
			leafST[c].add(1, 0, numLeaves-1, leafIn[u], leafOut[u], -1)
			target := leafST[c].query(1, 0, numLeaves-1, leafIn[u], leafOut[u])

			topNode := u
			curr := parent[u]
			for curr > 0 {
				if leafST[c].query(1, 0, numLeaves-1, leafIn[curr], leafOut[curr]) != target {
					break
				}
				head := hldHead[curr]
				if leafST[c].query(1, 0, numLeaves-1, leafIn[head], leafOut[head]) == target {
					topNode = head
					curr = parent[head]
				} else {
					low := hldPos[head]
					high := hldPos[curr]
					ansPos := high
					for low <= high {
						mid := low + (high-low)/2
						node := hldNode[mid]
						if leafST[c].query(1, 0, numLeaves-1, leafIn[node], leafOut[node]) == target {
							ansPos = mid
							high = mid - 1
						} else {
							low = mid + 1
						}
					}
					topNode = hldNode[ansPos]
					break
				}
			}

			globalST.add(1, 1, n, hldPosIn[u], hldPosOut[u], -1)
			if topNode != u {
				curr = parent[u]
				for depth[curr] >= depth[topNode] {
					head := hldHead[curr]
					if depth[head] < depth[topNode] {
						head = topNode
					}
					globalST.add(1, 1, n, hldPos[head], hldPos[curr], -1)
					curr = parent[head]
					if curr == 0 {
						break
					}
				}
			}
		}

		if newC == '?' {
			globalST.add(1, 1, n, hldPosIn[u], hldPosOut[u], 1)
		} else {
			c := int(newC - 'a')
			target := leafST[c].query(1, 0, numLeaves-1, leafIn[u], leafOut[u])

			topNode := u
			curr := parent[u]
			for curr > 0 {
				if leafST[c].query(1, 0, numLeaves-1, leafIn[curr], leafOut[curr]) != target {
					break
				}
				head := hldHead[curr]
				if leafST[c].query(1, 0, numLeaves-1, leafIn[head], leafOut[head]) == target {
					topNode = head
					curr = parent[head]
				} else {
					low := hldPos[head]
					high := hldPos[curr]
					ansPos := high
					for low <= high {
						mid := low + (high-low)/2
						node := hldNode[mid]
						if leafST[c].query(1, 0, numLeaves-1, leafIn[node], leafOut[node]) == target {
							ansPos = mid
							high = mid - 1
						} else {
							low = mid + 1
						}
					}
					topNode = hldNode[ansPos]
					break
				}
			}

			leafST[c].add(1, 0, numLeaves-1, leafIn[u], leafOut[u], 1)

			globalST.add(1, 1, n, hldPosIn[u], hldPosOut[u], 1)
			if topNode != u {
				curr = parent[u]
				for depth[curr] >= depth[topNode] {
					head := hldHead[curr]
					if depth[head] < depth[topNode] {
						head = topNode
					}
					globalST.add(1, 1, n, hldPos[head], hldPos[curr], 1)
					curr = parent[head]
					if curr == 0 {
						break
					}
				}
			}
		}
		edgeChar[u] = newC
	}

	for i := 2; i <= n; i++ {
		if initialChar[i] != '?' {
			updateEdge(i, initialChar[i])
		}
	}

	var results []string
	for i := 0; i < q; i++ {
		v := scanInt()
		cStr := scanString()
		c := cStr[0]

		updateEdge(v, c)

		if !possible {
			results = append(results, "Fou")
		} else if globalST.query(1, 1, n, 1, n) > 0 {
			results = append(results, "Fou")
		} else {
			S := 0
			var maxC [26]int
			for j := 0; j < 26; j++ {
				maxC[j] = leafST[j].query(1, 0, numLeaves-1, 0, numLeaves-1)
				S += maxC[j]
			}
			ans := 0
			for j := 0; j < 26; j++ {
				f := maxC[j] + D - S
				ans += f * (j + 1)
			}
			results = append(results, fmt.Sprintf("Shi %d", ans))
		}
	}
	return strings.Join(results, "\n")
}

// End of embedded solver

const testcasesRaw = `100
3 5 1 c 2 a 2 ? 3 b 2 ? 3 ? 2 b
3 5 1 a 1 b 2 c 2 c 3 ? 3 ? 3 b
4 1 1 b 2 b 2 ? 4 c
5 5 1 c 2 b 2 a 3 b 4 a 3 c 4 a 2 ? 5 a
4 1 1 b 1 c 2 ? 2 a
6 5 1 ? 2 c 3 b 1 c 1 a 2 a 3 ? 4 c 3 a 4 c
4 2 1 ? 2 ? 3 a 4 c 3 b
4 4 1 c 2 a 2 c 2 ? 4 b 2 c 3 c
4 5 1 ? 1 a 3 a 3 c 4 ? 3 c 2 c 2 c
4 5 1 c 2 a 1 b 3 b 4 c 2 c 2 ? 4 a
2 5 1 c 2 ? 2 a 2 b 2 c 2 a
2 5 1 c 2 c 2 a 2 b 2 c 2 ?
4 4 1 ? 2 a 2 b 2 a 3 ? 4 b 2 ?
6 3 1 b 1 c 1 b 1 a 5 b 5 a 2 ? 2 b
6 3 1 a 2 a 3 a 3 b 3 ? 2 c 3 b 2 a
3 2 1 b 1 ? 3 a 3 b
4 5 1 a 2 c 1 a 2 a 2 a 2 ? 2 a 4 ?
4 2 1 a 2 ? 3 ? 4 c 3 c
3 3 1 a 1 a 3 a 2 a 3 ?
6 5 1 a 2 a 2 ? 3 ? 4 ? 2 b 3 c 6 a 5 b 5 b
2 3 1 c 2 ? 2 ? 2 c
6 5 1 a 2 b 1 ? 1 a 5 a 5 b 2 c 2 a 2 ? 4 c
2 1 1 a 2 a
6 3 1 a 1 b 3 b 4 ? 3 c 6 ? 4 ? 2 ?
6 2 1 b 2 ? 1 ? 2 b 1 ? 5 ? 6 b
3 3 1 b 2 b 3 ? 3 b 3 a
4 4 1 b 1 c 1 c 3 b 3 ? 4 b 3 a
5 1 1 a 2 b 1 a 2 c 3 b
4 2 1 a 2 b 1 c 2 ? 2 a
2 1 1 c 2 c
5 5 1 a 1 c 2 ? 4 ? 2 b 5 ? 3 c 2 c 2 ?
2 4 1 a 2 c 2 c 2 a 2 a
5 5 1 a 2 b 3 b 2 c 5 a 2 b 4 ? 4 b 5 ?
4 2 1 a 1 ? 2 a 3 c 3 a
3 1 1 c 1 c 3 ?
5 5 1 c 2 a 1 c 1 a 2 b 3 ? 5 b 3 ? 3 b
6 2 1 c 2 c 1 ? 4 ? 3 c 5 c 5 a
6 2 1 a 1 ? 1 ? 2 b 5 b 2 ? 2 c
3 2 1 a 1 a 3 b 2 ?
6 1 1 c 2 c 3 b 1 a 1 c 3 a
3 4 1 ? 2 a 2 ? 2 b 2 ? 3 ?
2 5 1 b 2 c 2 ? 2 ? 2 b 2 a
4 3 1 ? 2 ? 1 a 2 c 3 b 3 b
3 2 1 b 2 b 2 b 2 b
5 4 1 a 1 ? 2 c 4 a 2 b 5 a 5 ? 2 b
3 1 1 ? 1 b 3 a
4 1 1 c 2 ? 2 ? 2 c
2 3 1 c 2 a 2 a 2 a
2 1 1 b 2 ?
2 2 1 b 2 c 2 c
2 4 1 ? 2 c 2 ? 2 ? 2 ?
3 4 1 b 2 ? 3 b 3 ? 3 c 2 ?
5 1 1 b 1 ? 1 a 3 b 2 a
5 4 1 c 1 a 1 ? 1 c 4 a 5 a 3 b 2 b
3 5 1 a 1 c 2 b 3 c 2 ? 2 a 2 a
6 1 1 b 2 a 1 ? 4 b 5 a 5 c
5 5 1 a 2 ? 2 b 1 a 5 c 2 c 4 b 5 a 4 ?
2 5 1 b 2 a 2 b 2 ? 2 c 2 c
3 2 1 ? 2 c 2 ? 3 ?
5 5 1 c 1 b 1 a 2 c 2 ? 2 b 3 c 5 ? 2 c
5 2 1 b 1 ? 2 ? 4 a 3 b 5 c
5 3 1 ? 2 ? 2 b 3 c 4 a 5 c 4 c
3 5 1 a 1 b 2 b 2 a 2 ? 3 c 2 ?
3 2 1 c 2 c 2 a 3 c
4 2 1 ? 2 c 1 a 3 a 4 b
6 2 1 ? 2 b 2 a 2 a 1 ? 2 ? 5 ?
4 5 1 c 1 a 3 b 4 a 3 c 3 ? 3 c 3 c
6 2 1 b 1 ? 2 b 4 c 1 c 6 a 5 b
4 3 1 ? 1 c 1 ? 4 ? 2 c 3 b
5 4 1 a 2 b 1 c 3 ? 4 a 3 ? 4 ? 2 ?
3 2 1 a 1 a 3 a 3 c
5 4 1 c 1 b 1 ? 2 b 3 b 5 ? 2 b 5 a
3 2 1 b 1 b 2 b 2 ?
2 4 1 b 2 c 2 b 2 b 2 b
5 3 1 b 1 c 1 ? 2 a 4 ? 2 c 5 a
4 5 1 ? 2 c 2 b 4 ? 3 b 4 c 4 b 3 b
4 3 1 b 1 a 2 b 2 a 2 a 3 c
6 3 1 ? 2 c 1 b 3 ? 5 c 4 c 6 c 4 c
4 3 1 a 1 c 3 b 3 ? 3 ? 2 ?
4 5 1 ? 1 a 2 ? 4 c 2 b 3 ? 2 a 3 b
2 5 1 a 2 ? 2 a 2 c 2 ? 2 c
2 5 1 b 2 b 2 c 2 a 2 a 2 c
5 3 1 a 1 c 3 ? 3 b 5 c 2 ? 3 ?
3 3 1 b 1 b 2 ? 2 c 3 a
5 2 1 c 1 ? 1 ? 3 ? 4 a 2 a
5 3 1 a 2 ? 2 a 3 ? 3 c 3 ? 2 a
6 1 1 a 1 b 2 ? 3 c 4 b 3 c
4 5 1 c 2 c 3 b 3 a 4 ? 3 b 3 b 2 a
5 1 1 b 2 b 1 c 4 b 5 a
3 4 1 ? 1 b 3 c 3 c 2 c 2 b
3 3 1 b 1 b 3 b 2 c 2 b
6 4 1 b 2 ? 2 a 1 ? 4 a 6 b 4 ? 5 b 5 ?
4 4 1 b 1 ? 2 b 4 c 2 c 4 b 3 ?
6 2 1 c 2 a 3 c 2 a 5 c 6 c 6 c
5 3 1 ? 2 ? 2 c 2 b 4 b 5 a 4 b
4 4 1 ? 2 b 2 c 2 a 4 ? 4 b 2 b
2 1 1 ? 2 ?
4 1 1 c 2 b 2 a 2 a
5 3 1 b 2 c 3 b 3 b 3 b 4 c 4 a
2 2 1 b 2 c 2 c`

type testcase struct {
	n   int
	q   int
	par []int
	ch  []byte
	qs  []tquery
}

type tquery struct {
	v int
	c byte
}

var testcases = mustParseTestcases(testcasesRaw)

func mustParseTestcases(raw string) []testcase {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	start := 0
	firstFields := strings.Fields(lines[0])
	if len(firstFields) == 1 {
		start = 1
	}
	res := make([]testcase, 0, len(lines)-start)
	for idx := start; idx < len(lines); idx++ {
		line := lines[idx]
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			panic(fmt.Sprintf("line %d: too short", idx+1))
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			panic(fmt.Sprintf("line %d: bad n: %v", idx+1, err))
		}
		q, err := strconv.Atoi(fields[1])
		if err != nil {
			panic(fmt.Sprintf("line %d: bad q: %v", idx+1, err))
		}
		expectedLen := 2 + 2*(n-1) + 2*q
		if len(fields) != expectedLen {
			panic(fmt.Sprintf("line %d: expected %d fields got %d", idx+1, expectedLen, len(fields)))
		}
		parents := make([]int, n+1)
		chars := make([]byte, n+1)
		pos := 2
		for v := 2; v <= n; v++ {
			p, _ := strconv.Atoi(fields[pos])
			parents[v] = p
			pos++
			cf := fields[pos]
			if len(cf) != 1 {
				panic(fmt.Sprintf("line %d: invalid char %q", idx+1, cf))
			}
			chars[v] = cf[0]
			pos++
		}
		qs := make([]tquery, q)
		for i := 0; i < q; i++ {
			v, _ := strconv.Atoi(fields[pos])
			pos++
			cf := fields[pos]
			if len(cf) != 1 {
				panic(fmt.Sprintf("line %d: invalid query char %q", idx+1, cf))
			}
			pos++
			qs[i] = tquery{v: v, c: cf[0]}
		}
		res = append(res, testcase{n: n, q: q, par: parents, ch: chars, qs: qs})
	}
	return res
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
	return out.String(), nil
}

func parseCandidateOutput(out string) []string {
	lines := []string{}
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func buildInput(tc testcase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
	for v := 2; v <= tc.n; v++ {
		sb.WriteString(fmt.Sprintf("%d %c\n", tc.par[v], tc.ch[v]))
	}
	for _, q := range tc.qs {
		sb.WriteString(fmt.Sprintf("%d %c\n", q.v, q.c))
	}
	return sb.String()
}

func checkCase(bin string, idx int, tc testcase) error {
	input := buildInput(tc)
	expectedStr := oracleSolve(input)
	expected := parseCandidateOutput(expectedStr)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	got := parseCandidateOutput(out)
	if len(got) != len(expected) {
		return fmt.Errorf("expected %d outputs, got %d", len(expected), len(got))
	}
	for i := range expected {
		if got[i] != expected[i] {
			return fmt.Errorf("output %d: expected %s got %s", i+1, expected[i], got[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i, tc := range testcases {
		if err := checkCase(bin, i, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
