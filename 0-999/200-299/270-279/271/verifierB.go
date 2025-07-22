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

func expected(matrix [][]int) int {
	n := len(matrix)
	m := len(matrix[0])
	maxA := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if matrix[i][j] > maxA {
				maxA = matrix[i][j]
			}
		}
	}
	buf := 500
	limit := maxA + buf
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}
	nextPrime := make([]int, limit+2)
	next := -1
	for i := limit; i >= 0; i-- {
		if isPrime[i] {
			next = i
		}
		nextPrime[i] = next
	}
	rowSum := make([]int, n)
	colSum := make([]int, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			v := matrix[i][j]
			np := nextPrime[v]
			if np < 0 {
				np = v
			}
			delta := np - v
			rowSum[i] += delta
			colSum[j] += delta
		}
	}
	ans := rowSum[0]
	for i := 0; i < n; i++ {
		if rowSum[i] < ans {
			ans = rowSum[i]
		}
	}
	for j := 0; j < m; j++ {
		if colSum[j] < ans {
			ans = colSum[j]
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		panic(err)
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
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Printf("test %d invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		if len(parts) != 2+n*m {
			fmt.Printf("test %d wrong number of values\n", idx)
			os.Exit(1)
		}
		matrix := make([][]int, n)
		pos := 2
		for i := 0; i < n; i++ {
			matrix[i] = make([]int, m)
			for j := 0; j < m; j++ {
				v, _ := strconv.Atoi(parts[pos])
				matrix[i][j] = v
				pos++
			}
		}
		exp := expected(matrix)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", matrix[i][j]))
			}
			sb.WriteByte('\n')
		}
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(sb.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != fmt.Sprintf("%d", exp) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
