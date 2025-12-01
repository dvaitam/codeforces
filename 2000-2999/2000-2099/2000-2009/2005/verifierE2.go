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

type matrix struct {
	l int
	n int
	m int
	a []int
	b [][]int
}

type testCase struct {
	input string
	data  []matrix
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/binary")
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
		refVals, err := parseOutputs(refOut, len(tc.data))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotVals, err := parseOutputs(gotOut, len(tc.data))
		if err != nil {
			fmt.Fprintf(os.Stderr, "participant output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		for caseIdx := 0; caseIdx < len(tc.data); caseIdx++ {
			if refVals[caseIdx] != gotVals[caseIdx] {
				fmt.Fprintf(os.Stderr, "test %d case %d mismatch: expected %s got %s\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n",
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
	tmp, err := os.CreateTemp("", "2005E2_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "2005E2.go")
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
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
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

func parseOutputs(out string, expected int) ([]string, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d tokens got %d", expected, len(fields))
	}
	for i := range fields {
		fields[i] = strings.ToUpper(fields[i])
		if fields[i] != "T" && fields[i] != "N" {
			return nil, fmt.Errorf("invalid token %q", fields[i])
		}
	}
	return fields, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, 30, 50, 50)...)
	tests = append(tests, randomTests(rng, 30, 100, 100)...)
	tests = append(tests, randomTests(rng, 25, 200, 200)...)
	tests = append(tests, stressTests()...)
	return tests
}

func manualTests() []testCase {
	cases := []matrix{
		{
			l: 2, n: 2, m: 3,
			a: []int{1, 2},
			b: [][]int{
				{1, 3, 6},
				{4, 6, 2},
			},
		},
		{
			l: 4, n: 2, m: 2,
			a: []int{1, 2, 1, 2},
			b: [][]int{
				{1, 1},
				{2, 2},
			},
		},
		{
			l: 1, n: 1, m: 1,
			a: []int{1},
			b: [][]int{{1}},
		},
	}
	return []testCase{makeTestCase(cases)}
}

func randomTests(rng *rand.Rand, batches, maxN, maxM int) []testCase {
	const limit = 3_000_000
	var tests []testCase
	for b := 0; b < batches; b++ {
		caseCnt := rng.Intn(3) + 1
		var cases []matrix
		sumNM := 0
		sumL := 0
		for len(cases) < caseCnt {
			n := rng.Intn(maxN) + 1
			m := rng.Intn(maxM) + 1
			if sumNM+n*m > limit {
				break
			}
			maxRemain := 1500 - sumL
			if maxRemain <= 0 {
				break
			}
			l := rng.Intn(minInt(maxRemain, 1500)) + 1
			k := n * m
			a := make([]int, l)
			for i := 0; i < l; i++ {
				a[i] = rng.Intn(k) + 1
			}
			b := make([][]int, n)
			for i := 0; i < n; i++ {
				row := make([]int, m)
				for j := 0; j < m; j++ {
					row[j] = rng.Intn(k) + 1
				}
				b[i] = row
			}
			cases = append(cases, matrix{l: l, n: n, m: m, a: a, b: b})
			sumNM += n * m
			sumL += l
		}
		if len(cases) == 0 {
			cases = append(cases, matrix{
				l: 1, n: 1, m: 1,
				a: []int{1},
				b: [][]int{{1}},
			})
		}
		tests = append(tests, makeTestCase(cases))
	}
	return tests
}

func stressTests() []testCase {
	largeA := make([]int, 1500)
	for i := range largeA {
		largeA[i] = i%100 + 1
	}
	largeB := make([][]int, 1500)
	for i := range largeB {
		row := make([]int, 2)
		for j := range row {
			row[j] = (i+j)%100 + 1
		}
		largeB[i] = row
	}
	return []testCase{
		makeTestCase([]matrix{
			{
				l: len(largeA),
				n: len(largeB),
				m: len(largeB[0]),
				a: largeA,
				b: largeB,
			},
		}),
	}
}

func makeTestCase(cases []matrix) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	stored := make([]matrix, len(cases))
	for i, c := range cases {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", c.l, c.n, c.m))
		for j := 0; j < c.l; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(c.a[j]))
		}
		sb.WriteByte('\n')
		for r := 0; r < c.n; r++ {
			for col := 0; col < c.m; col++ {
				if col > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(c.b[r][col]))
			}
			sb.WriteByte('\n')
		}
		copyA := append([]int(nil), c.a...)
		copyB := make([][]int, c.n)
		for r := 0; r < c.n; r++ {
			copyB[r] = append([]int(nil), c.b[r]...)
		}
		stored[i] = matrix{l: c.l, n: c.n, m: c.m, a: copyA, b: copyB}
	}
	return testCase{input: sb.String(), data: stored}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
