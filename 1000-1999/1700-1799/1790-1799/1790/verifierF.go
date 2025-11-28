package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Generate a random tree with n nodes and return edges
func genTree(n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	// Create a basic tree structure
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	
	// Shuffle node labels to avoid bias
	perm := rand.Perm(n)
	// map 1..n to perm[0]+1 .. perm[n-1]+1
	mapping := make([]int, n+1)
	for i, v := range perm {
		mapping[i+1] = v + 1
	}
	
	shuffledEdges := make([][2]int, len(edges))
	for i, e := range edges {
		shuffledEdges[i] = [2]int{mapping[e[0]], mapping[e[1]]}
	}
	
	return shuffledEdges
}

func genCase() string {
	n := rand.Intn(50) + 2 // Random n between 2 and 51
	c0 := rand.Intn(n) + 1

	// Generate permutation of other nodes
	perm := rand.Perm(n)
	c := make([]int, 0, n-1)
	for _, v := range perm {
		val := v + 1
		if val != c0 {
			c = append(c, val)
		}
	}

	edges := genTree(n)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, c0))
	for i, v := range c {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteString("\n")
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func runProgram(path string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func main() {
	rand.Seed(time.Now().UnixNano())

	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]

	_, src, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(src)
	refPath := filepath.Join(baseDir, "1790F.go")

	// Run 20 random test cases
	const numTests = 20
	for i := 1; i <= numTests; i++ {
		caseStr := genCase()
		// Wrap in "1\n" for t=1
		fullInput := "1\n" + caseStr

		refOut, err := runProgram(refPath, fullInput)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Reference failed on test %d: %v\n", i, err)
			os.Exit(1)
		}

		targetOut, err := runProgram(target, fullInput)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Target failed on test %d: %v\n", i, err)
			os.Exit(1)
		}

		// Simple whitespace-agnostic comparison
		refFields := strings.Fields(refOut)
		targetFields := strings.Fields(targetOut)

		if len(refFields) != len(targetFields) {
			fmt.Fprintf(os.Stderr, "Test %d failed: length mismatch. Expected %d, got %d\nInput:\n%s\nExpected:\n%s\nGot:\n%s\n", i, len(refFields), len(targetFields), fullInput, refOut, targetOut)
			os.Exit(1)
		}

		for j, rf := range refFields {
			if rf != targetFields[j] {
				fmt.Fprintf(os.Stderr, "Test %d failed at token %d. Expected %s, got %s\nInput:\n%s\n", i, j, rf, targetFields[j], fullInput)
				os.Exit(1)
			}
		}
		// fmt.Printf("Test %d passed\n", i)
	}
	fmt.Printf("All %d tests passed\n", numTests)
}