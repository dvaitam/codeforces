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

func buildReference() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "891A.go")
	ref := filepath.Join(os.TempDir(), "ref891A")
	cmd := exec.Command("go", "build", "-o", ref, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = new(bytes.Buffer)
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []string {
	rand.Seed(1)
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rand.Intn(9)+1)
		}
		sb.WriteByte('\n')
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildReference()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for i, t := range tests {
		exp, err := runBinary(ref, t)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, t)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("Test %d failed:\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
