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

func countTriangles(n, m int, grid []string) int64 {
	g := make([][]byte, n+2)
	for i := range g {
		g[i] = make([]byte, m+2)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			g[i][j] = grid[i-1][j-1]
		}
	}
	black := make([][]int, n+2)
	for i := range black {
		black[i] = make([]int, m+2)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if g[i][j] == '1' {
				black[i][j] = 1
			}
		}
	}
	rowPS := make([][]int, n+2)
	colPS := make([][]int, n+2)
	for i := range rowPS {
		rowPS[i] = make([]int, m+2)
		colPS[i] = make([]int, m+2)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			rowPS[i][j] = rowPS[i][j-1] + black[i][j]
			colPS[i][j] = colPS[i-1][j] + black[i][j]
		}
	}
	diag1 := make([][]int, n+2)
	diag2 := make([][]int, n+2)
	for i := range diag1 {
		diag1[i] = make([]int, m+3)
		diag2[i] = make([]int, m+3)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			diag1[i][j] = black[i][j] + diag1[i-1][j-1]
		}
		for j := m; j >= 1; j-- {
			diag2[i][j] = black[i][j] + diag2[i-1][j+1]
		}
	}
	hRun := make([][]int, n+2)
	hRunL := make([][]int, n+2)
	vRun := make([][]int, n+2)
	vRunU := make([][]int, n+2)
	for i := range hRun {
		hRun[i] = make([]int, m+2)
		hRunL[i] = make([]int, m+2)
		vRun[i] = make([]int, m+2)
		vRunU[i] = make([]int, m+2)
	}
	for i := n; i >= 1; i-- {
		for j := m; j >= 1; j-- {
			if g[i][j] == '0' {
				hRun[i][j] = hRun[i][j+1] + 1
				vRun[i][j] = vRun[i+1][j] + 1
			}
		}
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if g[i][j] == '0' {
				hRunL[i][j] = hRunL[i][j-1] + 1
				vRunU[i][j] = vRunU[i-1][j] + 1
			}
		}
	}
	var ans int64
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if g[i][j] != '0' {
				continue
			}
			maxK := hRun[i][j]
			if vRun[i][j] < maxK {
				maxK = vRun[i][j]
			}
			for k := 1; k < maxK; k++ {
				if g[i][j+k] != '0' || g[i+k][j] != '0' {
					continue
				}
				if rowPS[i][j+k-1]-rowPS[i][j] != 0 {
					continue
				}
				if colPS[i+k-1][j]-colPS[i][j] != 0 {
					continue
				}
				if diag2[i+k-1][j+1]-diag2[i][j+k] != 0 {
					continue
				}
				ans++
			}
			maxK = hRun[i][j]
			if vRunU[i][j] < maxK {
				maxK = vRunU[i][j]
			}
			for k := 1; k < maxK; k++ {
				if g[i][j+k] != '0' || g[i-k][j] != '0' {
					continue
				}
				if rowPS[i][j+k-1]-rowPS[i][j] != 0 {
					continue
				}
				if colPS[i-1][j]-colPS[i-k][j] != 0 {
					continue
				}
				if diag1[i-1][j+k-1]-diag1[i-k][j] != 0 {
					continue
				}
				ans++
			}
			maxK = hRunL[i][j]
			if vRun[i][j] < maxK {
				maxK = vRun[i][j]
			}
			for k := 1; k < maxK; k++ {
				if g[i][j-k] != '0' || g[i+k][j] != '0' {
					continue
				}
				if rowPS[i][j-1]-rowPS[i][j-k] != 0 {
					continue
				}
				if colPS[i+k-1][j]-colPS[i][j] != 0 {
					continue
				}
				if diag1[i+k-1][j-1]-diag1[i][j-k] != 0 {
					continue
				}
				ans++
			}
			maxK = hRunL[i][j]
			if vRunU[i][j] < maxK {
				maxK = vRunU[i][j]
			}
			for k := 1; k < maxK; k++ {
				if g[i][j-k] != '0' || g[i-k][j] != '0' {
					continue
				}
				if rowPS[i][j-1]-rowPS[i][j-k] != 0 {
					continue
				}
				if colPS[i-1][j]-colPS[i-k][j] != 0 {
					continue
				}
				if diag2[i-1][j-k+1]-diag2[i-k][j] != 0 {
					continue
				}
				ans++
			}
		}
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if g[i][j] != '0' {
				continue
			}
			maxL := hRun[i][j]
			for k := 2; k < maxL; k += 2 {
				mid := k / 2
				bj := j + k
				if bj > m || g[i][bj] != '0' {
					continue
				}
				if rowPS[i][bj-1]-rowPS[i][j] != 0 {
					continue
				}
				ci := i - mid
				cj := j + mid
				if ci >= 1 {
					if g[ci][cj] == '0' {
						if diag2[i-1][j+1]-diag2[ci][cj] == 0 && diag1[i-1][bj-1]-diag1[ci][cj] == 0 {
							ans++
						}
					}
				}
				ci = i + mid
				if ci <= n {
					if g[ci][cj] == '0' {
						if diag1[ci-1][cj-1]-diag1[i][j] == 0 && diag2[ci-1][cj+1]-diag2[i][bj] == 0 {
							ans++
						}
					}
				}
			}
		}
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if g[i][j] != '0' {
				continue
			}
			maxL := vRun[i][j]
			for k := 2; k < maxL; k += 2 {
				mid := k / 2
				bi := i + k
				if bi > n || g[bi][j] != '0' {
					continue
				}
				if colPS[bi-1][j]-colPS[i][j] != 0 {
					continue
				}
				ci := i + mid
				cj := j - mid
				if cj >= 1 && g[ci][cj] == '0' {
					if diag2[ci-1][cj+1]-diag2[i][j] == 0 && diag1[bi-1][j-1]-diag1[ci][cj] == 0 {
						ans++
					}
				}
				cj = j + mid
				if cj <= m && g[ci][cj] == '0' {
					if diag1[ci-1][cj-1]-diag1[i][j] == 0 && diag2[bi-1][j+1]-diag2[ci][cj] == 0 {
						ans++
					}
				}
			}
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	var pending string
	for {
		var line string
		if pending != "" {
			line = pending
			pending = ""
		} else {
			if !scanner.Scan() {
				break
			}
			line = strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "bad header on test %d\n", idx+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "unexpected EOF in test %d\n", idx+1)
				os.Exit(1)
			}
			grid[i] = strings.TrimSpace(scanner.Text())
		}
		expect := countTriangles(n, m, grid)
		input := fmt.Sprintf("%d %d\n", n, m)
		for i := 0; i < n; i++ {
			input += grid[i] + "\n"
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		var got int64
		if _, err := fmt.Fscan(&out, &got); err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx+1, out.String())
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx+1, expect, got)
			os.Exit(1)
		}
		idx++
		if scanner.Scan() {
			pending = strings.TrimSpace(scanner.Text())
		} else {
			pending = ""
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
