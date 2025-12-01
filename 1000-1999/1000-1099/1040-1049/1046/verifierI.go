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

const (
	refSource  = "./1046I.go"
	totalTests = 80
)

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refVal, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVal, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if candVal != refVal {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d, got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
				idx+1, tc.name, refVal, candVal, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "ref-1046I-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref1046I.bin")
	cmd := exec.Command("go", "build", "-o", bin, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return bin, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseOutput(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer, got %d tokens", len(fields))
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", fields[0])
	}
	return val, nil
}

func buildTests() []testCase {
	tests := []testCase{
		{
			name:  "static_close",
			input: "2\n1 2\n0 0 0 0\n0 0 0 0\n",
		},
		{
			name:  "move_in_out",
			input: "4\n1 3\n0 0 0 0\n0 0 0 4\n0 0 0 0\n0 3 0 0\n",
		},
		{
			name:  "never_meet",
			input: "3\n1 5\n0 0 5 5\n0 0 5 5\n0 0 5 5\n",
		},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-5 {
		tests = append(tests, randomTest(rng, len(tests)+1, 50))
	}
	tests = append(tests,
		randomTest(rand.New(rand.NewSource(1)), len(tests)+1, 200),
		randomTest(rand.New(rand.NewSource(2)), len(tests)+2, 500),
		randomTest(rand.New(rand.NewSource(3)), len(tests)+3, 1000),
		randomLarge(rand.New(rand.NewSource(4)), 5000),
		randomLarge(rand.New(rand.NewSource(5)), 10000),
	)
	return tests
}

func randomTest(rng *rand.Rand, idx int, maxN int) testCase {
	n := rng.Intn(maxN-2) + 2
	d1 := rng.Intn(50) + 1
	d2 := d1 + rng.Intn(50) + 1
	lines := make([]string, 0, n+2)
	lines = append(lines, fmt.Sprintf("%d", n))
	lines = append(lines, fmt.Sprintf("%d %d", d1, d2))
	ax, ay := rng.Intn(100), rng.Intn(100)
	bx, by := rng.Intn(100), rng.Intn(100)
	lines = append(lines, fmt.Sprintf("%d %d %d %d", ax, ay, bx, by))
	for i := 1; i < n; i++ {
		if rng.Intn(2) == 0 {
			ax += rng.Intn(11) - 5
			ay += rng.Intn(11) - 5
			if ax < 0 {
				ax = 0
			}
			if ay < 0 {
				ay = 0
			}
		}
		if rng.Intn(2) == 0 {
			bx += rng.Intn(11) - 5
			by += rng.Intn(11) - 5
			if bx < 0 {
				bx = 0
			}
			if by < 0 {
				by = 0
			}
		}
		if ax > 1000 {
			ax = 1000
		}
		if ay > 1000 {
			ay = 1000
		}
		if bx > 1000 {
			bx = 1000
		}
		if by > 1000 {
			by = 1000
		}
		lines = append(lines, fmt.Sprintf("%d %d %d %d", ax, ay, bx, by))
	}
	return testCase{
		name:  fmt.Sprintf("random_%d", idx),
		input: strings.Join(lines, "\n") + "\n",
	}
}

func randomLarge(rng *rand.Rand, n int) testCase {
	d1 := rng.Intn(400) + 1
	d2 := d1 + rng.Intn(400) + 1
	lines := make([]string, 0, n+2)
	lines = append(lines, fmt.Sprintf("%d", n))
	lines = append(lines, fmt.Sprintf("%d %d", d1, d2))
	ax, ay := rng.Intn(1000), rng.Intn(1000)
	bx, by := rng.Intn(1000), rng.Intn(1000)
	lines = append(lines, fmt.Sprintf("%d %d %d %d", ax, ay, bx, by))
	for i := 1; i < n; i++ {
		ax += rng.Intn(21) - 10
		ay += rng.Intn(21) - 10
		bx += rng.Intn(21) - 10
		by += rng.Intn(21) - 10
		if ax < 0 {
			ax = 0
		}
		if ay < 0 {
			ay = 0
		}
		if bx < 0 {
			bx = 0
		}
		if by < 0 {
			by = 0
		}
		if ax > 1000 {
			ax = 1000
		}
		if ay > 1000 {
			ay = 1000
		}
		if bx > 1000 {
			bx = 1000
		}
		if by > 1000 {
			by = 1000
		}
		lines = append(lines, fmt.Sprintf("%d %d %d %d", ax, ay, bx, by))
	}
	return testCase{
		name:  fmt.Sprintf("large_%d", n),
		input: strings.Join(lines, "\n") + "\n",
	}
}
