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
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &p[i])
		}
		minVal := int(1e9)
		minFalse := int(1e9)
		lucky := 0
		for _, val := range p {
			if minVal > val || minFalse < val {
				// current player has winning move or no moves at all
			} else {
				// losing state for current player
				lucky++
				if val < minFalse {
					minFalse = val
				}
			}
			if val < minVal {
				minVal = val
			}
		}
		fmt.Fprintln(writer, lucky)
	}
}
