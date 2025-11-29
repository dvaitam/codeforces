package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesRaw = `2 -3 3 -3 1
4 -5 0 -1 0 4 -5 -2 4
1 0 0
3 -4 4 5 5 0 0
1 1 -2
3 4 -1 -3 -2 3 -3
5 -3 3 -4 0 4 4 3 5 4 -4
5 -3 0 -1 2 -1 2 -3 1 1 0
1 0 -3
3 -3 0 2 -4 0 0
3 5 0 -3 1 -2 -2
2 4 2 4 1
5 -2 -5 0 -3 -5 2 1 4 -4 -5
4 -1 4 -2 1 -5 4 -2 1
5 2 2 4 1 1 1 -5 -1 -3 -3
3 -4 3 4 2 -3 4
4 -2 -3 -1 -5 2 0 2 0
5 -2 -1 4 -5 2 -3 -1 -2 0 -5
3 4 1 -2 -5 -2 4
4 -3 -1 5 4 3 -2 -3 1
5 5 2 0 5 2 4 -3 -5 4 1
4 0 -1 -3 -1 2 -1 1 4
2 2 -2 -3 2
4 -5 1 -2 3 -1 1 -4 -1
3 -3 -3 -2 -1 -3 -3
3 -2 4 4 -3 -5 -4
2 -1 4 -1 -1
2 -5 -1 4 0
1 1 -1
3 4 3 3 5 -1 -4
5 3 -5 -5 0 1 4 2 0 5 4
2 -4 1 1 5
1 4 -1
5 0 -3 -5 -5 -3 1 0 -5 0 0
5 4 4 1 3 -4 -4 4 3 -2 -2
3 5 3 1 -3 -5 -4
1 5 0
2 1 4 5 -4
4 -2 0 -4 -4 -3 -1 -2 5
2 -2 1 1 -2
2 -5 3 -1 1
2 3 -4 4 0
3 -1 3 -2 2 4 -2
5 -5 4 5 -4 -1 1 5 -4 -5 -4
3 0 -4 4 -1 -4 -2
3 4 -4 -4 1 2 -3
3 -2 -3 -2 -3 -5 4
1 2 2
5 3 -2 1 -2 -2 4 -4 0 -3 -1
1 2 2
2 5 2 -1 2
2 2 5 -3 4
2 4 1 0 1
3 1 -3 -3 4 3 2
1 -2 4
3 5 4 -1 -5 -1 2
2 1 -2 -1 4
3 -5 0 -5 1 -5 4
5 -5 -5 -5 1 4 -1 2 0 5 -5
5 1 -4 1 -1 -5 0 -5 -5 0 5
1 4 -5
5 5 -4 4 3 -5 -2 4 4 1 -1
1 -4 -3
1 -3 3
2 -2 2 4 -1
4 -2 -5 3 3 2 -4 -2 -3
4 5 2 -4 -2 -2 -4 -3 1
4 1 -1 4 -2 3 -5 3 5
4 -3 2 4 -3 -5 -1 1 4
4 5 2 4 1 0 5 3 -2
1 -1 3
5 -2 0 -1 -2 3 1 -1 -1 1 0
3 -2 -4 -1 -1 -2 -4
2 1 3 -3 0
2 1 2 -2 -1
5 0 1 -4 4 2 2 5 -4 -3 5
1 3 -1
3 -5 1 1 -2 -5 4
2 -5 1 1 0
4 1 -5 -1 2 5 4 4 -5
1 -5 -3
5 3 -2 -3 0 1 -3 3 -1 2 -2
1 -1 -3
5 -2 -2 -3 -5 -5 2 -1 -1 3 -5
4 2 4 5 -3 -5 -3 5 -5
5 1 0 -4 0 -5 5 2 3 -2 -1
4 -5 1 2 2 -3 -2 5 -4
5 3 -4 3 -2 5 3 1 0 -1 -2
5 -5 0 5 -1 -4 4 4 -3 -4 2
1 2 0
3 1 -3 -5 3 -3 -3
4 -2 -1 4 3 4 -5 -5 -2
1 4 5
3 -1 2 4 -1 -1 3
3 0 -3 -4 5 2 3
1 -4 -3
2 0 1 -2 3
3 0 5 0 -1 3 -5
5 -4 5 1 5 -4 4 1 2 -4 -4
4 5 -3 1 -3 1 3 -5 0
4 -4 -2 -1 3 2 -5 -2 1`

type Point struct{ x, y int }
type VerSeg struct{ x, y1, y2 int }
type HorSeg struct{ y, x1, x2 int }

// referenceSolve replicates 1054F.go to produce expected output.
func referenceSolve(k int, coords []int) (string, error) {
	if len(coords) != 2*k {
		return "", fmt.Errorf("expected %d coords, got %d", 2*k, len(coords))
	}
	pts := make([]Point, k)
	for i := 0; i < k; i++ {
		pts[i] = Point{coords[2*i], coords[2*i+1]}
	}
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].x != pts[j].x {
			return pts[i].x < pts[j].x
		}
		return pts[i].y < pts[j].y
	})
	ver := make([]VerSeg, 0)
	for i := 1; i < k; i++ {
		if pts[i].x == pts[i-1].x {
			ver = append(ver, VerSeg{pts[i].x, pts[i-1].y, pts[i].y})
		}
	}
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].y != pts[j].y {
			return pts[i].y < pts[j].y
		}
		return pts[i].x < pts[j].x
	})
	hor := make([]HorSeg, 0)
	for i := 1; i < k; i++ {
		if pts[i].y == pts[i-1].y {
			hor = append(hor, HorSeg{pts[i].y, pts[i-1].x, pts[i].x})
		}
	}

	n := len(ver)
	m := len(hor)
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if ver[i].y1 < hor[j].y && ver[i].y2 > hor[j].y &&
				hor[j].x1 < ver[i].x && hor[j].x2 > ver[i].x {
				adj[i] = append(adj[i], j)
			}
		}
	}

	pairedL := make([]int, n)
	pairedR := make([]int, m)
	for i := range pairedL {
		pairedL[i] = -1
	}
	for j := range pairedR {
		pairedR[j] = -1
	}
	used := make([]bool, n)
	var dfs func(u int) bool
	dfs = func(u int) bool {
		used[u] = true
		for _, j := range adj[u] {
			if pairedR[j] < 0 || (!used[pairedR[j]] && dfs(pairedR[j])) {
				pairedL[u] = j
				pairedR[j] = u
				return true
			}
		}
		return false
	}
	for u := 0; u < n; u++ {
		if pairedL[u] < 0 {
			for i := range used {
				used[i] = false
			}
			dfs(u)
		}
	}

	visL := make([]bool, n)
	visR := make([]bool, m)
	var dfsL func(u int)
	var dfsR func(j int)
	dfsL = func(u int) {
		if visL[u] {
			return
		}
		visL[u] = true
		for _, j := range adj[u] {
			if pairedL[u] != j {
				dfsR(j)
			}
		}
	}
	dfsR = func(j int) {
		if visR[j] {
			return
		}
		visR[j] = true
		if pairedR[j] >= 0 {
			dfsL(pairedR[j])
		}
	}
	for u := 0; u < n; u++ {
		if pairedL[u] < 0 {
			dfsL(u)
		}
	}
	maxIndL := make([]bool, n)
	maxIndR := make([]bool, m)
	for i := 0; i < n; i++ {
		if visL[i] {
			maxIndL[i] = true
		}
	}
	for j := 0; j < m; j++ {
		if !visR[j] {
			maxIndR[j] = true
		}
	}

	type Pair struct{ x, y int }
	versAns := make([]VerSeg, 0)
	was := make(map[Pair]bool)
	for i := 0; i < n; i++ {
		if !maxIndL[i] {
			continue
		}
		x := ver[i].x
		ind := n
		for j := i; j < n; j++ {
			if !maxIndL[j] || ver[j].x != x {
				ind = j
				break
			}
			was[Pair{x, ver[j].y1}] = true
			was[Pair{x, ver[j].y2}] = true
		}
		start := ver[i].y1
		end := ver[ind-1].y2
		versAns = append(versAns, VerSeg{x, start, end})
		i = ind - 1
	}
	for _, p := range pts {
		key := Pair{p.x, p.y}
		if !was[key] {
			versAns = append(versAns, VerSeg{p.x, p.y, p.y})
		}
	}

	horsAns := make([]HorSeg, 0)
	was = make(map[Pair]bool)
	for i := 0; i < m; i++ {
		if !maxIndR[i] {
			continue
		}
		y := hor[i].y
		ind := m
		for j := i; j < m; j++ {
			if !maxIndR[j] || hor[j].y != y {
				ind = j
				break
			}
			was[Pair{hor[j].x1, y}] = true
			was[Pair{hor[j].x2, y}] = true
		}
		start := hor[i].x1
		end := hor[ind-1].x2
		horsAns = append(horsAns, HorSeg{y, start, end})
		i = ind - 1
	}
	for _, p := range pts {
		key := Pair{p.x, p.y}
		if !was[key] {
			horsAns = append(horsAns, HorSeg{p.y, p.x, p.x})
		}
	}

	var out bytes.Buffer
	fmt.Fprintln(&out, len(horsAns))
	for _, s := range horsAns {
		fmt.Fprintf(&out, "%d %d %d %d\n", s.x1, s.y, s.x2, s.y)
	}
	fmt.Fprintln(&out, len(versAns))
	for _, s := range versAns {
		fmt.Fprintf(&out, "%d %d %d %d\n", s.x, s.y1, s.x, s.y2)
	}
	return strings.TrimSpace(out.String()), nil
}

func parseLine(line string) (int, []int, error) {
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return 0, nil, fmt.Errorf("empty test case")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid k: %v", err)
	}
	if len(fields) != 1+2*k {
		return 0, nil, fmt.Errorf("expected %d coords, got %d", 2*k, len(fields)-1)
	}
	coords := make([]int, 2*k)
	for i := 0; i < 2*k; i++ {
		coords[i], err = strconv.Atoi(fields[1+i])
		if err != nil {
			return 0, nil, fmt.Errorf("invalid coord: %v", err)
		}
	}
	return k, coords, nil
}

func runCase(bin string, input string, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		k, coords, err := parseLine(line)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
		input := line + "\n"
		expected, err := referenceSolve(k, coords)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if err := runCase(bin, input, expected); err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
