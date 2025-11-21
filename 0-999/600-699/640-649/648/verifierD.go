package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type testCase struct {
	name  string
	input string
}

type bowl struct {
	u int64
	t int64
}

var problemDir string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to determine verifier location")
	}
	problemDir = filepath.Dir(file)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expected, err := parseSingleInt(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "solution runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		got, err := parseSingleInt(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse solution output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d, got %d\ninput:\n%s", idx+1, tc.name, expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("Accepted (%d tests)\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	tmp, err := os.CreateTemp("", "cf-648D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "648D.go")
	cmd.Dir = problemDir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func parseSingleInt(output string) (int64, error) {
	reader := strings.NewReader(strings.TrimSpace(output))
	var val int64
	if _, err := fmt.Fscan(reader, &val); err != nil {
		return 0, err
	}
	if err := ensureNoExtraTokens(reader); err != nil {
		return 0, err
	}
	return val, nil
}

func ensureNoExtraTokens(reader *strings.Reader) error {
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return fmt.Errorf("unexpected extra token %q", extra)
	}
	return nil
}

func generateTests() []testCase {
	var tests []testCase
	add := func(name string, dogs []int64, bowls []bowl) {
		tests = append(tests, makeTestCase(name, dogs, bowls))
	}

	add("basic_mix",
		[]int64{-5, -1, 2, 6, 10},
		[]bowl{
			{u: -2, t: 2},
			{u: 2, t: 1},
			{u: 7, t: 3},
			{u: 11, t: 2},
		},
	)

	add("single_bowl",
		[]int64{-10, 10},
		[]bowl{
			{u: 0, t: 5},
		},
	)

	add("tight_deadlines",
		[]int64{-8, -3, 4, 9},
		[]bowl{
			{u: -7, t: 1},
			{u: -2, t: 2},
			{u: 5, t: 1},
			{u: 12, t: 3},
		},
	)

	add("many_bowls",
		[]int64{-15, -9, -1, 2, 20, 50},
		[]bowl{
			{u: -14, t: 3},
			{u: -5, t: 2},
			{u: 0, t: 5},
			{u: 10, t: 4},
			{u: 25, t: 10},
			{u: 45, t: 3},
			{u: 55, t: 1},
		},
	)

	rng := rand.New(rand.NewSource(648))
	tests = append(tests, randomTestCase(rng, "random_small", 20, 25, 100))
	tests = append(tests, randomTestCase(rng, "random_medium", 200, 250, 5000))
	tests = append(tests, randomTestCase(rng, "random_large", 5000, 7000, 100000))
	tests = append(tests, randomTestCase(rng, "random_huge", 200000, 200000, 1000000000))

	return tests
}

func makeTestCase(name string, dogs []int64, bowls []bowl) testCase {
	if len(dogs) == 0 {
		panic("need at least one dog")
	}
	if len(bowls) == 0 {
		panic("need at least one bowl")
	}
	if len(dogs) > 200000 || len(bowls) > 200000 {
		panic("test exceeds constraints")
	}
	dogSet := make(map[int64]struct{}, len(dogs))
	for _, x := range dogs {
		if x < -1_000_000_000 || x > 1_000_000_000 {
			panic("dog coordinate out of bounds")
		}
		if _, ok := dogSet[x]; ok {
			panic("duplicate dog coordinate")
		}
		dogSet[x] = struct{}{}
	}
	bowlSet := make(map[int64]struct{}, len(bowls))
	for _, b := range bowls {
		if b.u < -1_000_000_000 || b.u > 1_000_000_000 {
			panic("bowl coordinate out of bounds")
		}
		if b.t < 1 || b.t > 1_000_000_000 {
			panic("bowl time out of bounds")
		}
		if _, ok := bowlSet[b.u]; ok {
			panic("duplicate bowl coordinate")
		}
		bowlSet[b.u] = struct{}{}
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", len(dogs), len(bowls))
	for i, x := range dogs {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", x)
	}
	b.WriteByte('\n')
	for _, bw := range bowls {
		fmt.Fprintf(&b, "%d %d\n", bw.u, bw.t)
	}

	return testCase{name: name, input: b.String()}
}

func randomTestCase(rng *rand.Rand, name string, n, m int, coordLimit int64) testCase {
	dogs := randomUniqueCoords(rng, n, coordLimit)
	bowls := randomBowls(rng, m, coordLimit)
	return makeTestCase(name, dogs, bowls)
}

func randomUniqueCoords(rng *rand.Rand, n int, coordLimit int64) []int64 {
	if coordLimit <= 0 {
		coordLimit = 1
	}
	seen := make(map[int64]struct{}, n)
	res := make([]int64, 0, n)
	for len(res) < n {
		val := rng.Int63n(2*coordLimit+1) - coordLimit
		if _, ok := seen[val]; ok {
			continue
		}
		seen[val] = struct{}{}
		res = append(res, val)
	}
	sort.Slice(res, func(i, j int) bool { return res[i] < res[j] })
	return res
}

func randomBowls(rng *rand.Rand, m int, coordLimit int64) []bowl {
	if coordLimit <= 0 {
		coordLimit = 1
	}
	seen := make(map[int64]struct{}, m)
	res := make([]bowl, 0, m)
	for len(res) < m {
		u := rng.Int63n(2*coordLimit+1) - coordLimit
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		t := rng.Int63n(1_000_000_000) + 1
		res = append(res, bowl{u: u, t: t})
	}
	return res
}
