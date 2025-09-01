package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Case struct{ input string }

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(5) + 2 // 2..6
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Intn(6))
		}
		sb.WriteByte('\n')
		cases[i] = Case{sb.String()}
	}
	return cases
}

// Dinic max flow implementation
type edge struct{ to, rev, cap int }

type dinic struct {
	g         [][]edge
	level, it []int
}

func newDinic(n int) *dinic {
	return &dinic{g: make([][]edge, n), level: make([]int, n), it: make([]int, n)}
}

func (d *dinic) addEdge(u, v, c int) {
	du := edge{to: v, rev: len(d.g[v]), cap: c}
	dv := edge{to: u, rev: len(d.g[u]), cap: 0}
	d.g[u] = append(d.g[u], du)
	d.g[v] = append(d.g[v], dv)
}

func (d *dinic) bfs(s, t int) bool {
	for i := range d.level {
		d.level[i] = -1
	}
	q := make([]int, 0, len(d.g))
	d.level[s] = 0
	q = append(q, s)
	for head := 0; head < len(q); head++ {
		u := q[head]
		for _, e := range d.g[u] {
			if e.cap > 0 && d.level[e.to] < 0 {
				d.level[e.to] = d.level[u] + 1
				q = append(q, e.to)
			}
		}
	}
	return d.level[t] >= 0
}

func (d *dinic) dfs(u, t, f int) int {
	if u == t {
		return f
	}
	for ; d.it[u] < len(d.g[u]); d.it[u]++ {
		i := d.it[u]
		e := d.g[u][i]
		if e.cap > 0 && d.level[u] < d.level[e.to] {
			ret := d.dfs(e.to, t, min(f, e.cap))
			if ret > 0 {
				// update
				d.g[u][i].cap -= ret
				v := e.to
				r := d.g[u][i].rev
				d.g[v][r].cap += ret
				return ret
			}
		}
	}
	return 0
}

func (d *dinic) maxflow(s, t int) int {
	flow := 0
	for d.bfs(s, t) {
		for i := range d.it {
			d.it[i] = 0
		}
		for {
			f := d.dfs(s, t, 1<<30)
			if f == 0 {
				break
			}
			flow += f
		}
	}
	return flow
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func verifyOne(input, output string) error {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return fmt.Errorf("bad n: %v", err)
	}
	r := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(in, &r[i]); err != nil {
			return fmt.Errorf("bad r[%d]", i)
		}
	}

	scan := bufio.NewScanner(strings.NewReader(output))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return fmt.Errorf("missing R")
	}
	var R int
	if _, err := fmt.Sscan(scan.Text(), &R); err != nil {
		return fmt.Errorf("bad R")
	}
	if !scan.Scan() {
		return fmt.Errorf("missing K")
	}
	var K int
	if _, err := fmt.Sscan(scan.Text(), &K); err != nil {
		return fmt.Errorf("bad K")
	}
	if K < 0 {
		return fmt.Errorf("negative K")
	}

	needs := make([]int, n)
	sumNeeds := 0
	for i := 0; i < n; i++ {
		needs[i] = r[i] - R
		if needs[i] < 0 {
			return fmt.Errorf("R too large: r[%d]=%d < R=%d", i, r[i], R)
		}
		sumNeeds += needs[i]
	}
	// read steps
	steps := make([]string, K)
	caps := make([]int, K)
	for step := 0; step < K; step++ {
		if !scan.Scan() {
			return fmt.Errorf("step %d: missing bitstring", step+1)
		}
		line := scan.Text()
		if len(line) != n {
			return fmt.Errorf("step %d: bitstring length %d != n=%d", step+1, len(line), n)
		}
		ones := 0
		for i := 0; i < n; i++ {
			c := line[i]
			if c != '0' && c != '1' {
				return fmt.Errorf("step %d: invalid char", step+1)
			}
			if c == '1' {
				ones++
			}
		}
		if ones != 2 && ones != 3 {
			return fmt.Errorf("step %d: must select exactly 2 or 3 indices, got %d", step+1, ones)
		}
		steps[step] = line
		caps[step] = ones
	}
	if scan.Scan() {
		return fmt.Errorf("extra output: %s", scan.Text())
	}

	// quick necessary condition: each index must appear at least needs[i] times
	for i := 0; i < n; i++ {
		occ := 0
		for step := 0; step < K; step++ {
			if steps[step][i] == '1' {
				occ++
			}
		}
		if occ < needs[i] {
			return fmt.Errorf("index %d appears %d times < need %d", i+1, occ, needs[i])
		}
	}
	// Build flow: source -> indices (needs), indices -> steps (1 if '1'), steps -> sink (cap = ones)
	sz := n + K + 2
	s := n + K
	t := s + 1
	d := newDinic(sz)
	for i := 0; i < n; i++ {
		if needs[i] > 0 {
			d.addEdge(s, i, needs[i])
		}
	}
	for j := 0; j < K; j++ {
		for i := 0; i < n; i++ {
			if steps[j][i] == '1' {
				d.addEdge(i, n+j, 1)
			}
		}
		d.addEdge(n+j, t, caps[j])
	}
	flow := d.maxflow(s, t)
	if flow != sumNeeds {
		return fmt.Errorf("cannot assign decrements: need %d got %d", sumNeeds, flow)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, c := range cases {
		out, err := runBinary(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
		if err := verifyOne(c.input, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, c.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
