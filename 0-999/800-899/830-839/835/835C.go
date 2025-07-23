package main

import (
	"bufio"
	"fmt"
	"os"
)

const MaxCoord = 100

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q, c int
	if _, err := fmt.Fscan(reader, &n, &q, &c); err != nil {
		return
	}

	// brightness[t][x][y] stores total brightness of stars at (x,y) at time t
	brightness := make([][][]int, c+1)
	for t := 0; t <= c; t++ {
		brightness[t] = make([][]int, MaxCoord+1)
		for i := 0; i <= MaxCoord; i++ {
			brightness[t][i] = make([]int, MaxCoord+1)
		}
	}

	for i := 0; i < n; i++ {
		var x, y, s int
		fmt.Fscan(reader, &x, &y, &s)
		for t := 0; t <= c; t++ {
			brightness[t][x][y] += (s + t) % (c + 1)
		}
	}

	// build prefix sums for each t
	prefix := make([][][]int, c+1)
	for t := 0; t <= c; t++ {
		prefix[t] = make([][]int, MaxCoord+1)
		for i := 0; i <= MaxCoord; i++ {
			prefix[t][i] = make([]int, MaxCoord+1)
		}
		for i := 1; i <= MaxCoord; i++ {
			rowSum := 0
			for j := 1; j <= MaxCoord; j++ {
				rowSum += brightness[t][i][j]
				prefix[t][i][j] = prefix[t][i-1][j] + rowSum
			}
		}
	}

	for i := 0; i < q; i++ {
		var time, x1, y1, x2, y2 int
		fmt.Fscan(reader, &time, &x1, &y1, &x2, &y2)
		t := time % (c + 1)
		ans := prefix[t][x2][y2] - prefix[t][x1-1][y2] - prefix[t][x2][y1-1] + prefix[t][x1-1][y1-1]
		fmt.Fprintln(writer, ans)
	}
}
