package main

import (
	"bufio"
	"fmt"
	"os"
)

const limitValue = 1000000

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int
	if _, err := fmt.Fscan(in, &n, &m, &q); err != nil {
		return
	}

	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		grid[i] = []byte(s)
	}

	values := make([]int, q)
	index := make(map[int]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &values[i])
		index[values[i]] = i
	}

	counts := make([]int64, q)

	for r := 0; r < n; r++ {
		for c := 0; c < m; c++ {
			ch := grid[r][c]
			if ch >= '1' && ch <= '9' {
				if idxPos, ok := index[int(ch-'0')]; ok {
					counts[idxPos]++
				}
			}
		}
	}

	dirs := [][2]int{
		{0, 1}, {0, -1},
		{1, 0}, {-1, 0},
		{1, 1}, {-1, -1},
		{1, -1}, {-1, 1},
	}
	for _, d := range dirs {
		processDirection(grid, n, m, d[0], d[1], index, counts)
	}

	for i := 0; i < q; i++ {
		fmt.Fprintln(out, counts[i])
	}
}

func processDirection(grid [][]byte, n, m, dx, dy int, idx map[int]int, counts []int64) {
	for r := 0; r < n; r++ {
		for c := 0; c < m; c++ {
			prevR, prevC := r-dx, c-dy
			if prevR >= 0 && prevR < n && prevC >= 0 && prevC < m {
				continue
			}
			line := make([]byte, 0)
			for rr, cc := r, c; rr >= 0 && rr < n && cc >= 0 && cc < m; rr, cc = rr+dx, cc+dy {
				line = append(line, grid[rr][cc])
			}
			processLine(line, idx, counts)
		}
	}
}

func processLine(line []byte, idx map[int]int, counts []int64) {
	limit := int64(limitValue)
	L := len(line)
	for start := 0; start < L; start++ {
		if line[start] < '1' || line[start] > '9' {
			continue
		}
		curSum := int64(0)
		curTerm := int64(0)
		curNum := int64(0)
		lastOp := byte('+')
		for pos := start; pos < L; pos++ {
			ch := line[pos]
			if ch >= '1' && ch <= '9' {
				curNum = curNum*10 + int64(ch-'0')
				if curNum > limit {
					break
				}
				if pos == start {
					continue
				}
				var value int64
				if lastOp == '+' {
					temp := curSum + curTerm
					if temp > limit {
						break
					}
					value = temp + curNum
				} else {
					prod := curTerm * curNum
					if prod > limit {
						break
					}
					value = curSum + prod
				}
				if value > limit {
					break
				}
				if idxPos, ok := idx[int(value)]; ok {
					counts[idxPos]++
				}
			} else {
				if ch != '+' && ch != '*' {
					break
				}
				if curNum == 0 {
					break
				}
				if lastOp == '+' {
					curSum += curTerm
					if curSum > limit {
						break
					}
					curTerm = curNum
				} else {
					curTerm *= curNum
					if curTerm > limit {
						break
					}
				}
				curNum = 0
				lastOp = ch
			}
		}
	}
}
