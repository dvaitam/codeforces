package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const testCount = 80

type point struct {
	x int
	y int
}

type solution struct {
	impossible bool
	edges      [][2]int
}

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1045E.go")
	tmp, err := os.CreateTemp("", "oracle1045E")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()
	os.Remove(path)
	cmd := exec.Command("go", "build", "-o", path, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(out))
	}
	return path, nil
}

func runProgram(bin, input string) (string, error) {
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

func collinear(a, b, c point) bool {
	return int64(b.x-a.x)*int64(c.y-a.y) == int64(b.y-a.y)*int64(c.x-a.x)
}

func genPoints(n int, r *rand.Rand) []point {
	points := make([]point, 0, n)
	for len(points) < n {
		x := r.Intn(2001) - 1000
		y := r.Intn(2001) - 1000
		ok := true
		for _, p := range points {
			if p.x == x && p.y == y {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		if len(points) >= 2 {
			for i := 0; i < len(points) && ok; i++ {
				for j := i + 1; j < len(points); j++ {
					if collinear(points[i], points[j], point{x: x, y: y}) {
						ok = false
						break
					}
				}
			}
		}
		if ok {
			points = append(points, point{x: x, y: y})
		}
	}
	return points
}

func ensureColors(n int, r *rand.Rand) []int {
	colors := make([]int, n)
	hasZero, hasOne := false, false
	for i := 0; i < n; i++ {
		colors[i] = r.Intn(2)
		if colors[i] == 0 {
			hasZero = true
		} else {
			hasOne = true
		}
	}
	if n > 1 {
		if !hasZero {
			idx := r.Intn(n)
			colors[idx] = 0
			hasZero = true
		}
		if !hasOne {
			idx := r.Intn(n)
			for colors[idx] == 1 {
				idx = r.Intn(n)
			}
			colors[idx] = 1
		}
	}
	return colors
}

func genCase(r *rand.Rand) (string, []point, []int) {
	n := r.Intn(30) + 1
	points := genPoints(n, r)
	colors := ensureColors(n, r)

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d %d\n", points[i].x, points[i].y, colors[i])
	}
	return sb.String(), points, colors
}

func parseSolution(out string, n int) (solution, error) {
	var sol solution
	reader := bufio.NewReader(strings.NewReader(out))
	firstToken, err := nextToken(reader)
	if err != nil {
		return sol, fmt.Errorf("output missing")
	}
	if strings.EqualFold(firstToken, "impossible") {
		if extra, err := nextToken(reader); err == nil && extra != "" {
			return sol, fmt.Errorf("extra data after Impossible")
		}
		sol.impossible = true
		return sol, nil
	}
	m, err := strconv.Atoi(firstToken)
	if err != nil {
		return sol, fmt.Errorf("first token is not an integer")
	}
	if m < 0 {
		return sol, fmt.Errorf("negative number of roads")
	}
	sol.edges = make([][2]int, m)
	for i := 0; i < m; i++ {
		aStr, err := nextToken(reader)
		if err != nil {
			return sol, fmt.Errorf("missing endpoint for road %d", i+1)
		}
		bStr, err2 := nextToken(reader)
		if err2 != nil {
			return sol, fmt.Errorf("missing endpoint for road %d", i+1)
		}
		a, err := strconv.Atoi(aStr)
		if err != nil {
			return sol, fmt.Errorf("invalid vertex index in road %d", i+1)
		}
		b, err := strconv.Atoi(bStr)
		if err != nil {
			return sol, fmt.Errorf("invalid vertex index in road %d", i+1)
		}
		sol.edges[i] = [2]int{a, b}
	}
	if extra, err := nextToken(reader); err == nil && extra != "" {
		return sol, fmt.Errorf("extra data after listed roads")
	}
	return sol, nil
}

func nextToken(r *bufio.Reader) (string, error) {
	var sb strings.Builder
	for {
		ch, err := r.ReadByte()
		if err != nil {
			if sb.Len() == 0 {
				return "", err
			}
			break
		}
		if ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r' {
			if sb.Len() == 0 {
				continue
			}
			break
		}
		sb.WriteByte(ch)
	}
	return sb.String(), nil
}

func validateSolution(edges [][2]int, pts []point, colors []int) error {
	n := len(colors)
	adj := make([][]int, n)
	type key struct{ u, v int }
	used := make(map[key]struct{})
	colorEdgeCount := make(map[int]int)

	for idx, e := range edges {
		u, v := e[0], e[1]
		if u < 0 || u >= n || v < 0 || v >= n {
			return fmt.Errorf("road %d uses invalid vertex index", idx+1)
		}
		if u == v {
			return fmt.Errorf("road %d connects vertex %d to itself", idx+1, u)
		}
		if colors[u] != colors[v] {
			return fmt.Errorf("road %d connects different civilizations", idx+1)
		}
		a, b := u, v
		if a > b {
			a, b = b, a
		}
		k := key{a, b}
		if _, ok := used[k]; ok {
			return fmt.Errorf("road %d duplicates connection between %d and %d", idx+1, a, b)
		}
		used[k] = struct{}{}
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		colorEdgeCount[colors[u]]++
	}

	colorNodes := make(map[int][]int)
	for i, c := range colors {
		colorNodes[c] = append(colorNodes[c], i)
	}

	for c, nodes := range colorNodes {
		edgeCnt := colorEdgeCount[c]
		if len(nodes) <= 1 {
			if edgeCnt != 0 {
				return fmt.Errorf("civilization %d should have 0 roads but has %d", c, edgeCnt)
			}
			continue
		}
		if edgeCnt != len(nodes)-1 {
			return fmt.Errorf("civilization %d should have %d roads but has %d", c, len(nodes)-1, edgeCnt)
		}
		visited := make([]bool, n)
		queue := []int{nodes[0]}
		visited[nodes[0]] = true
		count := 0
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			count++
			for _, v := range adj[u] {
				if !visited[v] {
					visited[v] = true
					queue = append(queue, v)
				}
			}
		}
		if count != len(nodes) {
			return fmt.Errorf("civilization %d graph is not connected", c)
		}
	}

	for i := 0; i < len(edges); i++ {
		for j := i + 1; j < len(edges); j++ {
			u1, v1 := edges[i][0], edges[i][1]
			u2, v2 := edges[j][0], edges[j][1]
			if u1 == u2 || u1 == v2 || v1 == u2 || v1 == v2 {
				continue
			}
			if segmentsIntersect(pts[u1], pts[v1], pts[u2], pts[v2]) {
				return fmt.Errorf("roads (%d,%d) and (%d,%d) intersect", u1, v1, u2, v2)
			}
		}
	}

	return nil
}

func orientation(a, b, c point) int {
	val := int64(b.x-a.x)*int64(c.y-a.y) - int64(b.y-a.y)*int64(c.x-a.x)
	if val > 0 {
		return 1
	}
	if val < 0 {
		return -1
	}
	return 0
}

func segmentsIntersect(a1, a2, b1, b2 point) bool {
	o1 := orientation(a1, a2, b1)
	o2 := orientation(a1, a2, b2)
	o3 := orientation(b1, b2, a1)
	o4 := orientation(b1, b2, a2)
	return o1*o2 < 0 && o3*o4 < 0
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	r := rand.New(rand.NewSource(1))
	for t := 0; t < testCount; t++ {
		input, pts, colors := genCase(r)
		expectOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		expectSol, err := parseSolution(expectOut, len(colors))
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotOut, err := runProgram(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotSol, err := parseSolution(gotOut, len(colors))
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", t+1, input, err)
			os.Exit(1)
		}
		if expectSol.impossible {
			if gotSol.impossible {
				continue
			}
			fmt.Printf("test %d failed\ninput:\n%s\nerror: expected Impossible but got a plan\n", t+1, input)
			os.Exit(1)
		}
		if gotSol.impossible {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: unexpected Impossible\n", t+1, input)
			os.Exit(1)
		}
		if err := validateSolution(gotSol.edges, pts, colors); err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", t+1, input, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
