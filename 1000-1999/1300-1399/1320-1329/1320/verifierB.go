package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const refSource = "1000-1999/1300-1399/1320-1329/1320/1320B.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/candidate")
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
	for idx, input := range tests {
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fail("reference failed on test %d: %v", idx+1, err)
		}
		candOut, err := runProgram(candidate, input)
		if err != nil {
			fail("candidate failed on test %d: %v", idx+1, err)
		}
		if normalize(refOut) != normalize(candOut) {
			fail("mismatch on test %d\nInput:\n%sExpected: %sGot: %s", idx+1, input, refOut, candOut)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1320B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, buf.String())
	}
	return tmp.Name(), nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	switch filepath.Ext(target) {
	case ".go":
		cmd = exec.Command("go", "run", target)
	case ".py":
		cmd = exec.Command("python3", target)
	default:
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return stdout.String(), fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
		}
		return stdout.String(), err
	}
	return stdout.String(), nil
}

func normalize(s string) string {
	return strings.TrimSpace(s)
}

func buildTests() []string {
	samples := []string{
		"6 9\n1 5\n5 4\n1 2\n2 3\n3 4\n4 1\n2 6\n6 4\n4 2\n4\n1 2 3 4\n",
		"7 7\n1 2\n2 3\n3 4\n4 5\n5 6\n6 7\n7 1\n7\n1 2 3 4 5 6 7\n",
		"8 13\n8 7\n8 6\n7 5\n7 4\n6 5\n6 4\n5 3\n5 2\n4 3\n4 2\n3 1\n2 1\n1 8\n5\n8 7 5 2 1\n",
	}

	tests := make([]string, 0, 20)
	tests = append(tests, samples...)

	rng := rand.New(rand.NewSource(1))
	for len(tests) < 15 {
		tests = append(tests, generateRandomTest(rng))
	}

	return tests
}

func generateRandomTest(rng *rand.Rand) string {
	n := rng.Intn(5) + 4 // 4..8
	edgeSet := make(map[[2]int]struct{})
	addEdge := func(u, v int) {
		if u == v {
			return
		}
		edgeSet[[2]int{u, v}] = struct{}{}
	}

	for i := 0; i < n; i++ {
		addEdge(i, (i+1)%n)
	}

	extra := rng.Intn(n + 3)
	for i := 0; i < extra; i++ {
		u := rng.Intn(n)
		v := rng.Intn(n)
		addEdge(u, v)
	}

	edges := make([][2]int, 0, len(edgeSet))
	for e := range edgeSet {
		edges = append(edges, e)
	}
	sort.Slice(edges, func(i, j int) bool {
		if edges[i][0] == edges[j][0] {
			return edges[i][1] < edges[j][1]
		}
		return edges[i][0] < edges[j][0]
	})

	k := rng.Intn(n-1) + 2 // 2..n
	start := rng.Intn(n)
	path := make([]int, k)
	for i := 0; i < k; i++ {
		path[i] = (start + i) % n
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&b, "%d %d\n", e[0]+1, e[1]+1)
	}
	fmt.Fprintf(&b, "%d\n", k)
	for i, node := range path {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", node+1)
	}
	b.WriteByte('\n')
	return b.String()
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
