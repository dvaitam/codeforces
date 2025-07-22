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

func expected(n, m int, grid []string) string {
	row := make([][26]int, n)
	col := make([][26]int, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			c := grid[i][j] - 'a'
			row[i][c]++
			col[j][c]++
		}
	}
	var ans []byte
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			c := grid[i][j] - 'a'
			if row[i][c] == 1 && col[j][c] == 1 {
				ans = append(ans, grid[i][j])
			}
		}
	}
	return string(ans)
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
		if len(parts) < 2 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		if len(parts) != 2+n {
			fmt.Printf("test %d: wrong number of rows\n", idx)
			os.Exit(1)
		}
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			if len(parts[2+i]) != m {
				fmt.Printf("test %d: row %d has wrong length\n", idx, i+1)
				os.Exit(1)
			}
			grid[i] = parts[2+i]
		}
		expect := expected(n, m, grid)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			input.WriteString(grid[i])
			input.WriteByte('\n')
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
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
			fmt.Printf("test %d failed:\nexpected: %q\n   got: %q\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
