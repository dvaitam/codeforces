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

const refSource = "2000-2999/2000-2099/2010-2019/2014/2014H.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/candidate")
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
	tmp, err := os.CreateTemp("", "2014H-ref-*")
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
	return strings.Join(strings.Fields(strings.ToUpper(out)), "\n")
}

func buildTests() []string {
	tests := []string{
		"1\n3 3\n1 2 2\n1 2\n1 3\n2 3\n",
		"1\n5 3\n3 2 1 2 3\n1 5\n2 4\n3 5\n",
		"1\n2 1\n1000000 1000000\n1 2\n",
	}

	randomConfigs := []struct {
		t    int
		sumN int
		sumQ int
		seed int64
	}{
		{3, 50, 50, 1},
		{5, 200, 200, 2},
		{8, 2000, 2000, 3},
		{10, 5000, 5000, 4},
		{12, 200000, 200000, 5},
		{9, 200000, 200000, time.Now().UnixNano()},
	}
	for _, cfg := range randomConfigs {
		tests = append(tests, randomTest(cfg.t, cfg.sumN, cfg.sumQ, cfg.seed))
	}
	return tests
}

func randomTest(t, sumN, sumQ int, seed int64) string {
	if t < 1 {
		t = 1
	}
	r := rand.New(rand.NewSource(seed))
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	remN := sumN
	remQ := sumQ
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		casesLeft := t - caseIdx - 1
		minRemainN := casesLeft
		minRemainQ := casesLeft
		maxN := remN - minRemainN
		maxQ := remQ - minRemainQ
		if maxN < 1 {
			maxN = 1
		}
		if maxQ < 1 {
			maxQ = 1
		}
		if maxN > 200000 {
			maxN = 200000
		}
		if maxQ > 200000 {
			maxQ = 200000
		}
		n := r.Intn(maxN) + 1
		q := r.Intn(maxQ) + 1
		remN -= n
		remQ -= q
		sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
		for i := 0; i < n; i++ {
			val := randomValue(r, i)
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", val))
		}
		sb.WriteByte('\n')
		for i := 0; i < q; i++ {
			l := r.Intn(n) + 1
			rv := r.Intn(n-l+1) + l
			sb.WriteString(fmt.Sprintf("%d %d\n", l, rv))
		}
	}
	return sb.String()
}

func randomValue(r *rand.Rand, idx int) int {
	switch r.Intn(4) {
	case 0:
		return r.Intn(10) + 1
	case 1:
		return r.Intn(1000) + 1
	case 2:
		return r.Intn(1000000) + 1
	default:
		base := r.Intn(500000) + 1
		return base
	}
}
