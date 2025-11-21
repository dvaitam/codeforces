package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	refSourceE = "737E.go"
	refBinaryE = "ref737E.bin"
	totalTests = 60
)

type childReq struct {
	machine int
	time    int
}

type testCase struct {
	n int
	m int
	b int
	p []int
	k [][]childReq
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := generateTests()
	for idx, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}
		refAns, refMask, refSegments, err := parseOutput(refOut, tc)
		if err != nil {
			fmt.Printf("failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			printInput(input)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}
		candAns, candMask, candSegments, err := parseOutput(candOut, tc)
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			printInput(input)
			os.Exit(1)
		}

		if candAns != refAns {
			fmt.Printf("test %d failed: expected minimal time %d, got %d\n", idx+1, refAns, candAns)
			printInput(input)
			fmt.Println("Reference output:")
			fmt.Println(refOut)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}

		if len(candMask) != tc.m {
			fmt.Printf("test %d failed: expected mask length %d, got %d\n", idx+1, tc.m, len(candMask))
			printInput(input)
			os.Exit(1)
		}
		if !validateMaskBudget(candMask, tc.p, tc.b) {
			fmt.Printf("test %d failed: candidate mask exceeds budget\n")
			printInput(input)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}

		if err := validateSchedule(tc, candAns, candMask, candSegments); err != nil {
			fmt.Printf("test %d schedule invalid: %v\n", idx+1, err)
			printInput(input)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinaryE, refSourceE)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryE), nil
}

func runProgram(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.b))
	for i := 0; i < tc.m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(tc.p[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < tc.n; i++ {
		sb.WriteString(strconv.Itoa(len(tc.k[i])))
		for _, req := range tc.k[i] {
			sb.WriteString(fmt.Sprintf(" %d %d", req.machine, req.time))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, tc testCase) (int, []int, []segment, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	readLine := func() (string, error) {
		for {
			line, err := reader.ReadString('\n')
			if len(line) == 0 && err != nil {
				return "", err
			}
			line = strings.TrimSpace(line)
			if line == "" {
				if err != nil {
					return "", err
				}
				continue
			}
			return line, nil
		}
	}
	first, err := readLine()
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to read time line: %v", err)
	}
	ans, err := strconv.Atoi(first)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("invalid time value: %v", err)
	}
	maskLine, err := readLine()
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to read mask: %v", err)
	}
	if len(maskLine) != tc.m {
		return 0, nil, nil, fmt.Errorf("invalid mask length")
	}
	mask := make([]int, tc.m)
	for i := 0; i < tc.m; i++ {
		switch maskLine[i] {
		case '0':
			mask[i] = 0
		case '1':
			mask[i] = 1
		default:
			return 0, nil, nil, fmt.Errorf("mask contains invalid character")
		}
	}
	gLine, err := readLine()
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to read segment count: %v", err)
	}
	g, err := strconv.Atoi(gLine)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("invalid segment count: %v", err)
	}
	segs := make([]segment, 0, g)
	for i := 0; i < g; i++ {
		line, err := readLine()
		if err != nil {
			return 0, nil, nil, fmt.Errorf("failed to read segment line: %v", err)
		}
		fields := strings.Fields(line)
		if len(fields) != 4 {
			return 0, nil, nil, fmt.Errorf("segment line must have 4 integers")
		}
		child, err := strconv.Atoi(fields[0])
		if err != nil {
			return 0, nil, nil, fmt.Errorf("invalid child index: %v", err)
		}
		machine, err := strconv.Atoi(fields[1])
		if err != nil {
			return 0, nil, nil, fmt.Errorf("invalid machine index: %v", err)
		}
		start, err := strconv.Atoi(fields[2])
		if err != nil {
			return 0, nil, nil, fmt.Errorf("invalid start time: %v", err)
		}
		duration, err := strconv.Atoi(fields[3])
		if err != nil {
			return 0, nil, nil, fmt.Errorf("invalid duration: %v", err)
		}
		segs = append(segs, segment{
			child:   child,
			machine: machine,
			start:   start,
			length:  duration,
		})
	}
	return ans, mask, segs, nil
}

func validateMaskBudget(mask []int, prices []int, budget int) bool {
	total := 0
	for i := 0; i < len(mask); i++ {
		if mask[i] == 1 {
			total += prices[i]
		}
	}
	return total <= budget
}

type segment struct {
	child   int
	machine int
	start   int
	length  int
}

func validateSchedule(tc testCase, timeLimit int, mask []int, segs []segment) error {
	sumReq := make([]int, tc.n)
	for i := 0; i < tc.n; i++ {
		for _, req := range tc.k[i] {
			sumReq[i] += req.time
		}
	}
	need := make([]map[int]int, tc.n)
	for i := 0; i < tc.n; i++ {
		need[i] = make(map[int]int)
		for _, req := range tc.k[i] {
			need[i][req.machine] = req.time
		}
	}
	childSegs := make([][]segment, tc.n)
	for i := 0; i < tc.n; i++ {
		childSegs[i] = make([]segment, 0)
	}
	machineSegs := make([][]segment, tc.m)
	for j := 0; j < tc.m; j++ {
		machineSegs[j] = make([]segment, 0)
	}
	for _, seg := range segs {
		if seg.child < 1 || seg.child > tc.n {
			return fmt.Errorf("segment has invalid child index %d", seg.child)
		}
		if seg.machine < 1 || seg.machine > tc.m {
			return fmt.Errorf("segment has invalid machine index %d", seg.machine)
		}
		if seg.start < 0 || seg.start+seg.length > timeLimit {
			return fmt.Errorf("segment exceeds time limit")
		}
		if seg.length <= 0 {
			return fmt.Errorf("segment length must be positive")
		}
		childIdx := seg.child - 1
		machineIdx := seg.machine - 1
		need[childIdx][seg.machine] -= seg.length
		childSegs[childIdx] = append(childSegs[childIdx], seg)
		machineSegs[machineIdx] = append(machineSegs[machineIdx], seg)
	}
	for i := 0; i < tc.n; i++ {
		sort.Slice(childSegs[i], func(a, b int) bool {
			if childSegs[i][a].start == childSegs[i][b].start {
				return childSegs[i][a].machine < childSegs[i][b].machine
			}
			return childSegs[i][a].start < childSegs[i][b].start
		})
		curEnd := 0
		for _, seg := range childSegs[i] {
			if seg.start < curEnd {
				return fmt.Errorf("child %d overlaps segments", i+1)
			}
			curEnd = seg.start + seg.length
		}
	}
	for j := 0; j < tc.m; j++ {
		events := make(map[int]int)
		for _, seg := range machineSegs[j] {
			events[seg.start]++
			events[seg.start+seg.length]--
		}
		if len(events) == 0 {
			continue
		}
		times := make([]int, 0, len(events))
		for t := range events {
			times = append(times, t)
		}
		sort.Ints(times)
		running := 0
		copies := 1 + mask[j]
		for _, t := range times {
			running += events[t]
			if running > copies {
				return fmt.Errorf("machine %d exceeds available copies at time %d", j+1, t)
			}
		}
	}
	for i := 0; i < tc.n; i++ {
		for machine, rem := range need[i] {
			if rem != 0 {
				return fmt.Errorf("child %d missing %d minutes on machine %d", i+1, rem, machine)
			}
		}
	}
	return nil
}

func generateTests() []testCase {
	tests := []testCase{
		buildTestCase(1, 1, 0, []int{1}, [][]childReq{{{1, 3}}}),
		buildTestCase(2, 2, 10, []int{3, 4}, [][]childReq{{{1, 5}}, {{2, 6}}}),
		buildTestCase(3, 3, 5, []int{4, 3, 2}, [][]childReq{
			{{1, 3}, {2, 2}},
			{{2, 4}},
			{{3, 5}},
		}),
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-5 {
		tests = append(tests, randomTest(rnd))
	}
	tests = append(tests,
		randomTest(rand.New(rand.NewSource(1))),
		randomTest(rand.New(rand.NewSource(2))),
		randomTest(rand.New(rand.NewSource(3))),
		randomTest(rand.New(rand.NewSource(4))),
		randomTest(rand.New(rand.NewSource(5))),
	)
	return tests
}

func buildTestCase(n, m, b int, prices []int, reqs [][]childReq) testCase {
	return testCase{
		n: n,
		m: m,
		b: b,
		p: append([]int(nil), prices...),
		k: copyReqs(reqs),
	}
}

func copyReqs(reqs [][]childReq) [][]childReq {
	out := make([][]childReq, len(reqs))
	for i := range reqs {
		out[i] = append([]childReq(nil), reqs[i]...)
	}
	return out
}

func randomTest(rnd *rand.Rand) testCase {
	n := rnd.Intn(8) + 1
	m := rnd.Intn(4) + 1
	b := rnd.Intn(100)
	prices := make([]int, m)
	for i := 0; i < m; i++ {
		prices[i] = rnd.Intn(10) + 1
	}
	reqs := make([][]childReq, n)
	for i := 0; i < n; i++ {
		used := make(map[int]bool)
		count := rnd.Intn(m) + 1
		reqs[i] = make([]childReq, 0, count)
		for len(reqs[i]) < count {
			machine := rnd.Intn(m) + 1
			if used[machine] {
				continue
			}
			used[machine] = true
			reqs[i] = append(reqs[i], childReq{
				machine: machine,
				time:    rnd.Intn(5) + 1,
			})
		}
	}
	return buildTestCase(n, m, b, prices, reqs)
}

func printInput(in string) {
	fmt.Println("Input used:")
	fmt.Print(in)
}
