package main

import (
	"bufio"
	"fmt"
	"os"
)

type state struct {
	mask int
	val  int
}

type parent struct {
	prevKey int
	idx     int
}

func reduce(x, k int) int {
	for x%k == 0 {
		x /= k
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		sum += a[i]
	}

	maxVal := sum
	fullMask := 1<<n - 1

	// visited map from key to parent info
	visited := make(map[int]parent)
	type queueEntry struct {
		mask int
		val  int
	}
	queue := make([]queueEntry, 0)

	key := func(mask, val int) int { return mask*(maxVal+1) + val }

	for i := 0; i < n; i++ {
		m := 1 << i
		v := a[i]
		k0 := key(m, v)
		visited[k0] = parent{prevKey: -1, idx: i}
		queue = append(queue, queueEntry{mask: m, val: v})
	}

	var finalKey int = -1

	for front := 0; front < len(queue); front++ {
		cur := queue[front]
		if cur.mask == fullMask && cur.val == 1 {
			finalKey = key(cur.mask, cur.val)
			break
		}
		for j := 0; j < n; j++ {
			if cur.mask&(1<<j) != 0 {
				continue
			}
			newMask := cur.mask | (1 << j)
			newVal := reduce(cur.val+a[j], k)
			k2 := key(newMask, newVal)
			if _, ok := visited[k2]; !ok {
				visited[k2] = parent{prevKey: key(cur.mask, cur.val), idx: j}
				queue = append(queue, queueEntry{mask: newMask, val: newVal})
			}
		}
	}

	if finalKey == -1 {
		fmt.Fprintln(out, "NO")
		return
	}

	// reconstruct order of indices
	order := make([]int, 0, n)
	curKey := finalKey
	for {
		info := visited[curKey]
		order = append(order, info.idx)
		if info.prevKey == -1 {
			break
		}
		curKey = info.prevKey
	}

	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}

	// simulate operations to collect pairs
	res := a[order[0]]
	ops := make([][2]int, 0, n-1)
	for i := 1; i < len(order); i++ {
		x := res
		y := a[order[i]]
		ops = append(ops, [2]int{x, y})
		res = reduce(x+y, k)
	}

	fmt.Fprintln(out, "YES")
	for _, p := range ops {
		fmt.Fprintf(out, "%d %d\n", p[0], p[1])
	}
}
