package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const refSource = "1000-1999/1700-1799/1770-1779/1774/1774G.go"

type testCase struct {
	name  string
	input string
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

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refVals, err := parseInts(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVals, err := parseInts(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if len(candVals) != len(refVals) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d values, got %d\ninput:\n%s\n", idx+1, tc.name, len(refVals), len(candVals), tc.input)
			os.Exit(1)
		}
		for i := range refVals {
			if candVals[i] != refVals[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s) at query %d: expected %d, got %d\ninput:\n%s\n", idx+1, tc.name, i+1, refVals[i], candVals[i], tc.input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1774G-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func runBinary(path, input string) (string, error) {
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
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func parseInts(out string) ([]int, error) {
	lines := filterNonEmpty(strings.Split(out, "\n"))
	values := make([]int, len(lines))
	for i, line := range lines {
		val, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", line, err)
		}
		values[i] = val
	}
	return values, nil
}

func filterNonEmpty(lines []string) []string {
	var res []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			res = append(res, line)
		}
	}
	return res
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase("sample", []Seg{{1, 3}, {2, 4}, {3, 5}}, [][2]int{{1, 4}, {1, 5}}),
		buildCase("small", []Seg{{1, 2}, {2, 3}, {3, 4}}, [][2]int{{1, 3}, {2, 4}, {1, 4}}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		m := rng.Intn(5) + 1
		q := rng.Intn(5) + 1
		var segs []Seg
		used := make(map[Seg]bool)
		for len(segs) < m {
			l := rng.Intn(10) + 1
			r := l + rng.Intn(5) + 1
			seg := Seg{l, r}
			if !used[seg] {
				used[seg] = true
				segs = append(segs, seg)
			}
		}
		queries := make([][2]int, q)
		for j := 0; j < q; j++ {
			l := rng.Intn(10) + 1
			r := l + rng.Intn(5) + 1
			queries[j] = [2]int{l, r}
		}
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), segs, queries))
	}
	return tests
}

type Seg struct {
	l, r int
}

func buildCase(name string, segs []Seg, queries [][2]int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", len(segs), len(queries))
	for _, s := range segs {
		fmt.Fprintf(&sb, "%d %d\n", s.l, s.r)
	}
	for _, q := range queries {
		fmt.Fprintf(&sb, "%d %d\n", q[0], q[1])
	}
	return testCase{name: name, input: sb.String()}
}
