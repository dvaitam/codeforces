package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type pair struct{ x, y int }

func solveCase(n, m int, seats []pair) []pair {
	occ := make([][]bool, n+1)
	for i := 0; i <= n; i++ {
		occ[i] = make([]bool, m+1)
	}
	res := make([]pair, len(seats))
	for idx, p := range seats {
		x1, y1 := p.x, p.y
		if !occ[x1][y1] {
			occ[x1][y1] = true
			res[idx] = pair{x1, y1}
			continue
		}
		found := false
		var rx, ry int
		for d := 1; !found; d++ {
			for dx := 0; dx <= d && !found; dx++ {
				dy := d - dx
				rows := []int{}
				if dx == 0 {
					if x1 >= 1 && x1 <= n {
						rows = append(rows, x1)
					}
				} else {
					if x1-dx >= 1 {
						rows = append(rows, x1-dx)
					}
					if x1+dx <= n {
						rows = append(rows, x1+dx)
					}
				}
				for _, x2 := range rows {
					cols := []int{}
					if dy == 0 {
						if y1 >= 1 && y1 <= m {
							cols = append(cols, y1)
						}
					} else {
						if y1-dy >= 1 {
							cols = append(cols, y1-dy)
						}
						if y1+dy <= m {
							cols = append(cols, y1+dy)
						}
					}
					for _, y2 := range cols {
						if !occ[x2][y2] {
							rx, ry = x2, y2
							found = true
							break
						}
					}
					if found {
						break
					}
				}
			}
		}
		occ[rx][ry] = true
		res[idx] = pair{rx, ry}
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	var expected []pair
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		k, _ := strconv.Atoi(scan.Text())
		seats := make([]pair, k)
		for j := 0; j < k; j++ {
			scan.Scan()
			x, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			y, _ := strconv.Atoi(scan.Text())
			seats[j] = pair{x, y}
		}
		res := solveCase(n, m, seats)
		expected = append(expected, res...)
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < len(expected); i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for line %d\n", i+1)
			os.Exit(1)
		}
		x, _ := strconv.Atoi(outScan.Text())
		if !outScan.Scan() {
			fmt.Printf("missing output for line %d\n", i+1)
			os.Exit(1)
		}
		y, _ := strconv.Atoi(outScan.Text())
		if x != expected[i].x || y != expected[i].y {
			fmt.Printf("line %d mismatch: expected %d %d got %d %d\n", i+1, expected[i].x, expected[i].y, x, y)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
