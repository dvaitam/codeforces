package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "./1551B2.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB2.go /path/to/candidate")
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
	for idx, tc := range tests {
		wantOut, err := runProgram(refBin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s", idx+1, err, tc)
			os.Exit(1)
		}
		want, err := parseOutput(tc, wantOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, tc, wantOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", idx+1, err, tc)
			os.Exit(1)
		}
		got, err := parseOutput(tc, gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, tc, gotOut)
			os.Exit(1)
		}

		if err := validateSolution(got, want); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%sreference output:\n%s\ncandidate output:\n%s", idx+1, err, tc, wantOut, gotOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1551B2-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		abs, err := filepath.Abs(bin)
		if err != nil {
			return "", err
		}
		cmd = exec.Command("go", "run", abs)
	} else {
		cmd = exec.Command(bin)
	}
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

type testInput struct {
	t  int
	cs []caseData
}

type caseData struct {
	n, k int
	a    []int
}

type result struct {
	colors [][]int
}

func parseInput(input string) (testInput, error) {
	sc := bufio.NewScanner(strings.NewReader(input))
	sc.Split(bufio.ScanWords)
	nextInt := func() (int, error) {
		if !sc.Scan() {
			return 0, fmt.Errorf("unexpected EOF")
		}
		return strconv.Atoi(sc.Text())
	}
	t, err := nextInt()
	if err != nil {
		return testInput{}, err
	}
	cases := make([]caseData, t)
	for i := 0; i < t; i++ {
		n, err := nextInt()
		if err != nil {
			return testInput{}, err
		}
		k, err := nextInt()
		if err != nil {
			return testInput{}, err
		}
		a := make([]int, n)
		for j := 0; j < n; j++ {
			val, err := nextInt()
			if err != nil {
				return testInput{}, err
			}
			a[j] = val
		}
		cases[i] = caseData{n: n, k: k, a: a}
	}
	return testInput{t: t, cs: cases}, nil
}

func parseOutput(input string, out string) (result, error) {
	test, err := parseInput(input)
	if err != nil {
		return result{}, err
	}
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	results := make([][]int, test.t)
	for i := 0; i < test.t; i++ {
		if !sc.Scan() {
			return result{}, fmt.Errorf("missing output for case %d", i+1)
		}
		line := sc.Text()
		parts := strings.Fields(line)
		if len(parts) != test.cs[i].n {
			// allow multi-line per case
			cur := append([]string{}, parts...)
			for len(cur) < test.cs[i].n {
				if !sc.Scan() {
					return result{}, fmt.Errorf("incomplete output for case %d", i+1)
				}
				cur = append(cur, strings.Fields(sc.Text())...)
			}
			parts = cur
		}
		if len(parts) != test.cs[i].n {
			return result{}, fmt.Errorf("case %d: expected %d integers, got %d", i+1, test.cs[i].n, len(parts))
		}
		row := make([]int, test.cs[i].n)
		for j, p := range parts {
			v, err := strconv.Atoi(p)
			if err != nil {
				return result{}, fmt.Errorf("case %d: invalid integer %q", i+1, p)
			}
			row[j] = v
		}
		results[i] = row
	}
	return result{colors: results}, nil
}

func validateSolution(got result, want result) error {
	if len(got.colors) != len(want.colors) {
		return fmt.Errorf("number of cases mismatch")
	}
	for i := range got.colors {
		if err := validateCase(got.colors[i], want.colors[i]); err != nil {
			return fmt.Errorf("case %d: %w", i+1, err)
		}
	}
	return nil
}

func validateCase(candidate []int, reference []int) error {
	if len(candidate) != len(reference) {
		return fmt.Errorf("length mismatch")
	}
	kk := 0
	for _, v := range reference {
		if v > kk {
			kk = v
		}
	}
	if kk == 0 {
		kk = maxColor(candidate)
	}
	colorPos := make(map[int][]int)
	for idx, c := range candidate {
		if c < 0 {
			return fmt.Errorf("negative color at position %d", idx+1)
		}
		if c == 0 {
			continue
		}
		colorPos[c] = append(colorPos[c], idx)
	}
	colorCount := 0
	for _, positions := range colorPos {
		colorCount += len(positions)
	}
	refColorPos := make(map[int][]int)
	for idx, c := range reference {
		if c > 0 {
			refColorPos[c] = append(refColorPos[c], idx)
		}
	}
	refColorCount := 0
	for _, positions := range refColorPos {
		refColorCount += len(positions)
	}
	if colorCount != refColorCount {
		return fmt.Errorf("painted count mismatch: expected %d got %d", refColorCount, colorCount)
	}
	for color, positions := range colorPos {
		if len(positions) == 0 {
			continue
		}
		refPositions, ok := refColorPos[color]
		if !ok {
			return fmt.Errorf("color %d not used in reference but used in candidate", color)
		}
		if len(positions) != len(refPositions) {
			return fmt.Errorf("color %d count mismatch: expected %d got %d", color, len(refPositions), len(positions))
		}
	}
	return nil
}

func maxColor(arr []int) int {
	max := 0
	for _, v := range arr {
		if v > max {
			max = v
		}
	}
	return max
}

func generateTests() []string {
	var tests []string
	tests = append(tests, "1\n1 1\n1\n")
	tests = append(tests, "1\n5 2\n1 1 2 2 3\n")
	tests = append(tests, "1\n10 3\n3 1 1 1 1 10 3 10 10 2\n")
	tests = append(tests, randomTest(5, 3))
	tests = append(tests, randomTest(10, 4))
	tests = append(tests, randomTest(50, 5))
	tests = append(tests, randomTest(200, 6))
	tests = append(tests, randomTest(1000, 7))
	tests = append(tests, randomTest(200000, 8))
	return tests
}

func randomTest(n, k int) string {
	if k > n {
		k = n
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", 1+(i%k)))
	}
	sb.WriteByte('\n')
	return sb.String()
}
