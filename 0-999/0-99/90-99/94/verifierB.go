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

func runCandidate(bin, input string) (string, error) {
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

func expected(m int, edges [][2]int) string {
	var g [6][6]bool
	for _, e := range edges {
		g[e[0]][e[1]] = true
		g[e[1]][e[0]] = true
	}
	for i := 1; i <= 5; i++ {
		for j := i + 1; j <= 5; j++ {
			for k := j + 1; k <= 5; k++ {
				cnt := 0
				if g[i][j] {
					cnt++
				}
				if g[i][k] {
					cnt++
				}
				if g[j][k] {
					cnt++
				}
				if cnt == 3 || cnt == 0 {
					return "WIN"
				}
			}
		}
	}
	return "FAIL"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
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
		parts := strings.Fields(line)
		if len(parts) < 1 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		m, err := strconv.Atoi(parts[0])
		if err != nil || len(parts) != 1+2*m {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			a, _ := strconv.Atoi(parts[1+2*i])
			b, _ := strconv.Atoi(parts[2+2*i])
			edges[i] = [2]int{a, b}
		}
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", m))
		for _, e := range edges {
			input.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		want := expected(m, edges)
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
