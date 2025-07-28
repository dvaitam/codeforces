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

func runCandidate(bin string, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func hasValidSubset(n int, adj [][]bool) bool {
	comb := make([]int, 5)
	var dfs func(start, depth int) bool
	dfs = func(start, depth int) bool {
		if depth == 5 {
			clique := true
			independent := true
			for i := 0; i < 5; i++ {
				for j := i + 1; j < 5; j++ {
					if adj[comb[i]][comb[j]] {
						independent = false
					} else {
						clique = false
					}
				}
			}
			return clique || independent
		}
		for i := start; i <= n-(5-depth)+1; i++ {
			comb[depth] = i
			if dfs(i+1, depth+1) {
				return true
			}
		}
		return false
	}
	return dfs(1, 0)
}

func verifyOutput(out string, n int, adj [][]bool, feasible bool) error {
	out = strings.TrimSpace(out)
	if out == "-1" {
		if feasible {
			return fmt.Errorf("expected subset but got -1")
		}
		return nil
	}
	fields := strings.Fields(out)
	if len(fields) != 5 {
		return fmt.Errorf("expected 5 numbers got %d", len(fields))
	}
	nums := make([]int, 5)
	seen := make(map[int]bool)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid number: %v", err)
		}
		if v < 1 || v > n {
			return fmt.Errorf("number out of range")
		}
		if seen[v] {
			return fmt.Errorf("duplicate numbers")
		}
		seen[v] = true
		nums[i] = v
	}
	clique := true
	independent := true
	for i := 0; i < 5; i++ {
		for j := i + 1; j < 5; j++ {
			if adj[nums[i]][nums[j]] {
				independent = false
			} else {
				clique = false
			}
		}
	}
	if !(clique || independent) {
		return fmt.Errorf("numbers do not form clique or independent set")
	}
	if !feasible {
		return fmt.Errorf("output subset but none exist")
	}
	return nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]
	file, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		n, err1 := strconv.Atoi(fields[0])
		m, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil || len(fields) != 2+2*m {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		edges := make([][2]int, m)
		pos := 2
		for i := 0; i < m; i++ {
			u, _ := strconv.Atoi(fields[pos])
			pos++
			v, _ := strconv.Atoi(fields[pos])
			pos++
			edges[i] = [2]int{u, v}
		}
		adj := make([][]bool, n+1)
		for i := 0; i <= n; i++ {
			adj[i] = make([]bool, n+1)
		}
		for _, e := range edges {
			adj[e[0]][e[1]] = true
			adj[e[1]][e[0]] = true
		}
		feasible := hasValidSubset(n, adj)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		out, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if err := verifyOutput(out, n, adj, feasible); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
