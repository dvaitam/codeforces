package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	floors := make([][]byte, n)
	for i := 0; i < n; i++ {
		var line string
		fmt.Fscan(in, &line)
		floors[n-1-i] = []byte(line)
	}

	col := 0
	dir := 1
	var total int64

	for floor := n - 1; floor > 0; floor-- {
		spent, nextCol, nextDir, ok := processFloor(floors[floor], floors[floor-1], col, dir)
		if !ok {
			fmt.Println("Never")
			return
		}
		total += spent
		col = nextCol
		dir = nextDir
	}

	fmt.Println(total)
}

func processFloor(row []byte, below []byte, col int, dir int) (int64, int, int, bool) {
	m := len(row)
	compLeft := col
	for compLeft > 0 && row[compLeft-1] != '#' {
		compLeft--
	}
	compRight := col
	for compRight+1 < m && row[compRight+1] != '#' {
		compRight++
	}

	fallable := make([]bool, m)
	hasHole := false
	for i := compLeft; i <= compRight; i++ {
		if below[i] == '.' {
			fallable[i] = true
			hasHole = true
		}
	}
	if !hasHole {
		return 0, 0, 0, false
	}

	nextRight := make([]int, m)
	next := -1
	for i := m - 1; i >= 0; i-- {
		if fallable[i] {
			next = i
		}
		nextRight[i] = next
	}
	nextLeft := make([]int, m)
	next = -1
	for i := 0; i < m; i++ {
		if fallable[i] {
			next = i
		}
		nextLeft[i] = next
	}

	left, right := col, col
	pos := col
	var spent int64

	for {
		if fallable[pos] {
			spent++
			return spent, pos, dir, true
		}

		target := right
		if dir == -1 {
			target = left
		}

		if pos != target {
			if dir == 1 {
				cand := nextRight[pos]
				if cand != -1 && cand <= target {
					dist := cand - pos
					spent += int64(dist)
					pos = cand
					continue
				}
			} else {
				cand := nextLeft[pos]
				if cand != -1 && cand >= target {
					dist := pos - cand
					spent += int64(dist)
					pos = cand
					continue
				}
			}
			dist := target - pos
			if dist < 0 {
				dist = -dist
			}
			spent += int64(dist)
			pos = target
		}

		nextIdx := pos + dir
		if nextIdx < compLeft || nextIdx > compRight || row[nextIdx] == '#' {
			spent++
			dir = -dir
			continue
		}
		if row[nextIdx] == '+' {
			row[nextIdx] = '.'
			spent++
			dir = -dir
			continue
		}

		spent++
		pos = nextIdx
		if dir == 1 {
			right++
		} else {
			left--
		}
	}
}
