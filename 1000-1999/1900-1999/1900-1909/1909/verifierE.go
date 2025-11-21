package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const refSource = "1000-1999/1900-1999/1900-1909/1909/1909E.go"

type testCase struct {
	name  string
	input string
}

type graphTest struct {
	n, m  int
	edges [][2]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refAnswers, err := parseOutput(tc.input, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAnswers, err := parseOutput(tc.input, candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if len(refAnswers) != len(candAnswers) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d cases, got %d\ninput:\n%s\n", idx+1, tc.name, len(refAnswers), len(candAnswers), tc.input)
			os.Exit(1)
		}

		for caseIdx := range refAnswers {
			refAns := refAnswers[caseIdx]
			candAns := candAnswers[caseIdx]
			if refAns.status != candAns.status {
				fmt.Fprintf(os.Stderr, "wrong status on test %d case %d: expected %s, got %s\n", idx+1, caseIdx+1, refAns.status, candAns.status)
				os.Exit(1)
			}
			if refAns.status == -1 {
				continue
			}
			if !validAnswer(refAns.graph, candAns.nodes) {
				fmt.Fprintf(os.Stderr, "invalid solution on test %d case %d: nodes %v violate constraints\n", idx+1, caseIdx+1, candAns.nodes)
				os.Exit(1)
			}
			if len(candAns.nodes) != refAns.count {
				fmt.Fprintf(os.Stderr, "wrong size on test %d case %d: expected %d nodes, got %d\n", idx+1, caseIdx+1, refAns.count, len(candAns.nodes))
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1909E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

type answer struct {
	status int
	count  int
	nodes  []int
	graph  graphTest
}

func parseOutput(input, out string) ([]answer, error) {
	reader := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	tests := make([]graphTest, t)
	for i := 0; i < t; i++ {
		fmt.Fscan(reader, &tests[i].n, &tests[i].m)
		tests[i].edges = make([][2]int, tests[i].m)
		for j := 0; j < tests[i].m; j++ {
			fmt.Fscan(reader, &tests[i].edges[j][0], &tests[i].edges[j][1])
		}
	}
	lines := filterNonEmpty(strings.Split(out, "\n"))
	pos := 0
	ans := make([]answer, t)
	for i := 0; i < t; i++ {
		if pos >= len(lines) {
			return nil, fmt.Errorf("output ended early")
		}
		count, err := strconv.Atoi(strings.Fields(lines[pos])[0])
		if err != nil {
			return nil, fmt.Errorf("invalid count on case %d: %v", i+1, err)
		}
		pos++
		if count == -1 {
			ans[i] = answer{status: -1, graph: tests[i]}
			continue
		}
		if pos >= len(lines) {
			return nil, fmt.Errorf("missing node list for case %d", i+1)
		}
		fields := strings.Fields(lines[pos])
		pos++
		if len(fields) != count {
			return nil, fmt.Errorf("case %d expected %d nodes, got %d", i+1, count, len(fields))
		}
		nodes := make([]int, count)
		for j := 0; j < count; j++ {
			v, err := strconv.Atoi(fields[j])
			if err != nil {
				return nil, fmt.Errorf("invalid node id %q: %v", fields[j], err)
			}
			nodes[j] = v
		}
		ans[i] = answer{status: 1, count: count, nodes: nodes, graph: tests[i]}
	}
	return ans, nil
}

func validAnswer(g graphTest, nodes []int) bool {
	set := make(map[int]bool, len(nodes))
	for _, v := range nodes {
		if v < 1 || v > g.n {
			return false
		}
		if set[v] {
			return false
		}
		set[v] = true
	}
	if len(nodes) == 0 {
		return false
	}
	if len(nodes) != bits.OnesCount(uint(len(nodes))) {
		// len must be count of mask's 1 bits, but actual mask check is below
	}
	mask := 0
	for _, v := range nodes {
		mask |= 1 << (v - 1)
	}
	for _, e := range g.edges {
		x, y := e[0], e[1]
		if mask&(1<<(x-1)) != 0 && mask&(1<<(y-1)) == 0 {
			return false
		}
	}
	if len(nodes) > g.n/5 {
		return false
	}
	return true
}

func filterNonEmpty(lines []string) []string {
	var res []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			res = append(res, line)
		}
	}
	return res
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase("small1", graphTest{n: 4, m: 3, edges: [][2]int{{1, 2}, {2, 3}, {3, 4}}}),
		buildCase("small2", graphTest{n: 5, m: 5, edges: [][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 1}}}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		n := rng.Intn(15) + 5
		m := rng.Intn(n * 2)
		edges := make([][2]int, 0, m)
		seen := make(map[[2]int]bool)
		for len(edges) < m {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			if u > v {
				u, v = v, u
			}
			if seen[[2]int{u, v}] {
				continue
			}
			seen[[2]int{u, v}] = true
			edges = append(edges, [2]int{u, v})
		}
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), graphTest{n: n, m: m, edges: edges}))
	}
	return tests
}

func buildCase(name string, g graphTest) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", 1)
	fmt.Fprintf(&sb, "%d %d\n", g.n, g.m)
	for _, e := range g.edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	return testCase{name: name, input: sb.String()}
}
