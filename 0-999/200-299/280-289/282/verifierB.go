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

func feasible(a, g []int) bool {
	reachable := map[int]struct{}{0: {}}
	for i := range a {
		next := make(map[int]struct{}, len(reachable)*2)
		for d := range reachable {
			next[d+a[i]] = struct{}{}
			next[d-g[i]] = struct{}{}
		}
		reachable = next
	}
	for d := range reachable {
		if d >= -500 && d <= 500 {
			return true
		}
	}
	return false
}

func runCase(bin string, line string, idx int) error {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("test %d: invalid n", idx)
	}
	if len(fields) != 1+2*n {
		return fmt.Errorf("test %d: expected %d numbers got %d", idx, 1+2*n, len(fields))
	}
	A := make([]int, n)
	G := make([]int, n)
	for i := 0; i < n; i++ {
		ai, _ := strconv.Atoi(fields[1+2*i])
		gi, _ := strconv.Atoi(fields[2+2*i])
		if ai+gi != 1000 {
			return fmt.Errorf("test %d: ai+gi != 1000", idx)
		}
		A[i] = ai
		G[i] = gi
	}
	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&input, "%d %d\n", A[i], G[i])
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	ans := strings.TrimSpace(out.String())
	exists := feasible(A, G)
	if ans == "-1" {
		if exists {
			return fmt.Errorf("solution exists but got -1")
		}
		return nil
	}
	if len(ans) != n {
		return fmt.Errorf("expected string of length %d", n)
	}
	sa, sg := 0, 0
	for i, ch := range ans {
		if ch == 'A' {
			sa += A[i]
		} else if ch == 'G' {
			sg += G[i]
		} else {
			return fmt.Errorf("invalid character %c", ch)
		}
	}
	if abs(sa-sg) > 500 {
		return fmt.Errorf("difference %d exceeds 500", abs(sa-sg))
	}
	if !exists {
		return fmt.Errorf("no solution should exist, but output provided")
	}
	return nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open testcasesB.txt: %v\n", err)
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
		if err := runCase(bin, line, idx); err != nil {
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
