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
)

const refSource = "./542A.go"

type testCase struct {
	input    string
	videos   []interval
	channels []channel
}

type interval struct {
	l, r int64
}

type channel struct {
	a, b, c int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}

	tests := generateTests()

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	candidate := os.Args[1]

	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		bestVal, err := parseBestValue(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		gotOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
		if err := validateOutput(tc, gotOut, bestVal); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "542A-ref-*")
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

func parseBestValue(output string) (int64, error) {
	tokens := strings.Fields(output)
	if len(tokens) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func validateOutput(tc testCase, output string, bestVal int64) error {
	tokens := strings.Fields(output)
	if len(tokens) == 0 {
		return fmt.Errorf("empty output")
	}
	ansVal, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid number %q", tokens[0])
	}
	if bestVal == 0 {
		if ansVal != 0 {
			return fmt.Errorf("expected 0, got %d", ansVal)
		}
		if len(tokens) > 1 {
			return fmt.Errorf("unexpected extra tokens for zero answer")
		}
		return nil
	}
	if ansVal != bestVal {
		return fmt.Errorf("expected %d, got %d", bestVal, ansVal)
	}
	if len(tokens) < 3 {
		return fmt.Errorf("missing indices for positive answer")
	}
	vid, err := strconv.Atoi(tokens[1])
	if err != nil {
		return fmt.Errorf("invalid video index %q", tokens[1])
	}
	ch, err := strconv.Atoi(tokens[2])
	if err != nil {
		return fmt.Errorf("invalid channel index %q", tokens[2])
	}
	if vid < 1 || vid > len(tc.videos) {
		return fmt.Errorf("video index %d out of range", vid)
	}
	if ch < 1 || ch > len(tc.channels) {
		return fmt.Errorf("channel index %d out of range", ch)
	}
	length := intersectionLength(tc.videos[vid-1], tc.channels[ch-1])
	if length <= 0 {
		return fmt.Errorf("chosen pair has non-positive overlap")
	}
	eff := length * tc.channels[ch-1].c
	if eff != ansVal {
		return fmt.Errorf("efficiency mismatch: got %d, computed %d", ansVal, eff)
	}
	return nil
}

func intersectionLength(v interval, c channel) int64 {
	l := max64(v.l, c.a)
	r := min64(v.r, c.b)
	return r - l
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(5421))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, makeTest(
		[]interval{{7, 9}},
		[]channel{{2, 8, 2}},
	))

	for i := 0; i < 30; i++ {
		tests = append(tests, randomCase(rng, rng.Intn(4)+1, rng.Intn(4)+1))
	}

	tests = append(tests, limitCase())

	return tests
}

func sampleTest() testCase {
	videos := []interval{{7, 9}, {1, 4}}
	channels := []channel{{2, 8, 2}, {0, 4, 18}, {9, 9, 3}}
	return makeTest(videos, channels)
}

func randomCase(rng *rand.Rand, n, m int) testCase {
	videos := make([]interval, n)
	for i := 0; i < n; i++ {
		l := rng.Int63n(1_000_000)
		r := l + rng.Int63n(1_000_000)
		videos[i] = interval{l, r}
	}
	channels := make([]channel, m)
	for i := 0; i < m; i++ {
		a := rng.Int63n(1_000_000)
		b := a + rng.Int63n(1_000_000)
		c := rng.Int63n(1_000) + 1
		channels[i] = channel{a, b, c}
	}
	return makeTest(videos, channels)
}

func limitCase() testCase {
	videos := make([]interval, 5)
	channels := make([]channel, 5)
	for i := 0; i < 5; i++ {
		videos[i] = interval{int64(i) * 1_000_000_000, int64(i+1) * 1_000_000_000}
		channels[i] = channel{int64(i) * 1_000_000_000, int64(i+1)*1_000_000_000 + 500_000_000, 1_000_000_000}
	}
	return makeTest(videos, channels)
}

func makeTest(videos []interval, channels []channel) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", len(videos), len(channels))
	for _, v := range videos {
		fmt.Fprintf(&b, "%d %d\n", v.l, v.r)
	}
	for _, ch := range channels {
		fmt.Fprintf(&b, "%d %d %d\n", ch.a, ch.b, ch.c)
	}
	return testCase{input: b.String(), videos: append([]interval(nil), videos...), channels: append([]channel(nil), channels...)}
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
