package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solveCase(line string) string {
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return ""
	}
	idx := 0
	n, _ := strconv.Atoi(fields[idx])
	idx++
	g := make([][]bool, n)
	for i := 0; i < n; i++ {
		s := fields[idx]
		idx++
		g[i] = make([]bool, n)
		for j := 0; j < n && j < len(s); j++ {
			if s[j] == '1' {
				g[i][j] = true
			}
		}
	}
	reach := make([][]bool, n)
	for i := 0; i < n; i++ {
		reach[i] = make([]bool, n)
	}
	for i := 0; i < n; i++ {
		queue := []int{i}
		reach[i][i] = true
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			for u := 0; u < n; u++ {
				if g[v][u] && !reach[i][u] {
					reach[i][u] = true
					queue = append(queue, u)
				}
			}
		}
	}
	var sb strings.Builder
	sb.WriteString("3\n")
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if reach[i][j] {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('0')
			}
		}
		if i+1 < n {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idxCase := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idxCase++
		expected := solveCase(line)
		input := "1\n" + line + "\n"
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idxCase, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idxCase, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idxCase)
}
