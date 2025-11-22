package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	coordLimit = 1_000_000
	refBin     = "./2122C_ref.bin"
)

type point struct {
	x int64
	y int64
}

type pair struct {
	a int
	b int
}

type testCase struct {
	points []point
}

type testInput struct {
	name  string
	cases []testCase
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBin, "2122C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refBin, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(t testInput) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(t.cases)))
	sb.WriteByte('\n')
	for _, tc := range t.cases {
		sb.WriteString(strconv.Itoa(len(tc.points)))
		sb.WriteByte('\n')
		for _, p := range tc.points {
			sb.WriteString(strconv.FormatInt(p.x, 10))
			sb.WriteByte(' ')
			sb.WriteString(strconv.FormatInt(p.y, 10))
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func parseOutput(out string, t testInput) ([][]pair, error) {
	tokens := strings.Fields(out)
	idx := 0
	result := make([][]pair, 0, len(t.cases))

	readInt := func(tok string) (int, error) {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return 0, err
		}
		if v < 1 {
			return 0, fmt.Errorf("index %d out of range", v)
		}
		return int(v) - 1, nil
	}

	for caseIdx, tc := range t.cases {
		n := len(tc.points)
		if n%2 != 0 {
			return nil, fmt.Errorf("test %d: n is not even", caseIdx+1)
		}
		expect := n
		if idx+expect > len(tokens) {
			return nil, fmt.Errorf("test %d: expected %d integers, only %d tokens remain", caseIdx+1, expect, len(tokens)-idx)
		}
		pairs := make([]pair, n/2)
		used := make([]bool, n)
		for i := 0; i < n/2; i++ {
			aVal, err := readInt(tokens[idx+2*i])
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid index %q (%v)", caseIdx+1, tokens[idx+2*i], err)
			}
			bVal, err := readInt(tokens[idx+2*i+1])
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid index %q (%v)", caseIdx+1, tokens[idx+2*i+1], err)
			}
			if aVal < 0 || aVal >= n || bVal < 0 || bVal >= n {
				return nil, fmt.Errorf("test %d: index out of bounds in pair %d", caseIdx+1, i+1)
			}
			if aVal == bVal {
				return nil, fmt.Errorf("test %d: self-pairing at position %d", caseIdx+1, i+1)
			}
			if used[aVal] {
				return nil, fmt.Errorf("test %d: index %d used multiple times", caseIdx+1, aVal+1)
			}
			if used[bVal] {
				return nil, fmt.Errorf("test %d: index %d used multiple times", caseIdx+1, bVal+1)
			}
			used[aVal], used[bVal] = true, true
			pairs[i] = pair{a: aVal, b: bVal}
		}
		idx += expect
		for j, u := range used {
			if !u {
				return nil, fmt.Errorf("test %d: index %d unused", caseIdx+1, j+1)
			}
		}
		result = append(result, pairs)
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("extra output detected (%d tokens)", len(tokens)-idx)
	}
	return result, nil
}

func manhattanSum(tc testCase, pairs []pair) int64 {
	var res int64
	for _, pr := range pairs {
		p1 := tc.points[pr.a]
		p2 := tc.points[pr.b]
		dx := p1.x - p2.x
		if dx < 0 {
			dx = -dx
		}
		dy := p1.y - p2.y
		if dy < 0 {
			dy = -dy
		}
		res += dx + dy
	}
	return res
}

func compareOutputs(t testInput, refPairs, candPairs [][]pair) error {
	if len(refPairs) != len(candPairs) {
		return fmt.Errorf("expected %d test cases, got %d", len(refPairs), len(candPairs))
	}
	for i := range refPairs {
		refSum := manhattanSum(t.cases[i], refPairs[i])
		candSum := manhattanSum(t.cases[i], candPairs[i])
		if candSum != refSum {
			return fmt.Errorf("test %d: mismatched total distance (expected %d, got %d)", i+1, refSum, candSum)
		}
	}
	return nil
}

func sampleInput() testInput {
	points1 := []point{{1, 1}, {3, 0}, {4, 2}, {3, 4}}
	points2 := []point{{10, -1}, {-1, 2}, {-2, -2}, {-2, 0}, {0, 2}, {2, -3}, {-4, -4}, {-4, 0}, {0, 1}, {-4, -2}}
	return testInput{
		name:  "sample",
		cases: []testCase{{points: points1}, {points: points2}},
	}
}

func smallDegenerate() testInput {
	points1 := []point{{0, 0}, {0, 0}}
	points2 := []point{{5, 5}, {5, 5}, {5, 5}, {5, 5}}
	points3 := []point{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}
	return testInput{
		name:  "degenerate",
		cases: []testCase{{points: points1}, {points: points2}, {points: points3}},
	}
}

func randomPoints(rng *rand.Rand, n int) []point {
	pts := make([]point, n)
	for i := 0; i < n; i++ {
		pts[i] = point{
			x: rng.Int63n(2*coordLimit+1) - coordLimit,
			y: rng.Int63n(2*coordLimit+1) - coordLimit,
		}
	}
	return pts
}

func randomInput(rng *rand.Rand, name string, cases int, nMin, nMax int) testInput {
	all := make([]testCase, 0, cases)
	for i := 0; i < cases; i++ {
		n := rng.Intn(nMax-nMin+1) + nMin
		if n%2 == 1 {
			n++
		}
		all = append(all, testCase{points: randomPoints(rng, n)})
	}
	return testInput{name: name, cases: all}
}

func alignedInput(rng *rand.Rand) testInput {
	n := 20
	pts := make([]point, n)
	for i := 0; i < n; i++ {
		pts[i] = point{x: int64(i - 10), y: 0}
	}
	return testInput{name: "aligned", cases: []testCase{{points: pts}, {points: randomPoints(rng, n)}}}
}

func largeStress(rng *rand.Rand) testInput {
	n := 200000
	return testInput{name: "stress", cases: []testCase{{points: randomPoints(rng, n)}}}
}

func buildTests() []testInput {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []testInput{
		sampleInput(),
		smallDegenerate(),
		alignedInput(rng),
		randomInput(rng, "random-small", 5, 2, 12),
		randomInput(rng, "random-mid", 4, 50, 120),
		randomInput(rng, "random-large", 2, 500, 2000),
		largeStress(rng),
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := buildTests()
	for idx, t := range tests {
		input := buildInput(t)

		expRaw, err := runProgram(refPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\n", idx+1, t.name, err)
			os.Exit(1)
		}
		expPairs, err := parseOutput(expRaw, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, t.name, err, expRaw)
			os.Exit(1)
		}

		gotRaw, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s): runtime error: %v\ninput:\n%s", idx+1, t.name, err, input)
			os.Exit(1)
		}
		gotPairs, err := parseOutput(gotRaw, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s): invalid output: %v\ninput:\n%soutput:\n%s\n", idx+1, t.name, err, input, gotRaw)
			os.Exit(1)
		}

		if err := compareOutputs(t, expPairs, gotPairs); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, t.name, err, input, expRaw, gotRaw)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
