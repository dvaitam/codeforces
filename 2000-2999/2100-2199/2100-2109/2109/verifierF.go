package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ---------- embedded correct solver ----------

type solverItem struct {
	i, j int
	d    int64
}

type solverPQ []solverItem

func (pq solverPQ) Len() int           { return len(pq) }
func (pq solverPQ) Less(i, j int) bool { return pq[i].d < pq[j].d }
func (pq solverPQ) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *solverPQ) Push(x interface{}) {
	*pq = append(*pq, x.(solverItem))
}
func (pq *solverPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

type solverPair struct{ i, j int }

func solveCase(n, r, k int, a [][]int, c []string) (int, int) {
	dirs4 := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	dirs8 := [8][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}

	visited := make([][]bool, n)
	for i := 0; i < n; i++ {
		visited[i] = make([]bool, n)
	}

	q := make([]solverPair, 0, n*n)

	pathExists := func(limit int) bool {
		if a[0][0] > limit || a[r-1][n-1] > limit {
			return false
		}
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				visited[i][j] = false
			}
		}
		q = q[:0]
		q = append(q, solverPair{0, 0})
		visited[0][0] = true
		head := 0
		for head < len(q) {
			curr := q[head]
			head++
			if curr.i == r-1 && curr.j == n-1 {
				return true
			}
			for _, dir := range dirs4 {
				ni, nj := curr.i+dir[0], curr.j+dir[1]
				if ni >= 0 && ni < n && nj >= 0 && nj < n && !visited[ni][nj] && a[ni][nj] <= limit {
					visited[ni][nj] = true
					q = append(q, solverPair{ni, nj})
				}
			}
		}
		return false
	}

	lowM, highM := 1, 1000000
	DM := -1
	for lowM <= highM {
		mid := (lowM + highM) / 2
		if pathExists(mid) {
			DM = mid
			highM = mid - 1
		} else {
			lowM = mid + 1
		}
	}

	reach1 := make([][]bool, n)
	for i := 0; i < n; i++ {
		reach1[i] = make([]bool, n)
	}
	q = q[:0]
	q = append(q, solverPair{0, 0})
	reach1[0][0] = true
	head := 0
	for head < len(q) {
		curr := q[head]
		head++
		for _, dir := range dirs4 {
			ni, nj := curr.i+dir[0], curr.j+dir[1]
			if ni >= 0 && ni < n && nj >= 0 && nj < n && !reach1[ni][nj] && a[ni][nj] <= DM {
				reach1[ni][nj] = true
				q = append(q, solverPair{ni, nj})
			}
		}
	}

	reach2 := make([][]bool, n)
	for i := 0; i < n; i++ {
		reach2[i] = make([]bool, n)
	}
	q = q[:0]
	q = append(q, solverPair{r - 1, n - 1})
	reach2[r-1][n-1] = true
	head = 0
	for head < len(q) {
		curr := q[head]
		head++
		for _, dir := range dirs4 {
			ni, nj := curr.i+dir[0], curr.j+dir[1]
			if ni >= 0 && ni < n && nj >= 0 && nj < n && !reach2[ni][nj] && a[ni][nj] <= DM {
				reach2[ni][nj] = true
				q = append(q, solverPair{ni, nj})
			}
		}
	}

	inVpath := make([][]bool, n)
	for i := 0; i < n; i++ {
		inVpath[i] = make([]bool, n)
		for j := 0; j < n; j++ {
			if reach1[i][j] && reach2[i][j] {
				inVpath[i][j] = true
			}
		}
	}

	dist := make([][]int64, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int64, n)
	}
	pq := make(solverPQ, 0, n*n)

	check := func(X int) bool {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				dist[i][j] = 1e18
			}
		}
		pq = pq[:0]

		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if X > DM && inVpath[i][j] {
					continue
				}
				var w int64 = 0
				if a[i][j] < X {
					if c[i][j] == '1' {
						w = int64(X - a[i][j])
					} else {
						w = 1e18
					}
				}
				isP1 := (i == 0 || j == 0 || (j == n-1 && i <= r-1))
				if isP1 && w <= int64(k) {
					dist[i][j] = w
					heap.Push(&pq, solverItem{i, j, w})
				}
			}
		}

		for pq.Len() > 0 {
			curr := heap.Pop(&pq).(solverItem)
			ci, cj, d := curr.i, curr.j, curr.d
			if d > dist[ci][cj] {
				continue
			}

			isP2 := (ci == n-1 || (cj == n-1 && ci >= r-1))
			if isP2 {
				return true
			}

			for _, dir := range dirs8 {
				ni, nj := ci+dir[0], cj+dir[1]
				if ni >= 0 && ni < n && nj >= 0 && nj < n {
					if X > DM && inVpath[ni][nj] {
						continue
					}
					var w int64 = 0
					if a[ni][nj] < X {
						if c[ni][nj] == '1' {
							w = int64(X - a[ni][nj])
						} else {
							w = 1e18
						}
					}
					if w != 1e18 && d+w <= int64(k) && d+w < dist[ni][nj] {
						dist[ni][nj] = d + w
						heap.Push(&pq, solverItem{ni, nj, dist[ni][nj]})
					}
				}
			}
		}
		return false
	}

	ansF := 0
	lowF, highF := 1, 2000005
	for lowF <= highF {
		mid := (lowF + highF) / 2
		if check(mid) {
			ansF = mid
			lowF = mid + 1
		} else {
			highF = mid - 1
		}
	}

	return DM, ansF
}

// ---------- test infrastructure ----------

type testCase struct {
	n int
	r int
	k int
	a [][]int
	c []string
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutput(out string, t int) ([][2]int, error) {
	fields := strings.Fields(out)
	if len(fields) != 2*t {
		return nil, fmt.Errorf("expected %d tokens, got %d", 2*t, len(fields))
	}
	ans := make([][2]int, t)
	for i := 0; i < t; i++ {
		dm, err := strconv.Atoi(fields[2*i])
		if err != nil {
			return nil, fmt.Errorf("token %d invalid int %q: %v", 2*i+1, fields[2*i], err)
		}
		df, err := strconv.Atoi(fields[2*i+1])
		if err != nil {
			return nil, fmt.Errorf("token %d invalid int %q: %v", 2*i+2, fields[2*i+1], err)
		}
		ans[i] = [2]int{dm, df}
	}
	return ans, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.r, tc.k)
		for i := 0; i < tc.n; i++ {
			for j := 0; j < tc.n; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", tc.a[i][j])
			}
			sb.WriteByte('\n')
		}
		for i := 0; i < tc.n; i++ {
			sb.WriteString(tc.c[i])
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 2, r: 1, k: 30,
			a: [][]int{{2, 2}, {1, 1}},
			c: []string{"10", "11"},
		},
		{
			n: 3, r: 3, k: 5,
			a: [][]int{{9, 2, 2}, {2, 3, 2}, {2, 2, 2}},
			c: []string{"111", "110", "100"},
		},
		{
			n: 2, r: 2, k: 0,
			a: [][]int{{1, 1}, {1, 1}},
			c: []string{"00", "00"},
		},
	}
}

func randomTests() []testCase {
	const limit = 90000
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 40)
	sumSq := 0

	add := func(tc testCase) {
		if sumSq+tc.n*tc.n > limit {
			return
		}
		tests = append(tests, tc)
		sumSq += tc.n * tc.n
	}

	for len(tests) < 10 {
		add(randomCase(rng, 2, 8))
	}
	for len(tests) < 20 {
		add(randomCase(rng, 5, 20))
	}
	for len(tests) < 28 {
		add(randomCase(rng, 20, 80))
	}
	if sumSq+300*300 <= limit {
		add(randomCase(rng, 300, 300))
	} else if sumSq+200*200 <= limit {
		add(randomCase(rng, 200, 200))
	}

	return tests
}

func randomCase(rng *rand.Rand, lo, hi int) testCase {
	n := lo
	if hi > lo {
		n = rng.Intn(hi-lo+1) + lo
	}
	r := rng.Intn(n) + 1
	k := rng.Intn(1_000_001)
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, n)
		for j := 0; j < n; j++ {
			a[i][j] = rng.Intn(1_000_000) + 1
		}
	}
	c := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				row[j] = '0'
			} else {
				row[j] = '1'
			}
		}
		c[i] = string(row)
	}
	return testCase{n: n, r: r, k: k, a: a, c: c}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	tests := append(deterministicTests(), randomTests()...)
	input := buildInput(tests)

	// Compute expected answers using embedded solver.
	expAns := make([][2]int, len(tests))
	for i, tc := range tests {
		dm, df := solveCase(tc.n, tc.r, tc.k, tc.a, tc.c)
		expAns[i] = [2]int{dm, df}
	}

	gotOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\noutput:\n%s\n", err, gotOut)
		os.Exit(1)
	}

	gotAns, err := parseOutput(gotOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse target output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if expAns[i] != gotAns[i] {
			fmt.Fprintf(os.Stderr, "mismatch on test %d (n=%d): expected (%d %d), got (%d %d)\n", i+1, tests[i].n, expAns[i][0], expAns[i][1], gotAns[i][0], gotAns[i][1])
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
