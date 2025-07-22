package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCase struct {
	input  string
	output string
}

func solve(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(in, &n)
	angles := make([]float64, n)
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		deg := math.Atan2(float64(y), float64(x)) * 180.0 / math.Pi
		if deg < 0 {
			deg += 360.0
		}
		angles[i] = deg
	}
	sort.Float64s(angles)
	maxGap := 0.0
	for i := 1; i < n; i++ {
		gap := angles[i] - angles[i-1]
		if gap > maxGap {
			maxGap = gap
		}
	}
	if n > 0 {
		wrapGap := 360.0 - angles[n-1] + angles[0]
		if wrapGap > maxGap {
			maxGap = wrapGap
		}
	}
	result := 360.0 - maxGap
	return fmt.Sprintf("%.8f\n", result)
}

func generateTests() []testCase {
	rand.Seed(44)
	var tests []testCase
	fixed := []string{
		"1\n1 0\n",
		"2\n1 0\n0 1\n",
		"3\n1 0\n0 1\n-1 0\n",
		"4\n1 0\n0 1\n-1 0\n0 -1\n",
	}
	for _, f := range fixed {
		tests = append(tests, testCase{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rand.Intn(20) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		used := make(map[[2]int]bool)
		for i := 0; i < n; i++ {
			var x, y int
			for {
				x = rand.Intn(2001) - 1000
				y = rand.Intn(2001) - 1000
				if x != 0 || y != 0 {
					if !used[[2]int{x, y}] {
						used[[2]int{x, y}] = true
						break
					}
				}
			}
			fmt.Fprintf(&sb, "%d %d\n", x, y)
		}
		inp := sb.String()
		tests = append(tests, testCase{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.output) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected: %sGot: %s\n", i+1, t.input, strings.TrimSpace(t.output), strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
