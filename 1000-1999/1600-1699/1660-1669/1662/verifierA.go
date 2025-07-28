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
)

type TestCase struct {
	input string
}

func generateTests() []TestCase {
	r := rand.New(rand.NewSource(42))
	tests := make([]TestCase, 100)
	for i := range tests {
		n := r.Intn(10) + 1
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			b := r.Intn(10) + 1
			d := r.Intn(10) + 1
			sb.WriteString(fmt.Sprintf("%d %d\n", b, d))
		}
		tests[i] = TestCase{input: sb.String()}
	}
	return tests
}

func buildRef(dir string) (string, error) {
	ref := filepath.Join(dir, "refA.bin")
	cmd := exec.Command("go", "build", "-o", ref, "1662A.go")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("build ref failed: %v\n%s", err, string(out))
	}
	return ref, nil
}

func run(binary string, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	refBin, err := buildRef(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		exp, err := run(refBin, tc.input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(binary, tc.input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
