package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// ===== Embedded reference solver for 164D =====

type refPoint struct {
	x, y int64
}

func refDist2(p1, p2 refPoint) int64 {
	dx := p1.x - p2.x
	dy := p1.y - p2.y
	return dx*dx + dy*dy
}

type int64Slice []int64

func (x int64Slice) Len() int           { return len(x) }
func (x int64Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x int64Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func solveVC(adj [][]int, active []bool, k int) []int {
	n := len(adj)
	deg := make([]int, n)
	maxDeg := 0
	maxU := -1
	edgeCount := 0

	for i := 0; i < n; i++ {
		if !active[i] || len(adj[i]) == 0 {
			continue
		}
		d := 0
		for _, v := range adj[i] {
			if active[v] {
				d++
			}
		}
		deg[i] = d
		edgeCount += d
		if d > maxDeg {
			maxDeg = d
			maxU = i
		}
	}
	edgeCount /= 2

	if edgeCount == 0 {
		return []int{}
	}
	if k == 0 {
		return nil
	}
	if edgeCount > k*maxDeg {
		return nil
	}

	if maxDeg <= 2 {
		leaf := -1
		for i := 0; i < n; i++ {
			if active[i] && deg[i] == 1 {
				leaf = i
				break
			}
		}
		if leaf != -1 {
			parent := -1
			for _, v := range adj[leaf] {
				if active[v] {
					parent = v
					break
				}
			}
			active[parent] = false
			res := solveVC(adj, active, k-1)
			active[parent] = true
			if res != nil {
				return append(res, parent)
			}
			return nil
		}

		u := -1
		for i := 0; i < n; i++ {
			if active[i] && deg[i] == 2 {
				u = i
				break
			}
		}
		if u != -1 {
			active[u] = false
			res := solveVC(adj, active, k-1)
			active[u] = true
			if res != nil {
				return append(res, u)
			}
		}
		return nil
	}

	if maxDeg > k {
		active[maxU] = false
		res := solveVC(adj, active, k-1)
		active[maxU] = true
		if res != nil {
			return append(res, maxU)
		}
		return nil
	}

	active[maxU] = false
	res1 := solveVC(adj, active, k-1)
	active[maxU] = true
	if res1 != nil {
		return append(res1, maxU)
	}

	neighbors := make([]int, 0, deg[maxU])
	for _, v := range adj[maxU] {
		if active[v] {
			neighbors = append(neighbors, v)
		}
	}
	if k >= len(neighbors) {
		for _, v := range neighbors {
			active[v] = false
		}
		res2 := solveVC(adj, active, k-len(neighbors))
		for _, v := range neighbors {
			active[v] = true
		}
		if res2 != nil {
			return append(res2, neighbors...)
		}
	}

	return nil
}

func refCheck(D2 int64, pts []refPoint, k int) []int {
	n := len(pts)
	edgeCount := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if refDist2(pts[i], pts[j]) > D2 {
				edgeCount++
			}
		}
	}
	if edgeCount > k*(n-1) {
		return nil
	}

	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if refDist2(pts[i], pts[j]) > D2 {
				adj[i] = append(adj[i], j)
				adj[j] = append(adj[j], i)
			}
		}
	}

	active := make([]bool, n)
	for i := range active {
		active[i] = true
	}

	return solveVC(adj, active, k)
}

func refSolve(input string) string {
	fields := strings.Fields(input)
	pos := 0
	nextInt := func() int {
		v, _ := strconv.Atoi(fields[pos])
		pos++
		return v
	}

	n := nextInt()
	k := nextInt()
	if n == 0 {
		return ""
	}

	pts := make([]refPoint, n)
	for i := 0; i < n; i++ {
		pts[i].x = int64(nextInt())
		pts[i].y = int64(nextInt())
	}

	dists := make([]int64, 0, n*(n-1)/2)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dists = append(dists, refDist2(pts[i], pts[j]))
		}
	}
	sort.Sort(int64Slice(dists))

	uniqueDists := make([]int64, 0, len(dists)+1)
	uniqueDists = append(uniqueDists, 0)
	if len(dists) > 0 {
		if dists[0] != 0 {
			uniqueDists = append(uniqueDists, dists[0])
		}
		for i := 1; i < len(dists); i++ {
			if dists[i] != dists[i-1] {
				uniqueDists = append(uniqueDists, dists[i])
			}
		}
	}

	low := 0
	high := len(uniqueDists) - 1
	var bestRes []int

	for low <= high {
		mid := low + (high-low)/2
		D2 := uniqueDists[mid]

		res := refCheck(D2, pts, k)
		if res != nil {
			bestRes = res
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	deleted := make([]bool, n)
	for _, idx := range bestRes {
		deleted[idx] = true
	}

	var finalAns []int
	for _, idx := range bestRes {
		finalAns = append(finalAns, idx+1)
	}

	for i := 0; i < n && len(finalAns) < k; i++ {
		if !deleted[i] {
			finalAns = append(finalAns, i+1)
			deleted[i] = true
		}
	}

	parts := make([]string, len(finalAns))
	for i, v := range finalAns {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, " ")
}

// ===== Verifier infrastructure =====

type point struct{ x, y int }

type testCase struct {
	n int
	k int
	p []point
}

const encodedTestcases = `
H4sIAGQLK2kC/01WWXYjMQj89yl0BO3L/S821EJPnhPHaQGCoig8Systfm5ps8TbKQufWv+NOKkl3mvZ5cVJPOnxGjBqD49+YSj7pT+Xf2gI1xFG4VJ/O+wieC8HVr3Cv20GYoABI3jF84PIjcdw7HBsjLzkHglFhEqj9VuZ5oTLilwbw4RVpIwEGy/A+YQPymk3TyYuVSZ088FA/EPXKo8XN3Xfgqoafhu8UM5giY8OHZ5Vj/H+dE9AimqA9LMly3/M7sIyHjXBwGRxzQEKk3iqKGDdkMpgpkd+um4z1KB96wqBS+Sy2NL48PiCQdbEcpBp9ucJPPS70UvVxoEBb0wHZT0nFCkIuxNPRZU4G2WroiHqiF2blmwHwm+2gRUwx8HkzAFksoVEgKQbdfzxJBKiHwi4EWey2C7gb7wmydpF5k6+bLe/EoR8SgjY1kHGslxkVosowoQP+Xg9ES2bXU2f68QYkugIhMczcvcxneWZ4OBt3qQ0RfSbROaMDpXkbrLxHq1NO9EcQac/qZebNSyEQ2NX5ORJ7eq/qJ7NYd289uq4e1qr2wzkNu4cOi/VmPQs7n/ozTo0SDdzQDi3Y7CxBOkp3WxT/DPhdj2E5oyAvZ4yKclggWTR1+PmPJaJ0asRxRvtyZ4UhQavQQJkp66gRhjyr5M6SPkm1mwfwd6WuO6I1fycIg5VcYrEi+PRTrZ3imIwXGLO8ACAdYTRpH2a9WWVXgIEaHBEe049c5mco2laexYpXJeQPFFFHbofzbuV1hPB6WK8lvqiqjS8ZuKRshCJgefHPdmFqwLF1uwjlY7CnaL2WPi2DscNvDBD5vP2PCdNmlTV9Wp1NT8npFDLyiGAW9TCB+uPZH+avj3XT5NIolTJaJNgU1EwEMX96i7mahuY41VdI1vIwWmZHanam0RPjdjaWBqMGEvhPDihEjZp5iWxubP6X34sZZXL7HmLc70M8tLkX9RAEEwr6FCfG/Ev0thRVMBITVvWlEBFeTUxtEv9hJEZKUi62MOqq6PPMnKg2ehpkakC35v2JQrd68DTwnWnLeljrZhtJdrFkKdUiS7uyNC0UTq2+nw1lsuEG6aDW/qkxzWFmer03VizzY96pQVD2h6tVVGsPXWZ30QkvVKY861w79cAlghcFm0JsAC8r91e7c+AjqJtCIhXNvxInioJM9XjV7wj2O0qgWhwXDnPWlrScm5RjYCk7efh4Ze7hjpzN2sloYb5+wewoS9gPwoAAA==
`

func decodeTestcases() ([]testCase, error) {
	data, err := base64.StdEncoding.DecodeString(encodedTestcases)
	if err != nil {
		return nil, err
	}
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var out bytes.Buffer
	if _, err := out.ReadFrom(r); err != nil {
		return nil, err
	}
	fields := strings.Fields(out.String())
	if len(fields) == 0 {
		return nil, fmt.Errorf("no data")
	}
	pos := 0
	var cases []testCase
	for pos+2 <= len(fields) {
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, err
		}
		k, err := strconv.Atoi(fields[pos+1])
		if err != nil {
			return nil, err
		}
		pos += 2
		if pos+2*n > len(fields) {
			break
		}
		p := make([]point, n+1)
		for j := 1; j <= n; j++ {
			x, _ := strconv.Atoi(fields[pos])
			y, _ := strconv.Atoi(fields[pos+1])
			p[j] = point{x: x, y: y}
			pos += 2
		}
		cases = append(cases, testCase{n: n, k: k, p: p})
	}
	return cases, nil
}

func tcToInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
	for i := 1; i <= tc.n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", tc.p[i].x, tc.p[i].y)
	}
	return sb.String()
}

func runCandidate(bin string, tc testCase) (string, error) {
	input := tcToInput(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := decodeTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := tcToInput(tc)
		expect := refSolve(input)

		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
