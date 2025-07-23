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

	var n, q int
	fmt.Fscan(reader, &n, &q)

	t := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &t[i])
	}

	follow := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &follow[i])
		follow[i]--
	}

	followers := make([]map[int]struct{}, n)
	for i := 0; i < n; i++ {
		followers[i] = make(map[int]struct{})
	}
	for i := 0; i < n; i++ {
		followers[follow[i]][i] = struct{}{}
	}

	deg := make([]int, n)
	for i := 0; i < n; i++ {
		deg[i] = len(followers[i])
	}

	x := make([]int64, n)
	for i := 0; i < n; i++ {
		x[i] = t[i] / int64(deg[i]+2)
	}

	followerPart := make([]int64, n)
	for i := 0; i < n; i++ {
		followerPart[follow[i]] += x[i]
	}

	ownerPart := make([]int64, n)
	for i := 0; i < n; i++ {
		ownerPart[i] = t[i] - int64(deg[i]+1)*x[i]
	}

	income := func(i int) int64 {
		return ownerPart[i] + followerPart[i] + x[follow[i]]
	}

	for ; q > 0; q-- {
		var typ int
		fmt.Fscan(reader, &typ)
		if typ == 1 {
			var a, b int
			fmt.Fscan(reader, &a, &b)
			a--
			b--
			old := follow[a]
			if old == b {
				continue
			}
			// remove from old
			followerPart[old] -= x[a]
			delete(followers[old], a)
			oldX := x[old]
			deg[old]--
			x[old] = t[old] / int64(deg[old]+2)
			ownerPart[old] = t[old] - int64(deg[old]+1)*x[old]
			followerPart[follow[old]] += x[old] - oldX

			// add to new
			follow[a] = b
			followers[b][a] = struct{}{}
			followerPart[b] += x[a]
			prevX := x[b]
			deg[b]++
			x[b] = t[b] / int64(deg[b]+2)
			ownerPart[b] = t[b] - int64(deg[b]+1)*x[b]
			followerPart[follow[b]] += x[b] - prevX
		} else if typ == 2 {
			var i int
			fmt.Fscan(reader, &i)
			i--
			fmt.Fprintln(writer, income(i))
		} else if typ == 3 {
			minV := income(0)
			maxV := minV
			for i := 1; i < n; i++ {
				val := income(i)
				if val < minV {
					minV = val
				}
				if val > maxV {
					maxV = val
				}
			}
			fmt.Fprintln(writer, minV, maxV)
		}
	}
}
