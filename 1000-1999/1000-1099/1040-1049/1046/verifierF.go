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

type testCase struct {
	name  string
	input string
}

var problemDir string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to determine verifier location")
	}
	problemDir = filepath.Dir(file)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refAns := strings.TrimSpace(refOut)

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "solution runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		candAns := strings.TrimSpace(candOut)

		if refAns != candAns {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %s, got %s\ninput:\n%s", idx+1, tc.name, refAns, candAns, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests)\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	tmp, err := os.CreateTemp("", "cf-1046F-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "1046F.go")
	cmd.Dir = problemDir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func generateTests() []testCase {
	var tests []testCase

	add := func(name string, n int, arr []int64, x, f int64) {
		var b strings.Builder
		fmt.Fprintf(&b, "%d\n", n)
		for i, v := range arr {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
		fmt.Fprintf(&b, "%d %d\n", x, f)
		tests = append(tests, testCase{name: name, input: b.String()})
	}

	add("single_exact", 1, []int64{10}, 10, 1)
	add("single_large", 1, []int64{1_000_000_000}, 50_000_000, 1)
	add("mixed_small", 4, []int64{13, 7, 6, 2}, 8, 2)
	add("all_under", 5, []int64{3, 3, 3, 3, 3}, 10, 5)
	add("border_case", 3, []int64{20, 21, 22}, 20, 5)

	rng := rand.New(rand.NewSource(1046))
	for i := 0; i < 5; i++ {
		tests = append(tests, randomCase(rng, 20, 1000, 1_000_000_000, fmt.Sprintf("random_small_%d", i+1)))
	}
	tests = append(tests, randomCase(rng, 1000, 1_000_000_000, 1_000_000_000, "random_medium"))
	tests = append(tests, randomCase(rng, 200000, 1_000_000_000, 1_000_000_000, "random_large"))

	return tests
}

func randomCase(rng *rand.Rand, n int, maxA int64, maxX int64, name string) testCase {
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = rng.Int63n(maxA) + 1
	}
	x := rng.Int63n(maxX-1) + 2
	f := rng.Int63n(x-1) + 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i, v := range arr {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", v)
	}
	b.WriteByte('\n')
	fmt.Fprintf(&b, "%d %d\n", x, f)
	return testCase{name: name, input: b.String()}
}
