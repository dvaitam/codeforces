package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// ---------- embedded solver from accepted solution ----------

type Base struct {
	d int
	g int64
}

type FlowEdge struct {
	to  int
	rev int
	cap int64
}

type DinicFlow struct {
	g     [][]FlowEdge
	level []int
	it    []int
}

func NewDinicFlow(n int) *DinicFlow {
	return &DinicFlow{
		g:     make([][]FlowEdge, n),
		level: make([]int, n),
		it:    make([]int, n),
	}
}

func (d *DinicFlow) AddEdge(fr, to int, cap int64) {
	fwd := FlowEdge{to: to, rev: len(d.g[to]), cap: cap}
	rev := FlowEdge{to: fr, rev: len(d.g[fr]), cap: 0}
	d.g[fr] = append(d.g[fr], fwd)
	d.g[to] = append(d.g[to], rev)
}

func (d *DinicFlow) bfs(s, t int) bool {
	for i := range d.level {
		d.level[i] = -1
	}
	q := make([]int, 0, len(d.g))
	d.level[s] = 0
	q = append(q, s)
	for h := 0; h < len(q); h++ {
		v := q[h]
		for _, e := range d.g[v] {
			if e.cap > 0 && d.level[e.to] < 0 {
				d.level[e.to] = d.level[v] + 1
				q = append(q, e.to)
			}
		}
	}
	return d.level[t] >= 0
}

func flowMin64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func (d *DinicFlow) dfs(v, t int, f int64) int64 {
	if v == t {
		return f
	}
	for ; d.it[v] < len(d.g[v]); d.it[v]++ {
		i := d.it[v]
		e := &d.g[v][i]
		if e.cap > 0 && d.level[e.to] == d.level[v]+1 {
			ret := d.dfs(e.to, t, flowMin64(f, e.cap))
			if ret > 0 {
				e.cap -= ret
				re := &d.g[e.to][e.rev]
				re.cap += ret
				return ret
			}
		}
	}
	return 0
}

func (d *DinicFlow) MaxFlow(s, t int) int64 {
	var flow int64
	const INF int64 = 1 << 60
	for d.bfs(s, t) {
		for i := range d.it {
			d.it[i] = 0
		}
		for {
			f := d.dfs(s, t, INF)
			if f == 0 {
				break
			}
			flow += f
		}
	}
	return flow
}

func upperBound(a []int, x int) int {
	l, r := 0, len(a)
	for l < r {
		m := (l + r) >> 1
		if a[m] <= x {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}

func solveB3(input string) string {
	data := []byte(input)
	pos := 0
	nextInt := func() int {
		for pos < len(data) && (data[pos] < '0' || data[pos] > '9') && data[pos] != '-' {
			pos++
		}
		neg := false
		if pos < len(data) && data[pos] == '-' {
			neg = true
			pos++
		}
		v := 0
		for pos < len(data) && data[pos] >= '0' && data[pos] <= '9' {
			v = v*10 + int(data[pos]-'0')
			pos++
		}
		if neg {
			return -v
		}
		return v
	}

	n := nextInt()
	m := nextInt()

	const INFINT = 1000001000
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int, n)
		for j := 0; j < n; j++ {
			dist[i][j] = INFINT
		}
		dist[i][i] = 0
	}

	for i := 0; i < m; i++ {
		u := nextInt() - 1
		v := nextInt() - 1
		dist[u][v] = 1
		dist[v][u] = 1
	}

	for k := 0; k < n; k++ {
		rowk := dist[k]
		for i := 0; i < n; i++ {
			dik := dist[i][k]
			if dik == INFINT {
				continue
			}
			rowi := dist[i]
			for j := 0; j < n; j++ {
				nd := dik + rowk[j]
				if nd < rowi[j] {
					rowi[j] = nd
				}
			}
		}
	}

	s := nextInt()
	b := nextInt()
	kdep := nextInt()

	shipLoc := make([]int, s)
	shipAtk := make([]int, s)
	shipFuel := make([]int, s)
	shipPrice := make([]int64, s)

	for i := 0; i < s; i++ {
		shipLoc[i] = nextInt() - 1
		shipAtk[i] = nextInt()
		shipFuel[i] = nextInt()
		shipPrice[i] = int64(nextInt())
	}

	baseGroups := make([][]Base, n)
	for i := 0; i < b; i++ {
		x := nextInt() - 1
		d := nextInt()
		g := int64(nextInt())
		baseGroups[x] = append(baseGroups[x], Base{d: d, g: g})
	}

	defsByNode := make([][]int, n)
	prefByNode := make([][]int64, n)

	for i := 0; i < n; i++ {
		g := baseGroups[i]
		if len(g) == 0 {
			continue
		}
		sort.Slice(g, func(a, b int) bool {
			return g[a].d < g[b].d
		})
		defs := make([]int, len(g))
		pref := make([]int64, len(g))
		var mx int64 = -1
		for j := 0; j < len(g); j++ {
			defs[j] = g[j].d
			if g[j].g > mx {
				mx = g[j].g
			}
			pref[j] = mx
		}
		defsByNode[i] = defs
		prefByNode[i] = pref
	}

	values := make([]int64, s)
	available := make([]bool, s)

	for i := 0; i < s; i++ {
		best := int64(-1)
		x := shipLoc[i]
		atk := shipAtk[i]
		fuel := shipFuel[i]
		row := dist[x]
		for y := 0; y < n; y++ {
			if row[y] > fuel {
				continue
			}
			defs := defsByNode[y]
			if len(defs) == 0 {
				continue
			}
			p := upperBound(defs, atk) - 1
			if p >= 0 {
				gg := prefByNode[y][p]
				if gg > best {
					best = gg
				}
			}
		}
		if best >= 0 {
			available[i] = true
			values[i] = best - shipPrice[i]
		}
	}

	depA := make([]int, kdep)
	depB := make([]int, kdep)
	involved := make([]bool, s)

	for i := 0; i < kdep; i++ {
		a := nextInt() - 1
		bb := nextInt() - 1
		depA[i] = a
		depB[i] = bb
		involved[a] = true
		involved[bb] = true
	}

	idMap := make([]int, s)
	for i := 0; i < s; i++ {
		idMap[i] = -1
	}

	involvedList := make([]int, 0)
	var answer int64

	for i := 0; i < s; i++ {
		if involved[i] {
			idMap[i] = len(involvedList)
			involvedList = append(involvedList, i)
		} else if available[i] && values[i] > 0 {
			answer += values[i]
		}
	}

	cnt := len(involvedList)
	if cnt > 0 {
		src := cnt
		sink := cnt + 1
		din := NewDinicFlow(cnt + 2)
		const INF int64 = 1 << 60
		var totalPos int64

		for id, ship := range involvedList {
			if available[ship] {
				w := values[ship]
				if w > 0 {
					din.AddEdge(src, id, w)
					totalPos += w
				} else if w < 0 {
					din.AddEdge(id, sink, -w)
				}
			}
			if !available[ship] {
				din.AddEdge(id, sink, INF)
			}
		}

		for i := 0; i < kdep; i++ {
			u := idMap[depA[i]]
			v := idMap[depB[i]]
			din.AddEdge(u, v, INF)
		}

		answer += totalPos - din.MaxFlow(src, sink)
	}

	return strconv.FormatInt(answer, 10)
}

// ---------- verifier logic ----------

func run(bin, input string) (string, error) {
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

func genGraph(n, m int, r *rand.Rand) [][2]int {
	edges := make([][2]int, 0, m)
	for len(edges) < m {
		u := r.Intn(n) + 1
		v := r.Intn(n) + 1
		if u == v {
			continue
		}
		edges = append(edges, [2]int{u, v})
	}
	return edges
}

func genShips(s int, n int, r *rand.Rand) []string {
	ships := make([]string, s)
	for i := 0; i < s; i++ {
		x := r.Intn(n) + 1
		a := r.Intn(20)
		f := r.Intn(10)
		p := r.Intn(20)
		ships[i] = fmt.Sprintf("%d %d %d %d", x, a, f, p)
	}
	return ships
}

func genBases(b int, n int, r *rand.Rand) []string {
	bases := make([]string, b)
	for i := 0; i < b; i++ {
		x := r.Intn(n) + 1
		d := r.Intn(20)
		g := r.Intn(30)
		bases[i] = fmt.Sprintf("%d %d %d", x, d, g)
	}
	return bases
}

func genDeps(k, s int, r *rand.Rand) [][2]int {
	deps := make([][2]int, 0, k)
	for len(deps) < k {
		u := r.Intn(s) + 1
		v := r.Intn(s) + 1
		deps = append(deps, [2]int{u, v})
	}
	return deps
}

func genCase(r *rand.Rand) string {
	n := r.Intn(8) + 1
	m := r.Intn(10)
	edges := genGraph(n, m, r)
	s := r.Intn(6) + 1
	b := r.Intn(6) + 1
	k := r.Intn(5)
	ships := genShips(s, n, r)
	bases := genBases(b, n, r)
	deps := genDeps(k, s, r)

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	fmt.Fprintf(&sb, "%d %d %d\n", s, b, k)
	for _, ship := range ships {
		sb.WriteString(ship)
		sb.WriteByte('\n')
	}
	for _, base := range bases {
		sb.WriteString(base)
		sb.WriteByte('\n')
	}
	for _, dep := range deps {
		fmt.Fprintf(&sb, "%d %d\n", dep[0], dep[1])
	}
	return sb.String()
}

func parseOutput(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

const testCount = 100

func main() {
	// suppress unused import
	_ = io.Discard

	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB3.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	r := rand.New(rand.NewSource(1))
	for t := 0; t < testCount; t++ {
		input := genCase(r)
		expectStr := solveB3(input)
		gotStr, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		expectVal, err := parseOutput(expectStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "solver output invalid on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotVal, err := parseOutput(gotStr)
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", t+1, input, err)
			os.Exit(1)
		}
		if expectVal != gotVal {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected: %d\ngot: %d\n", t+1, input, expectVal, gotVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
