package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &grid[i])
		}
		visited := make([][]bool, n)
		for i := range visited {
			visited[i] = make([]bool, n)
		}
		ans := 0
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if visited[i][j] {
					continue
				}
				// compute orbit under 90 degree rotations
				coords := [][2]int{{i, j}, {j, n - 1 - i}, {n - 1 - i, n - 1 - j}, {n - 1 - j, i}}
				unique := make([][2]int, 0, 4)
				for _, c := range coords {
					found := false
					for _, u := range unique {
						if u == c {
							found = true
							break
						}
					}
					if !found {
						unique = append(unique, c)
					}
				}
				for _, c := range unique {
					visited[c[0]][c[1]] = true
				}
				ones := 0
				for _, c := range unique {
					if grid[c[0]][c[1]] == '1' {
						ones++
					}
				}
				size := len(unique)
				if ones < size-ones {
					ans += ones
				} else {
					ans += size - ones
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
