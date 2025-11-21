package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "0-999/900-999/930-939/936/936A.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for i, input := range tests {
		refOut, err := runExecutable(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		exp, err := parseFloat(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\n%s\n", i+1, err, refOut)
			os.Exit(1)
		}
		got, err := parseFloat(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\nOutput:\n%s\n", i+1, err, candOut)
			os.Exit(1)
		}
		if !closeEnough(exp, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected: %.12f\ngot: %.12f\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "936A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var combined bytes.Buffer
	cmd.Stdout = &combined
	cmd.Stderr = &combined
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, combined.String())
	}
	return tmp.Name(), nil
}

func runExecutable(path, input string) (string, error) {
	cmd := exec.Command(path)
	return execute(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return execute(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func execute(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseFloat(out string) (float64, error) {
	var val float64
	reader := strings.NewReader(out)
	if _, err := fmt.Fscan(reader, &val); err != nil {
		return 0, err
	}
	return val, nil
}

func closeEnough(expected, actual float64) bool {
	diff := math.Abs(expected - actual)
	if diff <= 1e-9 {
		return true
	}
	scale := math.Max(1.0, math.Abs(expected))
	return diff/scale <= 1e-9
}

func buildTests() []string {
	tests := []string{
		"3 2 6\n",
		"4 2 20\n",
		"5 10 100\n",
		"10 3 7\n",
		"1 1 1\n",
		"1000000000000000000 1 1000000000000000000\n",
		"1 1000000000000000000 1000000000000000000\n",
		"999999999999999999 999999999999999998 999999999999999997\n",
	}

	randomConfigs := []struct {
		seed int64
	}{
		{1},
		{2},
		{3},
		{time.Now().UnixNano()},
	}

	for _, cfg := range randomConfigs {
		tests = append(tests, randomTest(10, cfg.seed)...)
	}
	return tests
}

func randomTest(count int, seed int64) []string {
	r := rand.New(rand.NewSource(seed))
	tests := make([]string, 0, count)
	for i := 0; i < count; i++ {
		k := randRange(r, 1, 1_000_000_000_000_000_000)
		d := randRange(r, 1, 1_000_000_000_000_000_000)
		t := randRange(r, 1, 1_000_000_000_000_000_000)
		tests = append(tests, fmt.Sprintf("%d %d %d\n", k, d, t))
	}
	return tests
}

func randRange(r *rand.Rand, lo, hi int64) int64 {
	return r.Int63n(hi-lo+1) + lo
}
