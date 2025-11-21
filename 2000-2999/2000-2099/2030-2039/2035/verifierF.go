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

type testInput struct {
	text string
}

func buildReference() (string, error) {
	refDir := filepath.Join("2000-2999", "2000-2099", "2030-2039", "2035")
	tmp, err := os.CreateTemp("", "ref2035F")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	os.Remove(tmpPath)

	cmd := exec.Command("go", "build", "-o", tmpPath, "2035F.go")
	cmd.Dir = refDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return tmpPath, nil
}

func commandForPath(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	case ".js":
		return exec.Command("node", path)
	default:
		return exec.Command(path)
	}
}

func runBinary(path, input string) (string, error) {
	cmd := commandForPath(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func normalizeOutput(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return strings.Join(lines, "\n")
}

func fixedTests() []testInput {
	return []testInput{
		{"1\n1 1\n0\n"},
		{"1\n2 1\n1 2\n1 2\n"},
		{"1\n3 2\n1 2 3\n1 2\n2 3\n"},
	}
}

func randomTree(n int, rng *rand.Rand) [][2]int {
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		p := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{p, v})
	}
	return edges
}

func buildCase(n int, root int, values []int64, edges [][2]int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, root))
	for i, v := range values {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func randomTests() []testInput {
	tests := fixedTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Add a larger deterministic case to exercise deeper recursion while staying
	// within reference solver limits.
	nMax := 400
	values := make([]int64, nMax)
	edges := make([][2]int, 0, nMax-1)
	limitRng := rand.New(rand.NewSource(12345))
	for i := range values {
		values[i] = int64(i % 50)
		if i > 0 {
			edges = append(edges, [2]int{limitRng.Intn(i) + 1, i + 1})
		}
	}
	var sbLimit strings.Builder
	sbLimit.WriteString("1\n")
	sbLimit.WriteString(buildCase(nMax, nMax, values, edges))
	tests = append(tests, testInput{text: sbLimit.String()})

	for len(tests) < 60 {
		cases := rng.Intn(3) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", cases))
		for i := 0; i < cases; i++ {
			n := rng.Intn(10) + 1
			if rng.Intn(5) == 0 {
				n = rng.Intn(50) + 1
			}
			root := rng.Intn(n) + 1
			values := make([]int64, n)
			for j := 0; j < n; j++ {
				values[j] = rng.Int63n(1000)
			}
			edges := randomTree(n, rng)
			sb.WriteString(buildCase(n, root, values, edges))
			if i+1 < cases {
				sb.WriteByte('\n')
			}
		}
		tests = append(tests, testInput{text: sb.String()})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := randomTests()
	for idx, input := range tests {
		expect, err := runBinary(ref, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, input.text)
			os.Exit(1)
		}
		got, err := runBinary(bin, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input.text)
			os.Exit(1)
		}
		if normalizeOutput(expect) != normalizeOutput(got) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, input.text, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
