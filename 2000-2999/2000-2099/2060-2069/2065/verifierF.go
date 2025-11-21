package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const refSourceF = "2000-2999/2000-2099/2060-2069/2065/2065F.go"

type testCaseF struct {
	n     int
	a     []int
	edges [][2]int
	name  string
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReferenceF()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTestsF()
	input := buildInputF(tests)

	refOut, err := runProgramF(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseOutputF(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgramF(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutputF(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "case %d (%s) mismatch:\nexpected: %s\ngot:      %s\ninput:\n%s", i+1, tc.name, refAns[i], candAns[i], formatCaseF(tc))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReferenceF() (string, error) {
	outPath := "./ref_2065F.bin"
	cmd := exec.Command("go", "build", "-o", outPath, refSourceF)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgramF(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildTestsF() []testCaseF {
	var tests []testCaseF
	add := func(name string, n int, a []int, edges [][2]int) {
		cpA := append([]int(nil), a...)
		cpE := make([][2]int, len(edges))
		copy(cpE, edges)
		tests = append(tests, testCaseF{name: name, n: n, a: cpA, edges: cpE})
	}

	// Small handcrafted trees
	add("two_equal", 2, []int{1, 1}, [][2]int{{1, 2}})
	add("two_diff", 2, []int{1, 2}, [][2]int{{1, 2}})
	add("chain_three", 3, []int{1, 1, 2}, [][2]int{{1, 2}, {2, 3}})
	add("star_four", 4, []int{2, 2, 2, 3}, [][2]int{{1, 2}, {1, 3}, {1, 4}})
	add("mixed_five", 5, []int{1, 2, 2, 3, 2}, [][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalN := 0
	const maxTotalN = 120000
	for len(tests) < 120 && totalN < maxTotalN {
		n := rng.Intn(80) + 2 // 2..81 small/medium for speed
		if len(tests)%15 == 0 {
			n = rng.Intn(500) + 100 // occasional larger cases
		}
		if totalN+n > maxTotalN {
			break
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(n) + 1
		}
		edges := make([][2]int, 0, n-1)
		for v := 2; v <= n; v++ {
			p := rng.Intn(v-1) + 1
			edges = append(edges, [2]int{p, v})
		}
		add(fmt.Sprintf("random_%d", len(tests)), n, a, edges)
		totalN += n
	}
	return tests
}

func buildInputF(tests []testCaseF) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
	}
	return sb.String()
}

func parseOutputF(out string, tests []testCaseF) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != len(tests) {
		return nil, fmt.Errorf("expected %d outputs, got %d", len(tests), len(lines))
	}
	for i, s := range lines {
		if len(s) != tests[i].n {
			return nil, fmt.Errorf("case %d length mismatch: expected %d chars, got %d in %q", i+1, tests[i].n, len(s), s)
		}
		for _, ch := range s {
			if ch != '0' && ch != '1' {
				return nil, fmt.Errorf("case %d has invalid character %q", i+1, ch)
			}
		}
	}
	return lines, nil
}

func formatCaseF(tc testCaseF) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}
