package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x, y   int
	w      int
	indeg  int
	n1, n2 int // neighbors indices (right and up), -1 if none
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	points := make([]Point, n)
	pos := make(map[[2]int]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &points[i].x, &points[i].y)
		points[i].w = points[i].y - points[i].x
		points[i].n1, points[i].n2 = -1, -1
		pos[[2]int{points[i].x, points[i].y}] = i
	}

	// build edges
	for i := 0; i < n; i++ {
		x := points[i].x
		y := points[i].y
		if j, ok := pos[[2]int{x + 1, y}]; ok {
			points[i].n1 = j
			points[j].indeg++
		}
		if j, ok := pos[[2]int{x, y + 1}]; ok {
			points[i].n2 = j
			points[j].indeg++
		}
	}

	wseq := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &wseq[i])
	}

	// queue of available points per w
	avail := make(map[int][]int)
	for i := 0; i < n; i++ {
		if points[i].indeg == 0 {
			w := points[i].w
			avail[w] = append(avail[w], i)
		}
	}

	ans := make([][2]int, n)
	for i := 0; i < n; i++ {
		w := wseq[i]
		lst := avail[w]
		if len(lst) == 0 {
			fmt.Fprintln(writer, "NO")
			return
		}
		idx := lst[len(lst)-1]
		avail[w] = lst[:len(lst)-1]
		ans[i] = [2]int{points[idx].x, points[idx].y}

		// remove idx
		if points[idx].n1 != -1 {
			j := points[idx].n1
			points[j].indeg--
			if points[j].indeg == 0 {
				wj := points[j].w
				avail[wj] = append(avail[wj], j)
			}
		}
		if points[idx].n2 != -1 {
			j := points[idx].n2
			points[j].indeg--
			if points[j].indeg == 0 {
				wj := points[j].w
				avail[wj] = append(avail[wj], j)
			}
		}
	}

	fmt.Fprintln(writer, "YES")
	for i := 0; i < n; i++ {
		fmt.Fprintf(writer, "%d %d\n", ans[i][0], ans[i][1])
	}
}
