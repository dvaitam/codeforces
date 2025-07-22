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

func expected(lines []string) string {
	n, _ := strconv.Atoi(strings.TrimSpace(lines[0]))
	xs := make([]int, n)
	ys := make([]int, n)
	for i := 0; i < n; i++ {
		parts := strings.Fields(lines[1+i])
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		xs[i] = x
		ys[i] = y
	}
	visited := make([]bool, n)
	var dfs func(int)
	dfs = func(u int) {
		visited[u] = true
		for v := 0; v < n; v++ {
			if !visited[v] && (xs[u] == xs[v] || ys[u] == ys[v]) {
				dfs(v)
			}
		}
	}
	comp := 0
	for i := 0; i < n; i++ {
		if !visited[i] {
			comp++
			dfs(i)
		}
	}
	return fmt.Sprint(comp - 1)
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
	var lines []string
	idx := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			if len(lines) == 0 {
				continue
			}
			idx++
			expect := expected(lines)
			input := strings.Join(lines, "\n") + "\n"
			cmd := exec.Command(bin)
			cmd.Stdin = strings.NewReader(input)
			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			err := cmd.Run()
			if err != nil {
				fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
				os.Exit(1)
			}
			got := strings.TrimSpace(out.String())
			if got != expect {
				fmt.Printf("test %d failed: expected %s got %s\n", idx, expect, got)
				os.Exit(1)
			}
			lines = lines[:0]
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) > 0 {
		idx++
		expect := expected(lines)
		input := strings.Join(lines, "\n") + "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
