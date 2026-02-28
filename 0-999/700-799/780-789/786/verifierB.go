package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runExe(path string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildRef() (string, error) {
	ref := "./refB.bin"
	cmd := exec.Command("go", "build", "-o", ref, "786B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func generateTest(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	q := rng.Intn(10) + 1
	s := rng.Intn(n) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, q, s)
	for i := 0; i < q; i++ {
		t := rng.Intn(3) + 1
		if t == 1 {
			v := rng.Intn(n) + 1
			u := rng.Intn(n) + 1
			w := rng.Intn(100) + 1
			fmt.Fprintf(&sb, "1 %d %d %d\n", v, u, w)
		} else {
			v := rng.Intn(n) + 1
			l := rng.Intn(n) + 1
			r := rng.Intn(n) + 1
			if l > r {
				l, r = r, l
			}
			w := rng.Intn(100) + 1
			fmt.Fprintf(&sb, "%d %d %d %d %d\n", t, v, l, r, w)
		}
	}
	return sb.String()
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(42))
	tests := make([]string, 0, 100)
	// fixed test from problem statement
	tests = append(tests, "4 5 1\n2 1 1 3 1\n3 4 1 3 1\n1 1 4 1\n2 4 2 3 1\n3 3 1 2 1\n")
	for i := 0; i < 99; i++ {
		tests = append(tests, generateTest(rng))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for i, input := range tests {
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
