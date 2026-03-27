package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ── embedded reference solver (determines YES/NO only) ──────────────────────

type fEdge struct {
	to, rev int
	cap     int64
}

type fGraph struct {
	g     [][]fEdge
	level []int
	iter  []int
}

func newFGraph(n int) *fGraph {
	return &fGraph{g: make([][]fEdge, n), level: make([]int, n), iter: make([]int, n)}
}

func (gr *fGraph) addEdge(from, to int, cap int64) {
	gr.g[from] = append(gr.g[from], fEdge{to: to, rev: len(gr.g[to]), cap: cap})
	gr.g[to] = append(gr.g[to], fEdge{to: from, rev: len(gr.g[from]) - 1, cap: 0})
}

func (gr *fGraph) bfs(s int) {
	for i := range gr.level {
		gr.level[i] = -1
	}
	gr.level[s] = 0
	queue := []int{s}
	for qi := 0; qi < len(queue); qi++ {
		v := queue[qi]
		for _, e := range gr.g[v] {
			if e.cap > 0 && gr.level[e.to] < 0 {
				gr.level[e.to] = gr.level[v] + 1
				queue = append(queue, e.to)
			}
		}
	}
}

func (gr *fGraph) dfs(v, t int, f int64) int64 {
	if v == t {
		return f
	}
	for i := gr.iter[v]; i < len(gr.g[v]); i++ {
		gr.iter[v] = i
		e := &gr.g[v][i]
		if e.cap > 0 && gr.level[v] < gr.level[e.to] {
			mn := f
			if e.cap < mn {
				mn = e.cap
			}
			d := gr.dfs(e.to, t, mn)
			if d > 0 {
				e.cap -= d
				gr.g[e.to][e.rev].cap += d
				return d
			}
		}
	}
	return 0
}

func (gr *fGraph) maxFlow(s, t int) int64 {
	flow := int64(0)
	const inf = int64(1e18)
	for {
		gr.bfs(s)
		if gr.level[t] < 0 {
			break
		}
		for i := range gr.iter {
			gr.iter[i] = 0
		}
		for {
			f := gr.dfs(s, t, inf)
			if f == 0 {
				break
			}
			flow += f
		}
	}
	return flow
}

// refSolveYesNo returns true if the answer is YES for the given input.
func refSolveYesNo(n, m int, s, a []int, edges [][2]int) bool {
	totalNodes := n + m + 3
	S := n + m
	T := n + m + 1
	T2 := n + m + 2
	gr := newFGraph(totalNodes)

	cur := make([]int64, n)
	for i := 0; i < m; i++ {
		u := edges[i][0]
		v := edges[i][1]
		gr.addEdge(S, i, 1)
		gr.addEdge(i, m+u, 1)
		gr.addEdge(i, m+v, 1)
		cur[u]--
		cur[v]--
	}

	tmp := int64(m)
	for i := 0; i < n; i++ {
		if s[i] == 0 {
			gr.addEdge(m+i, T2, int64(1e18))
		} else if int64(a[i]) >= cur[i] {
			diff := int64(a[i]) - cur[i]
			if diff%2 != 0 {
				return false
			}
			cap := diff / 2
			gr.addEdge(m+i, T, cap)
			tmp -= cap
		} else {
			return false
		}
	}
	if tmp < 0 {
		return false
	}
	gr.addEdge(T2, T, tmp)
	res := gr.maxFlow(S, T)
	return res == int64(m)
}

// ── verifier harness ───────────────────────────────────────────────────────

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 2 // 2..6
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(5) + 1
	if m > maxEdges {
		m = maxEdges
	}
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	s := make([]int, n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
			s[i] = 0
		} else {
			sb.WriteByte('1')
			s[i] = 1
		}
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		if s[i] == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteString(fmt.Sprintf("%d", rng.Intn(2*m+1)-m))
		}
	}
	sb.WriteByte('\n')
	used := make(map[[2]int]bool)
	for i := 0; i < m; i++ {
		for {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			p1 := [2]int{u, v}
			p2 := [2]int{v, u}
			if used[p1] || used[p2] {
				continue
			}
			used[p1] = true
			sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
			break
		}
	}
	return sb.String()
}

func parseInput(input string) (int, int, []int, []int, [][2]int) {
	inLines := strings.Split(strings.TrimSpace(input), "\n")
	header := strings.Fields(inLines[0])
	n, _ := strconv.Atoi(header[0])
	m, _ := strconv.Atoi(header[1])
	sFields := strings.Fields(inLines[1])
	aFields := strings.Fields(inLines[2])
	s := make([]int, n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		s[i], _ = strconv.Atoi(sFields[i])
		a[i], _ = strconv.Atoi(aFields[i])
	}
	edges := make([][2]int, m)
	for i := 0; i < m; i++ {
		fields := strings.Fields(inLines[3+i])
		u, _ := strconv.Atoi(fields[0])
		v, _ := strconv.Atoi(fields[1])
		edges[i] = [2]int{u - 1, v - 1} // 0-indexed for the solver
	}
	return n, m, s, a, edges
}

func validateOutput(input, output string) error {
	inLines := strings.Split(strings.TrimSpace(input), "\n")
	header := strings.Fields(inLines[0])
	n, _ := strconv.Atoi(header[0])
	m, _ := strconv.Atoi(header[1])
	sFields := strings.Fields(inLines[1])
	aFields := strings.Fields(inLines[2])
	s := make([]int, n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		s[i], _ = strconv.Atoi(sFields[i])
		a[i], _ = strconv.Atoi(aFields[i])
	}
	type edge struct{ u, v int }
	inEdges := make([]edge, m)
	for i := 0; i < m; i++ {
		fields := strings.Fields(inLines[3+i])
		u, _ := strconv.Atoi(fields[0])
		v, _ := strconv.Atoi(fields[1])
		inEdges[i] = edge{u, v}
	}
	outLines := strings.Split(output, "\n")
	if len(outLines) < 1+m {
		return fmt.Errorf("expected %d lines of output after YES, got %d", m, len(outLines)-1)
	}
	b := make([]int, n+1)
	usedEdges := make(map[[2]int]bool)
	for i := 0; i < m; i++ {
		fields := strings.Fields(outLines[1+i])
		if len(fields) != 2 {
			return fmt.Errorf("line %d: expected 2 integers, got %q", i+1, outLines[1+i])
		}
		u, err1 := strconv.Atoi(fields[0])
		v, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("line %d: parse error", i+1)
		}
		foundIdx := -1
		for j, e := range inEdges {
			if (e.u == u && e.v == v) || (e.u == v && e.v == u) {
				foundIdx = j
				break
			}
		}
		if foundIdx == -1 {
			return fmt.Errorf("line %d: edge (%d,%d) not in input", i+1, u, v)
		}
		e1 := [2]int{u, v}
		e2 := [2]int{v, u}
		if usedEdges[e1] || usedEdges[e2] {
			return fmt.Errorf("line %d: duplicate edge (%d,%d)", i+1, u, v)
		}
		usedEdges[e1] = true
		b[u]--
		b[v]++
	}
	if len(usedEdges) != m {
		return fmt.Errorf("expected %d edges, got %d", m, len(usedEdges))
	}
	for i := 0; i < n; i++ {
		if s[i] == 1 && b[i+1] != a[i] {
			return fmt.Errorf("b[%d]=%d but a[%d]=%d (s[%d]=1)", i+1, b[i+1], i+1, a[i], i+1)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candPath := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		n, m, s, a, edges := parseInput(input)
		expectYes := refSolveYesNo(n, m, s, a, edges)

		got, err := runBinary(candPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		gotTrim := strings.TrimSpace(got)
		gotUpper := strings.ToUpper(strings.Split(gotTrim, "\n")[0])

		if !expectYes {
			if gotUpper != "NO" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected NO, got:\n%s\ninput:\n%s", i+1, got, input)
				os.Exit(1)
			}
			continue
		}
		// Reference says YES
		if gotUpper != "YES" {
			fmt.Fprintf(os.Stderr, "case %d failed: expected YES, got:\n%s\ninput:\n%s", i+1, got, input)
			os.Exit(1)
		}
		if err := validateOutput(input, gotTrim); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%sgot:\n%s\n", i+1, err, input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
