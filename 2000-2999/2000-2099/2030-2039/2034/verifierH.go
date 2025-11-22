package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const refSource = "2034H.go"

type testCase struct {
	name  string
	input string
	data  []caseData
}

type caseData struct {
	n   int
	arr []int
}

type parsedOutput struct {
	sets [][]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/candidate_binary_or_go_file")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanupRef, err := buildBinary(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin := candidate
	candCleanup := func() {}
	if strings.HasSuffix(candidate, ".go") {
		candBin, candCleanup, err = buildBinary(candidate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to build candidate: %v\n", err)
			os.Exit(1)
		}
	}
	defer candCleanup()

	tests := buildTests()

	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		candOut, err := runProgram(candBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		refParsed, err := parseOutput(refOut, tc.data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}
		candParsed, err := parseOutput(candOut, tc.data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		for i, cd := range tc.data {
			refSet := refParsed.sets[i]
			candSet := candParsed.sets[i]
			if len(candSet) != len(refSet) {
				fmt.Fprintf(os.Stderr, "test %d (%s) case %d: expected subset size %d, got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, i+1, len(refSet), len(candSet), tc.input, refOut, candOut)
				os.Exit(1)
			}
			if err := validateSubset(candSet, cd.arr); err != nil {
				fmt.Fprintf(os.Stderr, "test %d (%s) case %d: invalid subset: %v\ninput:\n%scandidate output:\n%s\n",
					idx+1, tc.name, i+1, err, tc.input, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildBinary(src string) (string, func(), error) {
	abs, err := filepath.Abs(src)
	if err != nil {
		return "", nil, err
	}
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot locate verifier path")
	}
	dir := filepath.Dir(file)

	tmpDir, err := os.MkdirTemp("", "2034H-bin-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "bin")

	cmd := exec.Command("go", "build", "-o", binPath, abs)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build failed: %v\n%s", err, stderr.String())
	}

	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(output string, cases []caseData) (parsedOutput, error) {
	tokens := strings.Fields(output)
	pos := 0
	res := parsedOutput{sets: make([][]int, 0, len(cases))}

	nextInt := func() (int, error) {
		if pos >= len(tokens) {
			return 0, fmt.Errorf("unexpected end of output")
		}
		v, err := strconv.Atoi(tokens[pos])
		if err != nil {
			return 0, fmt.Errorf("invalid integer %q", tokens[pos])
		}
		pos++
		return v, nil
	}

	for idx, cd := range cases {
		k, err := nextInt()
		if err != nil {
			return res, fmt.Errorf("case %d: %v", idx+1, err)
		}
		if k < 1 || k > cd.n {
			return res, fmt.Errorf("case %d: declared size %d out of range", idx+1, k)
		}
		if pos+k > len(tokens) {
			return res, fmt.Errorf("case %d: insufficient numbers for subset", idx+1)
		}
		set := make([]int, k)
		for i := 0; i < k; i++ {
			v, err := nextInt()
			if err != nil {
				return res, fmt.Errorf("case %d: %v", idx+1, err)
			}
			set[i] = v
		}
		res.sets = append(res.sets, set)
	}
	if pos != len(tokens) {
		return res, fmt.Errorf("extra tokens in output")
	}
	return res, nil
}

func validateSubset(subset []int, arr []int) error {
	orig := make(map[int]bool, len(arr))
	for _, v := range arr {
		orig[v] = true
	}
	seen := make(map[int]bool, len(subset))
	for _, v := range subset {
		if !orig[v] {
			return fmt.Errorf("value %d not in input set", v)
		}
		if seen[v] {
			return fmt.Errorf("duplicate value %d", v)
		}
		seen[v] = true
	}
	if len(subset) == 1 {
		return nil
	}

	for i, v := range subset {
		g := 0
		for j, u := range subset {
			if i == j {
				continue
			}
			g = gcd(g, abs(u))
		}
		if g != 0 && v%g == 0 {
			return fmt.Errorf("value %d is integer combination of others (gcd %d)", v, g)
		}
	}
	return nil
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func buildTests() []testCase {
	tests := []testCase{
		newManualTest("single_element", []caseData{
			{n: 1, arr: []int{5}},
		}),
		newManualTest("small_examples", []caseData{
			{n: 5, arr: []int{2, 4, 6, 8, 10}},
			{n: 5, arr: []int{12, 15, 21, 30, 35}},
		}),
		newManualTest("pair_not_divisible", []caseData{
			{n: 4, arr: []int{2, 3, 6, 7}},
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func newManualTest(name string, cases []caseData) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, cd := range cases {
		sb.WriteString(strconv.Itoa(cd.n))
		sb.WriteByte('\n')
		for i, v := range cd.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String(), data: cases}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(4) + 1
	cases := make([]caseData, t)
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t))
	sb.WriteByte('\n')
	for i := 0; i < t; i++ {
		n := rng.Intn(10) + 1
		arr := generateDistinct(rng, n)
		cases[i] = caseData{n: n, arr: arr}
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return testCase{
		name:  fmt.Sprintf("random_%d", idx+1),
		input: sb.String(),
		data:  cases,
	}
}

func generateDistinct(rng *rand.Rand, n int) []int {
	seen := make(map[int]bool)
	arr := make([]int, 0, n)
	for len(arr) < n {
		v := rng.Intn(100000) + 1
		if seen[v] {
			continue
		}
		seen[v] = true
		arr = append(arr, v)
	}
	return arr
}
