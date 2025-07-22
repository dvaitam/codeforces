package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
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
