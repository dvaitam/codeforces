package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesD = `100
2 5 9
2 4
1 2
2 1
1 5
1 1
1 4
2 2
2 5
1 3
3 1 1
1 1
2 4 4
2 3
1 1
2 4
2 2
3 6 1
1 1
3 5 11
2 4
2 1
1 5
3 1
1 1
1 4
2 3
3 3
3 2
1 3
3 5
5 1 4
3 1
1 1
2 1
5 1
2 6 1
1 1
3 3 6
1 2
1 1
3 3
2 2
3 2
1 3
5 3 7
1 2
2 1
1 1
4 2
2 2
5 3
5 2
6 6 22
4 3
4 6
5 1
2 5
1 3
6 2
4 2
3 3
5 6
3 6
5 3
1 5
6 1
6 4
4 1
3 5
5 2
1 1
2 3
2 6
6 6
6 3
4 3 4
1 1
3 3
1 3
2 2
3 5 3
2 3
1 1
3 3
4 3 11
1 3
1 2
2 1
4 3
3 1
1 1
4 2
3 3
2 2
3 2
4 1
3 2 3
1 1
2 1
3 2
1 5 2
1 1
1 3
4 5 18
4 4
2 4
1 2
3 4
2 1
1 5
3 1
4 3
1 1
1 4
4 2
4 5
3 3
2 2
3 2
2 5
4 1
3 5
2 5 9
2 4
1 2
2 1
1 5
1 1
1 4
2 2
2 5
1 3
2 2 1
1 1
3 2 1
1 1
3 2 4
1 1
1 2
2 1
2 2
2 6 10
2 4
1 2
2 1
1 5
1 1
1 4
2 6
1 6
2 5
1 3
4 3 3
2 3
1 1
4 2
1 3 2
1 1
1 3
3 2 6
1 2
2 1
3 1
1 1
2 2
3 2
3 3 1
1 1
4 4 13
1 3
2 4
1 2
2 1
4 3
3 1
1 1
4 2
1 4
3 3
2 2
3 2
4 1
1 4 1
1 1
2 6 3
1 1
2 6
2 4
1 5 4
1 1
1 3
1 4
1 5
3 5 7
1 2
1 5
1 1
2 2
3 2
2 5
1 3
4 4 6
4 4
2 4
3 4
1 1
1 4
4 1
2 2 1
1 1
4 4 7
4 3
3 1
1 1
4 2
2 3
3 3
1 3
3 6 12
1 2
2 1
3 4
1 5
3 1
1 1
1 4
2 3
3 6
2 2
1 6
3 5
5 3 9
1 2
4 3
3 1
1 1
4 2
2 3
5 3
3 2
4 1
1 6 6
1 2
1 5
1 1
1 4
1 6
1 3
6 1 2
6 1
1 1
5 1 1
1 1
2 3 2
1 1
1 3
5 4 17
4 4
2 4
1 2
4 3
3 1
4 1
1 1
5 4
5 1
4 2
1 4
2 3
2 2
5 3
3 2
1 3
5 2
6 5 15
4 4
6 2
1 2
2 1
6 5
4 3
3 1
3 4
1 1
5 1
6 4
6 3
2 5
4 1
3 5
2 1 2
1 1
2 1
4 1 2
1 1
2 1
6 3 13
6 2
2 1
4 3
3 1
6 1
1 1
5 1
4 2
2 3
2 2
3 2
6 3
4 1
3 4 6
2 1
3 4
1 1
1 4
2 3
2 2
1 6 6
1 2
1 5
1 1
1 4
1 6
1 3
4 5 20
3 4
4 3
3 1
2 2
2 5
1 3
4 2
4 5
3 3
2 4
1 2
2 1
1 5
3 2
4 1
3 5
4 4
1 1
1 4
2 3
2 2 1
1 1
5 1 4
1 1
2 1
5 1
4 1
5 6 25
3 4
4 3
3 1
5 4
4 6
5 1
1 3
4 2
4 5
3 3
5 6
3 6
5 3
1 2
2 1
1 5
3 2
3 5
5 2
4 4
5 5
1 1
1 4
2 3
2 6
4 2 6
1 2
2 1
3 1
1 1
4 2
3 2
2 4 3
2 3
1 1
2 1
3 4 1
1 1
4 2 8
1 2
2 1
3 1
1 1
4 2
2 2
3 2
4 1
3 5 3
1 1
2 5
3 5
5 3 6
1 1
5 1
4 2
2 2
5 3
4 1
2 2 1
1 1
4 1 4
3 1
1 1
4 1
2 1
4 2 3
1 1
1 2
3 2
6 1 2
1 1
2 1
4 2 3
1 1
2 1
2 2
6 3 8
4 3
3 1
1 1
2 3
5 3
3 2
1 3
5 2
4 4 1
1 1
5 2 5
1 1
4 2
3 2
4 1
5 2
5 4 6
1 1
1 4
5 3
3 2
1 3
5 2
3 4 11
2 4
2 1
3 4
3 1
1 1
1 4
2 3
3 3
2 2
3 2
1 3
6 1 3
3 1
1 1
5 1
2 5 7
2 4
2 1
1 5
1 1
1 4
2 2
2 5
3 1 2
1 1
2 1
2 3 1
1 1
2 1 2
1 1
2 1
2 3 3
1 1
1 2
1 3
2 2 3
1 1
1 2
2 2
3 1 2
1 1
2 1
3 4 12
2 4
1 2
3 4
2 1
3 1
1 1
1 4
2 3
3 3
2 2
3 2
1 3
3 5 2
1 1
3 5
2 3 1
1 1
1 4 1
1 1
5 3 4
5 3
1 1
3 3
2 2
5 4 7
2 4
3 1
1 1
5 4
5 1
4 2
3 2
3 5 3
1 1
3 3
1 3
1 1 1
1 1
1 4 2
1 1
1 4
1 1 1
1 1
6 3 10
6 2
1 2
2 1
6 1
1 1
4 2
2 2
3 2
6 3
1 3
5 6 25
3 4
4 3
3 1
5 4
4 6
5 1
2 2
1 6
2 5
4 2
4 5
3 3
5 6
3 6
1 2
2 1
1 5
3 2
4 1
3 5
5 2
4 4
5 5
1 1
2 3
5 2 7
2 1
3 1
1 1
5 1
4 2
3 2
4 1
5 3 6
2 1
3 1
1 1
2 2
4 1
5 2
2 6 2
1 1
2 5
3 2 5
1 2
2 1
3 1
1 1
2 2
4 4 9
2 4
1 2
4 3
3 1
1 1
1 4
2 3
3 3
2 2
3 1 1
1 1
3 5 9
2 4
2 1
3 4
3 1
1 1
1 4
2 3
2 5
1 3
5 3 8
1 2
2 1
3 1
1 1
2 3
3 2
1 3
5 2
1 3 3
1 1
1 2
1 3
4 6 18
4 4
2 4
1 2
2 2
2 1
3 4
4 1
3 1
1 5
1 1
4 6
4 2
2 3
3 3
3 6
1 6
1 3
3 5
5 4 7
4 4
4 3
4 1
1 1
5 4
2 3
1 3
4 6 6
2 4
3 1
1 1
2 3
2 2
2 5
6 5 4
1 1
3 3
6 2
5 1
1 3 2
1 1
1 2`

type point struct{ r, c int }

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func solve(n, m int, cells []point) int {
	cellMap := make(map[point]int, len(cells))
	for i, p := range cells {
		cellMap[p] = i
	}
	comp := make([]int, len(cells))
	for i := range comp {
		comp[i] = -1
	}
	compCount := 0
	dr := []int{-1, 1, 0, 0}
	dc := []int{0, 0, -1, 1}
	for i := 0; i < len(cells); i++ {
		if comp[i] != -1 {
			continue
		}
		queue := []int{i}
		comp[i] = compCount
		for head := 0; head < len(queue); head++ {
			u := queue[head]
			for dir := 0; dir < 4; dir++ {
				np := point{cells[u].r + dr[dir], cells[u].c + dc[dir]}
				if vIdx, ok := cellMap[np]; ok && comp[vIdx] == -1 {
					comp[vIdx] = compCount
					queue = append(queue, vIdx)
				}
			}
		}
		compCount++
	}

	compR := make([][]int, compCount)
	compC := make([][]int, compCount)
	for i, p := range cells {
		compR[comp[i]] = append(compR[comp[i]], p.r)
		compC[comp[i]] = append(compC[comp[i]], p.c)
	}
	for i := 0; i < compCount; i++ {
		sort.Ints(compR[i])
		sort.Ints(compC[i])
	}

	adj := make([][]int, compCount)
	for i := 0; i < compCount; i++ {
		for j := i + 1; j < compCount; j++ {
			can := false
			for _, r := range compR[i] {
				idx := sort.SearchInts(compR[j], r-2)
				if idx < len(compR[j]) && compR[j][idx] <= r+2 {
					can = true
					break
				}
			}
			if can {
				adj[i] = append(adj[i], j)
				adj[j] = append(adj[j], i)
				continue
			}
			for _, c := range compC[i] {
				idx := sort.SearchInts(compC[j], c-2)
				if idx < len(compC[j]) && compC[j][idx] <= c+2 {
					can = true
					break
				}
			}
			if can {
				adj[i] = append(adj[i], j)
				adj[j] = append(adj[j], i)
			}
		}
	}

	const inf = 1 << 30
	dist := make([]int, compCount)
	for i := range dist {
		dist[i] = inf
	}
	startIdx, ok := cellMap[point{1, 1}]
	if !ok {
		return -1
	}
	startComp := comp[startIdx]
	dist[startComp] = 0
	queue := []int{startComp}
	for head := 0; head < len(queue); head++ {
		u := queue[head]
		for _, v := range adj[u] {
			if dist[v] > dist[u]+1 {
				dist[v] = dist[u] + 1
				queue = append(queue, v)
			}
		}
	}

	ans := inf
	if targetIdx, ok := cellMap[point{n, m}]; ok {
		ans = min(ans, dist[comp[targetIdx]])
	}
	for i, p := range cells {
		if dist[comp[i]] == inf {
			continue
		}
		if abs(p.r-n) <= 1 || abs(p.c-m) <= 1 {
			ans = min(ans, dist[comp[i]]+1)
		}
	}
	if ans == inf {
		return -1
	}
	return ans
}

type testCase struct {
	n, m int
	pts  []point
}

func parseTests() ([]testCase, error) {
	fields := strings.Fields(testcasesD)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty testcases")
	}
	ptr := 0
	t, err := strconv.Atoi(fields[ptr])
	if err != nil {
		return nil, err
	}
	ptr++
	tests := make([]testCase, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if ptr+3 > len(fields) {
			return nil, fmt.Errorf("incomplete header for case %d", caseIdx+1)
		}
		n, _ := strconv.Atoi(fields[ptr])
		m, _ := strconv.Atoi(fields[ptr+1])
		k, _ := strconv.Atoi(fields[ptr+2])
		ptr += 3
		if ptr+2*k > len(fields) {
			return nil, fmt.Errorf("not enough points for case %d", caseIdx+1)
		}
		pts := make([]point, k)
		for i := 0; i < k; i++ {
			r, _ := strconv.Atoi(fields[ptr])
			c, _ := strconv.Atoi(fields[ptr+1])
			ptr += 2
			pts[i] = point{r, c}
		}
		tests = append(tests, testCase{n: n, m: m, pts: pts})
	}
	return tests, nil
}

func runCandidate(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTests()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		expected := strconv.Itoa(solve(tc.n, tc.m, tc.pts))
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, len(tc.pts)))
		for _, p := range tc.pts {
			input.WriteString(fmt.Sprintf("%d %d\n", p.r, p.c))
		}
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
