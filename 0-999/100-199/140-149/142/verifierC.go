package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	n, m int
	ans  int
	rem  int
	p    [9][9]byte
	a    = [3][12]byte{
		{'A', 'A', 'A', '.', 'A', '.', 'A', '.', '.', '.', '.', 'A'},
		{'.', 'A', '.', '.', 'A', '.', 'A', 'A', 'A', 'A', 'A', 'A'},
		{'.', 'A', '.', 'A', 'A', 'A', 'A', '.', '.', '.', '.', 'A'},
	}
)

func copyAns() {}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func saya(x, y, move int) {
	if x >= n-2 {
		ans = max(ans, move)
		return
	}
	if y >= m-2 {
		if p[x][y] == '.' {
			rem--
		}
		if p[x][y+1] == '.' {
			rem--
		}
		saya(x+1, 0, move)
		if p[x][y] == '.' {
			rem++
		}
		if p[x][y+1] == '.' {
			rem++
		}
		return
	}
	if rem/5 <= ans-move {
		return
	}
	if p[x][y] == '.' {
		rem--
	}
	for d := 0; d < 12; d += 3 {
		flag := false
		for i := 0; i < 3 && !flag; i++ {
			for j := 0; j < 3; j++ {
				if a[i][d+j] == 'A' && p[x+i][y+j] != '.' {
					flag = true
					break
				}
			}
		}
		if !flag {
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					if a[i][d+j] == 'A' {
						p[x+i][y+j] = 'A'
					}
				}
			}
			rem -= 5
			saya(x, y+1, move+1)
			rem += 5
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					if a[i][d+j] == 'A' {
						p[x+i][y+j] = '.'
					}
				}
			}
		}
	}
	saya(x, y+1, move)
	if p[x][y] == '.' {
		rem++
	}
}

func solveCase(nn, mm int) int {
	n, m = nn, mm
	ans = 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			p[i][j] = '.'
		}
	}
	rem = n * m
	saya(0, 0, 0)
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Println("could not open testcasesC.txt:", err)
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
		var nn, mm int
		if _, err := fmt.Sscan(line, &nn, &mm); err != nil {
			fmt.Printf("bad test case on line %d\n", idx)
			os.Exit(1)
		}
		expect := solveCase(nn, mm)
		input := fmt.Sprintf("%d %d\n", nn, mm)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		var errBuf bytes.Buffer
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Fscan(bytes.NewReader(out.Bytes()), &got); err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx, out.String())
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
