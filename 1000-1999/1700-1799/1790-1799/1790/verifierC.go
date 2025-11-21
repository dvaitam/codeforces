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

const refSource = "1000-1999/1700-1799/1790-1799/1790/1790C.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
		refAns, err := parsePermutationOutput(tc.input, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parsePermutationOutput(tc.input, candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if len(refAns) != len(candAns) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d lines, got %d\ninput:\n%s\n", idx+1, tc.name, len(refAns), len(candAns), tc.input)
			os.Exit(1)
		}

		reader := strings.NewReader(tc.input)
		var t int
		fmt.Fscan(reader, &t)
		for caseIdx := 0; caseIdx < t; caseIdx++ {
			var n int
			fmt.Fscan(reader, &n)
			seqs := make([][]int, n)
			for i := 0; i < n; i++ {
				seqs[i] = make([]int, n-1)
				for j := 0; j < n-1; j++ {
					fmt.Fscan(reader, &seqs[i][j])
				}
			}
			refPerm := refAns[caseIdx]
			candPerm := candAns[caseIdx]
			if !validPermutation(candPerm) || len(candPerm) != n {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d: invalid permutation output\n", idx+1, caseIdx+1)
				os.Exit(1)
			}
			if !equalPerms(refPerm, candPerm) {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d: expected %v got %v\ninput:\n%s\n", idx+1, caseIdx+1, refPerm, candPerm, tc.input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1790C-ref-*")
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

func parsePermutationOutput(input, out string) ([][]int, error) {
	reader := strings.NewReader(input)
	var t int
	fmt.Fscan(reader, &t)
	lines := filterNonEmpty(strings.Split(out, "\n"))
	if len(lines) != t {
		return nil, fmt.Errorf("expected %d lines, got %d", t, len(lines))
	}
	res := make([][]int, t)
	for idx, line := range lines {
		fields := strings.Fields(line)
		perm := make([]int, len(fields))
		for i, f := range fields {
			val, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q: %v", f, err)
			}
			perm[i] = val
		}
		res[idx] = perm
	}
	return res, nil
}

func validPermutation(perm []int) bool {
	seen := make(map[int]bool)
	for _, v := range perm {
		if v < 1 || v > len(perm) || seen[v] {
			return false
		}
		seen[v] = true
	}
	return true
}

func equalPerms(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
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
		buildCase("example", [][]int{
			{4, 2, 1},
			{4, 2, 3},
			{2, 1, 3},
			{4, 1, 3},
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		n := rng.Intn(10) + 3
		perm := randPermutation(rng, n)
		cases := make([][]int, n)
		for j := 0; j < n; j++ {
			seq := make([]int, 0, n-1)
			for k := 0; k < n; k++ {
				if k == j {
					continue
				}
				seq = append(seq, perm[k])
			}
			cases[j] = seq
		}
		rng.Shuffle(len(cases), func(a, b int) { cases[a], cases[b] = cases[b], cases[a] })
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), cases))
	}

	return tests
}

func randPermutation(rng *rand.Rand, n int) []int {
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) { perm[i], perm[j] = perm[j], perm[i] })
	return perm
}

func buildCase(name string, sequences [][]int) testCase {
	n := len(sequences)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n%d\n", 1, n)
	for _, seq := range sequences {
		for i, v := range seq {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}
