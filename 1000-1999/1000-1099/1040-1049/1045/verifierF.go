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

const (
	refSourceF   = "1000-1999/1000-1099/1040-1049/1045/1045F.go"
	randomTrials = 150
	maxN         = 200000
	maxCoord     = 1_000_000_000
)

type testCase struct {
	input string
	n     int
	pts   [][2]int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceF)
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := deterministicCases()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomTrials; i++ {
		tests = append(tests, randomCase(rng))
	}

	for idx, tc := range tests {
		expect, err := runProgram(refBin, tc.input)
		if err != nil {
			fail("reference failed on case %d: %v\ninput:\n%s", idx+1, err, tc.input)
		}
		expNorm := normalize(expect)
		if err := validateOutput(expNorm); err != nil {
			fail("reference emitted invalid output on case %d: %v\ninput:\n%s\noutput:\n%s", idx+1, err, tc.input, expect)
		}
		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fail("candidate failed on case %d: %v\ninput:\n%s", idx+1, err, tc.input)
		}
		gotNorm := normalize(got)
		if err := validateOutput(gotNorm); err != nil {
			fail("candidate invalid on case %d: %v\ninput:\n%s\noutput:\n%s", idx+1, err, tc.input, got)
		}
		if expNorm != gotNorm {
			fail("case %d mismatch\ninput:\n%s\nexpected: %s\ngot: %s", idx+1, tc.input, expNorm, gotNorm)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(src string) (string, error) {
	tmp, err := os.CreateTemp("", "1045F-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(src))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	return runCommand(cmd, input)
}

func runCandidate(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	return runCommand(cmd, input)
}

func runCommand(cmd *exec.Cmd, input string) (string, error) {
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

func validateOutput(out string) error {
	trim := strings.TrimSpace(out)
	if strings.EqualFold(trim, "Ani") {
		return nil
	}
	if strings.EqualFold(trim, "Borna") {
		return nil
	}
	return fmt.Errorf("unexpected verdict %q", out)
}

func deterministicCases() []testCase {
	var tests []testCase
	tests = append(tests, buildCase([][2]int64{{0, 0}, {2, 0}, {0, 2}, {2, 2}}))
	tests = append(tests, buildCase([][2]int64{{0, 1}, {2, 3}, {4, 5}}))
	tests = append(tests, buildCase([][2]int64{{1, 0}, {0, 1}, {2, 2}}))
	tests = append(tests, buildCase([][2]int64{{0, 0}, {0, 2}, {2, 0}, {2, 2}, {4, 4}}))
	tests = append(tests, buildCase([][2]int64{{1, 1}, {3, 3}, {5, 5}, {7, 8}}))
	tests = append(tests, buildCase([][2]int64{{0, 0}, {1, 1}}))
	tests = append(tests, buildCase([][2]int64{{10, 10}, {12, 14}, {14, 12}, {16, 18}, {18, 16}}))
	tests = append(tests, buildCase([][2]int64{{0, 0}, {4, 4}, {8, 8}, {12, 12}, {16, 16}, {20, 20}}))
	tests = append(tests, buildCase([][2]int64{{2, 2}, {6, 6}, {10, 10}, {3, 5}}))
	tests = append(tests, buildCase([][2]int64{{0, 0}, {1, 1}, {1, 0}, {0, 1}, {3, 3}}))
	return tests
}

func randomCase(rng *rand.Rand) testCase {
	n := randomN(rng)
	pts := make([][2]int64, n)
	mode := rng.Intn(4)
	for i := 0; i < n; i++ {
		switch mode {
		case 0:
			pts[i] = randomPoint(rng, maxCoord, false)
		case 1:
			pts[i] = randomPoint(rng, maxCoord, true)
		case 2:
			pts[i] = [2]int64{int64(rng.Intn(1000)), int64(rng.Intn(1000))}
		case 3:
			if rng.Intn(2) == 0 {
				pts[i] = [2]int64{int64(2 * rng.Intn(500)), int64(2 * rng.Intn(500))}
			} else {
				pts[i] = [2]int64{int64(2*rng.Intn(500) + 1), int64(2*rng.Intn(500) + 1)}
			}
		}
	}
	return buildCase(pts)
}

func randomN(rng *rand.Rand) int {
	switch rng.Intn(5) {
	case 0:
		return 2
	case 1:
		return rng.Intn(5) + 3
	case 2:
		return rng.Intn(50) + 10
	case 3:
		return rng.Intn(5000) + 200
	default:
		return rng.Intn(maxN-10000) + 10000
	}
}

func randomPoint(rng *rand.Rand, limit int, evenOnly bool) [2]int64 {
	var x, y int64
	if evenOnly {
		x = int64(2 * rng.Intn(limit/2+1))
		y = int64(2 * rng.Intn(limit/2+1))
	} else {
		x = int64(rng.Intn(limit + 1))
		y = int64(rng.Intn(limit + 1))
	}
	return [2]int64{x, y}
}

func buildCase(pts [][2]int64) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(pts))
	for _, p := range pts {
		fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
	}
	return testCase{
		input: sb.String(),
		n:     len(pts),
		pts:   pts,
	}
}

func normalize(out string) string {
	trim := strings.TrimSpace(out)
	switch {
	case strings.EqualFold(trim, "Ani"):
		return "Ani"
	case strings.EqualFold(trim, "Borna"):
		return "Borna"
	default:
		return trim
	}
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
