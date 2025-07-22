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

type test struct {
	input  string
	output string
}

func buildReference() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "433D.go")
	ref := filepath.Join(os.TempDir(), "ref433D")
	cmd := exec.Command("go", "build", "-o", ref, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = new(bytes.Buffer)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}

func generateTests() []string {
	rand.Seed(4)
	var tests []string
	for i := 0; i < 100; i++ {
		n := rand.Intn(4) + 1
		m := rand.Intn(4) + 1
		q := rand.Intn(4) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, m, q)
		for r := 0; r < n; r++ {
			for c := 0; c < m; c++ {
				if c > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteByte(byte('0' + rand.Intn(2)))
			}
			sb.WriteByte('\n')
		}
		for j := 0; j < q; j++ {
			op := rand.Intn(2) + 1
			x := rand.Intn(n) + 1
			y := rand.Intn(m) + 1
			fmt.Fprintf(&sb, "%d %d %d\n", op, x, y)
		}
		tests = append(tests, sb.String())
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
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
	for i, in := range tests {
		expected, err := runBinary(ref, in)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("Test %d failed:\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
