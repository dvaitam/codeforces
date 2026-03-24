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

type testCase struct {
	name  string
	input string
}

const refSource = `package main

import (
	"bufio"
	"os"
	"sort"
)

type FastScanner struct {
	r        *os.File
	buf      []byte
	idx, n   int
}

func NewFastScanner() *FastScanner {
	return &FastScanner{
		r:   os.Stdin,
		buf: make([]byte, 1<<20),
	}
}

func (fs *FastScanner) refill() bool {
	n, _ := fs.r.Read(fs.buf)
	fs.idx = 0
	fs.n = n
	return n > 0
}

func (fs *FastScanner) NextInt() int {
	sign := 1
	val := 0
	for {
		if fs.idx >= fs.n {
			if !fs.refill() {
				return 0
			}
		}
		c := fs.buf[fs.idx]
		if (c >= '0' && c <= '9') || c == '-' {
			if c == '-' {
				sign = -1
				fs.idx++
			}
			break
		}
		fs.idx++
	}
	for {
		if fs.idx >= fs.n {
			if !fs.refill() {
				break
			}
		}
		c := fs.buf[fs.idx]
		if c < '0' || c > '9' {
			break
		}
		val = val*10 + int(c-'0')
		fs.idx++
	}
	return sign * val
}

type Edge struct {
	u, v, w, id int32
}

type ByWeight []Edge

func (a ByWeight) Len() int           { return len(a) }
func (a ByWeight) Less(i, j int) bool { return a[i].w < a[j].w }
func (a ByWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func find(parent []int32, x int32) int32 {
	for parent[int(x)] != x {
		parent[int(x)] = parent[int(parent[int(x)])]
		x = parent[int(x)]
	}
	return x
}

func union(parent, size []int32, a, b int32) bool {
	ra := find(parent, a)
	rb := find(parent, b)
	if ra == rb {
		return false
	}
	if size[int(ra)] < size[int(rb)] {
		ra, rb = rb, ra
	}
	parent[int(rb)] = ra
	size[int(ra)] += size[int(rb)]
	return true
}

func query(u, v int32, up, mx [][]int32, depth []int, kmax int) int32 {
	var res int32
	if depth[int(u)] < depth[int(v)] {
		u, v = v, u
	}
	diff := depth[int(u)] - depth[int(v)]
	k := 0
	for diff > 0 {
		if diff&1 == 1 {
			t := mx[k][int(u)]
			if t > res {
				res = t
			}
			u = up[k][int(u)]
		}
		diff >>= 1
		k++
	}
	if u == v {
		return res
	}
	for k = kmax - 1; k >= 0; k-- {
		pu := up[k][int(u)]
		pv := up[k][int(v)]
		if pu != pv {
			t := mx[k][int(u)]
			if t > res {
				res = t
			}
			t = mx[k][int(v)]
			if t > res {
				res = t
			}
			u = pu
			v = pv
		}
	}
	t := mx[0][int(u)]
	if t > res {
		res = t
	}
	t = mx[0][int(v)]
	if t > res {
		res = t
	}
	return res
}

func appendInt32(buf []byte, x int32) []byte {
	if x == 0 {
		return append(buf, '0')
	}
	var tmp [11]byte
	i := len(tmp)
	for x > 0 {
		i--
		tmp[i] = byte(x%10) + '0'
		x /= 10
	}
	return append(buf, tmp[i:]...)
}

func main() {
	fs := NewFastScanner()
	n := fs.NextInt()
	m := fs.NextInt()

	edges := make([]Edge, m)
	for i := 0; i < m; i++ {
		u := fs.NextInt()
		v := fs.NextInt()
		w := fs.NextInt()
		edges[i] = Edge{int32(u), int32(v), int32(w), int32(i)}
	}

	sort.Sort(ByWeight(edges))

	parent := make([]int32, n+1)
	size := make([]int32, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = int32(i)
		size[i] = 1
	}

	isMST := make([]bool, m)

	head := make([]int32, n+1)
	for i := 1; i <= n; i++ {
		head[i] = -1
	}
	to := make([]int32, 2*(n-1))
	wt := make([]int32, 2*(n-1))
	next := make([]int32, 2*(n-1))
	var ec int32

	addEdge := func(u, v, w int32) {
		to[int(ec)] = v
		wt[int(ec)] = w
		next[int(ec)] = head[int(u)]
		head[int(u)] = ec
		ec++
	}

	added := 0
	for i := 0; i < m && added < n-1; i++ {
		e := edges[i]
		if union(parent, size, e.u, e.v) {
			isMST[int(e.id)] = true
			addEdge(e.u, e.v, e.w)
			addEdge(e.v, e.u, e.w)
			added++
		}
	}

	kmax := 0
	for (1 << kmax) <= n {
		kmax++
	}

	up := make([][]int32, kmax)
	mx := make([][]int32, kmax)
	for i := 0; i < kmax; i++ {
		up[i] = make([]int32, n+1)
		mx[i] = make([]int32, n+1)
	}

	depth := make([]int, n+1)
	queue := make([]int32, n)
	queue[0] = 1
	qh, qt := 0, 1

	for qh < qt {
		u := queue[qh]
		qh++
		pu := up[0][int(u)]
		for ei := head[int(u)]; ei != -1; ei = next[int(ei)] {
			v := to[int(ei)]
			if v == pu {
				continue
			}
			up[0][int(v)] = u
			mx[0][int(v)] = wt[int(ei)]
			depth[int(v)] = depth[int(u)] + 1
			queue[qt] = v
			qt++
		}
	}

	for k := 1; k < kmax; k++ {
		upPrev := up[k-1]
		mxPrev := mx[k-1]
		upCur := up[k]
		mxCur := mx[k]
		for v := 1; v <= n; v++ {
			p := upPrev[v]
			upCur[v] = upPrev[int(p)]
			a := mxPrev[v]
			b := mxPrev[int(p)]
			if b > a {
				a = b
			}
			mxCur[v] = a
		}
	}

	ans := make([]int32, m)
	for i := 0; i < m; i++ {
		id := int(edges[i].id)
		if !isMST[id] {
			ans[id] = query(edges[i].u, edges[i].v, up, mx, depth, kmax)
		}
	}

	out := make([]byte, 0, (m-n+1)*12)
	for i := 0; i < m; i++ {
		if !isMST[i] {
			out = appendInt32(out, ans[i])
			out = append(out, '\n')
		}
	}

	w := bufio.NewWriterSize(os.Stdout, 1<<20)
	w.Write(out)
	w.Flush()
}
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()

	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expected := strings.TrimSpace(refOut)

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "solution runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		got := strings.TrimSpace(candOut)

		if expected != got {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected:\n%s\nbut got:\n%s\ninput:\n%s", idx+1, tc.name, expected, got, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests)\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	srcFile, err := os.CreateTemp("", "cf-1184E2-src-*.go")
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

	tmp, err := os.CreateTemp("", "cf-1184E2-ref-*")
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

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func generateTests() []testCase {
	var tests []testCase

	tests = append(tests, manualCase("simple_tree", 3, [][]int{
		{1, 2, 1},
		{2, 3, 2},
	}))
	tests = append(tests, manualCase("line_plus_extra", 4, [][]int{
		{1, 2, 1},
		{2, 3, 2},
		{3, 4, 3},
		{1, 3, 4},
		{2, 4, 5},
	}))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomCase(rng, "random_small", 6, 10))
	tests = append(tests, randomCase(rng, "random_medium", 200, 500))
	tests = append(tests, randomCase(rng, "random_dense", 500, 2000))
	tests = append(tests, randomCase(rng, "random_big", 2000, 4000))
	tests = append(tests, randomCase(rng, "random_huge", 10000, 20000))

	return tests
}

func manualCase(name string, n int, edges [][]int) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&b, "%d %d %d\n", e[0], e[1], e[2])
	}
	return testCase{name: name, input: b.String()}
}

func randomCase(rng *rand.Rand, name string, n, m int) testCase {
	if m < n-1 {
		m = n - 1
	}
	type pair struct{ u, v int }
	edges := make([][3]int, 0, m)
	used := make(map[pair]struct{})
	// build random tree first
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		w := len(edges) + 1
		edges = append(edges, [3]int{u, v, w})
		used[pair{min(u, v), max(u, v)}] = struct{}{}
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		p := pair{min(u, v), max(u, v)}
		if _, ok := used[p]; ok {
			continue
		}
		used[p] = struct{}{}
		w := len(edges) + 1
		edges = append(edges, [3]int{u, v, w})
	}
	// shuffle edges
	for i := len(edges) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		edges[i], edges[j] = edges[j], edges[i]
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&b, "%d %d %d\n", e[0], e[1], e[2])
	}
	return testCase{name: name, input: b.String()}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
