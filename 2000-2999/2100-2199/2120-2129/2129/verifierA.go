package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "./2129A.go"

type testCase struct {
	cases []caseData
	input string
}

type caseData struct {
	a, b     []int
	maxCoord int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	candidate := os.Args[1]
	tests := generateTests()

	for i, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		refObjs, err := extractObjectives(tc, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, tc.input, refOut)
			os.Exit(1)
		}

		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, got)
			os.Exit(1)
		}

		userObjs, err := extractObjectives(tc, got)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid output on test %d: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, tc.input, got)
			os.Exit(1)
		}

		if len(refObjs) != len(userObjs) {
			fmt.Fprintf(os.Stderr, "wrong number of testcases in output on test %d\ninput:\n%s\n", i+1, tc.input)
			os.Exit(1)
		}
		for j := range refObjs {
			if refObjs[j] != userObjs[j] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d: expected objective %d, got %d\ninput:\n%s\n", i+1, j+1, refObjs[j], userObjs[j], tc.input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2129A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(21292129))
	var tests []testCase

	// Sample-like small cases.
	tests = append(tests, makeTest([]caseData{
		buildCase([][2]int{{1, 2}}),
		buildCase([][2]int{{1, 2}, {2, 3}, {1, 3}, {3, 5}}),
	}))

	tests = append(tests, makeTest([]caseData{
		randomCase(rng, 5),
		randomCase(rng, 6),
	}))

	// Mix of sizes.
	tests = append(tests, makeTest([]caseData{
		randomCaseWithSize(rng, 50),
		randomCaseWithSize(rng, 80),
	}))

	tests = append(tests, makeTest([]caseData{
		randomCaseWithSize(rng, 200),
	}))

	tests = append(tests, makeTest([]caseData{
		randomCaseWithSize(rng, 400),
	}))

	// Stress near limits while keeping sum n^2 manageable.
	tests = append(tests, makeTest([]caseData{
		randomCaseWithSize(rng, 600),
	}))

	return tests
}

func makeTest(cs []caseData) testCase {
	return testCase{cases: cs, input: buildInput(cs)}
}

func buildInput(cs []caseData) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(cs))
	for _, c := range cs {
		n := len(c.a)
		fmt.Fprintf(&b, "%d\n", n)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&b, "%d %d\n", c.a[i], c.b[i])
		}
	}
	return b.String()
}

func buildCase(intervals [][2]int) caseData {
	n := len(intervals)
	a := make([]int, n)
	b := make([]int, n)
	maxC := 0
	for i, p := range intervals {
		a[i], b[i] = p[0], p[1]
		if b[i] > maxC {
			maxC = b[i]
		}
	}
	return caseData{a: a, b: b, maxCoord: maxC}
}

func randomCase(rng *rand.Rand, n int) caseData {
	return randomCaseWithSize(rng, n)
}

func randomCaseWithSize(rng *rand.Rand, n int) caseData {
	type pair struct{ a, b int }
	seen := make(map[pair]bool)
	a := make([]int, n)
	b := make([]int, n)
	maxC := 0
	for i := 0; i < n; i++ {
		for {
			x := rng.Intn(2*n-1) + 1
			y := rng.Intn(2*n-x) + x + 1
			p := pair{x, y}
			if !seen[p] {
				seen[p] = true
				a[i], b[i] = x, y
				if y > maxC {
					maxC = y
				}
				break
			}
		}
	}
	return caseData{a: a, b: b, maxCoord: maxC}
}

func extractObjectives(tc testCase, out string) ([]int, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	tOut, err := nextInt(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read number of test cases: %v", err)
	}
	if tOut != len(tc.cases) {
		return nil, fmt.Errorf("expected %d test cases, got %d", len(tc.cases), tOut)
	}

	objs := make([]int, len(tc.cases))
	for idx, cs := range tc.cases {
		k, err := nextInt(reader)
		if err != nil {
			return nil, fmt.Errorf("test %d: failed to read k: %v", idx+1, err)
		}
		if k < 0 || k > len(cs.a) {
			return nil, fmt.Errorf("test %d: k=%d out of bounds", idx+1, k)
		}
		seen := make([]bool, len(cs.a))
		selected := make([]int, k)
		for i := 0; i < k; i++ {
			id, err := nextInt(reader)
			if err != nil {
				return nil, fmt.Errorf("test %d: failed to read index %d: %v", idx+1, i+1, err)
			}
			if id < 1 || id > len(cs.a) {
				return nil, fmt.Errorf("test %d: index %d out of range", idx+1, id)
			}
			if seen[id-1] {
				return nil, fmt.Errorf("test %d: duplicate index %d", idx+1, id)
			}
			seen[id-1] = true
			selected[i] = id - 1
		}
		obj := computeObjective(cs, selected)
		objs[idx] = obj
	}

	if extra, err := nextInt(reader); err == nil {
		return nil, fmt.Errorf("unexpected extra output, next int=%d", extra)
	}

	return objs, nil
}

func computeObjective(cs caseData, selected []int) int {
	maxC := cs.maxCoord
	diff := make([]int, maxC+2)
	adj := make([][]int, maxC+1)
	deg := make([]int, maxC+1)

	for _, idx := range selected {
		a, b := cs.a[idx], cs.b[idx]
		diff[a]++
		diff[b]--
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
		deg[a]++
		deg[b]++
	}

	cover := 0
	cur := 0
	for x := 1; x <= maxC; x++ {
		cur += diff[x]
		if cur > 0 {
			cover++
		}
	}

	queue := make([]int, 0)
	inQ := make([]bool, maxC+1)
	for v := 1; v <= maxC; v++ {
		if deg[v] > 0 && deg[v] <= 1 {
			queue = append(queue, v)
			inQ[v] = true
		}
	}
	removed := make([]bool, maxC+1)
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		if removed[v] {
			continue
		}
		removed[v] = true
		for _, to := range adj[v] {
			if removed[to] {
				continue
			}
			deg[to]--
			if deg[to] == 1 && !inQ[to] {
				queue = append(queue, to)
				inQ[to] = true
			}
		}
	}

	cycleNodes := 0
	for v := 1; v <= maxC; v++ {
		if deg[v] > 0 && !removed[v] {
			cycleNodes++
		}
	}

	return cover - cycleNodes
}

func nextInt(r *bufio.Reader) (int, error) {
	sign, val := 1, 0
	for {
		c, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		if c == '-' {
			sign = -1
			continue
		}
		if c >= '0' && c <= '9' {
			val = int(c - '0')
			for {
				c, err = r.ReadByte()
				if err != nil {
					return sign * val, nil
				}
				if c < '0' || c > '9' {
					if err := r.UnreadByte(); err != nil {
						return 0, err
					}
					break
				}
				val = val*10 + int(c-'0')
			}
			return sign * val, nil
		}
		if c > ' ' {
			return 0, fmt.Errorf("unexpected character %q", c)
		}
	}
}
