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

const refSource = "./314A.go"

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		want, err := runProgram(refBin, tc.input)
		if err != nil {
			fail("reference runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
		}
		got, err := runProgram(candidate, tc.input)
		if err != nil {
			fail("candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
		}
		if normalize(got) != normalize(want) {
			fail("wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s", i+1, tc.input, want, got)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "314A-ref-*")
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
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", filepath.Clean(bin))
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func normalize(s string) string {
	return strings.TrimSpace(s)
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(20240602))
	var tests []testCase

	tests = append(tests, makeTestCase(1, 0, []int64{1}))
	tests = append(tests, makeTestCase(2, -1, []int64{5, 3}))
	tests = append(tests, makeTestCase(4, -5, []int64{5, 3, 4, 1}))
	tests = append(tests, makeTestCase(8, 0, []int64{8, 7, 6, 5, 4, 3, 2, 1}))
	tests = append(tests, makeTestCase(8, -1000000000, []int64{1, 2, 3, 4, 5, 6, 7, 8}))

	for _, n := range []int{1, 2, 4, 8, 16} {
		for k := int64(0); k >= -5; k-- {
			arr := make([]int64, n)
			for i := range arr {
				arr[i] = int64(i + 1)
			}
			tests = append(tests, makeTestCase(n, k, arr))
		}
	}

	powers := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072}
	for i := 0; i < 120; i++ {
		n := powers[rng.Intn(len(powers))]
		tests = append(tests, randomTest(rng, n))
	}

	tests = append(tests, makeSpecialLargeTest(131072))
	tests = append(tests, makeSpecialLargeTest(65536))

	return tests
}

func randomTest(rng *rand.Rand, n int) testCase {
	k := -int64(rng.Int63n(1_000_000_001))
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = rng.Int63n(1_000_000_000) + 1
	}
	return makeTestCase(n, k, arr)
}

func makeSpecialLargeTest(n int) testCase {
	arr := make([]int64, n)
	for i := range arr {
		if i%2 == 0 {
			arr[i] = int64(1_000_000_000 - i)
		} else {
			arr[i] = int64(i + 1)
		}
	}
	return makeTestCase(n, 0, arr)
}

func makeTestCase(n int, k int64, arr []int64) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return testCase{input: sb.String()}
}
