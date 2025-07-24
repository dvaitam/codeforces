package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	maxQ = 300000 + 5
	maxV = 20
)

var (
	parent [maxQ]int
	val    [maxQ]int
	cnt    [maxQ][maxV]int
	sum    int64
)

func update(u int) {
	old := val[u]
	newVal := 0
	for newVal < maxV {
		if cnt[u][newVal] == 0 {
			break
		}
		newVal++
	}
	if newVal == old {
		return
	}
	val[u] = newVal
	sum += int64(newVal - old)
	p := parent[u]
	if p != 0 {
		cnt[p][old]--
		cnt[p][newVal]++
		update(p)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)

	current := 1
	for i := 1; i <= q; i++ {
		var x int
		fmt.Fscan(in, &x)
		current++
		node := current
		parent[node] = x
		// new node value is 0, sum unchanged
		cnt[x][0]++
		update(x)
		fmt.Fprintln(out, sum)
	}
}
