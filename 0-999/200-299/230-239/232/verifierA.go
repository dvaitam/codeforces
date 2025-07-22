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

func countTriangles(adj [][]bool) int {
	n := len(adj)
	cnt := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if !adj[i][j] {
				continue
			}
			for k := j + 1; k < n; k++ {
				if adj[i][k] && adj[j][k] {
					cnt++
				}
			}
		}
	}
	return cnt
}

func runCase(bin string, k int) error {
	input := fmt.Sprintf("%d\n", k)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(out.String())))
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	n64, err := strconv.ParseInt(strings.TrimSpace(scanner.Text()), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse n: %v", err)
	}
	n := int(n64)
	if n < 3 || n > 100 {
		return fmt.Errorf("invalid n %d", n)
	}
	adj := make([][]bool, n)
	for i := 0; i < n; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("not enough lines for adjacency matrix")
		}
		line := strings.TrimSpace(scanner.Text())
		if len(line) != n {
			return fmt.Errorf("line %d length mismatch", i+1)
		}
		adj[i] = make([]bool, n)
		for j := 0; j < n; j++ {
			c := line[j]
			if c != '0' && c != '1' {
				return fmt.Errorf("invalid char at row %d col %d", i+1, j+1)
			}
			if i == j && c != '0' {
				return fmt.Errorf("self loop at row %d", i+1)
			}
			adj[i][j] = c == '1'
		}
	}
	if scanner.Scan() {
		extra := strings.TrimSpace(scanner.Text())
		if extra != "" {
			return fmt.Errorf("extra output: %q", extra)
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if adj[i][j] != adj[j][i] {
				return fmt.Errorf("asymmetry at %d,%d", i, j)
			}
		}
	}
	if countTriangles(adj) != k {
		return fmt.Errorf("expected %d triangles", k)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
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
		k, err := strconv.Atoi(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid testcase line %d\n", idx)
			os.Exit(1)
		}
		if err := runCase(bin, k); err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
