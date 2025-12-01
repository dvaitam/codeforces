package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "./2018E2.go"
const totalLimit = 300000

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}

	candidate := os.Args[1]
	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	expected := tokenize(refOut)

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	got := tokenize(candOut)

	if len(expected) != len(got) {
		fmt.Fprintf(os.Stderr, "wrong number of tokens: expected %d got %d\n", len(expected), len(got))
		os.Exit(1)
	}
	for i := range expected {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "mismatch at token %d: expected %q got %q\n", i+1, expected[i], got[i])
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2018E2-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.WriteString(errBuf.String())
		return out.String(), err
	}
	if errBuf.Len() > 0 {
		out.WriteString(errBuf.String())
	}
	return out.String(), nil
}

func tokenize(s string) []string {
	return strings.Fields(strings.TrimSpace(s))
}

type segment struct {
	l int
	r int
}

type testCase struct {
	segs []segment
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		n := len(tc.segs)
		fmt.Fprintf(&b, "%d\n", n)
		for i, seg := range tc.segs {
			fmt.Fprintf(&b, "%d", seg.l)
			if i+1 == n {
				fmt.Fprintln(&b)
			} else {
				fmt.Fprintf(&b, " ")
			}
		}
		for i, seg := range tc.segs {
			fmt.Fprintf(&b, "%d", seg.r)
			if i+1 == n {
				fmt.Fprintln(&b)
			} else {
				fmt.Fprintf(&b, " ")
			}
		}
	}
	return b.String()
}

func generateTests() []testCase {
	var tests []testCase
	total := 0

	add := func(segs []segment) {
		if len(segs) == 0 {
			return
		}
		if total+len(segs) > totalLimit {
			return
		}
		tests = append(tests, testCase{segs: segs})
		total += len(segs)
	}

	add([]segment{{1, 1}})
	add([]segment{{1, 2}, {2, 3}, {3, 4}})
	add([]segment{{1, 5}, {2, 4}, {6, 9}, {8, 10}, {3, 7}})
	add([]segment{{1, 6}, {2, 5}, {3, 4}, {7, 9}, {8, 10}, {11, 11}})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for total < totalLimit {
		remaining := totalLimit - total
		n := rng.Intn(min(5000, remaining)) + 1
		segs := make([]segment, n)
		maxCoord := 2 * n
		if maxCoord < 10 {
			maxCoord = 10
		}
		for i := 0; i < n; i++ {
			l := rng.Intn(maxCoord) + 1
			r := rng.Intn(maxCoord-l+1) + l
			segs[i] = segment{l, r}
		}
		add(segs)
	}

	return tests
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
