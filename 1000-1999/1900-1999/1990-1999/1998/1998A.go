package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var xc, yc int64
		var k int
		fmt.Fscan(in, &xc, &yc, &k)

		points := make([][2]int64, 0, k)
		if k%2 == 1 {
			points = append(points, [2]int64{xc, yc})
		}

		for i := 1; len(points) < k; i++ {
			points = append(points, [2]int64{xc + int64(i), yc + int64(i)})
			if len(points) < k {
				points = append(points, [2]int64{xc - int64(i), yc - int64(i)})
			}
		}

		for _, p := range points {
			fmt.Fprintln(out, p[0], p[1])
		}
	}
}
