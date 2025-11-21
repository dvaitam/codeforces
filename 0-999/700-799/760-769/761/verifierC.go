package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refSrc, err := locateReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refBin, err := buildReference(refSrc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		want, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s\n", idx+1, tc.input, err)
			os.Exit(1)
		}
		got, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, tc.input, err)
			os.Exit(1)
		}
		if normalize(got) != normalize(want) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d\nInput:\n%sExpected:\n%sGot:\n%s\n", idx+1, tc.input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"761C.go",
		filepath.Join("0-999", "700-799", "760-769", "761", "761C.go"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("could not find 761C.go relative to working directory")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref761C_%d.bin", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", outPath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func normalize(out string) string {
	return strings.TrimSpace(out)
}

func generateTests() []testCase {
	tests := []testCase{
		newTest(3, 1, []string{"0", "a", "#"}),
		newTest(3, 2, []string{"09", "ab", "#*"}),
		newTest(4, 3, []string{"123", "abc", "##*", "&z1"}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 200 {
		n := rng.Intn(48) + 3
		m := rng.Intn(50) + 1
		rows := make([]string, n)
		for i := 0; i < n; i++ {
			var sb strings.Builder
			for j := 0; j < m; j++ {
				sb.WriteByte(randomChar(rng))
			}
			rows[i] = sb.String()
		}
		tests = append(tests, newTest(n, m, rows))
	}
	return tests
}

func randomChar(rng *rand.Rand) byte {
	switch rng.Intn(3) {
	case 0:
		return byte('0' + rng.Intn(10))
	case 1:
		return byte('a' + rng.Intn(26))
	default:
		return []byte{'#', '*', '&'}[rng.Intn(3)]
	}
}

func newTest(n, m int, rows []string) testCase {
	if len(rows) != n {
		panic("row count mismatch")
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, row := range rows {
		if len(row) != m {
			panic("row length mismatch")
		}
		b.WriteString(row)
		b.WriteByte('\n')
	}
	return testCase{input: b.String()}
}
