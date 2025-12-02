package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Pair represents an undirected edge
type Pair struct {
	u, v int
}

// DirectedEdge represents a directed edge
type DirectedEdge struct {
	u, v int
}

const testcases = `100
4 3
2 4
1 3
3 4
3 1
1 2
5 4
1 3
2 4
1 5
3 4
1 0
6 2
3 6
3 4
6 7
5 6
1 6
1 5
3 4
1 4
2 6
1 3
2 1
1 2
1 0
3 0
3 0
5 9
3 4
1 2
1 3
1 4
4 5
2 3
2 5
1 5
3 5
4 1
3 4
2 0
5 8
1 5
1 4
2 5
4 5
3 5
2 4
1 2
1 3
6 11
1 2
5 6
1 3
1 4
2 5
3 5
2 3
3 4
1 5
4 5
1 6
1 0
5 1
1 5
6 0
1 0
5 3
3 5
1 5
2 4
2 0
4 4
2 3
3 4
1 4
2 4
1 0
1 0
1 0
3 0
5 9
3 5
2 5
2 3
1 4
4 5
1 5
2 4
1 3
1 2
2 0
1 0
2 1
1 2
5 9
4 5
2 5
1 4
1 5
1 2
3 5
2 4
3 4
1 3
5 7
1 5
1 2
4 5
1 3
1 4
3 5
3 4
6 4
4 5
5 6
4 6
3 4
6 12
4 5
2 5
2 6
5 6
3 4
1 4
2 3
1 5
1 2
1 3
3 5
4 6
6 14
2 6
1 6
2 3
1 4
2 5
1 5
1 3
4 6
4 5
3 6
1 2
2 4
3 4
3 5
5 6
2 4
2 5
1 3
2 3
1 4
1 5
2 0
6 6
1 3
5 6
2 5
1 6
3 6
3 4
3 0
4 6
3 4
1 4
1 3
2 4
1 2
2 3
1 0
2 1
1 2
1 0
3 0
2 1
1 2
5 10
3 4
1 3
1 4
1 5
2 3
3 5
2 5
1 2
4 5
2 4
4 4
3 4
2 3
1 2
1 3
3 1
1 3
3 2
1 2
2 3
6 9
1 5
3 6
2 4
4 6
2 5
2 6
5 6
1 4
1 6
1 0
6 15
3 4
3 5
4 6
5 6
2 3
3 6
2 5
2 4
4 5
1 3
2 6
1 4
1 6
1 2
1 5
4 5
2 4
2 3
1 3
1 2
3 4
4 1
3 4
6 12
4 6
1 4
1 6
2 5
1 2
2 6
2 3
1 5
5 6
3 6
1 3
3 5
5 3
1 4
2 4
1 2
3 3
1 3
1 2
2 3
2 1
1 2
1 0
6 3
3 6
1 6
4 6
2 1
1 2
4 1
3 4
5 8
1 3
1 2
3 5
2 3
1 5
1 4
2 4
3 4
5 4
1 3
4 5
3 5
2 3
5 1
1 3
2 0
1 0
1 0
3 3
1 3
2 3
1 2
3 0
5 7
1 2
2 3
2 5
1 4
3 4
1 5
2 4
4 5
2 3
1 4
1 2
1 3
3 4
2 1
1 2
5 6
2 5
3 4
1 3
3 5
1 4
1 5
1 0
6 11
3 5
3 6
4 6
2 6
2 4
1 3
4 5
2 5
3 4
1 2
2 3
4 0
1 0
2 0
1 0
6 10
3 4
4 6
1 4
1 6
1 2
4 5
2 5
2 4
1 5
5 6
6 14
4 5
1 5
2 5
5 6
1 2
2 6
2 3
4 6
2 4
1 4
1 6
3 4
3 6
3 5
3 2
2 3
1 3
6 3
2 3
2 4
1 4
2 1
1 2
5 7
3 4
1 4
2 3
4 5
2 4
1 2
1 5
4 5
2 4
1 4
1 3
1 2
2 3
3 0
1 0
6 0
3 1
2 3
2 0
4 1
3 4
1 0
3 2
1 2
2 3
3 1
1 2
5 8
1 4
2 4
1 2
3 4
2 5
1 5
3 5
2 3
2 0
2 1
1 2
3 1
2 3
`

func solveOptimal(n int, edges []Pair) int {
	deg := make([]int, n+1)
	for _, e := range edges {
		deg[e.u]++
		deg[e.v]++
	}
	odd := 0
	for i := 1; i <= n; i++ {
		if deg[i]%2 != 0 {
			odd++
		}
	}
	return n - odd
}

func runCase(bin string, input string, n, m int, originalEdges []Pair) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out // Capture stderr too just in case

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}

	// Parse output
	scanner := bufio.NewScanner(&out)
	scanner.Split(bufio.ScanWords)

	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	userK, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return fmt.Errorf("invalid K: %v", err)
	}

	// We expect m edges
	userEdges := make([]DirectedEdge, 0, m)
	for i := 0; i < m; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("expected %d edges, got %d", m, i)
		}
		uStr := scanner.Text()
		u, err := strconv.Atoi(uStr)
		if err != nil {
			return fmt.Errorf("invalid u at edge %d: %v", i, err)
		}

		if !scanner.Scan() {
			return fmt.Errorf("incomplete edge %d", i)
		}
		vStr := scanner.Text()
		v, err := strconv.Atoi(vStr)
		if err != nil {
			return fmt.Errorf("invalid v at edge %d: %v", i, err)
		}
		userEdges = append(userEdges, DirectedEdge{u, v})
	}

	// 1. Verify K is optimal
	optimalK := solveOptimal(n, originalEdges)
	if userK != optimalK {
		return fmt.Errorf("claimed K=%d, optimal is %d", userK, optimalK)
	}

	// 2. Verify edges match input (multiset equality)
	// Map canonical undirected edge to count
	edgeCounts := make(map[Pair]int)
	for _, e := range originalEdges {
		u, v := e.u, e.v
		if u > v {
			u, v = v, u
		}
		edgeCounts[Pair{u, v}]++
	}

	for _, e := range userEdges {
		u, v := e.u, e.v
		if u > v {
			u, v = v, u
		}
		p := Pair{u, v}
		if edgeCounts[p] > 0 {
			edgeCounts[p]--
		} else {
			return fmt.Errorf("edge %d-%d returned by user not in input or used too many times", e.u, e.v)
		}
	}

	// 3. Calculate in/out degrees and verify K count
	inDeg := make([]int, n+1)
	outDeg := make([]int, n+1)
	for _, e := range userEdges {
		outDeg[e.u]++
		inDeg[e.v]++
	}

	actualK := 0
	for i := 1; i <= n; i++ {
		if inDeg[i] == outDeg[i] {
			actualK++
		}
	}

	if actualK != userK {
		return fmt.Errorf("user claimed K=%d, but actual edges give K=%d", userK, actualK)
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	// Read test cases from string
	scanner := bufio.NewScanner(strings.NewReader(testcases))
	scanner.Split(bufio.ScanWords)

	if !scanner.Scan() {
		fmt.Println("Empty test file")
		os.Exit(1)
	}
	tStr := scanner.Text()
	t, _ := strconv.Atoi(tStr)

	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			break
		}
		n, _ := strconv.Atoi(scanner.Text())
		if !scanner.Scan() {
			break
		}
		m, _ := strconv.Atoi(scanner.Text())

		edges := make([]Pair, m)
		for j := 0; j < m; j++ {
			scanner.Scan()
			u, _ := strconv.Atoi(scanner.Text())
			scanner.Scan()
			v, _ := strconv.Atoi(scanner.Text())
			edges[j] = Pair{u, v}
		}

		// Construct input string for the binary
		var sb strings.Builder
		sb.WriteString("1\n") // We feed one case at a time to the binary
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
		}

		if err := runCase(bin, sb.String(), n, m, edges); err != nil {
			fmt.Printf("Test case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed")
}
