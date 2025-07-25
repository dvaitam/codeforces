package main

import (
	"bufio"
	"fmt"
	"os"
)

type node struct {
	val  int
	prev int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		nodes := make([]node, 1, n+1) // nodes[0] is empty state
		mp := make(map[[2]int]int)
		head := 0
		cnt := map[int]int64{0: 1}
		var ans int64
		for _, x := range arr {
			if head != 0 && nodes[head].val == x {
				head = nodes[head].prev
			} else {
				key := [2]int{head, x}
				if id, ok := mp[key]; ok {
					head = id
				} else {
					nodes = append(nodes, node{val: x, prev: head})
					head = len(nodes) - 1
					mp[key] = head
				}
			}
			ans += cnt[head]
			cnt[head]++
		}
		fmt.Fprintln(out, ans)
	}
}
