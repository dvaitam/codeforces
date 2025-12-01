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

const refSource = "./538E.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
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
	tmp, err := os.CreateTemp("", "538E-ref-*")
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
		"5\n1 2\n1 3\n2 4\n2 5\n",
		"6\n1 2\n1 3\n3 4\n1 5\n5 6\n",
		"1\n",
		lineTree(8),
		starTree(10),
		balancedTree(15, 3),
	}

	randomConfigs := []struct {
		n    int
		seed int64
	}{
		{20, 1},
		{50, 2},
		{100, 3},
		{300, 4},
		{1000, 5},
		{5000, time.Now().UnixNano()},
	}
	for _, cfg := range randomConfigs {
		tests = append(tests, randomTree(cfg.n, cfg.seed))
	}
	return tests
}

func lineTree(n int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 2; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", i-1, i))
	}
	return sb.String()
}

func starTree(n int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 2; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("1 %d\n", i))
	}
	return sb.String()
}

func balancedTree(n, branching int) string {
	children := make([][]int, n+1)
	q := []int{1}
	cur := 2
	for len(q) > 0 && cur <= n {
		v := q[0]
		q = q[1:]
		for j := 0; j < branching && cur <= n; j++ {
			children[v] = append(children[v], cur)
			q = append(q, cur)
			cur++
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for v := 1; v <= n; v++ {
		for _, u := range children[v] {
			sb.WriteString(fmt.Sprintf("%d %d\n", v, u))
		}
	}
	return sb.String()
}

func randomTree(n int, seed int64) string {
	r := rand.New(rand.NewSource(seed))
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for v := 2; v <= n; v++ {
		parent := r.Intn(v-1) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", parent, v))
	}
	return sb.String()
}
