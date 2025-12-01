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

const (
	refSourcePath = "./2069C.go"
	mod           = 998244353
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, input := range tests {
		expected, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		got, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		if normalize(got) != normalize(expected) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, input, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2069C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourcePath))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
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

func normalize(s string) string {
	return strings.TrimSpace(s)
}

func generateTests() []string {
	var tests []string

	tests = append(tests, sampleInput())

	tests = append(tests, buildInput([][]int{
		{1, 2, 3},
	}))

	tests = append(tests, buildInput([][]int{
		{1, 1, 1},
		{3, 3, 3, 3, 3},
		{1, 2, 1, 2, 1},
	}))

	tests = append(tests, buildInput([][]int{
		repeatValue(300, 1),
		repeatValue(300, 2),
		repeatValue(300, 3),
	}))

	tests = append(tests, buildInput([][]int{
		alternatingArray(600, []int{1, 2, 3}),
		alternatingArray(600, []int{3, 2, 1}),
	}))

	rng := rand.New(rand.NewSource(2069))
	for i := 0; i < 5; i++ {
		tests = append(tests, randomTest(rng, 5000))
	}

	tests = append(tests, randomTest(rng, 50000))

	tests = append(tests, buildInput([][]int{
		randomArray(rng, 200000),
	}))

	return tests
}

func sampleInput() string {
	return strings.TrimSpace(`4
7
3 2 1 2 2 1 3
4
3 1 2 2
3
1 2 3
9
1 2 3 2 1 3 2 2 3
`) + "\n"
}

func buildInput(cases [][]int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, arr := range cases {
		if len(arr) < 3 {
			panic("array length must be at least 3")
		}
		b.WriteString(fmt.Sprintf("%d\n", len(arr)))
		for i, v := range arr {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(fmt.Sprintf("%d", v))
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func randomTest(rng *rand.Rand, totalN int) string {
	var cases [][]int
	remaining := totalN
	for remaining > 3 {
		maxLen := 1000
		if maxLen > remaining-3 {
			maxLen = remaining - 3
		}
		if maxLen < 3 {
			break
		}
		length := rng.Intn(maxLen-2) + 3
		cases = append(cases, randomArray(rng, length))
		remaining -= length
	}
	if remaining > 0 {
		if remaining < 3 && len(cases) > 0 {
			last := len(cases) - 1
			cases[last] = append(cases[last], randomArray(rng, remaining)...)
		} else if remaining >= 3 {
			cases = append(cases, randomArray(rng, remaining))
		}
	}
	return buildInput(cases)
}

func randomArray(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(3) + 1
	}
	return arr
}

func repeatValue(n, val int) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = val
	}
	return arr
}

func alternatingArray(n int, pattern []int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = pattern[i%len(pattern)]
	}
	return arr
}
