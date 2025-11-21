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

type testInput struct {
	text     string
	ansCount int
}

func buildReference() (string, error) {
	refDir := filepath.Join("2000-2999", "2000-2099", "2000-2009", "2005")
	tmp, err := os.CreateTemp("", "ref2005B2")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	os.Remove(tmpPath)

	cmd := exec.Command("go", "build", "-o", tmpPath, "2005B2.go")
	cmd.Dir = refDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return tmpPath, nil
}

func commandForPath(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	case ".js":
		return exec.Command("node", path)
	default:
		return exec.Command(path)
	}
}

func runBinary(path, input string) (string, error) {
	cmd := commandForPath(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func parseOutput(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
	}
	ans := make([]int64, expected)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		ans[i] = val
	}
	return ans, nil
}

func fixedTests() []testInput {
	return []testInput{
		{text: "1\n8 1 1\n6\n1\n", ansCount: 1},
		{text: "1\n20 2 3\n5 15\n1 10 20\n", ansCount: 3},
	}
}

func randomTest(rng *rand.Rand) testInput {
	t := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	ansCount := 0
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		n := int64(rng.Intn(1_000_000_000-2) + 3)
		m := rng.Intn(5) + 1
		q := rng.Intn(5) + 1
		ansCount += q
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
		teacherPositions := make(map[int64]struct{})
		for len(teacherPositions) < m {
			pos := int64(rng.Intn(int(n)) + 1)
			teacherPositions[pos] = struct{}{}
		}
		first := true
		for pos := range teacherPositions {
			if !first {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", pos))
			first = false
		}
		sb.WriteByte('\n')

		for i := 0; i < q; i++ {
			var pos int64
			for {
				pos = int64(rng.Intn(int(n)) + 1)
				if _, exists := teacherPositions[pos]; !exists {
					break
				}
			}
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", pos))
		}
		sb.WriteByte('\n')
	}
	return testInput{text: sb.String(), ansCount: ansCount}
}

func edgeTests() []testInput {
	var sb strings.Builder
	sb.WriteString("1\n1000000000 3 3\n1 500000000 1000000000\n2 999999999 500000001\n")
	return []testInput{
		{text: sb.String(), ansCount: 3},
	}
}

func generateTests() []testInput {
	tests := fixedTests()
	tests = append(tests, edgeTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		tests = append(tests, randomTest(rng))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for idx, input := range tests {
		expectRaw, err := runBinary(ref, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, preview(input.text))
			os.Exit(1)
		}
		expect, err := parseOutput(expectRaw, input.ansCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, expectRaw)
			os.Exit(1)
		}

		gotRaw, err := runBinary(bin, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, preview(input.text))
			os.Exit(1)
		}
		got, err := parseOutput(gotRaw, input.ansCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, gotRaw)
			os.Exit(1)
		}

		for i := range expect {
			if expect[i] != got[i] {
				fmt.Fprintf(os.Stderr, "mismatch on test %d answer %d: expected %d got %d\ninput:\n%s\n", idx+1, i+1, expect[i], got[i], preview(input.text))
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func preview(s string) string {
	if len(s) <= 400 {
		return s
	}
	return s[:400] + "...\n"
}
