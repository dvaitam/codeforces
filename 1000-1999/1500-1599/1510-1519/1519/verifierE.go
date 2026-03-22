package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE.txt so the verifier is self-contained.
const testcasesRaw = `100
1
1
3
6
2 6 6 3 3 5
7 20 2 19 6 14
7
7 6 7 5 3 5 4
17 9 2 1 12 15 11
7
4 5 2 5 2 2 2
1 6 11 6 5 17 17
6
5 6 5 2 4 4
17 12 19 12 12 15
3
2 3 3
15 17 8
8
5 8 6 8 8 6 8 8
8 11 6 20 9 16 10 10
7
3 6 2 4 5 3 6
20 3 11 1 7 4 2
1
1
19
4
1 2 3 2
7 2 14 2
1
1
12
3
1 3 1
3 4 3
1
1
1
6
3 2 2 6 2 5
1 13 19 2 8 5
1
1
12
2
2 2
16 1
5
4 5 5 1 3
13 20 5 16 8
2
2 1
1 15
3
3 3 2
16 17 11
3
2 2 2
20 14 1
3
3 1 2
2 5 6
3
1 2 3
8 17 2
4
2 4 1 3
3 19 8 20
6
3 6 4 3 5 1
5 2 13 14 6 4
2
1 1
4 1
3
1 1 1
1 17 15
8
5 7 4 4 7 7 1 1
14 17 19 6 4 16 12 1
2
2 2
12 10
1
1
4
2
2 1
1 15
1
1
16
8
4 2 1 5 1 6 5 2
8 16 7 4 19 12 13 15
3
2 2 1
9 4 4
2
2 2
7 4
1
1
2
8
5 6 8 3 6 5 8 8
14 16 10 13 8 6 16 20
5
5 4 1 5 5
4 3 12 6 18
3
2 1 1
2 5 10
7
2 6 6 6 3 4 2
17 10 4 5 18 14 4
6
5 2 6 5 3 2
6 15 8 13 12 19
3
2 2 3
1 20 13
3
2 3 1
16 9 13
5
4 4 3 5 3
3 8 18 20 7
7
7 6 4 6 1 3 4
17 15 6 4 1 13 7
7
2 1 4 7 5 7 7
7 9 19 19 7 16 20
3
1 3 3
14 16 9
3
2 3 1
3 12 1
8
2 8 6 8 5 8 1 2
20 12 6 13 9 5 2 6
8
7 8 5 3 1 5 8 1
12 2 18 13 19 15 7 10
8
3 8 5 2 5 6 5 6
10 13 17 3 17 7 13 20
3
3 3 1
10 2 8
8
4 5 1 2 2 7 6 4
11 12 3 11 15 12 6 16
8
5 8 3 8 4 5 6 3
4 8 16 7 12 6 12 5
3
1 2 3
13 13 11
5
5 5 5 3 4
10 18 20 3 12
5
4 4 2 3 3
15 16 3 6 11
7
2 1 1 3 2 3 1
14 1 18 11 8 20 13
5
4 2 3 3 2
16 4 5 7 11
5
2 4 3 3 1
11 7 8 8 20
1
1
12
1
1
6
2
2 2
9 5
6
5 5 1 3 6 6
20 13 8 2 13 16
8
6 2 8 7 8 3 7 7
17 15 2 4 15 19 5 4
3
1 2 2
15 1 9
2
2 1
6 1
3
2 3 1
11 15 2
8
4 2 8 3 1 3 1 1
7 18 1 17 11 17 8 5
6
4 1 2 5 1 2
4 15 7 2 20 7
7
3 5 6 7 4 6 5
17 6 17 4 5 7 6
7
2 3 3 4 2 4 2
13 11 10 4 18 4 16
5
3 5 4 3 2
14 5 18 4 1
4
2 2 4 1
5 1 9 16
1
1
5
6
1 6 2 1 2 6
18 6 3 15 10 7
3
2 3 2
17 19 3
7
4 6 6 1 4 3 6
4 9 1 7 14 11 9
7
5 5 6 2 4 7 2
6 15 15 12 13 16 20
5
5 2 5 4 4
7 16 19 11 10
2
1 2
20 16
4
2 3 2 3
4 1 1 7
6
1 3 5 3 6 6
11 15 3 14 16 1
5
5 5 2 2 2
6 20 13 3 19
8
5 2 8 8 4 3 5 4
7 20 11 19 20 13 17 14
4
2 1 3 2
5 20 13 14
2
2 2
13 16
7
3 2 2 2 1 5 5
3 20 18 1 2 13 14
7
2 5 3 1 3 5 3
17 16 19 3 15 8 9
1
1
16
1
1
5
4
3 2 1 2
10 4 18 18
2
1 2
5 2
5
5 3 4 1 5
12 11 4 20 12
2
2 2
9 16
5
5 5 2 1 1
11 14 1 12 18
1
1
18
7
4 4 2 7 2 2 5
2 1 19 12 6 10 1
1
1
19
4
4 1 3 1
20 3 8 8
4
1 1 4 1
17 9 19 8
4
1 2 2 1
1 4 5 3
5 5 1 1
4 1 1 2
2
4 4 5 5
5 3 2 5
4
3 2 1 2
4 2 2 5
5 3 2 1
5 1 3 2
2
4 3 2 2
1 1 1 1
2
3 3 3 3
5 4 2 2
2
2 5 1 5
5 5 4 1
1
5 5 5 5`

type Frac struct {
	n int64
	d int64
}

type Frame struct {
	v     int
	pEdge int
	idx   int
	pool  []int
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// computeExpected uses the correct iterative solver from the accepted solution.
func computeExpected() (string, error) {
	fields := strings.Fields(testcasesRaw)
	pos := 0
	nextInt := func() int64 {
		v, _ := strconv.ParseInt(fields[pos], 10, 64)
		pos++
		return v
	}

	n := int(nextInt())

	idMap := make(map[Frac]int, 2*n)
	adj := make([][]int, 0, 2*n)

	getID := func(f Frac) int {
		if id, ok := idMap[f]; ok {
			return id
		}
		id := len(adj)
		idMap[f] = id
		adj = append(adj, nil)
		return id
	}

	eu := make([]int, n+1)
	ev := make([]int, n+1)

	for i := 1; i <= n; i++ {
		a := nextInt()
		b := nextInt()
		c := nextInt()
		d := nextInt()

		num1 := c * b
		den1 := d * (a + b)
		g1 := gcd(absInt64(num1), absInt64(den1))
		if g1 == 0 {
			g1 = 1
		}
		f1 := Frac{num1 / g1, den1 / g1}

		num2 := b * (c + d)
		den2 := a * d
		g2 := gcd(absInt64(num2), absInt64(den2))
		if g2 == 0 {
			g2 = 1
		}
		f2 := Frac{num2 / g2, den2 / g2}

		u := getID(f1)
		v := getID(f2)

		eu[i] = u
		ev[i] = v
		adj[u] = append(adj[u], i)
		adj[v] = append(adj[v], i)
	}

	visV := make([]bool, len(adj))
	visE := make([]bool, n+1)
	ans := make([][2]int, 0, n/2)
	stack := make([]Frame, 0, 64)

	for root := 0; root < len(adj); root++ {
		if visV[root] {
			continue
		}
		visV[root] = true
		stack = append(stack[:0], Frame{v: root, pEdge: -1})

		for len(stack) > 0 {
			top := len(stack) - 1
			if stack[top].idx < len(adj[stack[top].v]) {
				eid := adj[stack[top].v][stack[top].idx]
				stack[top].idx++
				if visE[eid] {
					continue
				}
				visE[eid] = true

				to := eu[eid]
				if to == stack[top].v {
					to = ev[eid]
				}

				if !visV[to] {
					visV[to] = true
					stack = append(stack, Frame{v: to, pEdge: eid})
				} else {
					stack[top].pool = append(stack[top].pool, eid)
				}
			} else {
				pool := stack[top].pool
				for len(pool) >= 2 {
					e1 := pool[len(pool)-1]
					e2 := pool[len(pool)-2]
					pool = pool[:len(pool)-2]
					ans = append(ans, [2]int{e1, e2})
				}

				pEdge := stack[top].pEdge
				usedParent := false
				if pEdge != -1 && len(pool) == 1 {
					ans = append(ans, [2]int{pool[0], pEdge})
					usedParent = true
				}

				stack = stack[:top]
				if len(stack) > 0 && pEdge != -1 && !usedParent {
					parent := len(stack) - 1
					stack[parent].pool = append(stack[parent].pool, pEdge)
				}
			}
		}
	}

	var out bytes.Buffer
	fmt.Fprintln(&out, len(ans))
	for _, p := range ans {
		fmt.Fprintf(&out, "%d %d\n", p[0], p[1])
	}

	return strings.TrimSpace(out.String()), nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// validateOutput checks that the candidate output is a valid matching.
// It verifies: correct count, valid pairs, no repeated edges, and each pair is
// geometrically valid (the two points can be moved onto a line through the origin).
func validateOutput(expected, got string) bool {
	// First compare the count line
	expLines := strings.Split(expected, "\n")
	gotLines := strings.Split(got, "\n")
	if len(expLines) == 0 || len(gotLines) == 0 {
		return false
	}
	// The count (first line) must match
	if strings.TrimSpace(expLines[0]) != strings.TrimSpace(gotLines[0]) {
		return false
	}
	expCount, err := strconv.Atoi(strings.TrimSpace(expLines[0]))
	if err != nil {
		return false
	}
	if len(gotLines) < expCount+1 {
		return false
	}

	// Parse all input data to validate pairs
	fields := strings.Fields(testcasesRaw)
	pos := 0
	nextInt := func() int64 {
		v, _ := strconv.ParseInt(fields[pos], 10, 64)
		pos++
		return v
	}
	n := int(nextInt())

	type PointData struct {
		a, b, c, d int64
	}
	points := make([]PointData, n+1)
	for i := 1; i <= n; i++ {
		points[i] = PointData{nextInt(), nextInt(), nextInt(), nextInt()}
	}

	// Check each pair in got
	used := make([]bool, n+1)
	for i := 1; i <= expCount; i++ {
		parts := strings.Fields(gotLines[i])
		if len(parts) != 2 {
			return false
		}
		p1, err1 := strconv.Atoi(parts[0])
		p2, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return false
		}
		if p1 < 1 || p1 > n || p2 < 1 || p2 > n || p1 == p2 {
			return false
		}
		if used[p1] || used[p2] {
			return false
		}
		used[p1] = true
		used[p2] = true

		// Validate geometric condition:
		// Point i has coords (a_i/b_i, c_i/d_i).
		// After moving +1 to either x or y for each point, both must lie on a line through origin.
		// A point on line through origin satisfies y/x = const (slope).
		// For point i moved in x: new = (a_i/b_i + 1, c_i/d_i), slope = c_i*b_i / (d_i*(a_i+b_i))
		// For point i moved in y: new = (a_i/b_i, c_i/d_i + 1), slope = b_i*(c_i+d_i) / (a_i*d_i)
		// Two points can be paired if they share a slope option.
		canPair := func(i, j int) bool {
			pi := points[i]
			pj := points[j]
			// slopes for point i
			// slope_ix = c_i*b_i / (d_i*(a_i+b_i))
			// slope_iy = b_i*(c_i+d_i) / (a_i*d_i)
			// slopes for point j similarly
			// Compare as fractions: a/b == c/d iff a*d == b*c
			type slope struct{ num, den int64 }
			slopesI := [2]slope{
				{pi.c * pi.b, pi.d * (pi.a + pi.b)},
				{pi.b * (pi.c + pi.d), pi.a * pi.d},
			}
			slopesJ := [2]slope{
				{pj.c * pj.b, pj.d * (pj.a + pj.b)},
				{pj.b * (pj.c + pj.d), pj.a * pj.d},
			}
			for _, si := range slopesI {
				for _, sj := range slopesJ {
					if si.num*sj.den == si.den*sj.num {
						return true
					}
				}
			}
			return false
		}

		if !canPair(p1, p2) {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	expected, err := computeExpected()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	got, err := runCandidate(bin, testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if strings.TrimSpace(got) == expected {
		fmt.Println("All tests passed")
		return
	}

	// If not exact match, validate the output structurally
	if validateOutput(expected, got) {
		fmt.Println("All tests passed")
		return
	}

	fmt.Printf("verification failed\nexpected:\n%s\n\ngot:\n%s\n", expected, got)
	os.Exit(1)
}
