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

const refSourceG = "2000-2999/2000-2099/2000-2009/2004/2004G.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReferenceG()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTestsG()
	for i, input := range tests {
		refOut, err := runExecutableG(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}

		candOut, err := runCandidateG(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}

		if normalizeG(refOut) != normalizeG(candOut) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReferenceG() (string, error) {
	tmp, err := os.CreateTemp("", "2004G-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceG))
	var combined bytes.Buffer
	cmd.Stdout = &combined
	cmd.Stderr = &combined
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, combined.String())
	}
	return tmp.Name(), nil
}

func runExecutableG(path, input string) (string, error) {
	cmd := exec.Command(path)
	return executeG(cmd, input)
}

func runCandidateG(path, input string) (string, error) {
	cmd := commandForG(path)
	return executeG(cmd, input)
}

func commandForG(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func executeG(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func normalizeG(out string) string {
	return strings.Join(strings.Fields(out), " ")
}

func buildTestsG() []string {
	tests := []string{
		"4 4\n5999\n",
		"10 3\n1111111111\n",
		"11 4\n99998641312\n",
		"6 2\n123456\n",
		"6 5\n999999\n",
	}

	randomConfigs := []struct {
		n    int
		kmin int
		seed int64
	}{
		{10, 2, 1},
		{50, 3, 2},
		{100, 5, 3},
		{500, 10, 4},
		{2000, 30, 5},
		{200000, 1000, 6},
		{150000, 500, time.Now().UnixNano()},
	}
	for _, cfg := range randomConfigs {
		tests = append(tests, randomTestG(cfg.n, cfg.kmin, cfg.seed))
	}
	return tests
}

func randomTestG(n, kmin int, seed int64) string {
	if n < 2 {
		n = 2
	}
	if kmin < 2 {
		kmin = 2
	}
	if kmin > n {
		kmin = n
	}
	r := rand.New(rand.NewSource(seed))
	k := kmin + r.Intn(n-kmin+1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		digit := r.Intn(9) + 1
		sb.WriteByte(byte('0' + digit))
	}
	sb.WriteByte('\n')
	return sb.String()
}
