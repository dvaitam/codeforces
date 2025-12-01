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

const refSource = "./1916A.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
		refAns, err := parseOutput(tc.input, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parseOutput(tc.input, candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if len(refAns) != len(candAns) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d test cases, got %d\ninput:\n%s\n", idx+1, tc.name, len(refAns), len(candAns), tc.input)
			os.Exit(1)
		}

		for caseIdx := range refAns {
			refCase := refAns[caseIdx]
			candCase := candAns[caseIdx]

			if refCase.status != candCase.status {
				fmt.Fprintf(os.Stderr, "wrong status on test %d case %d: expected %s, got %s\n", idx+1, caseIdx+1, refCase.status, candCase.status)
				os.Exit(1)
			}

			if refCase.status == "NO" {
				continue
			}

			if len(candCase.sequence) == 0 {
				fmt.Fprintf(os.Stderr, "missing sequence on test %d case %d\n", idx+1, caseIdx+1)
				os.Exit(1)
			}
			prod := candCase.seqProduct()
			const target = 2023
			if prod != target {
				fmt.Fprintf(os.Stderr, "wrong sequence product on test %d case %d: expected product %d got %d (sequence %v)\n", idx+1, caseIdx+1, target, prod, candCase.sequence)
				os.Exit(1)
			}
			if len(candCase.sequence) != refCase.k {
				fmt.Fprintf(os.Stderr, "wrong sequence length on test %d case %d: expected %d numbers, got %d\n", idx+1, caseIdx+1, refCase.k, len(candCase.sequence))
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1916A-ref-*")
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

type caseAnswer struct {
	status   string
	sequence []int
	k        int
}

func (c caseAnswer) seqProduct() int {
	prod := 1
	for _, v := range c.sequence {
		prod *= v
	}
	return prod
}

func parseOutput(input, out string) ([]caseAnswer, error) {
	reader := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	cases := make([]struct {
		n int
		k int
	}, t)
	for i := 0; i < t; i++ {
		fmt.Fscan(reader, &cases[i].n, &cases[i].k)
		for j := 0; j < cases[i].n; j++ {
			var x int
			fmt.Fscan(reader, &x)
		}
	}

	lines := filterNonEmpty(strings.Split(out, "\n"))
	pos := 0
	answers := make([]caseAnswer, t)
	for i := 0; i < t; i++ {
		if pos >= len(lines) {
			return nil, fmt.Errorf("not enough output lines")
		}
		status := strings.TrimSpace(lines[pos])
		pos++
		if status != "YES" && status != "NO" {
			return nil, fmt.Errorf("invalid status %q", status)
		}
		if status == "NO" {
			answers[i] = caseAnswer{status: status}
			continue
		}
		if pos >= len(lines) {
			return nil, fmt.Errorf("missing sequence line")
		}
		fields := strings.Fields(lines[pos])
		pos++
		seq := make([]int, len(fields))
		for j, f := range fields {
			val, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q: %v", f, err)
			}
			seq[j] = val
		}
		answers[i] = caseAnswer{status: status, sequence: seq, k: len(seq)}
	}
	return answers, nil
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
		buildCase("simple", [][]int{
			{2, 1, 2},
			{3, 1, 1, 1},
			{4, 7, 1, 1, 1},
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		t := rng.Intn(5) + 1
		var cases [][]int
		for j := 0; j < t; j++ {
			n := rng.Intn(5) + 1
			k := rng.Intn(5) + 1
			row := make([]int, 0, 2+n)
			row = append(row, n, k)
			prod := 1
			for x := 0; x < n; x++ {
				val := rng.Intn(5) + 1
				row = append(row, val)
				if prod <= 2023 {
					prod *= val
				}
			}
			cases = append(cases, row)
		}
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), cases))
	}
	return tests
}

func buildCase(name string, cases [][]int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, c := range cases {
		n := c[0]
		k := c[1]
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(c[2+i]))
		}
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}
