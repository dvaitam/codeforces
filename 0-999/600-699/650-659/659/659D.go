package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct{ x, y int64 }

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	pts := make([]Point, n+2)
	for i := 0; i <= n; i++ {
		fmt.Fscan(reader, &pts[i].x, &pts[i].y)
	}
	pts[n+1] = pts[1]

	var area int64
	for i := 0; i < n; i++ {
		area += pts[i].x*pts[i+1].y - pts[i+1].x*pts[i].y
	}
	orient := int64(1)
	if area < 0 {
		orient = -1
	}

	cnt := 0
	for i := 0; i < n; i++ {
		v1x := pts[i+1].x - pts[i].x
		v1y := pts[i+1].y - pts[i].y
		v2x := pts[i+2].x - pts[i+1].x
		v2y := pts[i+2].y - pts[i+1].y
		cross := v1x*v2y - v1y*v2x
		if cross*orient < 0 {
			cnt++
		}
	}

	fmt.Fprintln(writer, cnt)
}
