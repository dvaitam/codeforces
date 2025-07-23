package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCaseB struct {
	n    int
	perm []int
	mat  []string
}

func parseTestcases(path string) ([]testCaseB, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return nil, err
	}
	cases := make([]testCaseB, T)
	for i := 0; i < T; i++ {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return nil, err
		}
		perm := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &perm[j])
		}
		mat := make([]string, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &mat[j])
		}
		cases[i] = testCaseB{n: n, perm: perm, mat: mat}
	}
	return cases, nil
}

func solveCase(tc testCaseB) string {
	n := tc.n
	p := append([]int(nil), tc.perm...)
	adj := tc.mat
	visited := make([]bool, n)
	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}
		queue := []int{i}
		visited[i] = true
		comp := []int{i}
		for q := 0; q < len(queue); q++ {
			u := queue[q]
			for v := 0; v < n; v++ {
				if adj[u][v] == '1' && !visited[v] {
					visited[v] = true
					queue = append(queue, v)
					comp = append(comp, v)
				}
			}
		}
		vals := make([]int, len(comp))
		for j, idx := range comp {
			vals[j] = p[idx]
		}
		sort.Ints(comp)
		sort.Ints(vals)
		for j, idx := range comp {
			p[idx] = vals[j]
		}
	}
	var sb strings.Builder
	for i, v := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for j, v := range tc.perm {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, row := range tc.mat {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
		expected := solveCase(tc)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
