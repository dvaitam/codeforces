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

	grid := make([][][]bool, n)
	for i := 0; i < n; i++ {
		grid[i] = make([][]bool, m)
		for j := 0; j < m; j++ {
			var s string
			fmt.Fscan(reader, &s)
			grid[i][j] = make([]bool, k)
			for t := 0; t < k && t < len(s); t++ {
				if s[t] == '1' {
					grid[i][j][t] = true
				}
			}
		}
	}

	// directions
	pred := [3][3]int{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}}
	succ := [3][3]int{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}
	count := 0

	for x := 0; x < n; x++ {
		for y := 0; y < m; y++ {
			for z := 0; z < k; z++ {
				if !grid[x][y][z] {
					continue
				}
				// gather predecessors and successors
				preds := [][3]int{}
				succs := [][3]int{}

				for i := 0; i < 3; i++ {
					px, py, pz := x+pred[i][0], y+pred[i][1], z+pred[i][2]
					if px >= 0 && py >= 0 && pz >= 0 && grid[px][py][pz] {
						preds = append(preds, [3]int{px, py, pz})
					}
				}
				for j := 0; j < 3; j++ {
					qx, qy, qz := x+succ[j][0], y+succ[j][1], z+succ[j][2]
					if qx < n && qy < m && qz < k && grid[qx][qy][qz] {
						succs = append(succs, [3]int{qx, qy, qz})
					}
				}

				if len(preds) == 0 || len(succs) == 0 {
					continue
				}

				critical := false

				for pi, p := range preds {
					for qi := range succs {
						if critical {
							break
						}
						if pi == qi {
							critical = true
							break
						}
						rx := p[0] + succ[qi][0]
						ry := p[1] + succ[qi][1]
						rz := p[2] + succ[qi][2]
						if rx < 0 || ry < 0 || rz < 0 || rx >= n || ry >= m || rz >= k || !grid[rx][ry][rz] {
							critical = true
						}
					}
				}

				if critical {
					count++
				}
			}
		}
	}

	fmt.Fprintln(writer, count)
}
