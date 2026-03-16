package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)
const testcasesERaw = `
100
2 3
2 7 8
3 2 2
7 2 8
3 2 2
3 4
6 2 6 7
9 4 3 4
8 5 2 9
6 2 6 7
5 4 9 4
8 2 3 9
3 4
5 7 8 3
4 5 5 1
2 1 8 5
5 7 8 3
4 8 5 1
2 1 5 5
3 3
3 6 7
6 9 4
6 2 1
3 6 7
6 1 4
2 9 6
3 2
6 5
6 3
7 2
6 3
2 5
7 6
2 4
1 9 8 4
6 1 2 9
1 9 9 4
6 1 8 2
3 3
1 4 7
8 4 7
4 8 4
4 7 4
8 1 7
4 8 4
3 3
2 7 4
3 6 5
8 6 7
2 7 4
3 6 5
8 6 7
3 4
1 3 8 8
7 7 4 1
4 3 1 5
1 3 8 8
7 7 1 5
4 3 1 4
4 2
8 9
5 3
1 6
2 9
9 3
8 5
1 6
2 9
3 2
2 2
6 9
7 7
2 2
6 7
7 9
2 2
3 6
8 6
8 3
6 6
4 2
8 8
9 2
9 4
7 5
8 8
9 2
9 4
7 5
4 3
6 3 7
2 3 3
5 6 4
6 2 2
6 3 7
2 3 2
5 6 3
6 2 4
2 2
3 8
2 2
3 8
2 2
3 4
9 8 7 9
7 6 5 8
7 1 5 3
9 8 7 9
7 6 5 5
7 1 3 8
4 2
3 4
2 5
3 1
3 7
3 4
2 5
3 1
7 3
3 4
2 7 3 3
5 9 7 9
5 7 6 3
9 7 3 3
2 7 5 9
5 7 6 3
4 4
7 6 3 9
4 2 8 6
5 7 1 4
5 5 8 6
7 4 3 9
5 2 8 6
5 7 1 4
6 5 8 6
3 2
3 1
8 9
9 4
1 3
8 9
9 4
2 4
5 2 9 9
1 2 4 5
2 5 9 9
1 2 4 5
2 2
6 4
6 5
6 5
6 4
3 2
9 9
5 9
5 4
9 9
4 9
5 5
3 4
4 5 5 9
2 1 7 5
5 3 6 2
4 5 5 9
2 1 7 5
5 3 6 2
4 3
1 8 4
3 2 4
9 2 2
6 5 3
1 8 2
4 3 4
9 2 2
6 5 3
4 3
6 6 3
3 4 1
1 3 1
5 8 2
6 6 3
4 3 1
3 8 1
5 1 2
2 2
6 4
5 4
4 4
5 6
2 3
8 6 6
1 8 3
8 6 6
1 8 3
4 4
6 3 6 9
6 5 6 9
1 9 3 6
4 8 8 2
6 6 6 9
1 5 6 9
3 9 3 6
4 8 8 2
2 4
1 3 9 3
5 1 5 8
1 9 3 3
5 5 1 8
3 2
4 7
2 7
1 5
7 7
2 5
1 4
2 2
9 6
3 4
3 6
9 4
4 3
5 3 2
2 1 3
1 5 9
7 5 6
5 3 2
2 1 9
1 5 3
7 5 6
4 4
1 9 4 5
8 6 3 5
3 6 5 6
4 9 8 5
1 9 4 5
8 5 3 5
3 6 6 6
4 9 8 5
4 2
3 8
2 5
1 9
1 1
2 8
1 5
3 9
1 1
4 3
5 9 2
6 1 2
3 1 7
1 7 3
5 9 2
6 1 2
3 1 3
1 7 7
3 2
7 4
8 6
7 9
7 4
6 7
8 9
2 4
4 3 8 4
9 6 7 8
4 9 7 4
3 6 8 8
3 3
3 5 9
8 3 8
7 8 6
3 5 9
8 6 8
7 8 3
3 4
5 3 1 1
4 6 3 6
6 8 4 5
6 3 1 1
4 5 3 6
6 8 4 5
4 2
4 8
9 2
9 3
3 4
9 8
9 2
4 3
4 3
4 2
7 2
4 4
1 8
7 6
4 2
8 4
1 7
7 6
3 4
7 2 3 3
1 3 8 9
9 8 8 6
7 2 3 3
1 3 8 9
9 8 8 6
4 3
2 3 9
4 3 7
9 3 3
4 9 2
2 9 3
4 3 7
9 3 3
4 9 2
2 4
9 1 5 5
4 4 2 5
9 5 5 5
4 1 4 2
4 4
4 8 4 3
3 6 1 4
4 9 8 9
6 3 3 6
4 8 4 3
3 6 1 6
4 9 4 8
6 3 3 9
4 3
8 9 8
8 2 5
6 9 2
5 1 3
8 9 8
8 2 5
6 3 2
5 9 1
2 4
4 7 9 3
4 8 8 5
4 7 5 3
4 9 8 8
2 4
4 1 9 8
6 5 1 3
6 4 9 8
1 5 1 3
3 3
3 8 9
6 6 2
3 3 5
3 8 9
6 6 3
3 5 2
4 4
5 4 1 7
7 5 8 3
1 2 3 7
5 9 8 2
5 4 1 7
7 5 8 3
1 2 2 7
5 9 3 8
2 3
5 2 4
2 2 8
5 2 4
2 8 2
4 3
3 8 9
3 5 6
5 2 8
9 9 4
3 8 9
3 5 6
8 2 4
9 5 9
4 3
2 5 7
8 5 7
9 3 7
4 7 5
2 5 7
7 8 7
9 3 5
4 7 5
4 3
3 2 5
1 3 6
7 6 6
3 2 9
3 5 2
1 3 6
7 6 6
3 2 9
4 2
9 8
8 9
8 3
3 6
9 8
8 3
3 9
8 6
2 4
4 6 3 9
4 6 7 2
3 4 6 9
4 6 7 2
2 4
6 6 5 8
8 3 2 9
6 6 5 8
8 3 9 2
3 4
8 6 1 9
4 3 7 5
6 8 1 5
8 6 1 5
4 3 9 7
6 8 1 5
2 4
7 5 4 8
2 3 2 2
7 3 5 8
2 4 2 2
3 2
2 7
7 8
5 9
2 7
8 7
5 9
4 3
5 5 6
4 5 4
8 4 9
4 8 6
5 5 6
4 5 4
8 4 9
4 8 6
4 3
4 3 2
4 7 4
3 3 2
8 5 3
4 3 2
7 3 4
4 3 2
8 5 3
3 3
5 4 3
6 2 5
9 1 5
5 2 3
6 5 4
9 1 5
3 2
7 4
5 6
8 9
5 7
4 6
8 9
2 2
2 9
2 7
7 2
2 9
2 4
7 8 9 1
4 5 5 8
5 8 9 1
4 7 5 8
4 2
7 9
3 8
8 7
9 5
9 8
3 7
8 7
9 5
2 3
4 7 5
5 2 1
4 2 5
5 7 1
3 3
7 3 3
7 3 6
9 9 5
7 5 3
7 3 3
9 9 6
2 4
8 1 9 6
9 7 4 7
8 1 9 6
9 4 7 7
2 2
2 6
5 1
5 6
2 1
4 2
5 5
5 4
1 9
4 9
5 5
5 4
1 9
4 9
2 4
1 3 6 3
3 7 6 5
1 3 6 6
3 7 3 5
4 2
8 5
6 3
9 3
8 1
8 5
6 8
3 9
3 1
2 3
4 4 2
4 1 8
4 2 1
4 4 8
4 2
4 5
2 6
2 8
3 9
4 5
8 6
2 2
3 9
2 4
3 5 3 7
3 2 5 3
5 3 3 7
3 2 5 3
4 2
2 4
1 7
6 3
5 9
2 4
1 7
5 9
3 6
2 4
9 5 2 8
5 7 3 3
9 3 2 8
5 5 7 3
3 3
3 9 3
3 8 3
5 1 2
3 8 9
3 3 3
5 1 2
3 4
3 6 1 8
9 3 4 6
6 3 3 9
3 6 3 8
9 1 4 6
6 3 3 9
4 4
6 6 9 1
9 4 6 6
7 9 7 2
9 3 3 1
6 6 9 1
9 9 4 6
7 6 7 2
9 3 3 1
2 2
4 2
2 1
2 1
4 2
2 2
9 5
5 1
1 5
9 5
2 2
5 8
4 6
5 8
4 6
2 3
6 8 8
8 7 6
6 8 8
6 8 7
3 2
9 7
8 8
5 9
9 8
7 8
5 9
3 4
4 4 4 5
4 3 1 8
3 9 3 6
4 1 4 5
4 3 4 8
3 9 3 6
4 3
8 5 7
9 7 2
8 3 9
5 2 7
8 5 7
9 3 7
8 9 2
5 2 7
4 2
5 6
3 9
1 5
5 6
9 5
3 6
1 5
5 6
3 4
3 2 6 3
9 1 2 1
7 2 8 7
3 2 6 3
2 1 2 1
7 9 8 7
4 2
2 4
8 6
1 7
2 3
2 4
8 6
3 7
2 1
4 4
6 9 6 5
1 3 4 8
7 8 1 9
5 1 7 1
6 1 6 5
9 3 4 8
7 8 1 9
5 1 7 1
4 2
9 3
1 4
5 3
6 3
9 4
1 3
5 3
6 3
4 4
5 2 3 9
8 3 6 5
2 3 5 2
2 9 2 6
3 2 5 9
8 3 6 2
2 3 5 5
2 9 2 6
4 2
1 5
4 7
2 6
7 8
1 5
4 7
8 7
6 2
2 2
3 4
1 1
4 3
1 1
2 3
5 4 4
5 7 6
5 4 4
5 7 6
4 2
8 9
8 6
6 7
5 5
8 8
9 6
6 7
5 5
`


func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data := []byte(testcasesERaw)
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	// parse all cases into structures
	type caseData struct {
		n, m int
		s, t [][]int
	}
	cases := make([]caseData, t)
	for idx := 0; idx < t; idx++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		s := make([][]int, n)
		for i := 0; i < n; i++ {
			row := make([]int, m)
			for j := 0; j < m; j++ {
				scan.Scan()
				v, _ := strconv.Atoi(scan.Text())
				row[j] = v
			}
			s[i] = row
		}
		tt := make([][]int, n)
		for i := 0; i < n; i++ {
			row := make([]int, m)
			for j := 0; j < m; j++ {
				scan.Scan()
				v, _ := strconv.Atoi(scan.Text())
				row[j] = v
			}
			tt[i] = row
		}
		cases[idx] = caseData{n, m, s, tt}
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for idx, c := range cases {
		if !outScan.Scan() {
			fmt.Printf("missing output for case %d\n", idx+1)
			os.Exit(1)
		}
		tok := outScan.Text()
		if tok == "-1" {
			fmt.Printf("case %d: reported impossible but solution exists\n", idx+1)
			os.Exit(1)
		}
		k, err := strconv.Atoi(tok)
		if err != nil {
			fmt.Printf("case %d: invalid k\n", idx+1)
			os.Exit(1)
		}
		if k < 0 || k > 1000000 {
			fmt.Printf("case %d: invalid k %d\n", idx+1, k)
			os.Exit(1)
		}
		moves := make([][2]int, k+1)
		for i := 0; i <= k; i++ {
			if !outScan.Scan() {
				fmt.Printf("case %d: not enough coordinates\n", idx+1)
				os.Exit(1)
			}
			x, _ := strconv.Atoi(outScan.Text())
			if !outScan.Scan() {
				fmt.Printf("case %d: not enough coordinates\n", idx+1)
				os.Exit(1)
			}
			y, _ := strconv.Atoi(outScan.Text())
			moves[i] = [2]int{x - 1, y - 1}
			if x < 1 || x > c.n || y < 1 || y > c.m {
				fmt.Printf("case %d: coord out of range\n", idx+1)
				os.Exit(1)
			}
			if i > 0 {
				px, py := moves[i-1][0], moves[i-1][1]
				if abs(x-1-px) > 1 || abs(y-1-py) > 1 {
					fmt.Printf("case %d: move %d not adjacent\n", idx+1, i)
					os.Exit(1)
				}
			}
		}
		board := make([][]int, c.n)
		for i := 0; i < c.n; i++ {
			board[i] = append([]int(nil), c.s[i]...)
		}
		for i := 1; i <= k; i++ {
			x1, y1 := moves[i-1][0], moves[i-1][1]
			x2, y2 := moves[i][0], moves[i][1]
			board[x1][y1], board[x2][y2] = board[x2][y2], board[x1][y1]
		}
		for i := 0; i < c.n; i++ {
			for j := 0; j < c.m; j++ {
				if board[i][j] != c.t[i][j] {
					fmt.Printf("case %d failed: final board mismatch\n", idx+1)
					os.Exit(1)
				}
			}
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
