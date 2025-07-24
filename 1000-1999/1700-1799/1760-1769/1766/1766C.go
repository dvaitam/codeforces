package main

import (
	"bufio"
	"fmt"
	"os"
)

func check(start int, row0, row1 string) bool {
	n := len(row0)
	visited := [2][]bool{make([]bool, n), make([]bool, n)}
	r := start
	c := 0
	for {
		if r == 0 {
			if row0[c] != 'B' || visited[0][c] {
				return false
			}
			visited[0][c] = true
		} else {
			if row1[c] != 'B' || visited[1][c] {
				return false
			}
			visited[1][c] = true
		}
		other := 1 - r
		if other == 0 {
			if row0[c] == 'B' && !visited[0][c] {
				r = other
				continue
			}
		} else {
			if row1[c] == 'B' && !visited[1][c] {
				r = other
				continue
			}
		}
		if c+1 < n {
			if r == 0 && row0[c+1] == 'B' {
				c++
				continue
			}
			if r == 1 && row1[c+1] == 'B' {
				c++
				continue
			}
		}
		break
	}
	for i := 0; i < n; i++ {
		if row0[i] == 'B' && !visited[0][i] {
			return false
		}
		if row1[i] == 'B' && !visited[1][i] {
			return false
		}
	}
	return true
}

func solve(row0, row1 string) bool {
	starts := []int{}
	if row0[0] == 'B' {
		starts = append(starts, 0)
	}
	if row1[0] == 'B' {
		starts = append(starts, 1)
	}
	for _, st := range starts {
		if check(st, row0, row1) {
			return true
		}
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var m int
		fmt.Fscan(reader, &m)
		var row0, row1 string
		fmt.Fscan(reader, &row0)
		fmt.Fscan(reader, &row1)
		if solve(row0, row1) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
