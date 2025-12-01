package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	// refSource points to the local reference solution to avoid GOPATH resolution.
	refSource      = "1998A.go"
	coordLimit     = int64(1_000_000_000)
	centerMin      = -100
	centerMax      = 100
	maxTotalKInput = 1000
)

type caseSpec struct {
	xc, yc int64
	k      int
}

type testSet struct {
	name  string
	cases []caseSpec
}

func main() {
	candPath, ok := parseBinaryArg()
	if !ok {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}

	refBin, cleanupRef, err := buildBinary(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := buildBinary(candPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := buildTests()
	for idx, ts := range tests {
		input := buildInput(ts)

		refOut, err := runBinary(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference solution failed on test %d (%s): %v\ninput:\n%s", idx+1, ts.name, err, input)
			os.Exit(1)
		}
		if err := validateOutput(refOut, ts); err != nil {
			fmt.Fprintf(os.Stderr, "internal error: reference output invalid on test %d (%s): %v\noutput:\n%s", idx+1, ts.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runBinary(candBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%s", idx+1, ts.name, err, input)
			os.Exit(1)
		}
		if err := validateOutput(candOut, ts); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%soutput:\n%s", idx+1, ts.name, err, input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func parseBinaryArg() (string, bool) {
	if len(os.Args) == 2 {
		return os.Args[1], true
	}
	if len(os.Args) == 3 && os.Args[1] == "--" {
		return os.Args[2], true
	}
	return "", false
}

func buildBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "verifier1998A-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(path))
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return abs, func() {}, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildTests() []testSet {
	tests := []testSet{
		sampleTest(),
		handcraftedTest(),
		randomTest("random_small", 10, 1, 10),
		randomTest("random_medium", 20, 1, 40),
		bigKTest(),
	}
	return tests
}

func sampleTest() testSet {
	return testSet{
		name: "sample",
		cases: []caseSpec{
			{xc: 10, yc: 10, k: 1},
			{xc: 0, yc: 0, k: 3},
			{xc: -5, yc: -8, k: 8},
			{xc: 4, yc: -5, k: 3},
		},
	}
}

func handcraftedTest() testSet {
	return testSet{
		name: "handcrafted",
		cases: []caseSpec{
			{xc: 0, yc: 0, k: 2},
			{xc: -1, yc: 1, k: 5},
			{xc: 100, yc: -100, k: 7},
		},
	}
}

func bigKTest() testSet {
	return testSet{
		name: "big_k",
		cases: []caseSpec{
			{xc: 37, yc: -42, k: 1000},
		},
	}
}

func randomTest(name string, maxCases, minK, maxK int) testSet {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []caseSpec
	totalK := 0
	for len(cases) < maxCases && totalK < maxTotalKInput {
		remaining := maxTotalKInput - totalK
		if remaining <= 0 {
			break
		}
		kRange := maxK
		if remaining < kRange {
			kRange = remaining
		}
		k := rng.Intn(kRange-minK+1) + minK
		xc := int64(rng.Intn(centerMax-centerMin+1) + centerMin)
		yc := int64(rng.Intn(centerMax-centerMin+1) + centerMin)
		cases = append(cases, caseSpec{xc: xc, yc: yc, k: k})
		totalK += k
	}
	if len(cases) == 0 {
		cases = append(cases, caseSpec{xc: 0, yc: 0, k: 1})
	}
	return testSet{name: name, cases: cases}
}

func buildInput(ts testSet) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(ts.cases))
	for _, c := range ts.cases {
		fmt.Fprintf(&sb, "%d %d %d\n", c.xc, c.yc, c.k)
	}
	return sb.String()
}

func validateOutput(out string, ts testSet) error {
	tokens := strings.Fields(out)
	idx := 0
	for i, c := range ts.cases {
		needed := 2 * c.k
		if idx+needed > len(tokens) {
			return fmt.Errorf("case %d: expected %d numbers, got %d", i+1, needed, len(tokens)-idx)
		}
		seen := make(map[[2]int64]struct{}, c.k)
		sumX, sumY := int64(0), int64(0)
		for j := 0; j < c.k; j++ {
			x, err := parseInt64(tokens[idx])
			if err != nil {
				return fmt.Errorf("case %d: invalid x coordinate: %v", i+1, err)
			}
			y, err := parseInt64(tokens[idx+1])
			if err != nil {
				return fmt.Errorf("case %d: invalid y coordinate: %v", i+1, err)
			}
			idx += 2
			if x < -coordLimit || x > coordLimit || y < -coordLimit || y > coordLimit {
				return fmt.Errorf("case %d: coordinates (%d,%d) out of bounds", i+1, x, y)
			}
			point := [2]int64{x, y}
			if _, exists := seen[point]; exists {
				return fmt.Errorf("case %d: duplicate point (%d,%d)", i+1, x, y)
			}
			seen[point] = struct{}{}
			sumX += x
			sumY += y
		}
		if sumX != int64(c.k)*c.xc || sumY != int64(c.k)*c.yc {
			return fmt.Errorf("case %d: center mismatch", i+1)
		}
	}
	if idx != len(tokens) {
		return fmt.Errorf("extra output detected (%d tokens unused)", len(tokens)-idx)
	}
	return nil
}

func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
