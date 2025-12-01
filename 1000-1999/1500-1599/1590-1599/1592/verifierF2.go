package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "./1592F2.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF2.go /path/to/binary")
		os.Exit(1)
	}

	if _, err := buildReference(); err != nil {
		// Reference solution prints a placeholder value; its output is not used.
		// We only check that it builds successfully.
		fmt.Fprintf(os.Stderr, "warning: failed to build reference: %v\n", err)
	}

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		expected, err := solveInput(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to solve test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}

		output, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, output)
			os.Exit(1)
		}
		got, err := parseInt(output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, output)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d, got %d\ninput:\n%s\n", idx+1, tc.name, expected, got, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1592F2-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	os.Remove(tmp.Name())
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func parseInt(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	return val, nil
}

func solveInput(input string) (int64, error) {
	reader := strings.NewReader(input)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return 0, err
	}
	rows := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &rows[i])
		if len(rows[i]) != m {
			return 0, fmt.Errorf("row %d has invalid length", i)
		}
	}
	cells := n * m
	if cells > 16 {
		return 0, fmt.Errorf("grid too large for verifier solver")
	}
	target := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if rows[i][j] == 'B' {
				target |= 1 << (i*m + j)
			} else if rows[i][j] != 'W' {
				return 0, fmt.Errorf("invalid character %q", rows[i][j])
			}
		}
	}
	ops := buildOperations(n, m)
	return dijkstra(target, ops, cells), nil
}

type operation struct {
	mask int
	cost int64
}

func buildOperations(n, m int) []operation {
	costs := make(map[int]int64)
	apply := func(mask int, cost int64) {
		if mask == 0 {
			return
		}
		if old, ok := costs[mask]; !ok || cost < old {
			costs[mask] = cost
		}
	}

	for x := 1; x <= n; x++ {
		for y := 1; y <= m; y++ {
			apply(rectMask(1, 1, x, y, n, m), 1) // top-left
			apply(rectMask(x, 1, n, y, n, m), 3) // bottom-left
			apply(rectMask(1, y, x, m, n, m), 4) // top-right
			apply(rectMask(x, y, n, m, n, m), 2) // bottom-right
		}
	}

	ops := make([]operation, 0, len(costs))
	for mask, cost := range costs {
		ops = append(ops, operation{mask: mask, cost: cost})
	}
	return ops
}

func rectMask(x1, y1, x2, y2, n, m int) int {
	mask := 0
	for i := x1; i <= x2; i++ {
		for j := y1; j <= y2; j++ {
			idx := (i-1)*m + (j - 1)
			mask |= 1 << idx
		}
	}
	return mask
}

type state struct {
	mask int
	dist int64
}

type priorityQueue []state

func (pq priorityQueue) Len() int { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}
func (pq priorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(state))
}
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func dijkstra(target int, ops []operation, cells int) int64 {
	const inf = int64(1 << 60)
	size := 1 << cells
	dist := make([]int64, size)
	for i := range dist {
		dist[i] = inf
	}
	dist[0] = 0
	pq := &priorityQueue{}
	heap.Push(pq, state{mask: 0, dist: 0})

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(state)
		if cur.dist != dist[cur.mask] {
			continue
		}
		if cur.mask == target {
			return cur.dist
		}
		for _, op := range ops {
			next := cur.mask ^ op.mask
			nd := cur.dist + op.cost
			if nd < dist[next] {
				dist[next] = nd
				heap.Push(pq, state{mask: next, dist: nd})
			}
		}
	}
	return dist[target]
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase("all-white-1x1", []string{"W"}),
		buildCase("single-black-1x1", []string{"B"}),
		buildCase("row", []string{"WB"}),
		buildCase("column", []string{"W", "B"}),
		buildCase("checker-2x2", []string{"WB", "BW"}),
	}

	rng := rand.New(rand.NewSource(1592))
	for i := 0; i < 40; i++ {
		n := rng.Intn(3) + 1
		m := rng.Intn(4) + 1
		for n*m > 12 {
			n = rng.Intn(3) + 1
			m = rng.Intn(4) + 1
		}
		rows := make([]string, n)
		for r := 0; r < n; r++ {
			var sb strings.Builder
			for c := 0; c < m; c++ {
				if rng.Intn(2) == 0 {
					sb.WriteByte('W')
				} else {
					sb.WriteByte('B')
				}
			}
			rows[r] = sb.String()
		}
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), rows))
	}
	return tests
}

func buildCase(name string, rows []string) testCase {
	if len(rows) == 0 {
		panic("empty grid")
	}
	n := len(rows)
	m := len(rows[0])
	for i := range rows {
		if len(rows[i]) != m {
			panic("inconsistent row lengths")
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for _, row := range rows {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}
