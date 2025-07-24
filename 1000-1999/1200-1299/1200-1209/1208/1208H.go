package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n         int
	k         int
	children  [][]int
	parent    []int
	order     []int
	revOrder  []int
	isLeaf    []bool
	leafColor []int
	color     []int
	blueCount []int
	redCount  []int
)

func buildOrder() {
	children = make([][]int, n+1)
	parent = make([]int, n+1)
	order = make([]int, 0, n)
	stack := []int{1}
	parent[1] = 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, u := range tempAdj[v] {
			if u != parent[v] {
				parent[u] = v
				children[v] = append(children[v], u)
				stack = append(stack, u)
			}
		}
	}
	revOrder = make([]int, len(order))
	for i, v := range order {
		revOrder[len(order)-1-i] = v
	}
}

var tempAdj [][]int

func recomputeAll() {
	for _, v := range revOrder {
		if isLeaf[v] {
			color[v] = leafColor[v]
			blueCount[v] = 0
			redCount[v] = 0
		} else {
			b := 0
			r := 0
			for _, c := range children[v] {
				if color[c] == 1 {
					b++
				} else {
					r++
				}
			}
			blueCount[v] = b
			redCount[v] = r
			if b-r >= k {
				color[v] = 1
			} else {
				color[v] = 0
			}
		}
	}
}

func changeLeaf(v int, c int) {
	if color[v] == c {
		leafColor[v] = c
		return
	}
	oldColor := color[v]
	color[v] = c
	leafColor[v] = c
	child := v
	newColor := c
	p := parent[child]
	for p != 0 {
		if oldColor == 1 {
			blueCount[p]--
		} else {
			redCount[p]--
		}
		if newColor == 1 {
			blueCount[p]++
		} else {
			redCount[p]++
		}
		oldParentColor := color[p]
		var newParentColor int
		if blueCount[p]-redCount[p] >= k {
			newParentColor = 1
		} else {
			newParentColor = 0
		}
		if newParentColor == oldParentColor {
			break
		}
		color[p] = newParentColor
		oldColor = oldParentColor
		newColor = newParentColor
		child = p
		p = parent[p]
	}
}

func updateK(newK int) {
	k = newK
	recomputeAll()
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &n, &k)

	tempAdj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		tempAdj[u] = append(tempAdj[u], v)
		tempAdj[v] = append(tempAdj[v], u)
	}

	s := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &s[i])
	}

	buildOrder()

	isLeaf = make([]bool, n+1)
	leafColor = make([]int, n+1)
	color = make([]int, n+1)
	blueCount = make([]int, n+1)
	redCount = make([]int, n+1)

	for i := 1; i <= n; i++ {
		if s[i] != -1 && i != 1 {
			isLeaf[i] = true
			leafColor[i] = s[i]
		} else {
			isLeaf[i] = false
		}
	}

	recomputeAll()

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var v int
			fmt.Fscan(reader, &v)
			fmt.Fprintln(writer, color[v])
		} else if t == 2 {
			var v, c int
			fmt.Fscan(reader, &v, &c)
			changeLeaf(v, c)
		} else if t == 3 {
			var h int
			fmt.Fscan(reader, &h)
			updateK(h)
		}
	}
}
