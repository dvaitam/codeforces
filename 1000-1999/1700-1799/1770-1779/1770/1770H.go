package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct{ x, y int }

func boundaryIndex(n, x, y int) int {
	if x == 1 {
		return y - 1
	} else if y == n {
		return (n - 1) + (x - 1)
	} else if x == n {
		return (n - 1) + (n - 1) + (n - y)
	} else { // y==1
		return (n - 1) + (n - 1) + (n - 1) + (n - x)
	}
}

func boundaryVertex(n, idx int) (int, int) {
	if idx < n {
		return 1, idx + 1
	} else if idx < n+(n-1) {
		return idx - n + 2, n
	} else if idx < n+2*(n-1) {
		k := idx - (n + (n - 1))
		return n, n - 1 - k
	}
	k := idx - (n + 2*(n-1))
	return n - 1 - k, 1
}

func boundaryPath(n int, sx, sy, ex, ey int) []Point {
	s := boundaryIndex(n, sx, sy)
	e := boundaryIndex(n, ex, ey)
	L := 4*n - 4
	path := []Point{{sx, sy}}
	i := s
	for i != e {
		i = (i + 1) % L
		x, y := boundaryVertex(n, i)
		path = append(path, Point{x, y})
	}
	return path
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	p := make([]int, n)
	q := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &p[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &q[i])
	}

	for i := 0; i < n; i++ {
		path := boundaryPath(n, 1, i+1, n, p[i])
		fmt.Fprint(out, len(path))
		for _, v := range path {
			fmt.Fprintf(out, " %d %d", v.x, v.y)
		}
		fmt.Fprintln(out)
	}
	for i := 0; i < n; i++ {
		path := boundaryPath(n, i+1, 1, q[i], n)
		fmt.Fprint(out, len(path))
		for _, v := range path {
			fmt.Fprintf(out, " %d %d", v.x, v.y)
		}
		fmt.Fprintln(out)
	}
}
