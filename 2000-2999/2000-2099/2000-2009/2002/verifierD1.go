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

type caseData struct {
	n       int
	q       int
	p       []int
	queries [][2]int
}

type testCase struct {
	input string
	cases []caseData
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/binary")
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
		refVals, err := parseOutputs(refOut, tc.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotVals, err := parseOutputs(gotOut, tc.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "participant output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		for caseIdx := range tc.cases {
			refCase := refVals[caseIdx]
			gotCase := gotVals[caseIdx]
			if len(refCase) != len(gotCase) {
				fmt.Fprintf(os.Stderr, "test %d case %d: expected %d answers got %d\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n",
					idx+1, caseIdx+1, len(refCase), len(gotCase), tc.input, refOut, gotOut)
				os.Exit(1)
			}
			for i := range refCase {
				if refCase[i] != gotCase[i] {
					fmt.Fprintf(os.Stderr, "test %d case %d query %d mismatch: expected %s got %s\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n",
						idx+1, caseIdx+1, i+1, refCase[i], gotCase[i], tc.input, refOut, gotOut)
					os.Exit(1)
				}
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
	tmp, err := os.CreateTemp("", "2002D1_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "2002D1.go")
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

func parseOutputs(out string, cases []caseData) ([][]string, error) {
	lines := strings.Fields(out)
	idx := 0
	results := make([][]string, len(cases))
	for i, c := range cases {
		if idx+c.q > len(lines) {
			return nil, fmt.Errorf("case %d: expected %d answers but output ended early", i+1, c.q)
		}
		res := make([]string, c.q)
		for j := 0; j < c.q; j++ {
			res[j] = strings.ToUpper(lines[idx+j])
			if res[j] != "YES" && res[j] != "NO" {
				return nil, fmt.Errorf("case %d query %d: invalid token %q", i+1, j+1, lines[idx+j])
			}
		}
		idx += c.q
		results[i] = res
	}
	if idx != len(lines) {
		return nil, fmt.Errorf("output has %d extra tokens", len(lines)-idx)
	}
	return results, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, 40, 15)...)
	tests = append(tests, randomTests(rng, 30, 63)...)
	tests = append(tests, randomTests(rng, 25, 255)...)
	tests = append(tests, stressTests()...)
	return tests
}

func manualTests() []testCase {
	cases := []caseData{
		{
			n: 3, q: 3,
			p:       []int{1, 3, 2},
			queries: [][2]int{{1, 2}, {2, 3}, {1, 3}},
		},
		{
			n: 7, q: 4,
			p:       []int{1, 2, 3, 4, 5, 6, 7},
			queries: [][2]int{{2, 5}, {3, 7}, {4, 6}, {5, 7}},
		},
	}
	return []testCase{makeTestCase(cases)}
}

func randomTests(rng *rand.Rand, batches int, maxN int) []testCase {
	var tests []testCase
	for b := 0; b < batches; b++ {
		caseCnt := rng.Intn(3) + 1
		var cases []caseData
		sumN := 0
		sumQ := 0
		for i := 0; i < caseCnt; i++ {
			n := (1 << uint(rng.Intn(log2(maxN)+1))) - 1
			if n < 3 {
				n = 3
			}
			if sumN+n > 65535 {
				break
			}
			q := rng.Intn(20) + 2
			if q > 100 {
				q = 100
			}
			if sumQ+q > 50000 {
				break
			}
			p := randPermutation(rng, n)
			queries := make([][2]int, q)
			for j := 0; j < q; j++ {
				x := rng.Intn(n) + 1
				y := rng.Intn(n) + 1
				for y == x {
					y = rng.Intn(n) + 1
				}
				queries[j] = [2]int{x, y}
			}
			sumN += n
			sumQ += q
			cases = append(cases, caseData{n: n, q: q, p: p, queries: queries})
		}
		if len(cases) == 0 {
			n := 3
			q := 2
			cases = append(cases, caseData{
				n:       n,
				q:       q,
				p:       randPermutation(rng, n),
				queries: [][2]int{{1, 2}, {2, 3}},
			})
		}
		tests = append(tests, makeTestCase(cases))
	}
	return tests
}

func stressTests() []testCase {
	large := make([]caseData, 1)
	n := 65535
	p := randPermutation(rand.New(rand.NewSource(42)), n)
	q := 50
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		x := rand.Intn(n) + 1
		y := rand.Intn(n) + 1
		for y == x {
			y = rand.Intn(n) + 1
		}
		queries[i] = [2]int{x, y}
	}
	large[0] = caseData{n: n, q: q, p: p, queries: queries}
	return []testCase{makeTestCase(large)}
}

func makeTestCase(cases []caseData) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, c := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", c.n, c.q))
		for i := 2; i <= c.n; i++ {
			if i > 2 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(i / 2))
		}
		sb.WriteByte('\n')
		for i, val := range c.p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
		for _, q := range c.queries {
			sb.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
		}
	}
	return testCase{input: sb.String(), cases: cases}
}

func randPermutation(rng *rand.Rand, n int) []int {
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) {
		perm[i], perm[j] = perm[j], perm[i]
	})
	return perm
}

func log2(x int) int {
	res := 0
	for (1 << uint(res)) <= x {
		res++
	}
	return res - 1
}
