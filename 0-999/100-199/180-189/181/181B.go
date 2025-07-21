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

	var n int
	fmt.Fscan(reader, &n)
	xs := make([]int, n)
	ys := make([]int, n)
	points := make(map[[2]int]struct{}, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &xs[i], &ys[i])
		points[[2]int{xs[i], ys[i]}] = struct{}{}
	}
	var ans int64
	for i := 0; i < n; i++ {
		xi, yi := xs[i], ys[i]
		for j := i + 1; j < n; j++ {
			xj, yj := xs[j], ys[j]
			sx, sy := xi+xj, yi+yj
			if sx%2 != 0 || sy%2 != 0 {
				continue
			}
			mx, my := sx/2, sy/2
			if _, ok := points[[2]int{mx, my}]; ok {
				ans++
			}
		}
	}
	fmt.Fprintln(writer, ans)
}
