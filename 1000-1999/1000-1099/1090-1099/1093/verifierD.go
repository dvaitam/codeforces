package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runExe(path, input string) (string, error) {
	if !strings.Contains(path, "/") {
		path = "./" + path
	}
	cmd := exec.Command(path)
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildRef() (string, error) {
	ref := "refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1093D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

type edge struct{ u, v int }

type caseD struct {
	n, m  int
	edges []edge
}

func genCase(r *rand.Rand) caseD {
	n := r.Intn(6) + 1 // 1..6
	maxM := n * (n - 1) / 2
	m := r.Intn(min(maxM, 6) + 1) // up to 6 edges
	used := make(map[[2]int]bool)
	edges := make([]edge, 0, m)
	for len(edges) < m {
		u := r.Intn(n) + 1
		v := r.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		key := [2]int{u, v}
		if used[key] {
			continue
		}
		used[key] = true
		edges = append(edges, edge{u, v})
	}
	return caseD{n: n, m: m, edges: edges}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("1\n%d %d\n", tc.n, tc.m))
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
		}
		input := sb.String()
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\n got:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
