package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	name  string
	input string
}

func solveRef(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var s, b int
	if _, err := fmt.Fscan(reader, &s, &b); err != nil {
		return "", err
	}
	a := make([]int, s)
	for i := 0; i < s; i++ {
		fmt.Fscan(reader, &a[i])
	}
	bases := make([][2]int64, b)
	for i := 0; i < b; i++ {
		var d, g int64
		fmt.Fscan(reader, &d, &g)
		bases[i][0] = d
		bases[i][1] = g
	}
	sort.Slice(bases, func(i, j int) bool {
		return bases[i][0] < bases[j][0]
	})
	defs := make([]int64, b)
	prefix := make([]int64, b)
	for i := 0; i < b; i++ {
		defs[i] = bases[i][0]
		if i == 0 {
			prefix[i] = bases[i][1]
		} else {
			prefix[i] = prefix[i-1] + bases[i][1]
		}
	}
	var sb strings.Builder
	for i := 0; i < s; i++ {
		attack := int64(a[i])
		idx := sort.Search(len(defs), func(j int) bool {
			return defs[j] > attack
		})
		var gold int64
		if idx > 0 {
			gold = prefix[idx-1]
		}
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", gold)
	}
	return sb.String(), nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	if stderr.Len() > 0 {
		return "", fmt.Errorf("unexpected stderr output: %s", stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func compareOutputs(expect, got string) bool {
	return strings.TrimSpace(expect) == strings.TrimSpace(got)
}

func makeCase(name string, s, b int, ships []int, bases [][2]int64) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", s, b)
	for i := 0; i < s; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", ships[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < b; i++ {
		fmt.Fprintf(&sb, "%d %d\n", bases[i][0], bases[i][1])
	}
	return testCase{name: name, input: sb.String()}
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("single_zero", 1, 1, []int{0}, [][2]int64{{0, 5}}),
		makeCase("simple", 3, 3, []int{1, 10, 3}, [][2]int64{{1, 5}, {5, 7}, {3, 2}}),
		makeCase("no_bases", 2, 0, []int{10, 20}, [][2]int64{}),
		makeCase("strong_ships", 2, 3, []int{100, 100}, [][2]int64{{10, 5}, {20, 5}, {30, 5}}),
		makeCase("weak_ships", 2, 3, []int{0, 1}, [][2]int64{{2, 10}, {3, 20}, {4, 30}}),
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for i := 0; i < 100; i++ {
		s := rng.Intn(10) + 1
		b := rng.Intn(10) + 1
		ships := make([]int, s)
		for j := 0; j < s; j++ {
			ships[j] = rng.Intn(50)
		}
		bases := make([][2]int64, b)
		for j := 0; j < b; j++ {
			bases[j][0] = int64(rng.Intn(50))
			bases[j][1] = int64(rng.Intn(100))
		}
		tests = append(tests, makeCase(fmt.Sprintf("rand_%d", i+1), s, b, ships, bases))
	}
	return tests
}

func parseOutput(output string, expectedCount int) ([]string, error) {
	fields := strings.Fields(output)
	if len(fields) != expectedCount {
		return nil, fmt.Errorf("expected %d numbers, got %d", expectedCount, len(fields))
	}
	return fields, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		expect, err := solveRef(tc.input)
		if err != nil {
			fmt.Printf("failed to compute expected result for %s: %v\n", tc.name, err)
			os.Exit(1)
		}
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expFields, _ := parseOutput(expect, len(strings.Fields(expect)))
		gotFields, err := parseOutput(out, len(expFields))
		if err != nil {
			fmt.Printf("test %d (%s) invalid output: %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
		if strings.Join(expFields, " ") != strings.Join(gotFields, " ") {
			fmt.Printf("test %d (%s) mismatch\ninput:\n%s\nexpect:%s\nactual:%s\n", idx+1, tc.name, tc.input, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
