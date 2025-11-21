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

const refSource = "0-999/500-599/530-539/533/533B.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
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

		if normalize(refOut) != normalize(candOut) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n",
				i+1, input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "533B-ref-*")
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

func normalize(out string) string {
	return strings.TrimSpace(out)
}

func buildTests() []string {
	tests := []string{
		"7\n-1 3\n1 2\n1 1\n1 4\n4 5\n4 3\n5 2\n", // sample
		"1\n-1 99999\n",
		"4\n-1 5\n1 4\n1 3\n1 2\n",
		"5\n-1 10\n1 8\n1 6\n2 4\n2 2\n",
	}

	randomConfigs := []struct {
		n    int
		seed int64
	}{
		{10, 1},
		{25, 2},
		{50, 3},
		{100, 4},
		{200, 5},
		{500, 6},
		{1000, 7},
		{5000, time.Now().UnixNano()},
	}
	for _, cfg := range randomConfigs {
		tests = append(tests, randomTest(cfg.n, cfg.seed))
	}
	return tests
}

func randomTest(n int, seed int64) string {
	r := rand.New(rand.NewSource(seed))
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		parent := -1
		if i > 1 {
			parent = r.Intn(i-1) + 1
		}
		val := r.Intn(100000) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", parent, val))
	}
	return sb.String()
}
