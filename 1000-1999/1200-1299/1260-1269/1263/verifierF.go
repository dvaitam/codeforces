package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const refSource = `package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const INF = int(1e9)

type Tree struct {
	head      []int
	to        []int
	next      []int
	edgeCount int
	depth     []int
	up        [][]int
}

func NewTree(n int) *Tree {
	t := &Tree{
		head:  make([]int, n+1),
		to:    make([]int, 1),
		next:  make([]int, 1),
		depth: make([]int, n+1),
		up:    make([][]int, n+1),
	}
	for i := range t.up {
		t.up[i] = make([]int, 12)
	}
	return t
}

func (t *Tree) AddEdge(u, v int) {
	t.edgeCount++
	t.to = append(t.to, v)
	t.next = append(t.next, t.head[u])
	t.head[u] = t.edgeCount
}

func (t *Tree) DFS(u, p, d int) {
	t.depth[u] = d
	t.up[u][0] = p
	for i := 1; i < 12; i++ {
		t.up[u][i] = t.up[t.up[u][i-1]][i-1]
	}
	for e := t.head[u]; e != 0; e = t.next[e] {
		v := t.to[e]
		if v != p {
			t.DFS(v, u, d+1)
		}
	}
}

func (t *Tree) LCA(u, v int) int {
	if t.depth[u] < t.depth[v] {
		u, v = v, u
	}
	diff := t.depth[u] - t.depth[v]
	for i := 0; i < 12; i++ {
		if (diff & (1 << i)) != 0 {
			u = t.up[u][i]
		}
	}
	if u == v {
		return u
	}
	for i := 11; i >= 0; i-- {
		if t.up[u][i] != t.up[v][i] {
			u = t.up[u][i]
			v = t.up[v][i]
		}
	}
	return t.up[u][0]
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 1024*1024)
	scanner.Split(bufio.ScanWords)

	scanInt := func() int {
		scanner.Scan()
		res, _ := strconv.Atoi(scanner.Text())
		return res
	}

	if !scanner.Scan() {
		return
	}
	n, _ := strconv.Atoi(scanner.Text())

	a := scanInt()
	tree1 := NewTree(a)
	for i := 1; i <= a-1; i++ {
		p := scanInt()
		tree1.AddEdge(p, i+1)
	}
	x := make([]int, n+1)
	for i := 1; i <= n; i++ {
		x[i] = scanInt()
	}

	b := scanInt()
	tree2 := NewTree(b)
	for i := 1; i <= b-1; i++ {
		q := scanInt()
		tree2.AddEdge(q, i+1)
	}
	y := make([]int, n+1)
	for i := 1; i <= n; i++ {
		y[i] = scanInt()
	}

	tree1.DFS(1, 1, 0)
	tree2.DFS(1, 1, 0)

	dp1 := make([]int, n+1)
	dp2 := make([]int, n+1)
	for i := 0; i <= n; i++ {
		dp1[i] = INF
		dp2[i] = INF
	}

	dp1[0] = tree1.depth[x[1]]
	dp2[0] = tree2.depth[y[1]]

	for k := 1; k < n; k++ {
		K := k + 1
		next_dp1 := make([]int, n+1)
		next_dp2 := make([]int, n+1)
		for i := 0; i <= n; i++ {
			next_dp1[i] = INF
			next_dp2[i] = INF
		}

		for j := 0; j < k; j++ {
			if dp1[j] != INF {
				cost1 := dp1[j] + tree1.depth[x[K]] - tree1.depth[tree1.LCA(x[k], x[K])]
				if cost1 < next_dp1[j] {
					next_dp1[j] = cost1
				}

				sub := 0
				if j > 0 {
					sub = tree2.depth[tree2.LCA(y[j], y[K])]
				}
				cost2 := dp1[j] + tree2.depth[y[K]] - sub
				if cost2 < next_dp2[k] {
					next_dp2[k] = cost2
				}
			}
		}

		for i := 0; i < k; i++ {
			if dp2[i] != INF {
				sub := 0
				if i > 0 {
					sub = tree1.depth[tree1.LCA(x[i], x[K])]
				}
				cost1 := dp2[i] + tree1.depth[x[K]] - sub
				if cost1 < next_dp1[k] {
					next_dp1[k] = cost1
				}

				cost2 := dp2[i] + tree2.depth[y[K]] - tree2.depth[tree2.LCA(y[k], y[K])]
				if cost2 < next_dp2[i] {
					next_dp2[i] = cost2
				}
			}
		}
		dp1 = next_dp1
		dp2 = next_dp2
	}

	minCost := INF
	for i := 0; i < n; i++ {
		if dp1[i] < minCost {
			minCost = dp1[i]
		}
		if dp2[i] < minCost {
			minCost = dp2[i]
		}
	}

	ans := (a - 1) + (b - 1) - minCost
	fmt.Println(ans)
}
`

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func buildReferenceBinary() (string, error) {
	srcFile, err := os.CreateTemp("", "cf-1263F-src-*.go")
	if err != nil {
		return "", err
	}
	if _, err := srcFile.WriteString(refSource); err != nil {
		srcFile.Close()
		os.Remove(srcFile.Name())
		return "", err
	}
	srcFile.Close()
	defer os.Remove(srcFile.Name())

	tmp, err := os.CreateTemp("", "cf-1263F-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), srcFile.Name())
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func generateTree(nLeaves int, rng *rand.Rand) (int, []int) {
	m := rng.Intn(2)
	a := 1 + m + nLeaves
	par := make([]int, a+1)
	for i := 2; i <= 1+m; i++ {
		par[i] = 1
	}
	for i := 0; i < nLeaves; i++ {
		id := m + 2 + i
		var parent int
		if i < 1+m {
			parent = 1 + i
			if parent > 1+m {
				parent = 1
			}
		} else {
			parent = rng.Intn(1+m) + 1
		}
		par[id] = parent
	}
	return a, par
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(3) + 1
	a, parA := generateTree(n, rng)
	x := make([]int, n+1)
	for i := 1; i <= n; i++ {
		x[i] = len(parA) - n - 1 + i
	}
	b, parB := generateTree(n, rng)
	y := make([]int, n+1)
	for i := 1; i <= n; i++ {
		y[i] = len(parB) - n - 1 + i
	}
	rng.Shuffle(n, func(i, j int) { y[i+1], y[j+1] = y[j+1], y[i+1] })
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", n))
	input.WriteString(fmt.Sprintf("%d\n", a))
	for i := 2; i <= a; i++ {
		input.WriteString(fmt.Sprintf("%d ", parA[i]))
	}
	input.WriteByte('\n')
	for i := 1; i <= n; i++ {
		input.WriteString(fmt.Sprintf("%d ", x[i]))
	}
	input.WriteByte('\n')
	input.WriteString(fmt.Sprintf("%d\n", b))
	for i := 2; i <= b; i++ {
		input.WriteString(fmt.Sprintf("%d ", parB[i]))
	}
	input.WriteByte('\n')
	for i := 1; i <= n; i++ {
		input.WriteString(fmt.Sprintf("%d ", y[i]))
	}
	input.WriteByte('\n')
	return input.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refBin, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		exp, err := run(refBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
