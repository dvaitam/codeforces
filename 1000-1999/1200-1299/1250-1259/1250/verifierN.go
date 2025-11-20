package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// DSU is a Disjoint Set Union data structure.
type DSU struct {
	parent []int
	count  int
}

func NewDSU(n int) *DSU {
	parent := make([]int, n)
	for i := range parent {
		parent[i] = i
	}
	return &DSU{parent: parent, count: n}
}

func (d *DSU) Find(i int) int {
	if d.parent[i] == i {
		return i
	}
	d.parent[i] = d.Find(d.parent[i])
	return d.parent[i]
}

func (d *DSU) Union(i, j int) {
	rootI := d.Find(i)
	rootJ := d.Find(j)
	if rootI != rootJ {
		d.parent[rootI] = rootJ
		d.count--
	}
}

type Wire struct {
	id, u, v int
}

type Move struct {
	wID, oldP, newP int
}

// validate checks if the user's output is a correct solution for a given testcase.
func validate(testCaseInput string, userOutput string) error {
	// --- 1. Parse the initial wires from the test case input ---
	reader := bufio.NewReader(strings.NewReader(testCaseInput))
	var n int
	fmt.Fscan(reader, &n)
	initialWires := make([]Wire, n)
	allPoints := make(map[int]bool)
	for i := 0; i < n; i++ {
		initialWires[i].id = i + 1
		fmt.Fscan(reader, &initialWires[i].u, &initialWires[i].v)
		allPoints[initialWires[i].u] = true
		allPoints[initialWires[i].v] = true
	}

	// --- 2. Parse the user's output (k and k moves) ---
	outReader := bufio.NewReader(strings.NewReader(userOutput))
	var k int
	_, err := fmt.Fscan(outReader, &k)
	if err != nil {
		return fmt.Errorf("could not parse k (number of moves): %v", err)
	}
	if k < 0 {
		return fmt.Errorf("number of moves k cannot be negative: %d", k)
	}
	moves := make([]Move, k)
	for i := 0; i < k; i++ {
		_, err := fmt.Fscan(outReader, &moves[i].wID, &moves[i].oldP, &moves[i].newP)
		if err == io.EOF {
			return fmt.Errorf("incorrect number of moves provided. Expected %d, got %d", k, i)
		}
		if err != nil {
			return fmt.Errorf("could not parse move %d: %v", i+1, err)
		}
		allPoints[moves[i].newP] = true // Add new points for coordinate compression
	}

	// --- 3. Coordinate Compression ---
	sortedPoints := make([]int, 0, len(allPoints))
	for p := range allPoints {
		sortedPoints = append(sortedPoints, p)
	}
	coordMap := make(map[int]int)
	for i, p := range sortedPoints {
		coordMap[p] = i
	}
	
	numPoints := len(coordMap)
    if numPoints == 0 && n == 0 { // Empty case
        if k == 0 {
            return nil
        }
        return fmt.Errorf("expected 0 moves for 0 wires, but got %d", k)
    }
	if numPoints == 0 && n > 0 {
		return fmt.Errorf("no points found for %d wires", n)
	}

	// --- 4. Calculate initial connected components ---
	dsu := NewDSU(numPoints)
	for _, w := range initialWires {
		u, v := coordMap[w.u], coordMap[w.v]
		dsu.Union(u, v)
	}
	initialComponents := dsu.count

	// --- 5. Verify number of moves ---
	if k != initialComponents-1 {
		return fmt.Errorf("wrong number of moves. Expected %d, got %d", initialComponents-1, k)
	}
	if k == 0 && initialComponents == 1 { // Already connected
		return nil
	}

	// --- 6. Simulate and validate moves ---
	currentWires := make(map[int]struct{ u, v int })
	for _, w := range initialWires {
		currentWires[w.id] = struct{ u, v int }{w.u, w.v}
	}

	finalDSU := NewDSU(numPoints)
	// Re-apply initial connections to the second DSU
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

	// --- 7. Check final connectivity ---
	if finalDSU.count != 1 {
		return fmt.Errorf("graph is not fully connected after moves. %d components remain", finalDSU.count)
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println(`usage: go run verifierN.go /path/to/binary`)
		os.Exit(1)
	}
	bin := os.Args[1]

	testCasesData, err := os.ReadFile("testcasesN.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read testcases: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(testCasesData))
	
	scanner.Scan()
	numTestCases, _ := strconv.Atoi(scanner.Text())

	fmt.Printf("Running %d test cases...\n", numTestCases)

	for i := 0; i < numTestCases; i++ {
		
		var testCaseBuffer bytes.Buffer
		
		scanner.Scan()
		nStr := scanner.Text()
		if nStr == "" { // Handle potential blank lines between test cases
			i--
			continue
		}
		
		n, _ := strconv.Atoi(nStr)
		testCaseBuffer.WriteString(nStr + "\n")

		for j := 0; j < n; j++ {
			scanner.Scan()
			testCaseBuffer.WriteString(scanner.Text() + "\n")
		}
		
		testCaseInput := testCaseBuffer.String()

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(testCaseInput)
		userOutput, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, `[FAIL] Test case %d: execution failed: %v\nOutput:\n%s\n`, i+1, err, userOutput)
			os.Exit(1)
		}

		err = validate(testCaseInput, string(userOutput))
		if err != nil {
			fmt.Fprintf(os.Stderr, `[FAIL] Test case %d: validation failed: %v\n`, i+1, err)
			fmt.Fprintf(os.Stderr, `---
Input ---
%s
--- Your Output ---
%s
`, testCaseInput, userOutput)
			os.Exit(1)
		}
		fmt.Printf("[PASS] Test case %d\n", i+1)
	}

	fmt.Println("\nAll tests passed")
}
