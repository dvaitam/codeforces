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

const refSource = "2000-2999/2000-2099/2000-2009/2002/2002G.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/candidate")
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
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2002G-ref-*")
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
		"2\n3\n1 0 2\n0 1 3\n2 1 0\n3 0 3\n1 2\n0 1\n2 0\n3\n1 2 0\n0 1 2\n2 0 1\n1 2\n0 1\n2 0\n",
		"1\n2\n0 1\n1 0\n0\n1\n",
	}

	randomConfigs := []struct {
		t    int
		maxN int
		seed int64
	}{
		{3, 4, 1},
		{5, 6, 2},
		{6, 8, 3},
		{8, 10, 4},
		{10, 12, 5},
		{12, 15, 6},
		{15, 18, time.Now().UnixNano()},
	}
	for _, cfg := range randomConfigs {
		tests = append(tests, randomTest(cfg.t, cfg.maxN, cfg.seed))
	}
	return tests
}

func randomTest(t, maxN int, seed int64) string {
	if t < 1 {
		t = 1
	}
	r := rand.New(rand.NewSource(seed))
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	sumCubes := 0
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		remain := t - caseIdx
		minFuture := 8 * (remain - 1)
		if minFuture < 0 {
			minFuture = 0
		}
		limit := 8000 - sumCubes - minFuture
		if limit < 8 {
			limit = 8
		}
		maxAllowed := int(math.Floor(math.Pow(float64(limit), 1.0/3.0)))
		if maxAllowed > maxN {
			maxAllowed = maxN
		}
		if maxAllowed < 2 {
			maxAllowed = 2
		}
		n := 2 + r.Intn(maxAllowed-1)
		sumCubes += n * n * n
		sb.WriteString(fmt.Sprintf("%d\n", n))
		maxVal := 2*n - 2
		for i := 0; i < n-1; i++ {
			for j := 0; j < n; j++ {
				val := r.Intn(maxVal + 1)
				if j+1 == n {
					sb.WriteString(fmt.Sprintf("%d\n", val))
				} else {
					sb.WriteString(fmt.Sprintf("%d ", val))
				}
			}
		}
		for i := 0; i < n; i++ {
			for j := 0; j < n-1; j++ {
				val := r.Intn(maxVal + 1)
				if j+1 == n-1 {
					sb.WriteString(fmt.Sprintf("%d\n", val))
				} else {
					sb.WriteString(fmt.Sprintf("%d ", val))
				}
			}
		}
	}
	return sb.String()
}
