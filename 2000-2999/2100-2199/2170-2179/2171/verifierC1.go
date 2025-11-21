package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const randomTests = 200

type scoreSet struct {
	g int
	c int
	l int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC1.go /path/to/binary")
		return
	}

	candidate, candCleanup, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("failed to prepare contestant binary:", err)
		return
	}
	if candCleanup != nil {
		defer candCleanup()
	}

	oracle, oracleCleanup, err := prepareOracle()
	if err != nil {
		fmt.Println("failed to prepare reference solution:", err)
		return
	}
	defer oracleCleanup()

	tests := deterministicTests()
	total := 0
	for idx, test := range tests {
		if err := runTest(test, candidate, oracle); err != nil {
			fmt.Printf("deterministic test %d failed: %v\ninput:\n%s", idx+1, err, test)
			return
		}
		total++
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomTests; i++ {
		test := randomTest(rng)
		if err := runTest(test, candidate, oracle); err != nil {
			fmt.Printf("random test %d failed: %v\ninput:\n%s", i+1, err, test)
			return
		}
		total++
	}

	fmt.Printf("All %d tests passed.\n", total)
}

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", nil, err
		}
		tmp := filepath.Join(os.TempDir(), fmt.Sprintf("candidate2171C1_%d", time.Now().UnixNano()))
		cmd := exec.Command("go", "build", "-o", tmp, abs)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return abs, nil, nil
}

func prepareOracle() (string, func(), error) {
	dir := sourceDir()
	src := filepath.Join(dir, "2171C1.go")
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("oracle2171C1_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", tmp, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("go build oracle failed: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
}

func runTest(input string, candidate, oracle string) error {
	oracleOut, err := runBinary(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle runtime error: %v", err)
	}
	candOut, err := runBinary(candidate, input)
	if err != nil {
		return fmt.Errorf("contestant runtime error: %v", err)
	}
	oracleTokens := extractTokens(oracleOut)
	candTokens := extractTokens(candOut)
	if len(candTokens) != len(oracleTokens) {
		return fmt.Errorf("expected %v tokens got %v", oracleTokens, candTokens)
	}
	for i := range candTokens {
		if !equalToken(candTokens[i], oracleTokens[i]) {
			return fmt.Errorf("expected %s got %s", oracleTokens[i], candTokens[i])
		}
	}
	return nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func extractTokens(output string) []string {
	fields := strings.Fields(output)
	return fields
}

func equalToken(a, b string) bool {
	return strings.EqualFold(a, b)
}

func deterministicTests() []string {
	cases := []scoreSet{
		{80, 80, 80},
		{90, 95, 100},
		{90, 85, 80},
		{99, 98, 99},
		{100, 90, 80},
		{85, 94, 90},
	}
	tests := make([]string, 0, len(cases))
	for _, c := range cases {
		tests = append(tests, fmt.Sprintf("%d %d %d\n", c.g, c.c, c.l))
	}
	return tests
}

func randomTest(rng *rand.Rand) string {
	g := 80 + rng.Intn(21)
	c := 80 + rng.Intn(21)
	l := 80 + rng.Intn(21)
	return fmt.Sprintf("%d %d %d\n", g, c, l)
}

func sourceDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Dir(file)
}
