package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func compileRef() (string, error) {
	out := filepath.Join(os.TempDir(), fmt.Sprintf("refC_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", out, "1641C.go")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(42))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		q := rng.Intn(5) + 1
		var b strings.Builder
		fmt.Fprintf(&b, "%d %d\n", n, q)
		for j := 0; j < q; j++ {
			t := rng.Intn(2)
			if t == 0 {
				l := rng.Intn(n) + 1
				r := rng.Intn(n-l+1) + l
				x := rng.Intn(2)
				fmt.Fprintf(&b, "0 %d %d %d\n", l, r, x)
			} else {
				jv := rng.Intn(n) + 1
				fmt.Fprintf(&b, "1 %d\n", jv)
			}
		}
		tests[i] = b.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	ref, err := compileRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for i, t := range tests {
		expect, err := runBinary(ref, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runCandidate(candidate, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, t)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, t, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
