package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	input string
	n     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierGP1.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refSrc, err := locateReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refBin, err := buildReference(refSrc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		// Ensure reference handles the test (sanity)
		if _, err := runProgram(refBin, tc.input); err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := validateOutput(tc, gotOut); err != nil {
			fmt.Fprintf(os.Stderr, "candidate invalid on test %d: %v\nInput:\n%s\nOutput:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"1723GP1.go",
		filepath.Join("1000-1999", "1700-1799", "1720-1729", "1723", "1723GP1.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 1723GP1.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1723GP1_%d.bin", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", outPath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func validateOutput(tc testCase, out string) error {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("empty output")
	}
	groupsCount, err := strconv.Atoi(fields[0])
	if err != nil || groupsCount <= 0 {
		return fmt.Errorf("invalid groups count %q", fields[0])
	}
	fields = fields[1:]
	n := tc.n
	limit := int(math.Sqrt(float64(n)))
	if limit == 0 {
		limit = 1
	}
	seen := make([]bool, n)
	idx := 0
	for g := 0; g < groupsCount; g++ {
		if idx >= len(fields) {
			return fmt.Errorf("not enough data for group %d", g+1)
		}
		size, err := strconv.Atoi(fields[idx])
		if err != nil || size <= 0 {
			return fmt.Errorf("invalid size for group %d", g+1)
		}
		idx++
		if size > limit {
			return fmt.Errorf("group %d has size %d exceeding limit %d", g+1, size, limit)
		}
		if idx+size-1 >= len(fields) {
			return fmt.Errorf("not enough vertices for group %d", g+1)
		}
		for j := 0; j < size; j++ {
			vertex, err := strconv.Atoi(fields[idx+j])
			if err != nil || vertex < 0 || vertex >= n {
				return fmt.Errorf("invalid vertex %q in group %d", fields[idx+j], g+1)
			}
			if seen[vertex] {
				return fmt.Errorf("vertex %d appears multiple times", vertex)
			}
			seen[vertex] = true
		}
		idx += size
	}
	for i, v := range seen {
		if !v {
			return fmt.Errorf("vertex %d not included in any group", i)
		}
	}
	return nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		buildTest(1, [][3]int{}),
		buildTest(2, [][3]int{{0, 1, 5}}),
		buildTest(3, [][3]int{{0, 1, 1}, {1, 2, 2}}),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 40 {
		n := rng.Intn(50) + 1
		tests = append(tests, randomGraph(rng, n))
	}
	tests = append(tests, randomGraph(rng, 2000))
	return tests
}

func buildTest(n int, edges [][3]int) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		b.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
	}
	return testCase{input: b.String(), n: n}
}

func randomGraph(rng *rand.Rand, n int) testCase {
	total := rng.Intn(n*(n-1)/2-n+1) + (n - 1)
	edges := make([][3]int, 0, total)
	used := make(map[[2]int]struct{})
	for v := 1; v < n; v++ {
		u := rng.Intn(v)
		addEdge(u, v, rng, used, &edges)
	}
	for len(edges) < total {
		u := rng.Intn(n)
		v := rng.Intn(n)
		if u == v {
			continue
		}
		addEdge(u, v, rng, used, &edges)
	}
	return buildTest(n, edges)
}

func addEdge(u, v int, rng *rand.Rand, used map[[2]int]struct{}, edges *[][3]int) {
	a, b := u, v
	if a > b {
		a, b = b, a
	}
	key := [2]int{a, b}
	if _, ok := used[key]; ok {
		return
	}
	used[key] = struct{}{}
	w := rng.Intn(100)
	*edges = append(*edges, [3]int{u, v, w})
}
