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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	file, err := os.Open("testcasesB.txt")
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
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Printf("invalid test %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		if len(parts) != 2+n*m {
			fmt.Printf("test %d: expected %d numbers got %d\n", idx, n*m, len(parts)-2)
			os.Exit(1)
		}
		grid := make([][]int, n)
		for i := 0; i < n; i++ {
			grid[i] = make([]int, m)
			for j := 0; j < m; j++ {
				v, _ := strconv.Atoi(parts[2+i*m+j])
				grid[i][j] = v
			}
		}
		neigh := make([][]int, n)
		possible := true
		for i := 0; i < n; i++ {
			neigh[i] = make([]int, m)
			for j := 0; j < m; j++ {
				cnt := 0
				if i > 0 {
					cnt++
				}
				if i < n-1 {
					cnt++
				}
				if j > 0 {
					cnt++
				}
				if j < m-1 {
					cnt++
				}
				neigh[i][j] = cnt
				if grid[i][j] > cnt {
					possible = false
				}
			}
		}
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

		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input.String())
		var outBuf, errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		outFields := strings.Fields(strings.TrimSpace(outBuf.String()))
		if len(outFields) == 0 {
			fmt.Printf("Test %d: empty output\n", idx)
			os.Exit(1)
		}
		if !possible {
			if strings.ToUpper(outFields[0]) != "NO" {
				fmt.Printf("Test %d: expected NO\n", idx)
				os.Exit(1)
			}
			continue
		}
		if strings.ToUpper(outFields[0]) != "YES" {
			fmt.Printf("Test %d: expected YES\n", idx)
			os.Exit(1)
		}
		if len(outFields)-1 != n*m {
			fmt.Printf("Test %d: expected %d numbers got %d\n", idx, n*m, len(outFields)-1)
			os.Exit(1)
		}
		k := 1
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				v, err := strconv.Atoi(outFields[k])
				if err != nil || v != neigh[i][j] {
					fmt.Printf("Test %d: invalid value at cell (%d,%d)\n", idx, i, j)
					os.Exit(1)
				}
				k++
			}
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
