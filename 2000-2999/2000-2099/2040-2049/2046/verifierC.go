package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSourceC = "2000-2999/2000-2099/2040-2049/2046/2046C.go"

type caseData struct {
	points [][2]int64
}

type testInput struct {
	input string
	cases []caseData
}

type result struct {
	k int64
	x int64
	y int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceC)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, test := range tests {
		refOut, err := runProgram(refBin, test.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", idx+1, err, test.input)
			os.Exit(1)
		}
		expected, err := parseResults(refOut, len(test.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		userOut, err := runCandidate(candidate, test.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", idx+1, err, test.input)
			os.Exit(1)
		}
		actual, err := parseResults(userOut, len(test.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, test.input, userOut)
			os.Exit(1)
		}

		for caseIdx, cd := range test.cases {
			exp := expected[caseIdx]
			act := actual[caseIdx]
			if act.k != exp.k {
				fmt.Fprintf(os.Stderr, "test %d case %d: expected k=%d, got %d\n", idx+1, caseIdx+1, exp.k, act.k)
				os.Exit(1)
			}
			if act.k < 0 || act.k > int64(len(cd.points)/4) {
				fmt.Fprintf(os.Stderr, "test %d case %d: invalid k=%d\n", idx+1, caseIdx+1, act.k)
				os.Exit(1)
			}
			minCount := quadrantMin(cd.points, act.x, act.y)
			if int64(minCount) != act.k {
				fmt.Fprintf(os.Stderr, "test %d case %d: point (%d,%d) yields min=%d, reported k=%d\n",
					idx+1, caseIdx+1, act.x, act.y, minCount, act.k)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(source string) (string, error) {
	tmp, err := os.CreateTemp("", "2046C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(source))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseResults(out string, cases int) ([]result, error) {
	reader := strings.NewReader(out)
	res := make([]result, cases)
	for i := 0; i < cases; i++ {
		if _, err := fmt.Fscan(reader, &res[i].k); err != nil {
			return nil, fmt.Errorf("case %d: failed to read k: %v", i+1, err)
		}
		if _, err := fmt.Fscan(reader, &res[i].x, &res[i].y); err != nil {
			return nil, fmt.Errorf("case %d: failed to read point: %v", i+1, err)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected (starts with %s)", extra)
	}
	return res, nil
}

func quadrantMin(points [][2]int64, x0, y0 int64) int {
	counts := [4]int{}
	for _, p := range points {
		switch {
		case p[0] >= x0 && p[1] >= y0:
			counts[0]++
		case p[0] < x0 && p[1] >= y0:
			counts[1]++
		case p[0] >= x0 && p[1] < y0:
			counts[2]++
		default:
			counts[3]++
		}
	}
	minVal := counts[0]
	for i := 1; i < 4; i++ {
		if counts[i] < minVal {
			minVal = counts[i]
		}
	}
	return minVal
}

func generateTests() []testInput {
	var tests []testInput
	tests = append(tests, sampleTest())
	tests = append(tests, smallCustomTest())
	rng := rand.New(rand.NewSource(2046))
	tests = append(tests, randomTest(rng, []int{10, 12, 15, 18}, 50))
	tests = append(tests, randomTest(rng, []int{200, 250, 300}, 1000))
	tests = append(tests, randomTest(rng, []int{5000, 6000, 7000}, 100000))
	tests = append(tests, randomTest(rng, []int{50000, 50000}, 1000000000))
	return tests
}

func sampleTest() testInput {
	cases := []caseData{
		newCase([][2]int64{{1, 1}, {1, 2}, {2, 1}, {2, 2}}),
		newCase([][2]int64{{0, 0}, {0, 0}, {0, 0}, {0, 0}}),
		newCase([][2]int64{{1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}, {-2, -1}, {-2, 1}, {-1, 2}}),
		newCase([][2]int64{{1, 1}, {1, 2}, {1, 3}, {1, 4}, {2, 1}, {3, 1}, {4, 1}}),
	}
	return buildTest(cases)
}

func smallCustomTest() testInput {
	cases := []caseData{
		newCase([][2]int64{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}, {2, 2}, {-2, -2}}),
		newCase([][2]int64{{0, 5}, {5, 0}, {-5, 0}, {0, -5}, {10, 10}, {-10, -10}, {3, 3}, {3, -3}}),
		newCase([][2]int64{{100, 100}, {100, -100}, {-100, 100}, {-100, -100}}),
	}
	return buildTest(cases)
}

func randomTest(rng *rand.Rand, sizes []int, coordLimit int64) testInput {
	cases := make([]caseData, len(sizes))
	for i, n := range sizes {
		cases[i] = randomCase(rng, n, coordLimit)
	}
	return buildTest(cases)
}

func randomCase(rng *rand.Rand, n int, coordLimit int64) caseData {
	points := make([][2]int64, n)
	for i := 0; i < n; i++ {
		points[i][0] = randCoord(rng, coordLimit)
		points[i][1] = randCoord(rng, coordLimit)
	}
	return newCase(points)
}

func randCoord(rng *rand.Rand, limit int64) int64 {
	if limit <= 0 {
		return 0
	}
	val := rng.Int63n(2*limit + 1)
	return val - limit
}

func newCase(points [][2]int64) caseData {
	return caseData{points: points}
}

func buildTest(cases []caseData) testInput {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, c := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", len(c.points)))
		for _, p := range c.points {
			sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}
	}
	return testInput{input: sb.String(), cases: cases}
}
