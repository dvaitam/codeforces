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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(reader, &a[i][j])
		}
	}
	up := make([]int, m)
	rowMin := make([]int, n)
	for j := 0; j < m; j++ {
		up[j] = 0
	}
	// compute earliest row for each column
	for i := 0; i < n; i++ {
		minVal := i
		if i == 0 {
			for j := 0; j < m; j++ {
				up[j] = 0
			}
			minVal = 0
		} else {
			for j := 0; j < m; j++ {
				if a[i][j] < a[i-1][j] {
					up[j] = i
				}
				if up[j] < minVal {
					minVal = up[j]
				}
			}
		}
		// for first row we didn't compute above; ensure rowMin[0]
		if i == 0 {
			minVal = 0
		}
		rowMin[i] = minVal
	}

	var k int
	fmt.Fscan(reader, &k)
	for ; k > 0; k-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		l--
		r--
		if rowMin[r] <= l {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
