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

const refSource = "1000-1999/1800-1899/1840-1849/1842/1842H.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/candidate")
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
			fail("candidate failed on test %d: %v", idx+1, err)
		}
		if normalize(refOut) != normalize(candOut) {
			fail("mismatch on test %d\nInput:\n%sExpected: %sGot: %s", idx+1, input, refOut, candOut)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1842H-ref-*")
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
		"1 0\n",
		"1 1\n0 1 1\n",
		"1 1\n1 1 1\n",
		"3 2\n0 1 2\n1 3 3\n",
		"2 3\n0 1 1\n1 1 1\n0 1 2\n",
	}

	rng := rand.New(rand.NewSource(1))
	for len(tests) < 40 {
		tests = append(tests, randomCase(rng))
	}
	return tests
}

func randomCase(rng *rand.Rand) string {
	n := rng.Intn(6) + 1
	maxM := n*n + n
	m := rng.Intn(maxM + 1)
	seen := make(map[[3]int]struct{})
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", n, m)
	for len(seen) < m {
		t := rng.Intn(2)
		i := rng.Intn(n) + 1
		j := rng.Intn(n-i+1) + i
		key := [3]int{t, i, j}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		fmt.Fprintf(&b, "%d %d %d\n", t, i, j)
	}
	return b.String()
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
