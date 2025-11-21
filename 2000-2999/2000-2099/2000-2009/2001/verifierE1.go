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

const refSource = "2000-2999/2000-2099/2000-2009/2001/2001E1.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE1.go /path/to/candidate")
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
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n", i+1, input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2001E1-ref-*")
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
		"1\n1 1 998244353\n",
		"1\n2 1 998244353\n",
		"1\n3 3 100000007\n",
		"3\n1 1 998244353\n2 1 998244353\n3 2 998244853\n",
	}

	randomConfigs := []struct {
		t    int
		seed int64
	}{
		{5, 1},
		{10, 2},
		{20, 3},
		{50, 4},
		{100, 5},
		{150, 6},
		{200, time.Now().UnixNano()},
	}

	for _, cfg := range randomConfigs {
		tests = append(tests, randomTest(cfg.t, cfg.seed))
	}
	return tests
}

func randomTest(t int, seed int64) string {
	if t < 1 {
		t = 1
	}
	if t > 500 {
		t = 500
	}
	r := rand.New(rand.NewSource(seed))
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	sumN, sumK := 0, 0
	for i := 0; i < t; i++ {
		remCases := t - i
		remN := 500 - sumN
		remK := 500 - sumK
		maxN := min(500, remN-(remCases-1))
		if maxN < 1 {
			maxN = 1
		}
		maxK := min(500, remK-(remCases-1))
		if maxK < 1 {
			maxK = 1
		}
		n := r.Intn(maxN) + 1
		k := r.Intn(maxK) + 1
		sumN += n
		sumK += k
		p := randomPrime(r)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, k, p))
	}
	return sb.String()
}

func randomPrime(r *rand.Rand) int64 {
	primes := []int64{
		998244353, 998244853, 100000007, 100000037, 100000039, 100000007, 1000003, 1000033, 1000037, 1000039,
	}
	return primes[r.Intn(len(primes))]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
