package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Embedded solver logic from the accepted solution for 48G.

type solverEdge struct {
	to int
	w  int64
}

func solveCase(input string) string {
	fs := &fastScanner{data: []byte(input), n: len(input)}
	n := fs.nextInt()

	g := make([][]solverEdge, n+1)
	for i := 0; i < n; i++ {
		a := fs.nextInt()
		b := fs.nextInt()
		t := int64(fs.nextInt())
		g[a] = append(g[a], solverEdge{b, t})
		g[b] = append(g[b], solverEdge{a, t})
	}

	deg := make([]int, n+1)
	for i := 1; i <= n; i++ {
		deg[i] = len(g[i])
	}

	removed := make([]bool, n+1)
	q := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if deg[i] == 1 {
			q = append(q, i)
		}
	}

	for head := 0; head < len(q); head++ {
		v := q[head]
		removed[v] = true
		for _, e := range g[v] {
			to := e.to
			if !removed[to] {
				deg[to]--
				if deg[to] == 1 {
					q = append(q, to)
				}
			}
		}
	}

	isCycle := make([]bool, n+1)
	k := 0
	start := 0
	for i := 1; i <= n; i++ {
		if !removed[i] {
			isCycle[i] = true
			k++
			if start == 0 {
				start = i
			}
		}
	}

	cycleNodes := make([]int, 0, k)
	cycleEdges := make([]int64, 0, k)
	prev := 0
	curr := start
	for len(cycleNodes) < k {
		cycleNodes = append(cycleNodes, curr)
		next := 0
		var w int64
		for _, e := range g[curr] {
			if isCycle[e.to] && e.to != prev {
				next = e.to
				w = e.w
				break
			}
		}
		cycleEdges = append(cycleEdges, w)
		prev, curr = curr, next
	}

	compID := make([]int, n+1)
	parent := make([]int, n+1)
	distRoot := make([]int64, n+1)
	subSize := make([]int64, n+1)
	down := make([]int64, n+1)
	treeAns := make([]int64, n+1)

	sz := make([]int64, k)
	rootSum := make([]int64, k)
	var totalRootSum int64

	for i, root := range cycleNodes {
		parent[root] = 0
		distRoot[root] = 0
		compID[root] = i

		stack := make([]int, 0)
		order := make([]int, 0)
		stack = append(stack, root)

		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			order = append(order, v)

			for _, e := range g[v] {
				to := e.to
				if to == parent[v] || isCycle[to] {
					continue
				}
				parent[to] = v
				distRoot[to] = distRoot[v] + e.w
				compID[to] = i
				stack = append(stack, to)
			}
		}

		for idx := len(order) - 1; idx >= 0; idx-- {
			v := order[idx]
			subSize[v] = 1
			down[v] = 0
			for _, e := range g[v] {
				to := e.to
				if parent[to] == v {
					subSize[v] += subSize[to]
					down[v] += down[to] + subSize[to]*e.w
				}
			}
		}

		sz[i] = subSize[root]
		rootSum[i] = down[root]
		totalRootSum += rootSum[i]
		treeAns[root] = down[root]

		for _, v := range order {
			for _, e := range g[v] {
				to := e.to
				if parent[to] == v {
					treeAns[to] = treeAns[v] + (sz[i]-2*subSize[to])*e.w
				}
			}
		}
	}

	pos := make([]int64, k+1)
	for i := 0; i < k; i++ {
		pos[i+1] = pos[i] + cycleEdges[i]
	}
	L := pos[k]

	x := make([]int64, 2*k)
	wts := make([]int64, 2*k)
	for i := 0; i < k; i++ {
		x[i] = pos[i]
		x[i+k] = pos[i] + L
		wts[i] = sz[i]
		wts[i+k] = sz[i]
	}

	prefW := make([]int64, 2*k+1)
	prefWP := make([]int64, 2*k+1)
	for i := 0; i < 2*k; i++ {
		prefW[i+1] = prefW[i] + wts[i]
		prefWP[i+1] = prefWP[i] + wts[i]*x[i]
	}

	G := make([]int64, k)
	r := 0
	for i := 0; i < k; i++ {
		if r < i {
			r = i
		}
		for r+1 < i+k && 2*(x[r+1]-x[i]) <= L {
			r++
		}
		sumW1 := prefW[r+1] - prefW[i+1]
		sumWP1 := prefWP[r+1] - prefWP[i+1]
		cw := sumWP1 - x[i]*sumW1

		sumW2 := prefW[i+k] - prefW[r+1]
		sumWP2 := prefWP[i+k] - prefWP[r+1]
		ccw := (L+x[i])*sumW2 - sumWP2

		G[i] = cw + ccw
	}

	ans := make([]int64, n+1)
	N := int64(n)
	for v := 1; v <= n; v++ {
		i := compID[v]
		ans[v] = treeAns[v] + distRoot[v]*(N-sz[i]) + totalRootSum - rootSum[i] + G[i]
	}

	out := &bytes.Buffer{}
	bw := bufio.NewWriterSize(out, 1<<20)
	for i := 1; i <= n; i++ {
		if i > 1 {
			bw.WriteByte(' ')
		}
		bw.WriteString(strconv.FormatInt(ans[i], 10))
	}
	bw.WriteByte('\n')
	bw.Flush()
	return strings.TrimSpace(out.String())
}

type fastScanner struct {
	data []byte
	idx  int
	n    int
}

func (fs *fastScanner) nextInt() int {
	for fs.idx < fs.n && (fs.data[fs.idx] < '0' || fs.data[fs.idx] > '9') {
		fs.idx++
	}
	val := 0
	for fs.idx < fs.n && fs.data[fs.idx] >= '0' && fs.data[fs.idx] <= '9' {
		val = val*10 + int(fs.data[fs.idx]-'0')
		fs.idx++
	}
	return val
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		all, _ := io.ReadAll(&stderr)
		return "", fmt.Errorf("%v\n%s", err, string(all))
	}
	return strings.TrimSpace(stdout.String()), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 3
	k := rng.Intn(n-2) + 3
	edges := make([][3]int, 0, n)
	for i := 1; i < k; i++ {
		edges = append(edges, [3]int{i, i + 1, rng.Intn(5) + 1})
	}
	edges = append(edges, [3]int{k, 1, rng.Intn(5) + 1})
	for v := k + 1; v <= n; v++ {
		to := rng.Intn(k) + 1
		edges = append(edges, [3]int{v, to, rng.Intn(5) + 1})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		expect := solveCase(in)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, in, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
