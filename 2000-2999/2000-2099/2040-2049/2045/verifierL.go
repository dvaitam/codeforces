package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const refSourceL = "2000-2999/2000-2099/2040-2049/2045/2045L.go"

type edge struct {
	u, v int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierL.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}
	K, err := parseK(input)
	if err != nil {
		fail("failed to parse input: %v", err)
	}

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runCommand(exec.Command(refBin), input)
	if err != nil {
		fail("reference execution failed: %v", err)
	}
	expImpossible, err := isImpossible(refOut)
	if err != nil {
		fail("failed to parse reference output: %v\n%s", err, refOut)
	}

	userOut, err := runCommand(commandFor(candidate), input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}

	if expImpossible {
		if !isExactImpossible(userOut) {
			fail("expected impossible (-1 -1) but candidate produced:\n%s", userOut)
		}
		fmt.Println("OK")
		return
	}

	n, m, edges, err := parseGraph(userOut)
	if err != nil {
		fail("failed to parse candidate output: %v", err)
	}
	if n < 1 || n > 32768 {
		fail("N out of bounds: %d", n)
	}
	if m < 0 || m > 65536 {
		fail("M out of bounds: %d", m)
	}
	if len(edges) != m {
		fail("expected %d edges, parsed %d", m, len(edges))
	}

	seen := make(map[edge]struct{}, m)
	for _, e := range edges {
		if e.u < 1 || e.u > n || e.v < 1 || e.v > n {
			fail("edge endpoint out of bounds: %d %d (n=%d)", e.u, e.v, n)
		}
		if e.u == e.v {
			fail("self-loop detected on node %d", e.u)
		}
		a, b := e.u, e.v
		if a > b {
			a, b = b, a
		}
		key := edge{a, b}
		if _, ok := seen[key]; ok {
			fail("duplicate edge detected: %d %d", a, b)
		}
		seen[key] = struct{}{}
	}

	counter := runBDFS(n, edges)
	if counter != K {
		fail("BDFS counter mismatch: expected %d got %d", K, counter)
	}

	fmt.Println("OK")
}

func parseK(input []byte) (int64, error) {
	fields := strings.Fields(string(input))
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected 1 integer, got %d tokens", len(fields))
	}
	return strconv.ParseInt(fields[0], 10, 64)
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2045L-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceL))
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
	}
	return tmp.Name(), nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runCommand(cmd *exec.Cmd, input []byte) (string, error) {
	cmd.Stdin = bytes.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := stderr.String()
		if msg == "" {
			msg = stdout.String()
		}
		return "", fmt.Errorf("%v\n%s", err, msg)
	}
	return stdout.String(), nil
}

func isImpossible(out string) (bool, error) {
	fields := strings.Fields(out)
	if len(fields) == 2 && fields[0] == "-1" && fields[1] == "-1" {
		return true, nil
	}
	if len(fields) < 2 {
		return false, fmt.Errorf("not enough tokens")
	}
	_, err1 := strconv.Atoi(fields[0])
	_, err2 := strconv.Atoi(fields[1])
	if err1 != nil || err2 != nil {
		return false, fmt.Errorf("invalid start tokens")
	}
	return false, nil
}

func isExactImpossible(out string) bool {
	fields := strings.Fields(out)
	return len(fields) == 2 && fields[0] == "-1" && fields[1] == "-1"
}

func parseGraph(out string) (int, int, []edge, error) {
	fields := strings.Fields(out)
	if len(fields) < 2 {
		return 0, 0, nil, fmt.Errorf("not enough output tokens")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, nil, fmt.Errorf("invalid N: %v", err)
	}
	m, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, nil, fmt.Errorf("invalid M: %v", err)
	}
	if len(fields) != 2+2*m {
		return 0, 0, nil, fmt.Errorf("expected %d tokens, got %d", 2+2*m, len(fields))
	}
	edges := make([]edge, m)
	idx := 2
	for i := 0; i < m; i++ {
		u, err := strconv.Atoi(fields[idx])
		if err != nil {
			return 0, 0, nil, fmt.Errorf("invalid u at edge %d: %v", i+1, err)
		}
		v, err := strconv.Atoi(fields[idx+1])
		if err != nil {
			return 0, 0, nil, fmt.Errorf("invalid v at edge %d: %v", i+1, err)
		}
		edges[i] = edge{u, v}
		idx += 2
	}
	return n, m, edges, nil
}

func runBDFS(n int, edges []edge) int64 {
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	for i := 1; i <= n; i++ {
		sort.Ints(adj[i])
	}

	flag := make([]bool, n+1)
	stack := []int{1}
	var counter int64

	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		flag[u] = true
		for _, v := range adj[u] {
			counter++
			if !flag[v] {
				stack = append(stack, v)
			}
		}
	}
	return counter
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
