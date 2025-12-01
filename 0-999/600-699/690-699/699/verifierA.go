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

const refSource = "699A.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		exp, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		got, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, exp, got, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-699A-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref699A.bin")
	source := filepath.Join(".", refSource)
	cmd := exec.Command("go", "build", "-o", binPath, source)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer, got %d tokens", len(fields))
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", fields[0])
	}
	return val, nil
}

func buildTests() []testCase {
	tests := []testCase{
		formatManual("no_collision", []byte("LR"), []int{0, 2}),
		formatManual("simple_collision", []byte("RL"), []int{0, 2}),
		formatManual("three_particles", []byte("RRL"), []int{0, 2, 4}),
		formatManual("all_left", []byte("LLL"), []int{0, 2, 4}),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 300; i++ {
		tests = append(tests, randomCase(rng, i))
	}
	tests = append(tests, stressCase())
	return tests
}

func formatManual(name string, dirs []byte, pos []int) testCase {
	var sb strings.Builder
	n := len(dirs)
	fmt.Fprintf(&sb, "%d\n", n)
	sb.Write(dirs)
	sb.WriteByte('\n')
	for i, p := range pos {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(p))
	}
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String()}
}

func randomCase(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(200000) + 1
	dirs := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			dirs[i] = 'L'
		} else {
			dirs[i] = 'R'
		}
	}
	pos := make([]int, n)
	cur := rng.Intn(10)
	for i := 0; i < n; i++ {
		cur += rng.Intn(1000)*2 + 2
		pos[i] = cur
	}
	return testCase{
		name:  fmt.Sprintf("random_%d_n%d", idx+1, n),
		input: formatManual("", dirs, pos).input,
	}
}

func stressCase() testCase {
	n := 200000
	dirs := make([]byte, n)
	pos := make([]int, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			dirs[i] = 'R'
		} else {
			dirs[i] = 'L'
		}
		pos[i] = i * 2
	}
	return testCase{name: "stress_max", input: formatManual("", dirs, pos).input}
}
