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

const referenceSource = "2023F.go"

type testCase struct {
	name string
	n    int
	q    int
	a    []int64
	qs   [][2]int
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/candidate")
		os.Exit(1)
	}
	candPath := os.Args[len(os.Args)-1]

	refBin, refCleanup, err := buildBinary(referencePath())
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer refCleanup()

	candBin, candCleanup, err := prepareCandidate(candPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to prepare candidate:", err)
		os.Exit(1)
	}
	defer candCleanup()

	tests := buildTests()
	input := buildInput(tests)

	expOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	actOut, err := runProgram(candBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\ninput preview:\n%s", err, previewInput(input))
		os.Exit(1)
	}

	totalAnswers := totalQueries(tests)
	expVals, err := parseOutputs(expOut, totalAnswers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output invalid: %v\n%s", err, expOut)
		os.Exit(1)
	}
	actVals, err := parseOutputs(actOut, totalAnswers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate output invalid: %v\n%s", err, actOut)
		fmt.Fprintln(os.Stderr, previewInput(input))
		os.Exit(1)
	}

	for i := 0; i < totalAnswers; i++ {
		if expVals[i] != actVals[i] {
			fmt.Fprintf(os.Stderr, "mismatch at answer %d: expected %d got %d\n", i+1, expVals[i], actVals[i])
			fmt.Fprintln(os.Stderr, previewInput(input))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d answers matched across %d tests.\n", totalAnswers, len(tests))
}

func referencePath() string {
	if _, file, _, ok := runtime.Caller(0); ok {
		return filepath.Join(filepath.Dir(file), referenceSource)
	}
	return referenceSource
}

func buildBinary(src string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "ref2023F")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(tmpDir, "ref")
	cmd := exec.Command("go", "build", "-o", bin, src)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("%v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
}

func prepareCandidate(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	tmpDir, err := os.MkdirTemp("", "cand2023F")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(tmpDir, "candidate")
	cmd := exec.Command("go", "build", "-o", bin, path)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("%v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return errBuf.String(), fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func parseOutputs(out string, want int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != want {
		return nil, fmt.Errorf("expected %d integers, got %d", want, len(tokens))
	}
	res := make([]int64, want)
	for i, tok := range tokens {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		res[i] = v
	}
	return res, nil
}

func buildTests() []testCase {
	var tests []testCase
	tests = append(tests, smallHandcrafted()...)
	tests = append(tests, randomTests(200)...)
	return tests
}

func smallHandcrafted() []testCase {
	return []testCase{
		{
			name: "single-positive",
			n:    3,
			q:    3,
			a:    []int64{5, -2, -3},
			qs:   [][2]int{{1, 3}, {2, 3}, {1, 1}},
		},
		{
			name: "alternating",
			n:    6,
			q:    4,
			a:    []int64{2, -1, 3, -4, 5, -5},
			qs:   [][2]int{{1, 6}, {2, 5}, {3, 3}, {4, 6}},
		},
	}
}

func randomTests(count int) []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)

	for len(tests) < count {
		n := rng.Intn(60) + 1
		q := rng.Intn(80) + 1
		a := make([]int64, n)
		for i := range a {
			val := rng.Int63n(2000000000) - 1000000000
			if val == 0 {
				val = 1
			}
			a[i] = val
		}
		qs := make([][2]int, q)
		for i := 0; i < q; i++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n) + 1
			if l > r {
				l, r = r, l
			}
			qs[i] = [2]int{l, r}
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("rnd-%d", len(tests)+1),
			n:    n,
			q:    q,
			a:    a,
			qs:   qs,
		})
	}
	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.Grow(len(tests) * 128)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(tc.q))
		sb.WriteByte('\n')
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for _, q := range tc.qs {
			sb.WriteString(strconv.Itoa(q[0]))
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(q[1]))
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func totalQueries(tests []testCase) int {
	total := 0
	for _, tc := range tests {
		total += len(tc.qs)
	}
	return total
}

func previewInput(input string) string {
	lines := strings.Split(input, "\n")
	if len(lines) > 20 {
		lines = lines[:20]
	}
	return strings.Join(lines, "\n")
}
