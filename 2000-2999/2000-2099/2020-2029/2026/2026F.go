package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxP = 2000
var zeroDP = make([]int, maxP+1)

type node struct {
	price int
	tasty int
	prev  *node
	dp    []int
}

type store struct {
	front *node
	back  *node
}

func push(top *node, price, tasty int) *node {
	var base []int
	if top != nil {
		base = top.dp
	} else {
		base = zeroDP
	}
	dp := make([]int, maxP+1)
	copy(dp, base)
	for i := price; i <= maxP; i++ {
		if cand := base[i-price] + tasty; cand > dp[i] {
			dp[i] = cand
		}
	}
	return &node{price: price, tasty: tasty, prev: top, dp: dp}
}

func transfer(back *node) *node {
	var front *node
	for cur := back; cur != nil; cur = cur.prev {
		front = push(front, cur.price, cur.tasty)
	}
	return front
}

func combine(front, back *node, limit int) int {
	dpF := zeroDP
	dpB := zeroDP
	if front != nil {
		dpF = front.dp
	}
	if back != nil {
		dpB = back.dp
	}
	best := 0
	for i := 0; i <= limit; i++ {
		if val := dpF[i] + dpB[limit-i]; val > best {
			best = val
		}
	}
	return best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}
	stores := make([]store, 2, q+2) // 1-based indexing with initial empty store 1

	for ; q > 0; q-- {
		var tp int
		fmt.Fscan(in, &tp)
		switch tp {
		case 1:
			var x int
			fmt.Fscan(in, &x)
			// copy by value; nodes are immutable so pointers are safe to share
			stores = append(stores, stores[x])
		case 2:
			var x, p, t int
			fmt.Fscan(in, &x, &p, &t)
			cur := stores[x]
			cur.back = push(cur.back, p, t)
			stores[x] = cur
		case 3:
			var x int
			fmt.Fscan(in, &x)
			cur := stores[x]
			if cur.front == nil {
				cur.front = transfer(cur.back)
				cur.back = nil
			}
			cur.front = cur.front.prev
			stores[x] = cur
		case 4:
			var x, p int
			fmt.Fscan(in, &x, &p)
			cur := stores[x]
			ans := combine(cur.front, cur.back, p)
			fmt.Fprintln(out, ans)
		}
	}
}
