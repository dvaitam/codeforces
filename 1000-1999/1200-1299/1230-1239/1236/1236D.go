package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}

	rowObs := make(map[int][]int)
	colObs := make(map[int][]int)
	for i := 0; i < k; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		rowObs[x] = append(rowObs[x], y)
		colObs[y] = append(colObs[y], x)
	}

	for _, v := range rowObs {
		sort.Ints(v)
	}
	for _, v := range colObs {
		sort.Ints(v)
	}

	total := int64(n)*int64(m) - int64(k)
	visited := int64(1)
	x, y := 1, 1
	dir := 0
	top, bottom := 1, n
	left, right := 1, m

	// Helper functions
	nextRight := func(r, c, limit int) int {
		arr := rowObs[r]
		idx := sort.SearchInts(arr, c+1)
		next := limit
		if idx < len(arr) && arr[idx] <= limit {
			next = arr[idx] - 1
		}
		if next > limit {
			next = limit
		}
		return next
	}
	nextLeft := func(r, c, limit int) int {
		arr := rowObs[r]
		idx := sort.SearchInts(arr, c)
		next := limit
		if idx > 0 && arr[idx-1] >= limit {
			next = arr[idx-1] + 1
		}
		if next < limit {
			next = limit
		}
		return next
	}
	nextDown := func(c, r, limit int) int {
		arr := colObs[c]
		idx := sort.SearchInts(arr, r+1)
		next := limit
		if idx < len(arr) && arr[idx] <= limit {
			next = arr[idx] - 1
		}
		if next > limit {
			next = limit
		}
		return next
	}
	nextUp := func(c, r, limit int) int {
		arr := colObs[c]
		idx := sort.SearchInts(arr, r)
		next := limit
		if idx > 0 && arr[idx-1] >= limit {
			next = arr[idx-1] + 1
		}
		if next < limit {
			next = limit
		}
		return next
	}

	// If grid has only starting cell
	if total == 1 {
		fmt.Fprintln(writer, "Yes")
		return
	}

	for {
		moved := false
		switch dir {
		case 0: // right
			ny := nextRight(x, y, right)
			if ny > y {
				visited += int64(ny - y)
				y = ny
				top = x + 1
				moved = true
			}
		case 1: // down
			nx := nextDown(y, x, bottom)
			if nx > x {
				visited += int64(nx - x)
				x = nx
				right = y - 1
				moved = true
			}
		case 2: // left
			ny := nextLeft(x, y, left)
			if ny < y {
				visited += int64(y - ny)
				y = ny
				bottom = x - 1
				moved = true
			}
		case 3: // up
			nx := nextUp(y, x, top)
			if nx < x {
				visited += int64(x - nx)
				x = nx
				left = y + 1
				moved = true
			}
		}
		if !moved {
			break
		}
		dir = (dir + 1) % 4
	}

	if visited == total {
		fmt.Fprintln(writer, "Yes")
	} else {
		fmt.Fprintln(writer, "No")
	}
}
