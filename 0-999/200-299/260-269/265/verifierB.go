package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveRef(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return ""
	}
	h := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &h[i])
	}
	total := int64(0)
	curH := 0
	for i := 0; i < n; i++ {
		if i > 0 {
			total++
		}
		if h[i] > curH {
			total += int64(h[i] - curH)
			curH = h[i]
		}
		total++
		if i < n-1 && curH > h[i+1] {
			total += int64(curH - h[i+1])
			curH = h[i+1]
		}
	}
	return fmt.Sprintf("%d", total)
}

type testCase struct {
	name   string
	input  string
	expect string
}

func makeCase(name string, heights []int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(heights)))
	for i, v := range heights {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	return testCase{name: name, input: input, expect: solveRef(input)}
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("single_tree", []int{1}),
		makeCase("equal_heights", []int{3, 3, 3}),
		makeCase("strictly_increasing", []int{1, 2, 3, 4}),
		makeCase("strictly_decreasing", []int{5, 4, 3, 2, 1}),
		makeCase("alternating", []int{1, 3, 2, 4, 2}),
		makeCase("tall_and_small", []int{10, 1}),
		makeCase("large_jump_down", []int{2, 10, 1}),
		makeCase("flat_then_peak", []int{1, 1, 1, 10}),
		makeCase("peak_middle", []int{2, 8, 2, 8, 2}),
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(265))
	var tests []testCase
	add := func(prefix string, count, minN, maxN, maxH int) {
		for i := 0; i < count; i++ {
			n := minN + rng.Intn(maxN-minN+1)
			heights := make([]int, n)
			for j := 0; j < n; j++ {
				heights[j] = rng.Intn(maxH) + 1
			}
			tests = append(tests, makeCase(fmt.Sprintf("%s_%d", prefix, i+1), heights))
		}
	}
	add("small", 120, 1, 6, 10)
	add("medium", 120, 7, 50, 100)
	add("large", 60, 100, 400, 10000)
	add("very_large", 10, 50000, 100000, 10000)
	return tests
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		got, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if got != tc.expect {
			fmt.Printf("test %d (%s) failed\ninput:\n%s\nexpect:%s\nactual:%s\n", idx+1, tc.name, tc.input, tc.expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
