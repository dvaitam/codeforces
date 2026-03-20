package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSourceE = `package main

import (
	"bufio"
	"fmt"
	"os"
)

type SegmentTree struct {
	n    int
	tree []int
	lazy []int
}

func NewSegmentTree(n int) *SegmentTree {
	return &SegmentTree{
		n:    n,
		tree: make([]int, 4*n),
		lazy: make([]int, 4*n),
	}
}

func (st *SegmentTree) push(node int) {
	if st.lazy[node] != 0 {
		st.tree[2*node] += st.lazy[node]
		st.lazy[2*node] += st.lazy[node]
		st.tree[2*node+1] += st.lazy[node]
		st.lazy[2*node+1] += st.lazy[node]
		st.lazy[node] = 0
	}
}

func (st *SegmentTree) update(node, l, r, ql, qr, val int) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.tree[node] += val
		st.lazy[node] += val
		return
	}
	st.push(node)
	mid := (l + r) / 2
	st.update(2*node, l, mid, ql, qr, val)
	st.update(2*node+1, mid+1, r, ql, qr, val)
	if st.tree[2*node] > st.tree[2*node+1] {
		st.tree[node] = st.tree[2*node]
	} else {
		st.tree[node] = st.tree[2*node+1]
	}
}

func (st *SegmentTree) query() int {
	node := 1
	l, r := 0, st.n-1
	for l != r {
		st.push(node)
		mid := (l + r) / 2
		if st.tree[2*node] >= st.tree[2*node+1] {
			node = 2 * node
			r = mid
		} else {
			node = 2*node + 1
			l = mid + 1
		}
	}
	return l
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	fmt.Fscan(reader, &n, &q)

	length := 2*n + 1
	p := make([]int, length+1)
	pos := make([]int, length+1)

	st := NewSegmentTree(length)

	addInterval := func(u, v, val int) {
		if u <= v {
			st.update(1, 0, length-1, u, v, val)
		} else {
			st.update(1, 0, length-1, u, length-1, val)
			st.update(1, 0, length-1, 0, v, val)
		}
	}

	updateVal := func(i, val int) {
		if p[i] > i {
			addInterval(i, p[i]-1, val)
		} else if p[i] < i {
			addInterval(i, length-1, val)
			addInterval(0, p[i]-1, val)
		}
	}

	for i := 1; i <= length; i++ {
		fmt.Fscan(reader, &p[i])
		pos[p[i]] = i
	}

	for i := 1; i <= length; i++ {
		updateVal(i, 1)
	}

	for i := 0; i < q; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)

		updateVal(u, -1)
		updateVal(v, -1)

		p[u], p[v] = p[v], p[u]
		pos[p[u]] = u
		pos[p[v]] = v

		updateVal(u, 1)
		updateVal(v, 1)

		maxVal := st.tree[1]
		if maxVal <= n {
			fmt.Fprintln(writer, st.query())
		} else {
			fmt.Fprintln(writer, -1)
		}
	}
}
`

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func generateTests() []string {
	r := rand.New(rand.NewSource(6))
	tests := make([]string, 0, 100)
	for len(tests) < 100 {
		n := r.Intn(3) + 1 // 1..3
		q := r.Intn(3) + 1 // 1..3
		m := 2*n + 1
		perm := r.Perm(m)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
		for i := 0; i < m; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(perm[i] + 1))
		}
		sb.WriteByte('\n')
		for i := 0; i < q; i++ {
			u := r.Intn(m) + 1
			v := r.Intn(m) + 1
			sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
		}
		tests = append(tests, sb.String())
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tmp, err := os.CreateTemp("", "refE_*.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create temp file: %v\n", err)
		os.Exit(1)
	}
	if _, err := tmp.WriteString(refSourceE); err != nil {
		tmp.Close()
		fmt.Fprintf(os.Stderr, "failed to write temp file: %v\n", err)
		os.Exit(1)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	ref := filepath.Join(os.TempDir(), "refE_1685.bin")
	cmd := exec.Command("go", "build", "-o", ref, tmp.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n%s", err, out)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for i, input := range tests {
		candOut, cErr := runBinary(candidate, input)
		refOut, rErr := runBinary(ref, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "case %d: candidate error: %v\n", i+1, cErr)
			os.Exit(1)
		}
		if rErr != nil {
			fmt.Fprintf(os.Stderr, "case %d: reference error: %v\n", i+1, rErr)
			os.Exit(1)
		}
		if strings.TrimSpace(candOut) != strings.TrimSpace(refOut) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%sactual:%s", i+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
