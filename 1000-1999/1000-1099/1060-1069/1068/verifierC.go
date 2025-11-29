package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `100
6 9
6 2 2 1 6 5 6 1 5 4 4 6 4 2 6 3 2 5
6 7
1 3 4 3 1 5 6 1 2 6 4 1 5 2
2 0

7 0

6 0

2 0

5 10
1 2 4 3 3 1 5 1 4 2 1 4 4 5 5 3 3 2 2 5
8 7
2 7 5 8 4 6 5 1 3 6 1 3 4 7
9 3
6 2 8 4 6 9
8 5
6 2 4 6 5 1 5 7 1 4
5 7
2 4 2 1 5 1 1 4 4 5 3 2 1 3
7 0

3 1
2 1
3 2
3 2 2 1
10 8
1 2 10 8 3 7 5 1 9 5 4 5 8 9 10 2
9 8
7 1 2 7 6 4 9 5 2 9 5 6 8 5 5 2
10 2
3 1 10 1
3 3
3 1 2 3 2 1
10 0

9 9
9 3 3 7 6 8 9 6 4 6 1 4 5 7 2 5 8 5
1 0

1 0

4 2
2 3 4 2
1 0

10 8
2 1 3 4 8 1 6 1 4 6 1 4 6 7 1 3
8 6
7 1 8 1 6 7 3 6 5 3 7 5
10 1
3 1
9 2
2 6 9 3
4 6
1 2 3 4 4 1 4 2 2 3 1 3
3 3
3 2 1 2 1 3
10 8
2 7 8 7 6 8 5 4 2 10 5 1 6 7 5 10
6 3
6 3 5 6 1 5
1 0

6 5
2 1 5 4 1 4 3 6 2 5
3 2
2 3 1 3
5 9
2 1 4 3 5 4 5 1 1 4 4 2 5 3 3 2 5 2
1 0

10 6
3 4 6 9 7 10 6 8 4 1 1 9
6 10
6 2 1 2 2 4 5 4 6 4 5 6 1 6 3 2 4 1 5 2
2 1
1 2
1 0

1 0

4 6
2 1 3 4 3 1 4 2 1 4 3 2
3 2
2 3 1 3
3 3
2 3 2 1 3 1
6 9
2 4 1 2 6 5 5 4 6 4 1 6 3 2 1 3 5 2
10 3
3 8 8 9 3 6
9 0

9 1
9 6
1 0

2 1
2 1
3 2
3 1 1 2
10 3
1 7 10 3 1 5
6 1
4 1
8 5
8 7 3 1 8 5 2 6 4 7
2 0

1 0

1 0

7 9
1 5 4 1 7 3 1 7 7 6 5 6 1 6 1 3 4 7
3 2
3 2 1 3
1 0

7 0

4 1
3 2
3 3
3 1 1 2 2 3
9 0

1 0

2 0

3 3
2 3 1 2 3 1
3 2
3 1 3 2
2 0

7 9
3 1 5 1 4 2 1 4 7 3 5 7 2 6 6 3 3 5
10 5
3 7 1 7 7 2 9 1 2 8
10 7
10 5 6 8 9 2 1 4 7 3 8 9 3 6
5 2
5 3 3 4
1 0

9 6
3 1 4 6 5 7 1 7 2 6 8 2
2 1
1 2
5 2
2 4 3 5
9 4
6 7 8 9 4 7 5 7
7 1
6 2
3 0

10 9
4 10 9 2 2 6 7 2 10 3 10 9 8 6 2 5 2 8
6 9
6 2 2 4 2 1 1 5 4 3 5 4 6 4 5 6 6 3
1 0

7 2
4 5 1 4
2 1
2 1
9 7
3 4 9 6 1 8 9 8 2 5 9 7 8 5
3 2
3 1 2 1
2 1
1 2
8 3
2 3 2 5 3 6
1 0

10 8
9 10 2 1 5 4 1 7 8 9 8 2 4 1 8 5
1 0

6 8
3 4 3 1 5 4 6 4 2 6 3 2 6 3 5 2
2 1
1 2
6 6
4 3 3 1 4 6 1 4 2 3 2 6
3 2
1 3 2 1
7 6
3 4 2 3 4 5 2 6 5 6 7 5
1 0

8 5
7 4 7 1 4 6 3 2 2 5`

type Rook struct {
	x, y  int
	color int
}

type testCase struct {
	n     int
	edges [][2]int
}

func parseOutput(n int, out string) ([]Rook, [][]int, error) {
	r := strings.NewReader(out)
	var rooks []Rook
	colorIdx := make([][]int, n+1)
	coords := make(map[[2]int]bool)
	total := 0
	for c := 1; c <= n; c++ {
		var cnt int
		if _, err := fmt.Fscan(r, &cnt); err != nil {
			return nil, nil, fmt.Errorf("failed to read count for color %d: %v", c, err)
		}
		if cnt <= 0 || cnt > 5000 {
			return nil, nil, fmt.Errorf("invalid count for color %d", c)
		}
		for i := 0; i < cnt; i++ {
			var x, y int
			if _, err := fmt.Fscan(r, &x, &y); err != nil {
				return nil, nil, fmt.Errorf("failed to read rook for color %d: %v", c, err)
			}
			if x < 1 || x > 1_000_000_000 || y < 1 || y > 1_000_000_000 {
				return nil, nil, fmt.Errorf("coordinates out of range: %d %d", x, y)
			}
			key := [2]int{x, y}
			if coords[key] {
				return nil, nil, fmt.Errorf("duplicate cell %d %d", x, y)
			}
			coords[key] = true
			rooks = append(rooks, Rook{x: x, y: y, color: c})
			colorIdx[c] = append(colorIdx[c], len(rooks)-1)
		}
		total += cnt
		if total > 5000 {
			return nil, nil, fmt.Errorf("total rooks exceed limit")
		}
	}
	if _, err := fmt.Fscan(r, new(int)); err == nil {
		return nil, nil, fmt.Errorf("extra output")
	}
	return rooks, colorIdx, nil
}

func isConnected(indices []int, rooks []Rook, rowMap, colMap map[int][]int) bool {
	if len(indices) == 0 {
		return false
	}
	allowed := make(map[int]bool, len(indices))
	for _, idx := range indices {
		allowed[idx] = true
	}
	visited := make(map[int]bool)
	q := []int{indices[0]}
	visited[indices[0]] = true
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		row := rooks[u].x
		for _, v := range rowMap[row] {
			if allowed[v] && !visited[v] {
				visited[v] = true
				q = append(q, v)
			}
		}
		col := rooks[u].y
		for _, v := range colMap[col] {
			if allowed[v] && !visited[v] {
				visited[v] = true
				q = append(q, v)
			}
		}
	}
	return len(visited) == len(indices)
}

func verifyCase(n int, edges [][2]int, output string) error {
	rooks, colorIdx, err := parseOutput(n, output)
	if err != nil {
		return err
	}
	rowMap := make(map[int][]int)
	colMap := make(map[int][]int)
	for i, r := range rooks {
		rowMap[r.x] = append(rowMap[r.x], i)
		colMap[r.y] = append(colMap[r.y], i)
	}
	for c := 1; c <= n; c++ {
		if !isConnected(colorIdx[c], rooks, rowMap, colMap) {
			return fmt.Errorf("color %d rooks not connected", c)
		}
	}
	adj := make([][]bool, n+1)
	for i := range adj {
		adj[i] = make([]bool, n+1)
	}
	for _, e := range edges {
		a, b := e[0], e[1]
		adj[a][b] = true
		adj[b][a] = true
	}
	for a := 1; a <= n; a++ {
		for b := a + 1; b <= n; b++ {
			inds := append(append([]int{}, colorIdx[a]...), colorIdx[b]...)
			conn := isConnected(inds, rooks, rowMap, colMap)
			if adj[a][b] && !conn {
				return fmt.Errorf("colors %d and %d should be connected", a, b)
			}
			if !adj[a][b] && conn {
				return fmt.Errorf("colors %d and %d should not be connected", a, b)
			}
		}
	}
	return nil
}

// referenceSolve embeds the logic from 1068C.go to produce one valid placement.
func referenceSolve(tc testCase) string {
	n := tc.n
	sol := make([][]int, n+1)
	cnt := 0
	for _, e := range tc.edges {
		cnt++
		x, y := e[0], e[1]
		sol[x] = append(sol[x], cnt)
		sol[y] = append(sol[y], cnt)
	}
	for i := 1; i <= n; i++ {
		if len(sol[i]) == 0 {
			cnt++
			sol[i] = append(sol[i], cnt)
		}
	}
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		sb.WriteString(strconv.Itoa(len(sol[i])))
		sb.WriteByte('\n')
		for _, id := range sol[i] {
			sb.WriteString(fmt.Sprintf("%d %d\n", i, id))
		}
	}
	return sb.String()
}

func runSolution(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func parseTestcases(raw string) ([]testCase, error) {
	tokens := strings.Fields(raw)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("invalid test file")
	}
	idx := 0
	nextInt := func() (int, error) {
		if idx >= len(tokens) {
			return 0, fmt.Errorf("unexpected end of test data")
		}
		v, err := strconv.Atoi(tokens[idx])
		idx++
		return v, err
	}
	t, err := nextInt()
	if err != nil {
		return nil, err
	}
	tests := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		n, err := nextInt()
		if err != nil {
			return nil, err
		}
		m, err := nextInt()
		if err != nil {
			return nil, err
		}
		edges := make([][2]int, m)
		for j := 0; j < m; j++ {
			a, err := nextInt()
			if err != nil {
				return nil, err
			}
			b, err := nextInt()
			if err != nil {
				return nil, err
			}
			edges[j] = [2]int{a, b}
		}
		tests = append(tests, testCase{n: n, edges: edges})
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("extra data in testcases")
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		// Ensure reference solution is valid for this case (sanity check).
		if err := verifyCase(tc.n, tc.edges, referenceSolve(tc)); err != nil {
			fmt.Printf("internal reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.edges)))
		for _, e := range tc.edges {
			input.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		out, err := runSolution(bin, input.String())
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := verifyCase(tc.n, tc.edges, out); err != nil {
			fmt.Printf("test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
