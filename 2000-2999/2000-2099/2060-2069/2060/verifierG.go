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

const refSource2060G = "2000-2999/2000-2099/2060-2069/2060/2060G.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refVals, err := parseOutput(tc.input, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutput(tc.input, candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if len(refVals) != len(candVals) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d answers got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, len(refVals), len(candVals), tc.input, refOut, candOut)
			os.Exit(1)
		}
		for i := range refVals {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d: expected %s got %s\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, i+1, refVals[i], candVals[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2060G-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2060G.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2060G)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
	return binPath, cleanup, nil
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(input, output string) ([]string, error) {
	inFields := strings.Fields(input)
	if len(inFields) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	t, err := strconv.Atoi(inFields[0])
	if err != nil || t < 1 || t > 10000 {
		return nil, fmt.Errorf("invalid test count %q", inFields[0])
	}

	outFields := strings.Fields(output)
	if len(outFields) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(outFields))
	}
	res := make([]string, t)
	for i, tok := range outFields {
		val := strings.ToLower(tok)
		if val == "y" {
			val = "yes"
		} else if val == "n" {
			val = "no"
		}
		if val != "yes" && val != "no" {
			return nil, fmt.Errorf("answer %q is not yes/no", tok)
		}
		res[i] = val
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{
		makeManual("already_sorted", buildSingleCase([]int{1, 2, 3}, []int{4, 5, 6})),
		makeManual("statement_no", "1\n3\n2 1 3\n4 6 5\n"),
		makeManual("statement_yes", "1\n3\n2 1 5\n4 3 6\n"),
		makeManual("reverse_blocks", buildSingleCase([]int{6, 5, 4}, []int{3, 2, 1})),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 150; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	tests = append(tests, largeTest(200000, rng))

	return tests
}

func buildSingleCase(a, b []int) string {
	if len(a) != len(b) {
		panic("arrays must be same length")
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(a)))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func makeManual(name string, content string) testCase {
	return testCase{
		name:  name,
		input: content,
	}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(30) + 3
	if idx%10 == 0 {
		n = rng.Intn(200) + 3
	}
	input := randomCaseInput(n, rng)
	return testCase{
		name:  fmt.Sprintf("random_%d", idx+1),
		input: input,
	}
}

func largeTest(n int, rng *rand.Rand) testCase {
	return testCase{
		name:  "large_max",
		input: randomCaseInput(n, rng),
	}
}

func randomCaseInput(n int, rng *rand.Rand) string {
	perm := make([]int, 2*n)
	for i := 0; i < 2*n; i++ {
		perm[i] = i + 1
	}
	rng.Shuffle(len(perm), func(i, j int) {
		perm[i], perm[j] = perm[j], perm[i]
	})
	a := perm[:n]
	b := perm[n:]

	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}
