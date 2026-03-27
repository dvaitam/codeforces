package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ---------- embedded correct solver for 822F ----------

type edgeF struct {
	to int
	id int
}

func solve822F(input string) (int, [][]int, [][2]int, []float64) {
	fields := strings.Fields(input)
	pos := 0
	readInt := func() int {
		v, _ := strconv.Atoi(fields[pos])
		pos++
		return v
	}

	n := readInt()
	adj := make([][]edgeF, n+1)
	type rawEdge struct{ u, v int }
	rawEdges := make([]rawEdge, n) // 1-indexed

	for i := 1; i < n; i++ {
		u := readInt()
		v := readInt()
		adj[u] = append(adj[u], edgeF{v, i})
		adj[v] = append(adj[v], edgeF{u, i})
		rawEdges[i] = rawEdge{u, v}
	}

	// BFS traversal, assign times to edges
	type edgeTime struct {
		fromU float64 // time point arrives at u-side of edge (going u->v)
	}
	target := make([]map[int]float64, n) // target[edgeId][node] = time point is at that node
	for i := 1; i < n; i++ {
		target[i] = make(map[int]float64)
	}

	visited := make([]bool, n+1)
	visited[1] = true
	queue := []int{1}

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]

		du := len(adj[u])
		if du == 0 {
			continue
		}

		var fixedTime float64
		fixedEdgeId := -1

		for _, e := range adj[u] {
			if val, ok := target[e.id][u]; ok {
				fixedTime = val
				fixedEdgeId = e.id
				break
			}
		}

		assignIdx := 1
		if fixedEdgeId == -1 {
			assignIdx = 0
		}

		for _, e := range adj[u] {
			if e.id == fixedEdgeId {
				continue
			}
			t := fixedTime + float64(assignIdx)*2.0/float64(du)
			for t >= 2.0 {
				t -= 2.0
			}
			target[e.id][u] = t

			tTo := t + 1.0
			for tTo >= 2.0 {
				tTo -= 2.0
			}
			target[e.id][e.to] = tTo

			assignIdx++
		}

		for _, e := range adj[u] {
			if !visited[e.to] {
				visited[e.to] = true
				queue = append(queue, e.to)
			}
		}
	}

	numPaths := n - 1
	pathEdges := make([][]int, numPaths)
	pathEndpoints := make([][2]int, numPaths)
	pathX := make([]float64, numPaths)

	for i := 1; i < n; i++ {
		u := rawEdges[i].u
		v := rawEdges[i].v
		T := target[i][u]

		var printU, printV int
		var x float64

		if T <= 1.0 {
			printU = v
			printV = u
			x = 1.0 - T
		} else {
			printU = u
			printV = v
			x = 2.0 - T
		}

		if x < 0.0 {
			x = 0.0
		} else if x > 1.0 {
			x = 1.0
		}

		pathEdges[i-1] = []int{i}
		pathEndpoints[i-1] = [2]int{printU, printV}
		pathX[i-1] = x
	}

	return numPaths, pathEdges, pathEndpoints, pathX
}

// Simulate the res array from a solution output for a given tree
func simulateRes(n int, edgesUV [][2]int, output string) ([]float64, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty output")
	}

	numPaths, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("cannot parse numPaths: %v", err)
	}

	// Build adjacency for the tree (edge index -> endpoints)
	adj := make([]map[int]int, n+1) // adj[u][v] = edgeId
	for i := 1; i <= n; i++ {
		adj[i] = make(map[int]int)
	}
	for i, e := range edgesUV {
		adj[e[0]][e[1]] = i + 1
		adj[e[1]][e[0]] = i + 1
	}

	// For each vertex, track the period of visits. Each path has period = 2*len.
	// A point on path of length L bounces back and forth with period 2L.
	// res[v] = max time stopwatch shows = half the period between consecutive visits.
	// Actually res[v] = the max gap between consecutive arrivals at v (considering all points).

	// For each vertex, collect all arrival times modulo their periods.
	type arrival struct {
		offset float64 // first arrival time
		period float64 // period of recurring arrivals
	}
	arrivals := make([][]arrival, n+1)

	lineIdx := 1
	for p := 0; p < numPaths; p++ {
		if lineIdx >= len(lines) {
			return nil, fmt.Errorf("not enough lines for path %d", p+1)
		}
		fields := strings.Fields(lines[lineIdx])
		lineIdx++
		if len(fields) < 1 {
			return nil, fmt.Errorf("empty path line")
		}
		pathLen, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("cannot parse path len: %v", err)
		}
		if len(fields) < 1+pathLen+2+1 {
			return nil, fmt.Errorf("path %d: not enough fields, got %d need %d", p+1, len(fields), 1+pathLen+3)
		}

		edgeIds := make([]int, pathLen)
		for i := 0; i < pathLen; i++ {
			edgeIds[i], err = strconv.Atoi(fields[1+i])
			if err != nil {
				return nil, fmt.Errorf("cannot parse edge id: %v", err)
			}
		}

		uStr := fields[1+pathLen]
		vStr := fields[1+pathLen+1]
		xStr := fields[1+pathLen+2]

		u, _ := strconv.Atoi(uStr)
		v, _ := strconv.Atoi(vStr)
		x, _ := strconv.ParseFloat(xStr, 64)

		// Reconstruct the path as a sequence of vertices
		// The edges form a path. We need to find the vertex sequence.
		pathVertices := make([]int, 0, pathLen+1)

		// Find the path vertex sequence from edge list
		edgeEndpoints := make([][2]int, pathLen)
		for i, eid := range edgeIds {
			edgeEndpoints[i] = edgesUV[eid-1]
		}

		if pathLen == 1 {
			eu := edgeEndpoints[0][0]
			ev := edgeEndpoints[0][1]
			pathVertices = append(pathVertices, eu, ev)
		} else {
			// Build vertex sequence from edge sequence
			// First edge and second edge share a vertex
			e0 := edgeEndpoints[0]
			e1 := edgeEndpoints[1]
			var start int
			if e0[0] == e1[0] || e0[0] == e1[1] {
				start = e0[1]
			} else {
				start = e0[0]
			}
			pathVertices = append(pathVertices, start)
			for i := 0; i < pathLen; i++ {
				eu := edgeEndpoints[i][0]
				ev := edgeEndpoints[i][1]
				last := pathVertices[len(pathVertices)-1]
				if eu == last {
					pathVertices = append(pathVertices, ev)
				} else {
					pathVertices = append(pathVertices, eu)
				}
			}
		}

		// The point starts on edge (u, v) at distance x from u, moving toward v.
		// Find the position of (u,v) edge in the path
		pointEdgeIdx := -1
		pointForward := true // true if point moves in increasing path index direction
		for i := 0; i < pathLen; i++ {
			pu := pathVertices[i]
			pv := pathVertices[i+1]
			if pu == u && pv == v {
				pointEdgeIdx = i
				pointForward = true
				break
			}
			if pu == v && pv == u {
				pointEdgeIdx = i
				pointForward = false
				break
			}
		}
		if pointEdgeIdx == -1 {
			return nil, fmt.Errorf("path %d: edge (%d,%d) not found in path", p+1, u, v)
		}

		// Distance from path start (vertex 0) along the path
		// Edge i connects pathVertices[i] and pathVertices[i+1], each of length 1
		// Point is on edge pointEdgeIdx at distance x from u
		var distFromStart float64
		if pointForward {
			// u = pathVertices[pointEdgeIdx], v = pathVertices[pointEdgeIdx+1]
			// distance from start = pointEdgeIdx + x
			distFromStart = float64(pointEdgeIdx) + x
		} else {
			// u = pathVertices[pointEdgeIdx+1], v = pathVertices[pointEdgeIdx]
			// distance from u means distance from pathVertices[pointEdgeIdx+1]
			// distance from start = pointEdgeIdx + (1-x)
			distFromStart = float64(pointEdgeIdx) + (1.0 - x)
		}

		L := float64(pathLen)
		period := 2.0 * L

		// Moving toward v means:
		// If pointForward: moving toward higher index (toward end of path)
		// If !pointForward: moving toward lower index (toward start of path)
		movingTowardEnd := pointForward

		// Compute arrival times at each vertex on the path
		for vi := 0; vi <= pathLen; vi++ {
			vertexDist := float64(vi)
			node := pathVertices[vi]

			// Two types of visits: when going forward and when going backward
			// The point bounces at 0 and L.
			// Position at time t: depends on initial position and direction

			// Map to a "phase" position in [0, 2L) that increases with time
			// If moving toward end: phase = distFromStart, increases with time
			// If moving toward start: phase = 2L - distFromStart, increases with time
			var phase0 float64
			if movingTowardEnd {
				phase0 = distFromStart
			} else {
				phase0 = period - distFromStart
			}

			// The point visits vertexDist at phases: vertexDist + 2kL and (2L - vertexDist) + 2kL
			// First arrival going forward: phase = vertexDist
			// First arrival going backward: phase = 2L - vertexDist

			t1 := vertexDist - phase0
			for t1 < 0 {
				t1 += period
			}
			t2 := (period - vertexDist) - phase0
			for t2 < 0 {
				t2 += period
			}

			arrivals[node] = append(arrivals[node], arrival{offset: t1, period: period})
			if vi > 0 && vi < pathLen {
				arrivals[node] = append(arrivals[node], arrival{offset: t2, period: period})
			} else {
				// Endpoint: the bounce means both forward and backward visit coincide
				// Actually at endpoints the point arrives, bounces, and leaves.
				// It visits the endpoint once per period.
				// Already added as t1; t2 would be the same arrival.
				// But let's add t2 as well for correctness - duplicates are harmless
				arrivals[node] = append(arrivals[node], arrival{offset: t2, period: period})
			}
		}
	}

	// For each vertex, compute the maximum gap between consecutive arrivals
	res := make([]float64, n+1)
	for v := 1; v <= n; v++ {
		arrs := arrivals[v]
		if len(arrs) == 0 {
			res[v] = math.Inf(1)
			continue
		}

		// Find LCM of all periods to determine the overall period
		// For small n (<=100), periods are small. Compute arrivals in one full LCM period.
		// Actually, since periods can be different, we need to find the overall period.
		// For simplicity with small n, compute all arrival times in a large enough window.

		// Find LCM of periods
		lcmVal := arrs[0].period
		for _, a := range arrs[1:] {
			lcmVal = lcm(lcmVal, a.period)
		}

		// Collect all arrival times in [0, lcmVal)
		times := make([]float64, 0)
		for _, a := range arrs {
			t := a.offset
			for t >= lcmVal {
				t -= lcmVal
			}
			for t < 0 {
				t += a.period
			}
			for t < lcmVal-1e-12 {
				times = append(times, t)
				t += a.period
			}
		}

		if len(times) == 0 {
			res[v] = math.Inf(1)
			continue
		}

		// Sort and find max gap
		sortFloats(times)
		// Remove near-duplicates
		unique := []float64{times[0]}
		for i := 1; i < len(times); i++ {
			if times[i]-unique[len(unique)-1] > 1e-12 {
				unique = append(unique, times[i])
			}
		}
		times = unique

		maxGap := 0.0
		for i := 1; i < len(times); i++ {
			gap := times[i] - times[i-1]
			if gap > maxGap {
				maxGap = gap
			}
		}
		// Wrap-around gap
		wrapGap := lcmVal - times[len(times)-1] + times[0]
		if wrapGap > maxGap {
			maxGap = wrapGap
		}
		res[v] = maxGap
	}

	return res, nil
}

func gcdFloat(a, b float64) float64 {
	// For our purposes, periods are always integers or simple fractions
	// Convert to rational: period = 2*pathLen, which is always an integer
	// So just use integer GCD
	ai := int(math.Round(a * 1e6))
	bi := int(math.Round(b * 1e6))
	for bi != 0 {
		ai, bi = bi, ai%bi
	}
	return float64(ai) / 1e6
}

func lcm(a, b float64) float64 {
	g := gcdFloat(a, b)
	if g < 1e-12 {
		return a
	}
	return a / g * b
}

func sortFloats(a []float64) {
	// Simple sort for small arrays
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && a[j] < a[j-1]; j-- {
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
}

// ---------- test generation and verification ----------

func genCase(rng *rand.Rand) (string, int, [][2]int) {
	n := rng.Intn(8) + 2
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String(), n, edges
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, n, edges := genCase(rng)

		// Get reference res array from embedded solver
		numPaths, pathEdges, pathEndpoints, pathX := solve822F(input)
		var refOutput strings.Builder
		fmt.Fprintln(&refOutput, numPaths)
		for j := 0; j < numPaths; j++ {
			fmt.Fprintf(&refOutput, "%d", len(pathEdges[j]))
			for _, eid := range pathEdges[j] {
				fmt.Fprintf(&refOutput, " %d", eid)
			}
			fmt.Fprintf(&refOutput, " %d %d %.10f\n", pathEndpoints[j][0], pathEndpoints[j][1], pathX[j])
		}

		edgesUV := make([][2]int, len(edges))
		copy(edgesUV, edges)

		refRes, err := simulateRes(n, edgesUV, refOutput.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: reference simulation error: %v\n", i+1, err)
			os.Exit(1)
		}

		// Run candidate
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}

		gotRes, err := simulateRes(n, edgesUV, got)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: candidate simulation error: %v\noutput:\n%s\ninput:\n%s", i+1, err, got, input)
			os.Exit(1)
		}

		// Compare res arrays with tolerance
		for v := 1; v <= n; v++ {
			diff := math.Abs(refRes[v] - gotRes[v])
			denom := math.Max(1.0, refRes[v])
			if diff/denom > 1e-6 {
				fmt.Fprintf(os.Stderr, "case %d: vertex %d: expected res=%.10f got res=%.10f\ninput:\n%s", i+1, v, refRes[v], gotRes[v], input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
