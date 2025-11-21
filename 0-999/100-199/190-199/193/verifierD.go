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
	pos := make([]int, n+2)
	for i := 1; i <= n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		pos[x] = i
	}
	breaks := []int{0}
	for v := 1; v < n; v++ {
		if abs(pos[v+1]-pos[v]) != 1 {
			breaks = append(breaks, v)
		}
	}
	breaks = append(breaks, n)
	var ans int64
	for i := 0; i+1 < len(breaks); i++ {
		run := breaks[i+1] - breaks[i] - 1
		if run > 0 {
			ans += int64(run) * int64(run+1) / 2
		}
	}
	for i := 1; i+1 < len(breaks); i++ {
		left := breaks[i] - breaks[i-1]
		right := breaks[i+1] - breaks[i]
		ans += int64(left) * int64(right)
	}
	return fmt.Sprintf("%d", ans)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type testCase struct {
	name   string
	input  string
	expect string
}

func makeCase(name string, perm []int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(perm)))
	for i, v := range perm {
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
		makeCase("n1", []int{1}),
		makeCase("sorted_small", []int{1, 2}),
		makeCase("reverse_small", []int{2, 1}),
		makeCase("sample_like1", []int{1, 2, 3}),
		makeCase("sample_like2", []int{2, 1, 3}),
		makeCase("reverse_mid", []int{5, 4, 3, 2, 1}),
		makeCase("zigzag", []int{2, 4, 1, 3, 5}),
		makeCase("shifted", []int{3, 4, 5, 1, 2}),
		makeCase("random_fixed", []int{4, 1, 5, 2, 3, 6}),
		makeCase("longer", []int{3, 1, 4, 2, 5, 7, 6, 8}),
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(193))
	var tests []testCase
	add := func(prefix string, count, minN, maxN int) {
		for i := 0; i < count; i++ {
			n := minN + rng.Intn(maxN-minN+1)
			base := rng.Perm(n)
			for j := range base {
				base[j]++
			}
			tests = append(tests, makeCase(fmt.Sprintf("%s_%d", prefix, i+1), base))
		}
	}
	add("tiny", 120, 1, 6)
	add("small", 120, 7, 12)
	add("medium", 80, 13, 40)
	add("large", 40, 120, 300)
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
