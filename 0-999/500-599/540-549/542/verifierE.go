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

type testCase struct {
	n     int
	edges [][2]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, tc := range tests {
		input := serialize(tc)
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(refOut) != strings.TrimSpace(candOut) {
			fmt.Fprintf(os.Stderr, "Mismatch on test %d\nInput:\n%sExpected:\n%sGot:\n%s", idx+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_542E.bin"
	cmd := exec.Command("go", "build", "-o", refName, "542E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
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
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func buildTests() []testCase {
	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 200 {
		n := rng.Intn(25) + 1
		maxEdges := n * (n - 1) / 2
		m := rng.Intn(maxEdges + 1)
		edges := randomEdges(rng, n, m)
		tests = append(tests, testCase{n: n, edges: edges})
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, edges: nil},
		{n: 2, edges: [][2]int{{1, 2}}},
		{n: 4, edges: [][2]int{{1, 2}, {2, 3}, {3, 4}}},
		{n: 5, edges: [][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}}},
	}
}

func randomEdges(rng *rand.Rand, n, m int) [][2]int {
	type pair struct{ u, v int }
	edges := make([][2]int, 0, m)
	used := make(map[pair]struct{}, m)
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n-1) + 1
		if v >= u {
			v++
		}
		if u > v {
			u, v = v, u
		}
		key := pair{u, v}
		if _, ok := used[key]; ok {
			continue
		}
		used[key] = struct{}{}
		edges = append(edges, [2]int{u, v})
	}
	return edges
}

func serialize(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.edges)))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}
