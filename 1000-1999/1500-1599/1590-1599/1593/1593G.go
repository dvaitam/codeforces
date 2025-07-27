package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	length int
	start  byte
	countS int
}

func (n Node) end() byte {
	return n.start ^ byte((n.length-1)&1)
}

func countSPrefix(start byte, t int) int {
	if start == 1 {
		return (t + 1) / 2
	}
	return t / 2
}

func countSSuffix(start byte, length, t int) int {
	pos := length - t
	if pos%2 == 1 {
		start ^= 1
	}
	return countSPrefix(start, t)
}

func merge(a, b Node) Node {
	if a.length == 0 {
		return b
	}
	if b.length == 0 {
		return a
	}
	if a.end() != b.start {
		return Node{length: a.length + b.length, start: a.start, countS: a.countS + b.countS}
	}
	t := a.length
	if b.length < t {
		t = b.length
	}
	removedA := countSSuffix(a.start, a.length, t)
	removedB := countSPrefix(b.start, t)
	if a.length == b.length {
		return Node{}
	} else if a.length > b.length {
		return Node{length: a.length - t, start: a.start, countS: a.countS - removedA}
	}
	return Node{length: b.length - t, start: b.start ^ byte(t&1), countS: b.countS - removedB}
}

var s string
var tree []Node
var size int

func build(arr []byte) {
	n := len(arr)
	size = 1
	for size < n {
		size <<= 1
	}
	tree = make([]Node, size*2)
	for i := 0; i < n; i++ {
		t := byte(0)
		if arr[i] == '[' || arr[i] == ']' {
			t = 1
		}
		tree[size+i] = Node{length: 1, start: t, countS: int(t)}
	}
	for i := size - 1; i >= 1; i-- {
		tree[i] = merge(tree[i<<1], tree[i<<1|1])
	}
}

func query(l, r int) Node {
	l += size
	r += size
	left := Node{}
	right := Node{}
	for l < r {
		if l&1 == 1 {
			left = merge(left, tree[l])
			l++
		}
		if r&1 == 1 {
			r--
			right = merge(tree[r], right)
		}
		l >>= 1
		r >>= 1
	}
	return merge(left, right)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		fmt.Fscan(reader, &s)
		arr := []byte(s)
		build(arr)
		var q int
		fmt.Fscan(reader, &q)
		for i := 0; i < q; i++ {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			node := query(l-1, r)
			fmt.Fprintln(writer, node.countS)
		}
	}
}
