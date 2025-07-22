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
	src := filepath.Join(dir, "433E.go")
	ref := filepath.Join(os.TempDir(), "ref433E")
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

func toBase(x, base int) []int {
	if x == 0 {
		return []int{0}
	}
	var digits []int
	for x > 0 {
		digits = append([]int{x % base}, digits...)
		x /= base
	}
	return digits
}

func generateTests() []string {
	rand.Seed(5)
	var tests []string
	for i := 0; i < 100; i++ {
		m := rand.Intn(4) + 2
		r := rand.Intn(50) + 1
		l := rand.Intn(r) + 1
		n := rand.Intn(3) + 1
		k := rand.Intn(5) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
		ld := toBase(l, m)
		fmt.Fprintf(&sb, "%d", len(ld))
		for _, d := range ld {
			fmt.Fprintf(&sb, " %d", d)
		}
		sb.WriteByte('\n')
		rd := toBase(r, m)
		fmt.Fprintf(&sb, "%d", len(rd))
		for _, d := range rd {
			fmt.Fprintf(&sb, " %d", d)
		}
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			ln := rand.Intn(3) + 1
			fmt.Fprintf(&sb, "%d", ln)
			for t := 0; t < ln; t++ {
				fmt.Fprintf(&sb, " %d", rand.Intn(m))
			}
			val := rand.Intn(3) + 1
			fmt.Fprintf(&sb, " %d\n", val)
		}
		tests = append(tests, sb.String())
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
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
