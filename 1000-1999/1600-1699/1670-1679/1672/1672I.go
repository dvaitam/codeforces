package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &p[i])
	}
	ans := 0
	for len(p) > 0 {
		minDiff := 1 << 30
		idx := 0
		for i, v := range p {
			d := abs(i + 1 - v)
			if d < minDiff {
				minDiff = d
				idx = i
			}
		}
		if minDiff > ans {
			ans = minDiff
		}
		val := p[idx]
		p = append(p[:idx], p[idx+1:]...)
		for i := 0; i < len(p); i++ {
			if p[i] > val {
				p[i]--
			}
		}
	}
	fmt.Fprintln(writer, ans)
}
