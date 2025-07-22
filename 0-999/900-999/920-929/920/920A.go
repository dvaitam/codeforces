package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		taps := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &taps[i])
		}
		ans := 0
		for i := 1; i <= n; i++ {
			minDist := n + 1
			for _, x := range taps {
				d := x - i
				if d < 0 {
					d = -d
				}
				if d < minDist {
					minDist = d
				}
			}
			time := minDist + 1
			if time > ans {
				ans = time
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
