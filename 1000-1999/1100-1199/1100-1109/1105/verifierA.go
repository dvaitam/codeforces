package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	input string
	n     int
	a     []int
}

func genTests() []Test {
	r := rand.New(rand.NewSource(0))
	tests := make([]Test, 0, 105)
	// random tests
	for i := 0; i < 100; i++ {
		n := r.Intn(20) + 1
		a := make([]int, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			val := r.Intn(100) + 1
			a[j] = val
			fmt.Fprintf(&sb, "%d", val)
		}
		sb.WriteByte('\n')
		tests = append(tests, Test{sb.String(), n, a})
	}
	// edge cases
	tests = append(tests, 
		Test{"1\n1\n", 1, []int{1}},
		Test{"1\n100\n", 1, []int{100}},
		Test{"3\n1 1 1\n", 3, []int{1, 1, 1}},
		Test{"4\n10 20 30 40\n", 4, []int{10, 20, 30, 40}},
		Test{"5\n100 99 98 97 96\n", 5, []int{100, 99, 98, 97, 96}},
		Test{"1\n12\n", 1, []int{12}},
	)
	return tests
}

func solve(a []int) (int, int) {
	minCost := -1
	bestT := -1
	// Range of t can be [1, 100] as a_i are in [1, 100].
	// Costs are convex-ish, checking 1 to 100 is safe and sufficient.
	for t := 1; t <= 100; t++ {
		cost := 0
		for _, val := range a {
			dist := val - t
			if dist < 0 {
				dist = -dist
			}
			if dist > 1 {
				cost += dist - 1
			}
		}
		if minCost == -1 || cost < minCost {
			minCost = cost
			bestT = t
		}
	}
	return bestT, minCost
}

func calculateCost(a []int, t int) int {
	cost := 0
	for _, val := range a {
		dist := val - t
		if dist < 0 {
			dist = -dist
		}
		if dist > 1 {
			cost += dist - 1
		}
	}
	return cost
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	tests := genTests()
	for i, tc := range tests {
		// Determine expected minimum cost
		_, minCost := solve(tc.a)

		// Run candidate
		gotOutput, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}

		// Parse candidate output
		var gotT, gotCost int
		_, err = fmt.Fscan(strings.NewReader(gotOutput), &gotT, &gotCost)
		if err != nil {
			fmt.Printf("Test %d failed: invalid output format\nInput:\n%sGot:\n%s\n", i+1, tc.input, gotOutput)
			os.Exit(1)
		}

		// Check if candidate cost is optimal
		if gotCost != minCost {
			fmt.Printf("Test %d failed: non-optimal cost\nInput:\n%sExpected Min Cost: %d\nGot Cost: %d\nOutput: %s\n", i+1, tc.input, minCost, gotCost, gotOutput)
			os.Exit(1)
		}

		// Check if candidate t yields the claimed cost
		realCost := calculateCost(tc.a, gotT)
		if realCost != gotCost {
			fmt.Printf("Test %d failed: claimed cost doesn't match t\nInput:\n%sCandidate T: %d\nClaimed Cost: %d\nCalculated Cost: %d\n", i+1, tc.input, gotT, gotCost, realCost)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}