package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	name string
	k    int
	t    int64
	s    string
}

func buildOracle() (string, func(), error) {
	srcPath := os.Getenv("REFERENCE_SOURCE_PATH")
	if srcPath == "" {
		srcPath = "2172H_oracle.cpp"
	}
	tmpDir, err := os.MkdirTemp("", "oracle-2172H-")
	if err != nil {
		return "", nil, err
	}
	outPath := tmpDir + "/oracleH"
	// Detect language from content (worker saves all as .go)
	content, _ := os.ReadFile(srcPath)
	actualSrc := srcPath
	var cmd *exec.Cmd
	if strings.Contains(string(content), "#include") {
		// Copy to .cpp so g++ doesn't try the Go frontend
		actualSrc = tmpDir + "/ref.cpp"
		os.WriteFile(actualSrc, content, 0644)
		cmd = exec.Command("g++", "-O2", "-o", outPath, actualSrc)
	} else {
		cmd = exec.Command("go", "build", "-o", outPath, srcPath)
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.k, tc.t))
	sb.WriteString(tc.s)
	sb.WriteByte('\n')
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{name: "sample1", k: 4, t: 2, s: "baaabaaabaaabaaa"},
		{name: "sample2", k: 4, t: 999999999, s: "abcdefghijklmnop"},
		{name: "sample3", k: 4, t: 17, s: "bbcttckrdezzzbcd"},
		{name: "k1_zero", k: 1, t: 0, s: "ba"},
		{name: "all_same", k: 3, t: 100, s: strings.Repeat("z", 8)},
		{name: "alternating", k: 3, t: 5, s: "abababab"},
	}
}

func randomString(n int, rng *rand.Rand) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 15)
	for i := 0; i < 10; i++ {
		k := rng.Intn(8) + 1
		n := 1 << k
		t := rng.Int63n(1_000_000_000)
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_small_%d", i+1),
			k:    k,
			t:    t,
			s:    randomString(n, rng),
		})
	}
	for i := 0; i < 5; i++ {
		k := rng.Intn(4) + 10 // up to 13
		n := 1 << k
		t := rng.Int63n(1_000_000_000)
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_medium_%d", i+1),
			k:    k,
			t:    t,
			s:    randomString(n, rng),
		})
	}
	return tests
}

func stressTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano() + 123))
	return []testCase{
		{
			name: "stress_max_k",
			k:    14,
			t:    1_000_000_000,
			s:    randomString(1<<14, rng),
		},
		{
			name: "stress_repeating",
			k:    14,
			t:    987654321,
			s: func() string {
				n := 1 << 14
				b := make([]byte, n)
				pat := []byte("abcxyz")
				for i := 0; i < n; i++ {
					b[i] = pat[i%len(pat)]
				}
				return string(b)
			}(),
		},
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	tests = append(tests, stressTests()...)

	for idx, tc := range tests {
		input := buildInput(tc)
		expected, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}
		actual, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) runtime error: %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}
		if expected != actual {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %q got %q\ninput:\n%s", idx+1, tc.name, expected, actual, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
