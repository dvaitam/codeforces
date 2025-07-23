package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	counts := make([]int, n)
	points := 0
	for i := 0; i < m; i++ {
		var c int
		fmt.Fscan(reader, &c)
		counts[c-1]++
		if counts[c-1] == points+1 {
			ok := true
			for _, v := range counts {
				if v <= points {
					ok = false
					break
				}
			}
			if ok {
				points++
			}
		}
	}
	fmt.Println(points)
}
