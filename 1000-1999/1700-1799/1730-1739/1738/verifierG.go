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

func solveG(n, k int, grid [][]int) (bool, [][]int) {
	f := make([][]int, n+2)
	vst := make([][]bool, n+2)
	for i := 0; i < n+2; i++ {
		f[i] = make([]int, n+2)
		vst[i] = make([]bool, n+2)
	}
	mx := make([][]int, k)
	for i := 0; i < k; i++ {
		mx[i] = make([]int, n+2)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			if grid[i-1][j-1] == 0 {
				f[i][j] = 1
			} else {
				f[i][j] = 0
			}
		}
	}
	isNo := false
	for i := n; i >= 1 && !isNo; i-- {
		for j := n; j >= 1; j-- {
			if f[i+1][j+1] > 0 {
				f[i][j] += f[i+1][j+1]
			}
			if f[i+1][j] > f[i][j] {
				f[i][j] = f[i+1][j]
			}
			if f[i][j+1] > f[i][j] {
				f[i][j] = f[i][j+1]
			}
			if f[i][j] == k {
				isNo = true
				break
			}
			if mx[f[i][j]][j] == 0 {
				mx[f[i][j]][j] = i
			}
		}
	}
	if isNo {
		return false, nil
	}
	for level := k - 1; level >= 1; level-- {
		for j := n - 1; j >= 1; j-- {
			if mx[level][j+1] > mx[level][j] {
				mx[level][j] = mx[level][j+1]
			}
		}
		x, y := n, 1
		for y <= n && vst[x][y] {
			y++
		}
		for y <= n {
			vst[x][y] = true
			if (y == n || x != mx[level][y+1]) && x > 1 && !vst[x-1][y] {
				x--
			} else {
				y++
			}
		}
	}
	res := make([][]int, n)
	for i := range res {
		res[i] = make([]int, n)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			if vst[i][j] {
				res[i-1][j-1] = 1
			}
		}
	}
	return true, res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesG.txt")
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
		k, _ := strconv.Atoi(parts[1])
		if len(parts) != 2+n*n {
			fmt.Printf("test %d: wrong number of values\n", idx)
			os.Exit(1)
		}
		grid := make([][]int, n)
		idxp := 2
		for i := 0; i < n; i++ {
			grid[i] = make([]int, n)
			for j := 0; j < n; j++ {
				v, _ := strconv.Atoi(parts[idxp])
				grid[i][j] = v
				idxp++
			}
		}
		ok, expectGrid := solveG(n, k, grid)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				input.WriteByte(byte('0' + grid[i][j]))
			}
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
		gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
		if !ok {
			if strings.TrimSpace(gotLines[0]) != "NO" {
				fmt.Printf("test %d failed: expected NO got %s\n", idx, gotLines[0])
				os.Exit(1)
			}
			continue
		}
		if strings.TrimSpace(gotLines[0]) != "YES" || len(gotLines)-1 != n {
			fmt.Printf("test %d failed: malformed output\n", idx)
			os.Exit(1)
		}
		for i := 0; i < n; i++ {
			if len(gotLines[i+1]) != n {
				fmt.Printf("test %d line %d length mismatch\n", idx, i+1)
				os.Exit(1)
			}
			for j := 0; j < n; j++ {
				if int(gotLines[i+1][j]-'0') != expectGrid[i][j] {
					fmt.Printf("test %d grid mismatch\n", idx)
					os.Exit(1)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
