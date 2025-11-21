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

const refSource = "2000-2999/2000-2099/2010-2019/2013/2013D.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
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
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2013D-ref-*")
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
		"5\n1\n131\n3\n1 2 3\n4\n4 1 2 3\n4\n4 2 3 1\n5\n1 4 4 10 2\n",
		"1\n1\n1000000000000\n",
		"2\n2\n1 1\n2\n1 1000000000000\n",
	}

	randomConfigs := []struct {
		t      int
		maxSum int
		seed   int64
	}{
		{5, 50, 1},
		{10, 200, 2},
		{20, 2000, 3},
		{40, 200000, 4},
		{60, 200000, 5},
		{80, 200000, time.Now().UnixNano()},
	}
	for _, cfg := range randomConfigs {
		tests = append(tests, randomTest(cfg.t, cfg.maxSum, cfg.seed))
	}
	return tests
}

func randomTest(t, maxSum int, seed int64) string {
	if t < 1 {
		t = 1
	}
	if t > 100000 {
		t = 100000
	}
	r := rand.New(rand.NewSource(seed))
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	sumRemain := maxSum
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		need := t - caseIdx - 1
		if need < 0 {
			need = 0
		}
		if sumRemain <= need {
			sumRemain = need + 1
		}
		maxN := sumRemain - need
		if maxN < 1 {
			maxN = 1
		}
		if maxN > 200000 {
			maxN = 200000
		}
		n := r.Intn(maxN) + 1
		sumRemain -= n
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			val := randomValue(r, i)
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", val))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func randomValue(r *rand.Rand, idx int) int64 {
	switch r.Intn(4) {
	case 0:
		return int64(r.Intn(1000) + 1)
	case 1:
		return int64(r.Intn(1_000_000) + 1)
	case 2:
		return int64(r.Intn(1_000_000_000) + 1)
	default:
		base := int64(r.Intn(1_000_000) + 1)
		return base * 1000000
	}
}
