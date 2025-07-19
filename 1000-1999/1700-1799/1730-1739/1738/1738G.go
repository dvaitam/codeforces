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

	var t, n, k int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for tc := 0; tc < t; tc++ {
		fmt.Fscan(reader, &n, &k)
		// allocate arrays with padding
		f := make([][]int, n+2)
		vst := make([][]bool, n+2)
		for i := 0; i < n+2; i++ {
			f[i] = make([]int, n+2)
			vst[i] = make([]bool, n+2)
		}
		mx := make([][]int, k)
		for i := 0; i < k; i++ {
			mx[i] = make([]int, n+2)
		}
		// read grid and init f
		for i := 1; i <= n; i++ {
			for j := 1; j <= n; j++ {
				var c byte
				// read next non-space byte
				for {
					b, err := reader.ReadByte()
					if err != nil {
						break
					}
					if b == '0' || b == '1' {
						c = b
						break
					}
				}
				if c == '0' {
					f[i][j] = 1
				} else {
					f[i][j] = 0
				}
			}
		}
		// DP phase
		isNo := false
		for i := n; i >= 1 && !isNo; i-- {
			for j := n; j >= 1; j-- {
				// extend diagonal
				if f[i+1][j+1] > 0 {
					f[i][j] += f[i+1][j+1]
				}
				// take max of right and down
				if f[i+1][j] > f[i][j] {
					f[i][j] = f[i+1][j]
				}
				if f[i][j+1] > f[i][j] {
					f[i][j] = f[i][j+1]
				}
				if f[i][j] == k {
					writer.WriteString("NO\n")
					isNo = true
					break
				}
				if mx[f[i][j]][j] == 0 {
					mx[f[i][j]][j] = i
				}
			}
		}
		if isNo {
			continue
		}
		// Build paths
		for level := k - 1; level >= 1; level-- {
			// suffix max for mx[level]
			for j := n - 1; j >= 1; j-- {
				if mx[level][j+1] > mx[level][j] {
					mx[level][j] = mx[level][j+1]
				}
			}
			x, y := n, 1
			for y <= n && vst[x][y] {
				y++
			}
			for y <= n {
				vst[x][y] = true
				if (y == n || x != mx[level][y+1]) && x > 1 && !vst[x-1][y] {
					x--
				} else {
					y++
				}
			}
		}
		// Output YES and matrix
		writer.WriteString("YES\n")
		for i := 1; i <= n; i++ {
			for j := 1; j <= n; j++ {
				if vst[i][j] {
					writer.WriteByte('1')
				} else {
					writer.WriteByte('0')
				}
			}
			writer.WriteByte('\n')
		}
	}
}
