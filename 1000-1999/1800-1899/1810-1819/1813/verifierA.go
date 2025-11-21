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

const (
	refSource    = "1000-1999/1800-1899/1810-1819/1813/1813A.go"
	numTestCases = 150
)

type test struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary-or-source")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for i, tc := range tests {
		expect, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\n", i+1, tc.name, err)
			os.Exit(1)
		}
		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%s", i+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if normalize(expect) != normalize(got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s)\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n", i+1, tc.name, tc.input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	tmp, err := os.CreateTemp("", "ref1813A-*")
	if err != nil {
		return "", nil, fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmpPath, refSource)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmpPath)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, string(out))
	}
	cleanup := func() { os.Remove(tmpPath) }
	return tmpPath, cleanup, nil
}

func runProgram(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func runCandidate(target string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		abs, err := filepath.Abs(target)
		if err != nil {
			return "", fmt.Errorf("resolve candidate: %w", err)
		}
		cmd = exec.Command("go", "run", abs)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func normalize(out string) string {
	lines := strings.Fields(out)
	return strings.Join(lines, "\n")
}

func generateTests() []test {
	var tests []test
	tests = append(tests, test{
		name:  "minimal",
		input: buildInput(1, 1, 1, []int{1}, []int{1}, []triplet{{1, 1, 1}}, []op{{1, 1}}),
	})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < numTestCases {
		n := rng.Intn(5) + 1
		qMin := make([]int, n)
		sumMin := 0
		for i := 0; i < n; i++ {
			qMin[i] = rng.Intn(3) + 1
			sumMin += qMin[i]
		}
		q := sumMin + rng.Intn(5) + 1
		m := rng.Intn(10) + 1
		priority := make([]int, n)
		dbSize := make([]int, n)
		qTrip := make([]triplet, n)
		for i := 0; i < n; i++ {
			priority[i] = rng.Intn(10) + 1
			dbSize[i] = rng.Intn(20) + 1
			base := qMin[i] + rng.Intn(3)
			max := base + rng.Intn(3)
			qTrip[i] = triplet{qMin[i], base, max}
		}
		ops := make([]op, m)
		for i := 0; i < m; i++ {
			u := rng.Intn(n) + 1
			p := rng.Intn(dbSize[u-1]) + 1
			ops[i] = op{u, p}
		}
		tests = append(tests, test{
			name:  fmt.Sprintf("rand_%d", len(tests)),
			input: buildInput(n, q, m, priority, dbSize, qTrip, ops),
		})
	}
	return tests
}

type triplet struct {
	lo, base, hi int
}

type op struct {
	u, p int
}

func buildInput(n, q, m int, pri, size []int, triplets []triplet, ops []op) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, q, m))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(pri[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(size[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d %d %d", triplets[i].lo, triplets[i].base, triplets[i].hi))
	}
	sb.WriteByte('\n')
	for _, op := range ops {
		sb.WriteString(fmt.Sprintf("%d %d\n", op.u, op.p))
	}
	return sb.String()
}
