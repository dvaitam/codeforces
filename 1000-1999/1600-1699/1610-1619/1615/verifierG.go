package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSourceG = `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	hasExistingPair := make([]bool, 605)
	for i := 0; i < n-1; i++ {
		if a[i] == a[i+1] && a[i] > 0 {
			hasExistingPair[a[i]] = true
		}
	}

	inC := make([]bool, 605)
	for i := 0; i < n; i++ {
		if a[i] > 0 && !hasExistingPair[a[i]] {
			inC[a[i]] = true
		}
	}

	type Block struct {
		id    int
		start int
		end   int
		L     int
		u     int
		v     int
	}

	var oddBlocks []Block
	var evenBlocks []Block

	for i := 0; i < n; i++ {
		if a[i] == 0 {
			j := i
			for j < n && a[j] == 0 {
				j++
			}
			L := j - i
			u, v := -1, -1
			if i > 0 && inC[a[i-1]] {
				u = a[i-1]
			}
			if j < n && inC[a[j]] {
				v = a[j]
			}

			if L%2 != 0 {
				oddBlocks = append(oddBlocks, Block{id: len(oddBlocks) + 1, start: i, end: j - 1, L: L, u: u, v: v})
			} else {
				evenBlocks = append(evenBlocks, Block{id: len(evenBlocks) + 1, start: i, end: j - 1, L: L, u: u, v: v})
			}
			i = j - 1
		}
	}

	adjOdd := make([][]int, 605)
	selfLoops := make([]int, 605)
	for _, b := range oddBlocks {
		if b.u != -1 && b.v != -1 {
			if b.u != b.v {
				adjOdd[b.u] = append(adjOdd[b.u], b.v)
				adjOdd[b.v] = append(adjOdd[b.v], b.u)
			} else {
				selfLoops[b.u]++
			}
		} else if b.u != -1 {
			selfLoops[b.u]++
		} else if b.v != -1 {
			selfLoops[b.v]++
		}
	}

	visited := make([]bool, 605)
	compID := make([]int, 605)
	isCycle := make([]bool, 605)
	currComp := 0

	for i := 1; i <= 600; i++ {
		if inC[i] && !visited[i] {
			currComp++
			var q []int
			q = append(q, i)
			visited[i] = true
			compID[i] = currComp

			nodes := 0
			edges := 0

			for head := 0; head < len(q); head++ {
				u := q[head]
				nodes++
				edges += selfLoops[u]
				for _, v := range adjOdd[u] {
					edges++
					if !visited[v] {
						visited[v] = true
						compID[v] = currComp
						q = append(q, v)
					}
				}
			}
			edges /= 2
			for _, u := range q {
				edges += selfLoops[u]
			}

			if edges >= nodes {
				isCycle[currComp] = true
			}
		}
	}

	treeMapping := make([]int, currComp+1)
	numTrees := 0
	for c := 1; c <= currComp; c++ {
		if !isCycle[c] {
			numTrees++
			treeMapping[c] = numTrees
		}
	}

	edgeID := make(map[string]int)
	for _, b := range evenBlocks {
		if b.u != -1 && b.v != -1 && b.u != b.v {
			cu, cv := compID[b.u], compID[b.v]
			if isCycle[cu] && isCycle[cv] {
				continue
			}
			if isCycle[cu] {
				tv := treeMapping[cv]
				edgeID[fmt.Sprintf("%d,%d", tv, tv)] = b.id
			} else if isCycle[cv] {
				tu := treeMapping[cu]
				edgeID[fmt.Sprintf("%d,%d", tu, tu)] = b.id
			} else {
				tu, tv := treeMapping[cu], treeMapping[cv]
				if tu == tv {
					edgeID[fmt.Sprintf("%d,%d", tu, tu)] = b.id
				} else {
					minT, maxT := tu, tv
					if minT > maxT {
						minT, maxT = maxT, minT
					}
					edgeID[fmt.Sprintf("%d,%d", minT, maxT)] = b.id
				}
			}
		}
	}

	V := 2 * numTrees
	adjMeta := make([][]int, V+1)
	for k := range edgeID {
		var u, v int
		fmt.Sscanf(k, "%d,%d", &u, &v)
		if u == v {
			adjMeta[u] = append(adjMeta[u], numTrees+u)
			adjMeta[numTrees+u] = append(adjMeta[numTrees+u], u)
		} else {
			adjMeta[u] = append(adjMeta[u], v)
			adjMeta[v] = append(adjMeta[v], u)
		}
	}

	matchMeta := blossom(V, adjMeta)

	usedEven := make(map[int]bool)
	satisfied := make([]bool, 605)

	for i := 1; i <= numTrees; i++ {
		m := matchMeta[i]
		if m > 0 && i < m {
			var k string
			if m > numTrees {
				k = fmt.Sprintf("%d,%d", i, i)
			} else {
				k = fmt.Sprintf("%d,%d", i, m)
			}
			bID := edgeID[k]
			usedEven[bID] = true
			b := evenBlocks[bID-1]
			satisfied[b.u] = true
			satisfied[b.v] = true
			a[b.start] = b.u
			a[b.end] = b.v
		}
	}

	var leftNodes []int
	leftMap := make([]int, 605)
	for i := 1; i <= 600; i++ {
		if inC[i] && !satisfied[i] {
			leftNodes = append(leftNodes, i)
			leftMap[i] = len(leftNodes)
		}
	}

	nLeft := len(leftNodes)
	mRight := len(oddBlocks)
	adjBip := make([][]int, nLeft+1)
	for j, b := range oddBlocks {
		if b.u != -1 && !satisfied[b.u] {
			uIdx := leftMap[b.u]
			adjBip[uIdx] = append(adjBip[uIdx], j+1)
		}
		if b.v != -1 && b.v != b.u && !satisfied[b.v] {
			vIdx := leftMap[b.v]
			adjBip[vIdx] = append(adjBip[vIdx], j+1)
		}
	}

	matchBip := hopcroftKarp(nLeft, mRight, adjBip)
	for i := 1; i <= nLeft; i++ {
		j := matchBip[i]
		if j > 0 {
			x := leftNodes[i-1]
			b := oddBlocks[j-1]
			if b.u == x {
				a[b.start] = x
			} else {
				a[b.end] = x
			}
			satisfied[x] = true
		}
	}

	hasPair := make([]bool, n+1)
	for i := 0; i < n-1; i++ {
		if a[i] == a[i+1] && a[i] > 0 {
			hasPair[a[i]] = true
		}
	}

	val := 1
	for i := 0; i < n; i++ {
		if a[i] == 0 {
			j := i
			for j < n && a[j] == 0 {
				j++
			}
			for k := i; k+1 < j; k += 2 {
				for val <= n && hasPair[val] {
					val++
				}
				v := val
				if v > n {
					v = 1
				} else {
					hasPair[v] = true
				}
				a[k] = v
				a[k+1] = v
			}
			if (j-i)%2 != 0 {
				a[j-1] = 1
			}
			i = j - 1
		}
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, a[i])
	}
	fmt.Fprintln(out)
}

func hopcroftKarp(n int, m int, edges [][]int) []int {
	match := make([]int, n+1)
	matchR := make([]int, m+1)
	dist := make([]int, n+1)

	var bfs func() bool
	bfs = func() bool {
		var q []int
		for i := 1; i <= n; i++ {
			if match[i] == 0 {
				dist[i] = 0
				q = append(q, i)
			} else {
				dist[i] = 1e9
			}
		}
		dist[0] = 1e9
		for len(q) > 0 {
			u := q[0]
			q = q[1:]
			if dist[u] < dist[0] {
				for _, v := range edges[u] {
					if dist[matchR[v]] == 1e9 {
						dist[matchR[v]] = dist[u] + 1
						q = append(q, matchR[v])
					}
				}
			}
		}
		return dist[0] != 1e9
	}

	var dfs func(int) bool
	dfs = func(u int) bool {
		if u != 0 {
			for _, v := range edges[u] {
				if dist[matchR[v]] == dist[u]+1 {
					if dfs(matchR[v]) {
						match[u] = v
						matchR[v] = u
						return true
					}
				}
			}
			dist[u] = 1e9
			return false
		}
		return true
	}

	for bfs() {
		for i := 1; i <= n; i++ {
			if match[i] == 0 {
				dfs(i)
			}
		}
	}
	return match
}

func blossom(n int, adj [][]int) []int {
	match := make([]int, n+1)
	p := make([]int, n+1)
	base := make([]int, n+1)
	q := make([]int, n+1)
	inq := make([]bool, n+1)
	inb := make([]bool, n+1)

	lca := func(u, v int) int {
		for i := 1; i <= n; i++ {
			inb[i] = false
		}
		for {
			u = base[u]
			inb[u] = true
			if match[u] == 0 {
				break
			}
			u = p[match[u]]
		}
		for {
			v = base[v]
			if inb[v] {
				return v
			}
			v = p[match[v]]
		}
	}

	mark := func(u, v, b int, head *int, tail *int) {
		for base[u] != b {
			p[u] = v
			v = match[u]
			inb[base[u]] = true
			inb[base[v]] = true
			u = p[v]
			if base[v] != b {
				q[*tail] = v
				*tail++
				inq[v] = true
			}
		}
	}

	for i := 1; i <= n; i++ {
		if match[i] == 0 {
			for j := 1; j <= n; j++ {
				base[j] = j
				p[j] = 0
				inq[j] = false
			}
			head, tail := 0, 0
			q[tail] = i
			tail++
			inq[i] = true
			found := false
			for head < tail && !found {
				u := q[head]
				head++
				for _, v := range adj[u] {
					if base[u] != base[v] && match[u] != v {
						if v == i || (match[v] > 0 && p[match[v]] > 0) {
							b := lca(u, v)
							for j := 1; j <= n; j++ {
								inb[j] = false
							}
							mark(u, v, b, &head, &tail)
							mark(v, u, b, &head, &tail)
							for j := 1; j <= n; j++ {
								if inb[base[j]] {
									base[j] = b
									if !inq[j] {
										q[tail] = j
										tail++
										inq[j] = true
									}
								}
							}
						} else if p[v] == 0 {
							p[v] = u
							if match[v] > 0 {
								q[tail] = match[v]
								tail++
								inq[match[v]] = true
							} else {
								curr := v
								for curr > 0 {
									nxt := p[curr]
									nxtMatch := match[nxt]
									match[curr] = nxt
									match[nxt] = curr
									curr = nxtMatch
								}
								found = true
								break
							}
						}
					}
				}
			}
		}
	}
	return match
}
`

type Test struct {
	input string
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("%v: %s", err, errBuf.String())
	}
	return out.String(), err
}

func genTests() []Test {
	rand.Seed(6)
	tests := make([]Test, 0, 20)
	for t := 0; t < 19; t++ {
		n := rand.Intn(5) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("%d ", rand.Intn(4)))
		}
		sb.WriteByte('\n')
		tests = append(tests, Test{sb.String()})
	}
	tests = append(tests, Test{"1\n0\n"})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	tmp, err := os.CreateTemp("", "refG_*.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create temp file: %v\n", err)
		os.Exit(1)
	}
	if _, err := tmp.WriteString(refSourceG); err != nil {
		tmp.Close()
		fmt.Fprintf(os.Stderr, "failed to write temp file: %v\n", err)
		os.Exit(1)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	ref := filepath.Join(os.TempDir(), "refG_1615.bin")
	cmd := exec.Command("go", "build", "-o", ref, tmp.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("build reference: %v: %s", err, string(out)))
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		want, err := runExe(ref, tc.input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
