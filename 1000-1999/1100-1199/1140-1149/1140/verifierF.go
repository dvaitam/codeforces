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

const embeddedSolverF = `package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
)

type Edge struct{ u, v int }

type History struct {
	u, v int
	ans  int64
}

var (
	tree    [][]Edge
	parent  []int
	sz      []int
	L, R    []int64
	ans     int64
	history []History
	res     []int64
)

func find(i int) int {
	for parent[i] != i {
		i = parent[i]
	}
	return i
}

func union(u, v int) {
	rootU := find(u)
	rootV := find(v)
	if rootU == rootV {
		history = append(history, History{-1, -1, 0})
		return
	}
	if sz[rootU] > sz[rootV] {
		rootU, rootV = rootV, rootU
	}
	history = append(history, History{rootU, rootV, ans})

	ans -= L[rootU]*R[rootU] + L[rootV]*R[rootV]
	parent[rootU] = rootV
	sz[rootV] += sz[rootU]
	L[rootV] += L[rootU]
	R[rootV] += R[rootU]
	ans += L[rootV]*R[rootV]
}

func rollback(target int) {
	for len(history) > target {
		last := history[len(history)-1]
		history = history[:len(history)-1]
		u, v := last.u, last.v
		if u != -1 {
			L[v] -= L[u]
			R[v] -= R[u]
			sz[v] -= sz[u]
			parent[u] = u
			ans = last.ans
		}
	}
}

func add(node, tl, tr, l, r, u, v int) {
	if l > r {
		return
	}
	if l == tl && r == tr {
		tree[node] = append(tree[node], Edge{u, v})
		return
	}
	tm := (tl + tr) / 2
	add(2*node, tl, tm, l, min(r, tm), u, v)
	add(2*node+1, tm+1, tr, max(l, tm+1), r, u, v)
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

func dfs(node, tl, tr int) {
	historySize := len(history)
	for _, e := range tree[node] {
		union(e.u, e.v)
	}
	if tl == tr {
		res[tl] = ans
	} else {
		tm := (tl + tr) / 2
		dfs(2*node, tl, tm)
		dfs(2*node+1, tm+1, tr)
	}
	rollback(historySize)
}

func main() {
	buffer, _ := io.ReadAll(os.Stdin)
	cursor := 0

	nextInt := func() int {
		for cursor < len(buffer) && buffer[cursor] <= ' ' {
			cursor++
		}
		if cursor >= len(buffer) {
			return 0
		}
		res := 0
		for cursor < len(buffer) && buffer[cursor] > ' ' {
			res = res*10 + int(buffer[cursor]-'0')
			cursor++
		}
		return res
	}

	q := nextInt()
	if q == 0 {
		return
	}

	parent = make([]int, 600005)
	sz = make([]int, 600005)
	L = make([]int64, 600005)
	R = make([]int64, 600005)

	for i := 1; i <= 300000; i++ {
		parent[i] = i
		sz[i] = 1
		L[i] = 1
	}
	for i := 300001; i <= 600000; i++ {
		parent[i] = i
		sz[i] = 1
		R[i] = 1
	}

	tree = make([][]Edge, 4*q+1)
	res = make([]int64, q+1)
	active := make(map[int64]int)

	for i := 1; i <= q; i++ {
		x := nextInt()
		y := nextInt()
		key := (int64(x) << 32) | int64(y)

		if start, ok := active[key]; ok {
			add(1, 1, q, start, i-1, x, y+300000)
			delete(active, key)
		} else {
			active[key] = i
		}
	}

	for key, start := range active {
		x := int(key >> 32)
		y := int(key & 0xFFFFFFFF)
		add(1, 1, q, start, q, x, y+300000)
	}

	dfs(1, 1, q)

	out := bufio.NewWriter(os.Stdout)
	for i := 1; i <= q; i++ {
		out.WriteString(strconv.FormatInt(res[i], 10))
		out.WriteByte('\n')
	}
	out.Flush()
}
`

func buildEmbeddedOracle() (string, func(), error) {
	tmpSrc, err := os.CreateTemp("", "oracle1140F-*.go")
	if err != nil {
		return "", nil, err
	}
	if _, err := tmpSrc.WriteString(embeddedSolverF); err != nil {
		tmpSrc.Close()
		os.Remove(tmpSrc.Name())
		return "", nil, err
	}
	tmpSrc.Close()

	tmpBin, err := os.CreateTemp("", "oracle1140F-bin-*")
	if err != nil {
		os.Remove(tmpSrc.Name())
		return "", nil, err
	}
	tmpBin.Close()

	cmd := exec.Command("go", "build", "-o", tmpBin.Name(), tmpSrc.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmpSrc.Name())
		os.Remove(tmpBin.Name())
		return "", nil, fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	os.Remove(tmpSrc.Name())
	return tmpBin.Name(), func() { os.Remove(tmpBin.Name()) }, nil
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	q := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		x := rng.Intn(5) + 1
		y := rng.Intn(5) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	// Handle .go source files
	if strings.HasSuffix(bin, ".go") {
		tmp, err := os.CreateTemp("", "cand1140F-*")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), bin)
		if out, err := cmd.CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			fmt.Fprintf(os.Stderr, "build failed: %v\n%s", err, out)
			os.Exit(1)
		}
		bin = tmp.Name()
		defer os.Remove(bin)
	}

	ref, refCleanup, err := buildEmbeddedOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer refCleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		exp, err := runProgram(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected: %s\ngot: %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
