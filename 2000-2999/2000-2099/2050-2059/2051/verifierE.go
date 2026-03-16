package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidatePath := os.Args[1]

	refBin, cleanupRef, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := prepareCandidate(candidatePath)
	if err != nil {
		fail("failed to prepare candidate: %v", err)
	}
	defer cleanupCand()

	tests := generateTests()
	for idx, tc := range tests {
		refOut, refErr, err := runProgram(refBin, []byte(tc.input))
		if err != nil {
			fail("reference runtime error on test %d: %v\n%s", idx+1, err, refErr)
		}
		expected, err := parseOutputs(refOut, tc.t)
		if err != nil {
			fail("failed to parse reference output on test %d: %v\noutput:\n%s", idx+1, err, refOut)
		}

		candOut, candErr, err := runProgram(candBin, []byte(tc.input))
		if err != nil {
			fail("candidate runtime error on test %d: %v\n%s", idx+1, err, candErr)
		}
		got, err := parseOutputs(candOut, tc.t)
		if err != nil {
			fail("invalid candidate output on test %d: %v\noutput:\n%s", idx+1, err, candOut)
		}

		for i := 0; i < tc.t; i++ {
			if got[i] != expected[i] {
				fail("test %d case %d: wrong answer, expected %d got %d\ninput:\n%s", idx+1, i+1, expected[i], got[i], tc.input)
			}
		}
	}

	fmt.Println("OK")
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

type testCase struct {
	input string
	t     int
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase

	// Manual test from problem statement
	tests = append(tests, testCase{
		input: "5\n3 1\n4 2 3 2 1\n5 5 3\n3 1 2 3\n5 2 1 4 2\n3 2\n2 1 2\n4 3 2 2 1\n1 1\n1\n1\n4 2\n2 1 2 1\n3 2 1 3\n",
		t:     5,
	})

	// Random tests
	for i := 0; i < 50; i++ {
		t := rng.Intn(5) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", t)
		for j := 0; j < t; j++ {
			n := rng.Intn(20) + 1
			k := rng.Intn(n) + 1
			fmt.Fprintf(&sb, "%d %d\n", n, k)
			a := make([]int, n)
			for x := 0; x < n; x++ {
				a[x] = rng.Intn(n) + 1
			}
			for x := 0; x < n; x++ {
				if x > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", a[x])
			}
			sb.WriteByte('\n')
			b := make([]int, n)
			for x := 0; x < n; x++ {
				b[x] = rng.Intn(n) + 1
			}
			for x := 0; x < n; x++ {
				if x > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", b[x])
			}
			sb.WriteByte('\n')
		}
		tests = append(tests, testCase{input: sb.String(), t: t})
	}
	return tests
}

func buildReference() (string, func(), error) {
	refPath := os.Getenv("REFERENCE_SOURCE_PATH")
	if refPath == "" {
		return "", nil, fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	tmp, err := os.CreateTemp("", "2051E-ref-*")
	if err != nil {
		return "", nil, err
	}
	tmp.Close()

	content, err := os.ReadFile(refPath)
	if err != nil {
		os.Remove(tmp.Name())
		return "", nil, fmt.Errorf("read reference: %v", err)
	}
	if strings.Contains(string(content), "#include") {
		cppPath := filepath.Join(os.TempDir(), "ref2051E.cpp")
		if err := os.WriteFile(cppPath, content, 0644); err != nil {
			os.Remove(tmp.Name())
			return "", nil, err
		}
		cmd := exec.Command("g++", "-O2", "-o", tmp.Name(), cppPath)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build ref (c++): %v\n%s", err, out.String())
		}
	} else {
		cmd := exec.Command("go", "build", "-o", tmp.Name(), refPath)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
	}
	cleanup := func() {
		os.Remove(tmp.Name())
	}
	return tmp.Name(), cleanup, nil
}

func prepareCandidate(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", nil, err
		}
		tmp, err := os.CreateTemp("", "2051E-cand-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), abs)
		cmd.Dir = filepath.Dir(abs)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runProgram(path string, input []byte) (string, string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return out.String(), errBuf.String(), err
}

func parseOutputs(out string, t int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d integers, got %d", t, len(tokens))
	}
	res := make([]int64, t)
	for i := 0; i < t; i++ {
		v, err := strconv.ParseInt(tokens[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tokens[i])
		}
		res[i] = v
	}
	return res, nil
}
