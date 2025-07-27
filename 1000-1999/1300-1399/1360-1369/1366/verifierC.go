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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveCase(n, m int, grid [][]int) int {
	zeros := make([]int, n+m-1)
	ones := make([]int, n+m-1)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 0 {
				zeros[i+j]++
			} else {
				ones[i+j]++
			}
		}
	}
	total := n + m - 2
	ans := 0
	for l, r := 0, total; l < r; l, r = l+1, r-1 {
		ans += min(zeros[l]+zeros[r], ones[l]+ones[r])
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		panic(err)
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
			fmt.Printf("invalid test %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1])
		if len(fields) != 2+n*m {
			fmt.Printf("test %d wrong count\n", idx)
			os.Exit(1)
		}
		grid := make([][]int, n)
		pos := 2
		for i := 0; i < n; i++ {
			grid[i] = make([]int, m)
			for j := 0; j < m; j++ {
				v, _ := strconv.Atoi(fields[pos])
				pos++
				grid[i][j] = v
			}
		}
		expected := solveCase(n, m, grid)

		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if j > 0 {
					input.WriteByte(' ')
				}
				input.WriteString(fmt.Sprintf("%d", grid[i][j]))
			}
			input.WriteByte('\n')
		}

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\n%s", idx, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(out.String())
		if outStr != fmt.Sprintf("%d", expected) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, expected, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
