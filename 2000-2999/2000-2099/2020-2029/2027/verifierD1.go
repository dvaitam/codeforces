package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const refSource = "2000-2999/2000-2099/2020-2029/2027/2027D1.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD1.go /path/to/candidate")
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
	tmp, err := os.CreateTemp("", "2027D1-ref-*")
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
		"5\n4 2\n9 3 4 3\n11 7\n1 2\n2019 18\n10 22 5 2 1 10 3 2 9 9 6\n17 9\n10 11 2 2 2 2 2 2 2 2\n20 18 16 14 12 10 8 6 4\n11 6\n1032 16 8 4 2 1\n",
	}

	rng := rand.New(rand.NewSource(1))
	for len(tests) < 40 {
		tests = append(tests, randomCase(rng))
	}
	return tests
}

func randomCase(rng *rand.Rand) string {
	t := rng.Intn(4) + 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		n := rng.Intn(7) + 1
		m := rng.Intn(7) + 1
		fmt.Fprintf(&b, "%d %d\n", n, m)
		for i := 0; i < n; i++ {
			val := rng.Intn(50) + 1
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", val)
		}
		b.WriteByte('\n')

		bVals := make([]int, m)
		for i := 0; i < m; i++ {
			bVals[i] = rng.Intn(50) + 1 + (m-i)*5
		}
		sort.Slice(bVals, func(i, j int) bool { return bVals[i] > bVals[j] })
		for i := 1; i < m; i++ {
			if bVals[i] >= bVals[i-1] {
				bVals[i] = bVals[i-1] - 1
			}
			if bVals[i] < 1 {
				bVals[i] = 1
			}
		}
		for i := 0; i < m; i++ {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", bVals[i])
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
