package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct {
	to     int
	parity int
}

func solveBit(bit int, n, m, k int, a [][]int) (int, bool) {
	adj := make([][]edge, n)
	limit := k - 1
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			c00, c01, c10, c11 := 0, 0, 0, 0
			for col := 0; col < m; col++ {
				b1 := (a[i][col] >> bit) & 1
				b2 := (a[j][col] >> bit) & 1
				switch {
				case b1 == 0 && b2 == 0:
					c00++
				case b1 == 0 && b2 == 1:
					c01++
				case b1 == 1 && b2 == 0:
					c10++
				default:
					c11++
				}
			}
			allowEqual := c00 <= limit && c11 <= limit
			allowDiff := c01 <= limit && c10 <= limit
			if !allowEqual && !allowDiff {
				return 0, false
			}
			if allowEqual && !allowDiff {
				adj[i] = append(adj[i], edge{j, 0})
				adj[j] = append(adj[j], edge{i, 0})
			} else if !allowEqual && allowDiff {
				adj[i] = append(adj[i], edge{j, 1})
				adj[j] = append(adj[j], edge{i, 1})
			}
		}
	}

	color := make([]int, n)
	for i := range color {
		color[i] = -1
	}
	totalOnes := 0
	for i := 0; i < n; i++ {
		if color[i] != -1 {
			continue
		}
		queue := []int{i}
		color[i] = 0
		compSize := 0
		ones := 0
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			compSize++
			if color[v] == 1 {
				ones++
			}
			for _, e := range adj[v] {
				desired := color[v] ^ e.parity
				if color[e.to] == -1 {
					color[e.to] = desired
					queue = append(queue, e.to)
				} else if color[e.to] != desired {
					return 0, false
				}
			}
		}
		if ones > compSize-ones {
			ones = compSize - ones
		}
		totalOnes += ones
	}
	return totalOnes * (1 << bit), true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m, k int
		fmt.Fscan(reader, &n, &m, &k)
		a := make([][]int, n)
		for i := 0; i < n; i++ {
			a[i] = make([]int, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(reader, &a[i][j])
			}
		}

		total := 0
		ok := true
		for bit := 0; bit < 5; bit++ {
			cost, possible := solveBit(bit, n, m, k, a)
			if !possible {
				ok = false
				break
			}
			total += cost
		}

		if !ok {
			fmt.Fprintln(writer, -1)
		} else {
			fmt.Fprintln(writer, total)
		}
	}
}
