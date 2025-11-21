package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "0-999/700-799/720-729/720/720D.go"

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
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
	tmp, err := os.CreateTemp("", "720D-ref-*")
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
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	tests = append(tests, manualTest(3, 3, [][4]int64{}))
	tests = append(tests, manualTest(3, 3, [][4]int64{{1, 2, 2, 2}}))
	tests = append(tests, manualTest(5, 4, [][4]int64{{2, 2, 3, 3}, {4, 1, 4, 2}}))
	for i := 0; i < 50; i++ {
		tests = append(tests, randomTest(rng, 10, 10, 8))
	}
	for i := 0; i < 30; i++ {
		tests = append(tests, randomTest(rng, 1_000_000, 1_000_000, 100000))
	}
	return tests
}

func manualTest(n, m int64, obs [][4]int64) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, len(obs))
	for _, o := range obs {
		fmt.Fprintf(&sb, "%d %d %d %d\n", o[0], o[1], o[2], o[3])
	}
	return testCase{input: sb.String()}
}

func randomTest(rng *rand.Rand, maxN, maxM int64, maxK int) testCase {
	n := rng.Int63n(maxN-2) + 3
	m := rng.Int63n(maxM-2) + 3
	k := rng.Intn(maxK + 1)
	if k > int(n*m) {
		k = int(n*m) / 2
	}
	obs := make([][4]int64, 0, k)
	for i := 0; i < k; i++ {
		if rng.Intn(5) == 0 {
			// large rectangle
			x1 := rng.Int63n(n-2) + 1
			x2 := rng.Int63n(n-x1) + x1
			y1 := rng.Int63n(m-2) + 1
			y2 := rng.Int63n(m-y1) + y1
			if x2 >= n {
				x2 = n - 1
			}
			if y2 >= m {
				y2 = m - 1
			}
			obs = append(obs, [4]int64{x1, y1, x2, y2})
		} else {
			x := rng.Int63n(n-2) + 2
			y := rng.Int63n(m-2) + 2
			obs = append(obs, [4]int64{x, y, x, y})
		}
	}
	return manualTest(n, m, obs)
}
