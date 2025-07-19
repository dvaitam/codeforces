package main

import (
	"bufio"
	"fmt"
	"os"
)

func getChar(x int) byte {
	if x < 26 {
		return byte('a' + x)
	} else if x < 52 {
		return byte('A' + (x - 26))
	}
	return byte('0' + (x - 52))
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for T > 0 {
		T--
		var n, m, k int
		fmt.Fscan(reader, &n, &m, &k)
		grid := make([][]byte, n)
		rCount := 0
			for i := 0; i < n; i++ {
           var s string
           fmt.Fscan(reader, &s)
           row := []byte(s)
           for j := 0; j < m; j++ {
               if row[j] == 'R' {
                   rCount++
               }
           }
           grid[i] = row
			}
		// distribute counts
		counts := make([]int, k)
		for i := 0; i < rCount; i++ {
			counts[i%k]++
		}
		res := make([][]byte, n)
		for i := range res {
			res[i] = make([]byte, m)
		}
		chi := 0
		for i := 0; i < n; i++ {
			// determine traversal order
			if i%2 == 0 {
				for j := 0; j < m; j++ {
					for chi < k-1 && counts[chi] == 0 {
						chi++
					}
					if grid[i][j] == 'R' {
						counts[chi]--
					}
					res[i][j] = getChar(chi)
				}
			} else {
				for jj := m - 1; jj >= 0; jj-- {
					for chi < k-1 && counts[chi] == 0 {
						chi++
					}
					if grid[i][jj] == 'R' {
						counts[chi]--
					}
					res[i][jj] = getChar(chi)
				}
			}
		}
		// output
		for i := 0; i < n; i++ {
			writer.Write(res[i])
			writer.WriteByte('\n')
		}
	}
}
