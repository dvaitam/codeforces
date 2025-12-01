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

const (
	refSourceF2 = "./1909F2.go"
	modF2       = 998244353
)

type testCaseF2 struct {
	name  string
	input string
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReferenceF2()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTestsF2()
	for idx, tc := range tests {
		refOut, err := runProgramF2(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refVals, err := parseOutputF2(refOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgramF2(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutputF2(candOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		for i := 0; i < tc.t; i++ {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, i+1, refVals[i], candVals[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReferenceF2() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-1909F2-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref1909F2.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSourceF2)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
	return binPath, cleanup, nil
}

func runProgramF2(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputF2(output string, expected int) ([]int64, error) {
	fields := strings.Fields(output)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d tokens", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, tok := range fields {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		if val < 0 {
			val = (val%modF2 + modF2) % modF2
		} else {
			val %= modF2
		}
		res[i] = val
	}
	return res, nil
}

func buildTestsF2() []testCaseF2 {
	tests := []testCaseF2{
		makeManualTestF2("all_free", [][]int{
			{5, -1, -1, -1, -1, -1},
		}),
		makeManualTestF2("identity_only", [][]int{
			{5, 0, 1, 2, 3, 4},
		}),
		makeManualTestF2("sample_like", [][]int{
			{6, 0, 2, 2, 2, -1, -1},
		}),
		makeManualTestF2("mixed_cases", [][]int{
			{3, -1, 1, -1},
			{4, 0, 0, 0, 0},
			{4, 1, 2, 3, 4},
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTestF2(rng, i))
	}
	return tests
}

// arr format: first element n, followed by n values of a.
func makeManualTestF2(name string, arrs [][]int) testCaseF2 {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(arrs)))
	for _, data := range arrs {
		n := data[0]
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(data[i+1]))
		}
		sb.WriteByte('\n')
	}
	return testCaseF2{
		name:  name,
		t:     len(arrs),
		input: sb.String(),
	}
}

func randomTestF2(rng *rand.Rand, idx int) testCaseF2 {
	caseCnt := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", caseCnt))
	totalN := 0
	for i := 0; i < caseCnt; i++ {
		n := rng.Intn(10) + 1
		totalN += n
		if totalN > 200 {
			n = 1
		}
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			val := rng.Intn(n+2) - 1 // range [-1, n]
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
	}
	return testCaseF2{
		name:  fmt.Sprintf("random_%d", idx+1),
		t:     caseCnt,
		input: sb.String(),
	}
}
