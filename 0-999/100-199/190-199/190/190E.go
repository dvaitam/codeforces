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

	readInt := func() int {
		var x int
		fmt.Fscan(reader, &x)
		return x
	}

	n := readInt()
	m := readInt()
	neighbors := make([][]int, n+1)
	for i := 0; i < m; i++ {
		u := readInt()
		v := readInt()
		neighbors[u] = append(neighbors[u], v)
		neighbors[v] = append(neighbors[v], u)
	}

	num := make([]int, n+2)
	bel := make([]int, n+2)
	wz := make([]int, n+2)
	dl := make([]int, n+2)
	for i := 1; i <= n; i++ {
		dl[i] = i
		wz[i] = i
	}

	top, now := 0, 0
	for i := 1; i <= n; i++ {
		if now < i {
			now = i
			top++
		}
		x := dl[i]
		bel[x] = top
		num[top]++
		td := n
		for _, j := range neighbors[x] {
			if wz[j] <= now {
				continue
			}
			if wz[j] > td {
				continue
			}
			y := dl[td]
			wz[y], wz[j] = wz[j], wz[y]
			dl[wz[y]], dl[wz[j]] = y, j
			td--
		}
		now = td
	}

	a := make([]int, n)
	for i := 1; i <= n; i++ {
		a[i-1] = i
	}
	sort.Slice(a, func(i, j int) bool {
		if bel[a[i]] != bel[a[j]] {
			return bel[a[i]] < bel[a[j]]
		}
		return a[i] < a[j]
	})

	fmt.Fprintln(writer, top)
	idx := 0
	for comp := 1; comp <= top; comp++ {
		writer.WriteString(fmt.Sprintf("%d", num[comp]))
		for k := 0; k < num[comp]; k++ {
			writer.WriteByte(' ')
			writer.WriteString(fmt.Sprintf("%d", a[idx]))
			idx++
		}
		writer.WriteByte('\n')
	}
}
