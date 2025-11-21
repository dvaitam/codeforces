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

const refSource = "2000-2999/2100-2199/2160-2169/2165/2165F.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
	for i, input := range tests {
		expected, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := runProgramWithCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if normalize(got) != normalize(expected) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2165F-ref-*")
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

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runProgramWithCandidate(path, input string) (string, error) {
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

func normalize(out string) string {
	return strings.TrimSpace(out)
}

func generateTests() []string {
	var tests []string

	tests = append(tests, sampleInput())

	tests = append(tests, buildInput([][]int{
		{1},
		{1, 2},
		{2, 1},
		{1, 2, 3},
	}))

	tests = append(tests, buildInput([][]int{
		{3, 1, 4, 2, 5},
		{5, 4, 3, 2, 1},
		{1, 3, 5, 2, 4},
	}))

	rng := rand.New(rand.NewSource(2165))
	tests = append(tests, randomTest(rng, []int{6, 8, 10, 12}))
	tests = append(tests, randomTest(rng, []int{20, 25, 30, 35, 40}))

	tests = append(tests, buildInput([][]int{
		increasingPerm(60),
		decreasingPerm(60),
		increasingPerm(61),
	}))

	tests = append(tests, buildInput([][]int{
		increasingPerm(300),
		decreasingPerm(300),
		randomPerm(rng, 300),
	}))

	tests = append(tests, largeRandomTest(rng, []int{2000, 4000, 6000}))

	tests = append(tests, largeRandomTest(rng, []int{50000}))

	tests = append(tests, largeRandomTest(rng, []int{100000}))

	tests = append(tests, buildInput([][]int{
		randomPerm(rng, 300000),
	}))

	return tests
}

func sampleInput() string {
	return strings.TrimSpace(`5
5
2 1 4 3 5
10
10 3 5 2 1 4 9 8 6 7
15
3 9 15 6 11 10 5 13 12 7 4 8 14 1 2
12
10 7 12 5 4 1 2 9 3 8 6 11
30
22 30 7 17 4 13 26 28 24 20 2 11 27 21 5 19 9 10 23 14 1 25 6 8 3 18 29 12 16 15
`) + "\n"
}

func buildInput(perms [][]int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", len(perms)))
	for _, p := range perms {
		b.WriteString(fmt.Sprintf("%d\n", len(p)))
		for i, v := range p {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(fmt.Sprintf("%d", v))
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func randomTest(rng *rand.Rand, sizes []int) string {
	perms := make([][]int, len(sizes))
	for i, sz := range sizes {
		perms[i] = randomPerm(rng, sz)
	}
	return buildInput(perms)
}

func largeRandomTest(rng *rand.Rand, sizes []int) string {
	perms := make([][]int, len(sizes))
	for i, sz := range sizes {
		perms[i] = randomPerm(rng, sz)
	}
	return buildInput(perms)
}

func randomPerm(rng *rand.Rand, n int) []int {
	p := increasingPerm(n)
	for i := n - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}
	return p
}

func increasingPerm(n int) []int {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i + 1
	}
	return p
}

func decreasingPerm(n int) []int {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = n - i
	}
	return p
}
