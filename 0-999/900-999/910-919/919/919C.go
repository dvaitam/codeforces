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

	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &grid[i])
	}

	if k == 1 {
		count := 0
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if grid[i][j] == '.' {
					count++
				}
			}
		}
		fmt.Fprintln(writer, count)
		return
	}

	ans := 0
	// horizontal
	for i := 0; i < n; i++ {
		length := 0
		for j := 0; j < m; j++ {
			if grid[i][j] == '.' {
				length++
			} else {
				if length >= k {
					ans += length - k + 1
				}
				length = 0
			}
		}
		if length >= k {
			ans += length - k + 1
		}
	}

	// vertical
	for j := 0; j < m; j++ {
		length := 0
		for i := 0; i < n; i++ {
			if grid[i][j] == '.' {
				length++
			} else {
				if length >= k {
					ans += length - k + 1
				}
				length = 0
			}
		}
		if length >= k {
			ans += length - k + 1
		}
	}

	fmt.Fprintln(writer, ans)
}
