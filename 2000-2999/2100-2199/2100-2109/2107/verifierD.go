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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
		want := strings.TrimSpace(wantOut)

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		got := strings.TrimSpace(gotOut)

		if want != got {
			fmt.Fprintf(os.Stderr, "test %d mismatch:\nInput:\n%s\nExpected:\n%s\nGot:\n%s\n", idx+1, tc.input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"2107D.go",
		filepath.Join("2000-2999", "2100-2199", "2100-2109", "2107", "2107D.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 2107D.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref2107D_%d.bin", time.Now().UnixNano()))
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

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		buildTest([][]int{{1, 2}}),
		buildTest([][]int{{1, 2}, {2, 3}}),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 40 {
		n := rng.Intn(8) + 1
		tests = append(tests, randomTreeTest(rng, n))
	}
	return tests
}

func buildTest(edges [][]int) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("1\n%d\n", len(edges)+1))
	for _, e := range edges {
		b.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return testCase{input: b.String()}
}

func randomTreeTest(rng *rand.Rand, n int) testCase {
	if n == 1 {
		return buildTest(nil)
	}
	if n == 2 {
		return buildTest([][]int{{1, 2}})
	}
	prufer := make([]int, n-2)
	for i := 0; i < n-2; i++ {
		prufer[i] = rng.Intn(n) + 1
	}
	deg := make([]int, n+1)
	for i := 1; i <= n; i++ {
		deg[i] = 1
	}
	for _, v := range prufer {
		deg[v]++
	}
	edges := make([][]int, 0, n-1)
	for _, v := range prufer {
		u := 1
		for deg[u] != 1 {
			u++
		}
		edges = append(edges, []int{u, v})
		deg[u]--
		deg[v]--
	}
	var u, v int
	for i := 1; i <= n; i++ {
		if deg[i] == 1 {
			if u == 0 {
				u = i
			} else {
				v = i
				break
			}
		}
	}
	edges = append(edges, []int{u, v})
	return buildTest(edges)
}
