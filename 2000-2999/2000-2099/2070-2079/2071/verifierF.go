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

var refSource = func() string {
	if p := os.Getenv("REFERENCE_SOURCE_PATH"); p != "" {
		return p
	}
	return "./2071F.go"
}()

type testCase struct {
	n int
	k int
	a []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := buildTestSuite(rng)
	input := serializeTests(tests)

	expected, err := runAndParse(refBin, input, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}

	got, err := runAndParse(candidate, input, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	for i := range expected {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "Mismatch on test %d: expected %d got %d\nInput:\n%s", i+1, expected[i], got[i], input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

// isJavaSource returns true if content looks like Java source code.
func isJavaSource(content []byte) bool {
	s := string(content)
	return strings.Contains(s, "public class") || strings.Contains(s, "import java.")
}

// buildReference reads the reference source, detects language, compiles, and
// returns a command-runner string plus a cleanup function.
func buildReference() (string, func(), error) {
	src := filepath.Clean(refSource)
	content, err := os.ReadFile(src)
	if err != nil {
		return "", func() {}, fmt.Errorf("read reference source %s: %v", src, err)
	}

	if isJavaSource(content) {
		return buildJavaReference(content)
	}

	if strings.Contains(string(content), "#include") {
		return buildCppReference(content)
	}

	// Default: Go source
	return buildGoReference(src)
}

func buildGoReference(src string) (string, func(), error) {
	tmp, err := os.CreateTemp("", "ref_2071F_*.bin")
	if err != nil {
		return "", func() {}, err
	}
	tmp.Close()
	binPath := tmp.Name()

	cmd := exec.Command("go", "build", "-o", binPath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(binPath)
		return "", func() {}, fmt.Errorf("failed to build Go reference: %v\n%s", err, string(out))
	}
	return binPath, func() { os.Remove(binPath) }, nil
}

func buildCppReference(content []byte) (string, func(), error) {
	tmp, err := os.CreateTemp("", "ref_2071F_*.bin")
	if err != nil {
		return "", func() {}, err
	}
	tmp.Close()
	binPath := tmp.Name()

	cppSrc := binPath + ".cpp"
	if err := os.WriteFile(cppSrc, content, 0644); err != nil {
		os.Remove(binPath)
		return "", func() {}, err
	}

	cmd := exec.Command("g++", "-O2", "-o", binPath, cppSrc)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(binPath)
		os.Remove(cppSrc)
		return "", func() {}, fmt.Errorf("failed to build C++ reference: %v\n%s", err, string(out))
	}
	os.Remove(cppSrc)
	return binPath, func() { os.Remove(binPath) }, nil
}

func buildJavaReference(content []byte) (string, func(), error) {
	// Create a temp directory for Java compilation.
	tmpDir, err := os.MkdirTemp("", "ref_2071F_java_*")
	if err != nil {
		return "", func() {}, err
	}

	javaSrc := filepath.Join(tmpDir, "Main.java")
	if err := os.WriteFile(javaSrc, content, 0644); err != nil {
		os.RemoveAll(tmpDir)
		return "", func() {}, err
	}

	cmd := exec.Command("javac", javaSrc)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", func() {}, fmt.Errorf("failed to compile Java reference: %v\n%s", err, out.String())
	}

	// Return a special marker; we handle Java execution in runProgram.
	// The "binary" path is the tmpDir with a sentinel suffix.
	marker := tmpDir + "/.java_main"
	os.WriteFile(marker, []byte(tmpDir), 0644)

	return marker, func() { os.RemoveAll(tmpDir) }, nil
}

func runAndParse(target, input string, tests int) ([]int64, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return nil, err
	}
	return parseAnswers(out, tests)
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd

	if strings.HasSuffix(target, "/.java_main") {
		// Java reference: run with java -cp <dir> Main
		classDir := filepath.Dir(target)
		cmd = exec.Command("java", "-cp", classDir, "Main")
	} else if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}

	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func parseAnswers(out string, tests int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != tests {
		return nil, fmt.Errorf("expected %d answers, got %d (output: %q)", tests, len(fields), out)
	}
	ans := make([]int64, tests)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		ans[i] = val
	}
	return ans, nil
}

func buildTestSuite(rng *rand.Rand) []testCase {
	tests := deterministicTests()
	totalN := 0
	for _, tc := range tests {
		totalN += tc.n
	}

	for len(tests) < 200 && totalN < 200000 {
		n := rng.Intn(1000) + 1
		if totalN+n > 200000 {
			n = 200000 - totalN
		}
		k := rng.Intn(n)
		if rng.Intn(5) == 0 {
			k = n - 1
		}
		tc := randomTest(rng, n, k)
		tests = append(tests, tc)
		totalN += n
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, k: 0, a: []int{1}},
		{n: 2, k: 1, a: []int{1, 1}},
		{n: 5, k: 0, a: []int{2, 1, 4, 5, 2}},
		{n: 5, k: 3, a: []int{2, 1, 4, 5, 2}},
		{n: 6, k: 1, a: []int{1, 2, 3, 4, 5, 11}},
		{n: 11, k: 6, a: []int{6, 3, 8, 5, 8, 3, 2, 1, 2, 7, 11}},
		{n: 14, k: 3, a: []int{3, 2, 3, 5, 5, 2, 6, 7, 4, 8, 10, 1, 8, 9}},
		{n: 2, k: 0, a: []int{1, 1000000000}},
		{n: 3, k: 1, a: []int{1000000000, 1, 1}},
	}
}

func randomTest(rng *rand.Rand, n, k int) testCase {
	a := make([]int, n)
	maxVal := 1_000_000_000
	for i := 0; i < n; i++ {
		switch rng.Intn(6) {
		case 0:
			a[i] = 1
		case 1:
			a[i] = maxVal
		default:
			a[i] = rng.Intn(maxVal) + 1
		}
	}
	return testCase{n: n, k: k, a: a}
}

func serializeTests(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for idx, tc := range tests {
		if idx > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
