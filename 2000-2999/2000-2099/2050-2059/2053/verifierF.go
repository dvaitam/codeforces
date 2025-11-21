package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type matrixTest struct {
	n, m, k int
	data    [][]int
}

type testCase struct {
	input   string
	caseCnt int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(ref, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, tc.caseCnt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotVals, err := parseOutputs(gotOut, tc.caseCnt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "participant output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		for caseIdx := 0; caseIdx < tc.caseCnt; caseIdx++ {
			if refVals[caseIdx] != gotVals[caseIdx] {
				fmt.Fprintf(os.Stderr, "test %d case %d mismatch: expected %d got %d\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n",
					idx+1, caseIdx+1, refVals[caseIdx], gotVals[caseIdx], tc.input, refOut, gotOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "2053F_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "2053F.go")
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return path, nil
}

func verifierDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine verifier directory")
	}
	return filepath.Dir(file), nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	switch {
	case strings.HasSuffix(path, ".go"):
		cmd = exec.Command("go", "run", path)
	default:
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutputs(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers got %d", expected, len(fields))
	}
	vals := make([]int64, expected)
	for i, token := range fields {
		v, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", token)
		}
		vals[i] = v
	}
	return vals, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, 60, 8)...)
	tests = append(tests, randomTests(rng, 60, 60)...)
	tests = append(tests, randomTests(rng, 40, 200)...)
	tests = append(tests, stressTests()...)
	return tests
}

func manualTests() []testCase {
	cases := []matrixTest{
		{
			n: 3, m: 3, k: 3,
			data: [][]int{
				{1, 2, 2},
				{3, 1, 3},
				{3, 2, 1},
			},
		},
		{
			n: 2, m: 3, k: 3,
			data: [][]int{
				{-1, 3, 3},
				{2, 2, -1},
			},
		},
	}
	return []testCase{makeTestCase(cases)}
}

func randomTests(rng *rand.Rand, batches int, maxNM int) []testCase {
	const limit = 600000
	var tests []testCase
	for b := 0; b < batches; b++ {
		caseCnt := rng.Intn(3) + 1
		var cases []matrixTest
		sumNM := 0
		for len(cases) < caseCnt {
			n := rng.Intn(maxNM-1) + 2
			m := rng.Intn(maxNM-1) + 2
			if n*m > maxNM || sumNM+n*m > limit {
				break
			}
			k := rng.Intn(n*m) + 1
			mat := make([][]int, n)
			for i := 0; i < n; i++ {
				row := make([]int, m)
				for j := 0; j < m; j++ {
					if rng.Intn(5) == 0 {
						row[j] = -1
					} else {
						row[j] = rng.Intn(k) + 1
					}
				}
				mat[i] = row
			}
			cases = append(cases, matrixTest{n: n, m: m, k: k, data: mat})
			sumNM += n * m
		}
		if len(cases) == 0 {
			cases = append(cases, matrixTest{
				n: 2, m: 2, k: 4,
				data: [][]int{
					{rng.Intn(4) + 1, -1},
					{-1, rng.Intn(4) + 1},
				},
			})
		}
		tests = append(tests, makeTestCase(cases))
	}
	return tests
}

func stressTests() []testCase {
	return []testCase{
		makeTestCase([]matrixTest{
			fullMatrix(300, 2, 10),
		}),
		makeTestCase([]matrixTest{
			stripedMatrix(200, 3, 20),
		}),
	}
}

func fullMatrix(n, m, k int) matrixTest {
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		row := make([]int, m)
		for j := 0; j < m; j++ {
			row[j] = (i+j)%k + 1
		}
		mat[i] = row
	}
	return matrixTest{n: n, m: m, k: k, data: mat}
}

func stripedMatrix(n, m, k int) matrixTest {
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		row := make([]int, m)
		for j := 0; j < m; j++ {
			if j%2 == 0 {
				row[j] = -1
			} else {
				row[j] = (i + j) % k
				if row[j] == 0 {
					row[j] = k
				}
			}
		}
		mat[i] = row
	}
	return matrixTest{n: n, m: m, k: k, data: mat}
}

func makeTestCase(cases []matrixTest) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, c := range cases {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", c.n, c.m, c.k))
		for i := 0; i < c.n; i++ {
			for j := 0; j < c.m; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(c.data[i][j]))
			}
			sb.WriteByte('\n')
		}
	}
	return testCase{input: sb.String(), caseCnt: len(cases)}
}
