package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// DSU structure for union find (mirrors 1250N solution).
type DSU struct {
	parent []int
	rank   []int
	count  int
}

func NewDSU(n int) *DSU {
	d := &DSU{
		parent: make([]int, n),
		rank:   make([]int, n),
		count:  n,
	}
	for i := 0; i < n; i++ {
		d.parent[i] = i
	}
	return d
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(a, b int) {
	a = d.Find(a)
	b = d.Find(b)
	if a == b {
		return
	}
	if d.rank[a] < d.rank[b] {
		a, b = b, a
	}
	d.parent[b] = a
	if d.rank[a] == d.rank[b] {
		d.rank[a]++
	}
	d.count--
}

type Wire struct {
	id, u, v int
}

type Move struct {
	wID, oldP, newP int
}

// solve is the 1250N reference logic (unused in validation but embedded for completeness).
func solve(n int, x, y []int64) [][3]int64 {
	contactFirst := make(map[int64]int)
	dsu := NewDSU(n)
	contactMap := make(map[int64][]int)

	for i := 0; i < n; i++ {
		if j, ok := contactFirst[x[i]]; ok {
			dsu.Union(i, j)
		} else {
			contactFirst[x[i]] = i
		}
		if j, ok := contactFirst[y[i]]; ok {
			dsu.Union(i, j)
		} else {
			contactFirst[y[i]] = i
		}
		contactMap[x[i]] = append(contactMap[x[i]], i)
		contactMap[y[i]] = append(contactMap[y[i]], i)
	}

	compWires := make(map[int][]int)
	for i := 0; i < n; i++ {
		r := dsu.Find(i)
		compWires[r] = append(compWires[r], i)
	}

	rootRep := dsu.Find(0)
	rootContact := x[0]

	operations := make([][3]int64, 0)

	for rep, wires := range compWires {
		if rep == rootRep {
			continue
		}
		chosenWire := wires[0]
		oldContact := x[chosenWire]

		for _, w := range wires {
			if len(contactMap[x[w]]) == 1 {
				chosenWire = w
				oldContact = x[w]
				break
			}
			if len(contactMap[y[w]]) == 1 {
				chosenWire = w
				oldContact = y[w]
				break
			}
		}

		operations = append(operations, [3]int64{int64(chosenWire + 1), oldContact, rootContact})
	}
	return operations
}

// validate checks if the user's output is a correct solution for a given testcase.
func validate(testCaseInput string, userOutput string) error {
	reader := bufio.NewReader(strings.NewReader(testCaseInput))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return fmt.Errorf("failed to read n: %v", err)
	}

	initialWires := make([]Wire, n)
	allPoints := make(map[int]bool)
	for i := 0; i < n; i++ {
		initialWires[i].id = i + 1
		if _, err := fmt.Fscan(reader, &initialWires[i].u, &initialWires[i].v); err != nil {
			return fmt.Errorf("failed to read wire %d: %v", i+1, err)
		}
		allPoints[initialWires[i].u] = true
		allPoints[initialWires[i].v] = true
	}

	outReader := bufio.NewReader(strings.NewReader(userOutput))
	var k int
	if _, err := fmt.Fscan(outReader, &k); err != nil {
		return fmt.Errorf("could not parse k (number of moves): %v", err)
	}
	if k < 0 {
		return fmt.Errorf("number of moves k cannot be negative: %d", k)
	}
	moves := make([]Move, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(outReader, &moves[i].wID, &moves[i].oldP, &moves[i].newP); err != nil {
			if err == io.EOF {
				return fmt.Errorf("incorrect number of moves provided. Expected %d, got %d", k, i)
			}
			return fmt.Errorf("could not parse move %d: %v", i+1, err)
		}
		allPoints[moves[i].newP] = true
	}

	sortedPoints := make([]int, 0, len(allPoints))
	for p := range allPoints {
		sortedPoints = append(sortedPoints, p)
	}
	coordMap := make(map[int]int)
	for i, p := range sortedPoints {
		coordMap[p] = i
	}

	numPoints := len(coordMap)
	if numPoints == 0 && n == 0 {
		if k == 0 {
			return nil
		}
		return fmt.Errorf("expected 0 moves for 0 wires, but got %d", k)
	}
	if numPoints == 0 && n > 0 {
		return fmt.Errorf("no points found for %d wires", n)
	}

	dsu := NewDSU(numPoints)
	for _, w := range initialWires {
		u, v := coordMap[w.u], coordMap[w.v]
		dsu.Union(u, v)
	}
	initialComponents := dsu.count

	if k != initialComponents-1 {
		return fmt.Errorf("wrong number of moves. Expected %d, got %d", initialComponents-1, k)
	}
	if k == 0 && initialComponents == 1 {
		return nil
	}

	currentWires := make(map[int]struct{ u, v int })
	for _, w := range initialWires {
		currentWires[w.id] = struct{ u, v int }{w.u, w.v}
	}

	finalDSU := NewDSU(numPoints)
	for _, w := range initialWires {
		finalDSU.Union(coordMap[w.u], coordMap[w.v])
	}

	for _, m := range moves {
		wire, ok := currentWires[m.wID]
		if !ok {
			return fmt.Errorf("move specifies invalid wire ID %d", m.wID)
		}

		var otherP int
		if wire.u == m.oldP {
			otherP = wire.v
		} else if wire.v == m.oldP {
			otherP = wire.u
		} else {
			return fmt.Errorf("wire %d does not have an endpoint at %d", m.wID, m.oldP)
		}

		if m.newP == otherP {
			return fmt.Errorf("move for wire %d creates a self-loop from %d to %d", m.wID, otherP, m.newP)
		}

		currentWires[m.wID] = struct{ u, v int }{otherP, m.newP}
		finalDSU.Union(coordMap[otherP], coordMap[m.newP])
	}

	if finalDSU.count != 1 {
		return fmt.Errorf("graph is not fully connected after moves. %d components remain", finalDSU.count)
	}

	return nil
}

func runCandidate(bin string, input string) ([]byte, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	return cmd.CombinedOutput()
}

func loadTestcases() ([]string, error) {
	reader := bufio.NewReader(strings.NewReader(testcaseData))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("failed to read test count: %v", err)
	}
	tests := make([]string, 0, t)
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return nil, fmt.Errorf("test %d: failed to read n: %v", i+1, err)
		}
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			var u, v int
			if _, err := fmt.Fscan(reader, &u, &v); err != nil {
				return nil, fmt.Errorf("test %d: failed to read wire %d: %v", i+1, j+1, err)
			}
			sb.WriteString(strconv.Itoa(u))
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(v))
			sb.WriteByte('\n')
		}
		tests = append(tests, sb.String())
	}
	if len(tests) != t {
		return nil, fmt.Errorf("testcase count mismatch: expected %d, got %d", t, len(tests))
	}
	return tests, nil
}

const testcaseData = `
100
5
20 11
4 3
13 3
18 6
1 16
5
8 6
14 14
10 20
13 14
16 3
2
10 3
4 15
4
10 15
19 11
18 8
9 6
4
1 14
4 6
8 1
16 16
5
12 12
17 19
19 15
11 9
8 2
5
20 11
3 6
15 11
5 20
18 14
4
3 5
16 5
3 14
9 15
1
18 11
5
14 2
17 10
19 9
19 1
19 5
3
5 10
15 15
20 16
3
1 16
4 5
6 4
1
8 14
1
7 13
2
3 17
7 10
5
17 20
18 7
20 13
5 3
4 20
4
19 8
10 19
1 4
15 19
4
12 20
14 2
2 7
17 4
5
2 5
12 11
7 1
7 9
5 4
2
12 4
1 1
5
3 9
9 14
19 16
11 5
12 4
5
17 7
7 18
6 5
19 10
5 5
3
20 2
7 7
6 13
3
10 18
8 2
5 5
2
13 1
19 15
1
14 4
5
2 12
14 8
4 4
18 15
5 5
5
14 15
9 18
2 20
11 13
20 11
2
12 12
18 14
4
7 15
12 20
13 13
7 5
4
8 8
7 10
1 2
14 18
4
4 11
2 10
20 7
7 8
1
9 18
1
3 2
2
10 13
4 14
2
15 2
6 17
3
20 19
6 9
8 1
5
12 15
1 4
2 19
11 4
17 12
4
13 13
9 9
13 7
19 1
4
5 13
9 10
17 9
20 5
4
17 3
19 2
12 16
10 8
1
15 18
2
20 18
5 6
3
11 1
16 16
7 3
2
19 8
14 5
4
1 13
6 7
20 17
3 13
1
18 19
1
6 17
4
10 10
18 18
16 4
19 1
1
7 20
3
2 6
17 19
15 20
1
17 11
3
17 12
11 9
20 12
4
18 7
14 2
11 12
14 4
2
5 16
18 15
3
18 12
18 6
2 8
3
11 3
15 5
14 17
1
7 3
5
6 3
9 13
7 17
17 3
11 4
1
19 6
2
1 5
5 6
3
6 7
19 20
2 20
2
11 11
14 17
5
17 6
19 11
17 10
5 13
15 4
2
9 6
12 15
3
8 2
15 13
1 10
2
17 3
20 8
3
8 6
9 3
11 6
5
19 3
4 3
18 3
19 7
13 8
5
12 4
3 20
17 11
18 11
20 14
5
13 8
15 2
3 3
1 15
20 12
1
10 4
1
5 7
2
18 2
16 4
2
7 15
8 19
2
7 8
7 6
2
3 8
6 15
1
18 10
4
12 5
6 6
9 9
20 10
5
10 8
5 6
1 2
15 6
16 5
1
11 12
4
7 3
9 7
6 18
16 3
4
19 15
10 2
13 12
1 9
2
7 3
8 5
4
15 4
3 20
13 10
19 18
3
19 20
18 16
14 16
5
6 8
16 2
4 20
11 6
4 14
1
16 5
4
6 2
17 8
17 6
3 20
4
6 6
11 9
12 14
8 10
2
17 19
20 14
1
8 6
2
17 11
10 1
5
14 7
10 11
20 16
6 20
8 17
1
3 18
5
18 14
2 9
11 2
1 16
13 13
2
14 20
12 5
4
10 12
12 3
12 6
3 12
5
17 8
13 2
4 7
5 13
6 7
1
16 8
`

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println(`usage: go run verifierN.go /path/to/binary`)
		os.Exit(1)
	}
	bin := args[0]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Running %d test cases...\n", len(tests))
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%s", tc)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[FAIL] Test case %d: execution failed: %v\nOutput:\n%s\n", i+1, err, out)
			os.Exit(1)
		}
		if err := validate(tc, string(out)); err != nil {
			fmt.Fprintf(os.Stderr, "[FAIL] Test case %d: validation failed: %v\n---\nInput:\n%s---\nYour Output:\n%s\n", i+1, err, tc, out)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed")
}
