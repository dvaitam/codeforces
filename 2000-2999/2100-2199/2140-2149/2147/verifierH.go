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

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
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
		wantOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		if !compareOutputs(tc.input, wantOut, gotOut) {
			fmt.Fprintf(os.Stderr, "test %d mismatch.\nInput:\n%s\nExpected output:\n%s\nCandidate output:\n%s\n", idx+1, tc.input, wantOut, gotOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"2147H.go",
		filepath.Join("2000-2999", "2100-2199", "2140-2149", "2147", "2147H.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 2147H.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref2147H_%d.bin", time.Now().UnixNano()))
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

func compareOutputs(input, wantOut, gotOut string) bool {
	want := strings.Fields(wantOut)
	got := strings.Fields(gotOut)
	if len(want) != len(got) {
		return false
	}
	for i := range want {
		if want[i] != got[i] {
			return false
		}
	}
	return true
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		buildTest(2, [][3]int{{1, 2, 5}}),
		buildTest(3, [][3]int{{1, 2, 4}, {2, 3, 6}}),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 30 {
		n := rng.Intn(6) + 1
		m := rng.Intn(n*(n-1)/2 + 1)
		tests = append(tests, randomGraph(rng, n, m))
	}
	return tests
}

func buildTest(n int, edges [][3]int) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("1\n%d %d\n", n, len(edges)))
	for _, e := range edges {
		b.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
	}
	return testCase{input: b.String()}
}

func randomGraph(rng *rand.Rand, n, m int) testCase {
	type edge struct{ u, v int }
	used := make(map[edge]struct{})
	edges := make([][3]int, 0, m)
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		e := edge{u, v}
		if u > v {
			e.u, e.v = v, u
		}
		if _, ok := used[e]; ok {
			continue
		}
		used[e] = struct{}{}
		w := rng.Intn(100) + 1
		edges = append(edges, [3]int{u, v, w})
	}
	return buildTest(n, edges)
}
