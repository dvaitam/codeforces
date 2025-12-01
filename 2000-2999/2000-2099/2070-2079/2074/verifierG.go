package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "./2074G.go"
const maxN = 400
const maxA = 1000
const maxCubeSum = maxN * maxN * maxN

type testCase struct {
	input string
}

type caseSpec struct {
	n   int
	arr []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	candidate := os.Args[1]
	tests := generateTests()

	for i, tc := range tests {
		expect, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}

		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, got)
			os.Exit(1)
		}

		if !equalTokens(expect, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc.input, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2074G-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}

	return tmp.Name(), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func equalTokens(a, b string) bool {
	ta := strings.Fields(a)
	tb := strings.Fields(b)
	if len(ta) != len(tb) {
		return false
	}
	for i := range ta {
		if ta[i] != tb[i] {
			return false
		}
	}
	return true
}

func generateTests() []testCase {
	var tests []testCase
	rng := rand.New(rand.NewSource(20742074))

	tests = append(tests, buildInput([]caseSpec{
		{n: 3, arr: []int{1, 1, 1}},
		{n: 4, arr: []int{1, 2, 3, 4}},
		{n: 5, arr: []int{2, 2, 2, 2, 2}},
	}))

	tests = append(tests, buildInput([]caseSpec{
		{n: 6, arr: []int{1, 2, 1, 3, 1, 5}},
		{n: 7, arr: []int{9, 9, 8, 2, 4, 4, 3}},
	}))

	tests = append(tests, buildInput([]caseSpec{
		{n: 9, arr: []int{9, 9, 3, 2, 4, 4, 8, 5, 3}},
	}))

	tests = append(tests, buildInput([]caseSpec{
		{n: 10, arr: seqAscending(10, 1)},
		{n: 10, arr: seqDescending(10, 10)},
	}))

	tests = append(tests, buildInput([]caseSpec{
		{n: 50, arr: constantArray(50, 7)},
	}))

	tests = append(tests, buildInput([]caseSpec{
		{n: 400, arr: randomArray(rng, 400)},
	}))

	for i := 0; i < 10; i++ {
		tests = append(tests, randomBatch(rng, 5))
	}

	return tests
}

func buildInput(cases []caseSpec) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(cases))
	for _, cs := range cases {
		fmt.Fprintf(&b, "%d\n", cs.n)
		for i, v := range cs.arr {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return testCase{input: b.String()}
}

func randomBatch(rng *rand.Rand, maxCases int) testCase {
	var specs []caseSpec
	remainingCube := maxCubeSum
	t := rng.Intn(maxCases) + 1
	for i := 0; i < t; i++ {
		minRemaining := t - i - 1
		minReserve := minRemaining * 27 // at least n=3 -> 27
		if remainingCube < minReserve {
			break
		}
		maxAllowed := remainingCube - minReserve
		n := pickN(rng, maxAllowed)
		remainingCube -= n * n * n
		specs = append(specs, caseSpec{n: n, arr: randomArray(rng, n)})
	}
	if len(specs) == 0 {
		specs = append(specs, caseSpec{n: 3, arr: []int{1, 1, 1}})
	}
	return buildInput(specs)
}

func pickN(rng *rand.Rand, cubeBudget int) int {
	maxPossible := int(math.Cbrt(float64(cubeBudget)))
	if maxPossible < 3 {
		return 3
	}
	if maxPossible > maxN {
		maxPossible = maxN
	}
	return rng.Intn(maxPossible-2) + 3
}

func randomArray(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(maxA) + 1
	}
	return arr
}

func seqAscending(n, start int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = start + i
	}
	return arr
}

func seqDescending(n, start int) []int {
	arr := make([]int, n)
	cur := start
	for i := 0; i < n; i++ {
		arr[i] = cur
		if cur > 1 {
			cur--
		}
	}
	return arr
}

func constantArray(n, val int) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = val
	}
	return arr
}
