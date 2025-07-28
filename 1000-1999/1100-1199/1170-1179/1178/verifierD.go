package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const numTestsD = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verifD_bin")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return strings.TrimSpace(buf.String()), err
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	limit := int(math.Sqrt(float64(n))) + 1
	for i := 2; i < limit; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func nextPrime(n int) int {
	for !isPrime(n) {
		n++
	}
	return n
}

type edge struct{ u, v int }

func parseOutput(out string) (int, []edge, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		return 0, nil, fmt.Errorf("empty output")
	}
	m, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return 0, nil, fmt.Errorf("invalid first line: %v", err)
	}
	edges := make([]edge, 0, len(lines)-1)
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return 0, nil, fmt.Errorf("invalid edge line: %q", line)
		}
		u, err1 := strconv.Atoi(parts[0])
		v, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return 0, nil, fmt.Errorf("invalid numbers in line: %q", line)
		}
		if u > v {
			u, v = v, u
		}
		edges = append(edges, edge{u, v})
	}
	if len(edges) != m {
		return 0, nil, fmt.Errorf("edge count mismatch: %d vs %d", m, len(edges))
	}
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].u == edges[j].u {
			return edges[i].v < edges[j].v
		}
		return edges[i].u < edges[j].u
	})
	return m, edges, nil
}

func solveD(n int) string {
	m := nextPrime(n)
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(m))
	sb.WriteByte('\n')
	for i := 1; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", i, i+1))
	}
	sb.WriteString(fmt.Sprintf("1 %d\n", n))
	extra := m - n
	half := n / 2
	for i := 1; i <= extra; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", i, i+half))
	}
	return strings.TrimSpace(sb.String())
}

func genCaseD(rng *rand.Rand) int {
	return rng.Intn(50) + 3
}

func runCaseD(bin string, n int) error {
	input := fmt.Sprintf("%d\n", n)
	expectedStr := solveD(n)
	outStr, err := run(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	_, expEdges, err := parseOutput(expectedStr)
	if err != nil {
		return fmt.Errorf("bad expected output: %v", err)
	}
	_, gotEdges, err := parseOutput(outStr)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if len(expEdges) != len(gotEdges) {
		return fmt.Errorf("edge count mismatch: expected %d got %d", len(expEdges), len(gotEdges))
	}
	for i := range expEdges {
		if expEdges[i] != gotEdges[i] {
			return fmt.Errorf("expected:\n%s\n got:\n%s", expectedStr, outStr)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin, cleanup, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if cleanup != nil {
		defer cleanup()
	}
	rng := rand.New(rand.NewSource(4))
	for t := 0; t < numTestsD; t++ {
		n := genCaseD(rng)
		if err := runCaseD(bin, n); err != nil {
			fmt.Printf("case %d failed: %v\n", t+1, err)
			return
		}
	}
	fmt.Println("All tests passed")
}
