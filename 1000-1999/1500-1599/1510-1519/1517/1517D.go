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

	hor := make([][]int, n)
	for i := 0; i < n; i++ {
		hor[i] = make([]int, m-1)
		for j := 0; j < m-1; j++ {
			fmt.Fscan(reader, &hor[i][j])
		}
	}

	ver := make([][]int, n-1)
	for i := 0; i < n-1; i++ {
		ver[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(reader, &ver[i][j])
		}
	}

	if k%2 == 1 {
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if j > 0 {
					writer.WriteByte(' ')
				}
				writer.WriteString("-1")
			}
			writer.WriteByte('\n')
		}
		return
	}

	half := k / 2
	const INF = int(1e9)
	dpPrev := make([][]int, n)
	dpCurr := make([][]int, n)
	for i := 0; i < n; i++ {
		dpPrev[i] = make([]int, m)
		dpCurr[i] = make([]int, m)
	}

	for step := 1; step <= half; step++ {
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				best := INF
				if i > 0 {
					w := ver[i-1][j]
					if dpPrev[i-1][j]+w < best {
						best = dpPrev[i-1][j] + w
					}
				}
				if i+1 < n {
					w := ver[i][j]
					if dpPrev[i+1][j]+w < best {
						best = dpPrev[i+1][j] + w
					}
				}
				if j > 0 {
					w := hor[i][j-1]
					if dpPrev[i][j-1]+w < best {
						best = dpPrev[i][j-1] + w
					}
				}
				if j+1 < m {
					w := hor[i][j]
					if dpPrev[i][j+1]+w < best {
						best = dpPrev[i][j+1] + w
					}
				}
				dpCurr[i][j] = best
			}
		}
		dpPrev, dpCurr = dpCurr, dpPrev
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				writer.WriteByte(' ')
			}
			ans := dpPrev[i][j] * 2
			fmt.Fprint(writer, ans)
		}
		writer.WriteByte('\n')
	}
}
