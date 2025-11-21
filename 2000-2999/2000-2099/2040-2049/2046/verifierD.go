package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "2000-2999/2000-2099/2040-2049/2046/2046D.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, input := range tests {
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fail("reference failed on test %d: %v", idx+1, err)
		}
		candOut, err := runProgram(candidate, input)
		if err != nil {
			fail("candidate crashed on test %d: %v", idx+1, err)
		}
		if normalize(refOut) != normalize(candOut) {
			fail("mismatch on test %d\nInput:\n%sExpected: %sGot: %s", idx+1, input, refOut, candOut)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2046D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
	}
	return tmp.Name(), nil
}

func runProgram(path, input string) (string, error) {
	cmd := commandFor(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return stdout.String(), fmt.Errorf("%v\nstderr:\n%s", err, stderr.String())
		}
		return stdout.String(), err
	}
	return stdout.String(), nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func normalize(s string) string {
	return strings.TrimSpace(s)
}

func buildTests() []string {
	tests := []string{
		"2\n7 6\n2 1 0 1 2 3 4\n1 2\n1 3\n2 4\n2 5\n3 6\n3 7\n4 4\n1 1 1 1\n1 2\n1 3\n2 4\n3 4\n",
	}

	rng := rand.New(rand.NewSource(1))
	for len(tests) < 35 {
		tests = append(tests, randomCase(rng))
	}
	return tests
}

func randomCase(rng *rand.Rand) string {
	t := rng.Intn(4) + 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		n := rng.Intn(8) + 2
		m := rng.Intn(n*(n-1)) + 1
		fmt.Fprintf(&b, "%d %d\n", n, m)
		for i := 0; i < n; i++ {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", rng.Intn(n+1))
		}
		b.WriteByte('\n')
		for i := 0; i < m; i++ {
			u := rng.Intn(n) + 1
			v := rng.Intn(n-1) + 1
			if v >= u {
				v++
			}
			fmt.Fprintf(&b, "%d %d\n", u, v)
		}
	}
	return b.String()
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
