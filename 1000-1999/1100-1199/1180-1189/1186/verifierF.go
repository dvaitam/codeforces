package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Edge struct{ u, v int }

func norm(u, v int) Edge {
	if u > v {
		u, v = v, u
	}
	return Edge{u, v}
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
	}
	return strings.TrimSpace(out.String()), nil
}

func validate(n, m int, edges []Edge, out string) error {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("empty output")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid k")
	}
	if k < 0 || k > (n+m+1)/2 {
		return fmt.Errorf("k out of range")
	}
	if len(fields) != 1+2*k {
		return fmt.Errorf("expected %d edges got %d", k, (len(fields)-1)/2)
	}
	orig := make(map[Edge]bool)
	degOrig := make([]int, n+1)
	for _, e := range edges {
		orig[norm(e.u, e.v)] = true
		degOrig[e.u]++
		degOrig[e.v]++
	}
	used := make(map[Edge]bool)
	deg := make([]int, n+1)
	for i := 0; i < k; i++ {
		u, err1 := strconv.Atoi(fields[1+2*i])
		v, err2 := strconv.Atoi(fields[1+2*i+1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("bad edge value")
		}
		if u < 1 || u > n || v < 1 || v > n {
			return fmt.Errorf("edge out of range")
		}
		e := norm(u, v)
		if !orig[e] {
			return fmt.Errorf("edge %d %d not in graph", u, v)
		}
		if used[e] {
			return fmt.Errorf("duplicate edge")
		}
		used[e] = true
		deg[u]++
		deg[v]++
	}
	for i := 1; i <= n; i++ {
		need := (degOrig[i] + 1) / 2
		if deg[i] < need {
			return fmt.Errorf("vertex %d degree %d < required %d", i, deg[i], need)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcasesF.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("bad testcase %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1])
		if len(fields) != 2+2*m {
			fmt.Printf("bad testcase %d\n", idx)
			os.Exit(1)
		}
		edges := make([]Edge, m)
		for i := 0; i < m; i++ {
			u, _ := strconv.Atoi(fields[2+2*i])
			v, _ := strconv.Atoi(fields[2+2*i+1])
			edges[i] = Edge{u, v}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
		}
		input := sb.String()
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if err := validate(n, m, edges, out); err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
