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

const refSource = "1000-1999/1800-1899/1840-1849/1842/1842E.go"

type testCase struct {
	name  string
	input string
}

type point struct {
	x int
	y int
	c int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refAns, err := parseAnswer(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parseAnswer(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if refAns != candAns {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d, got %d\ninput:\n%s\n", idx+1, tc.name, refAns, candAns, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1842E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
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
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func parseAnswer(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected 1 integer, got %d tokens", len(fields))
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	return val, nil
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase("single-point", 1, 5, 3, []point{{x: 1, y: 1, c: 10}}),
		buildCase("triangle-cheap", 3, 8, 1, []point{
			{x: 0, y: 0, c: 100},
			{x: 3, y: 2, c: 5},
			{x: 1, y: 1, c: 4},
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 60; i++ {
		n := rng.Intn(6) + 1
		k := rng.Intn(20) + n + 5
		A := rng.Intn(50) + 1
		pts := make([]point, n)
		for j := 0; j < n; j++ {
			x := rng.Intn(k)
			maxY := k - x - 1
			if maxY < 0 {
				maxY = 0
			}
			y := 0
			if maxY > 0 {
				y = rng.Intn(maxY + 1)
			}
			if x+y >= k {
				y = k - x - 1
				if y < 0 {
					y = 0
				}
			}
			pts[j] = point{x: x, y: y, c: rng.Intn(100) + 1}
		}
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), n, k, A, pts))
	}
	return tests
}

func buildCase(name string, n, k, A int, pts []point) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, k, A)
	for _, p := range pts {
		if p.x+p.y >= k {
			// ensure valid input, clamp y if necessary
			p.y = k - p.x - 1
			if p.y < 0 {
				p.y = 0
			}
		}
		fmt.Fprintf(&sb, "%d %d %d\n", p.x, p.y, p.c)
	}
	return testCase{name: name, input: sb.String()}
}
