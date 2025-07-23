package main

import (
	"bufio"
	"fmt"
	"os"
)

func checkHorizontal(n, m int, g [][]byte) bool {
	if n%3 != 0 {
		return false
	}
	h := n / 3
	colors := make(map[byte]bool)
	for k := 0; k < 3; k++ {
		col := g[k*h][0]
		for i := k * h; i < (k+1)*h; i++ {
			for j := 0; j < m; j++ {
				if g[i][j] != col {
					return false
				}
			}
		}
		if colors[col] {
			return false
		}
		colors[col] = true
	}
	return len(colors) == 3
}

func checkVertical(n, m int, g [][]byte) bool {
	if m%3 != 0 {
		return false
	}
	w := m / 3
	colors := make(map[byte]bool)
	for k := 0; k < 3; k++ {
		col := g[0][k*w]
		for j := k * w; j < (k+1)*w; j++ {
			for i := 0; i < n; i++ {
				if g[i][j] != col {
					return false
				}
			}
		}
		if colors[col] {
			return false
		}
		colors[col] = true
	}
	return len(colors) == 3
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		grid[i] = []byte(s)
	}
	if checkHorizontal(n, m, grid) || checkVertical(n, m, grid) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
