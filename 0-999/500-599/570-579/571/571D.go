package main

// See problemD.txt for the statement.

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n int

	uParent []int
	uSize   []int
	uAdd    []int64
	uDiff   []int64

	mParent []int
	mSize   []int
	members [][]int

	base   []int64
	accumU []int64
)

func findU(x int) int {
	if uParent[x] != x {
		root := findU(uParent[x])
		uDiff[x] += uDiff[uParent[x]]
		uParent[x] = root
	}
	return uParent[x]
}

func findM(x int) int {
	for mParent[x] != x {
		x = mParent[x]
	}
	return x
}

func syncDorm(i int) {
	ru := findU(i)
	val := uAdd[ru] + uDiff[i]
	base[i] += val - accumU[i]
	accumU[i] = val
}

func unionU(a, b int) {
	ra := findU(a)
	rb := findU(b)
	if ra == rb {
		return
	}
	if uSize[ra] < uSize[rb] {
		ra, rb = rb, ra
	}
	uParent[rb] = ra
	uDiff[rb] = uAdd[rb] - uAdd[ra]
	uSize[ra] += uSize[rb]
}

func unionM(c, d int) {
	rc := findM(c)
	rd := findM(d)
	if rc == rd {
		return
	}
	if len(members[rc]) < len(members[rd]) {
		rc, rd = rd, rc
	}
	mParent[rd] = rc
	mSize[rc] += mSize[rd]
	members[rc] = append(members[rc], members[rd]...)
	members[rd] = nil
}

func raid(y int) {
	ry := findM(y)
	for _, d := range members[ry] {
		syncDorm(d)
		base[d] = 0
		ru := findU(d)
		accumU[d] = uAdd[ru] + uDiff[d]
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	uParent = make([]int, n+1)
	uSize = make([]int, n+1)
	uAdd = make([]int64, n+1)
	uDiff = make([]int64, n+1)

	mParent = make([]int, n+1)
	mSize = make([]int, n+1)
	members = make([][]int, n+1)

	base = make([]int64, n+1)
	accumU = make([]int64, n+1)

	for i := 1; i <= n; i++ {
		uParent[i] = i
		uSize[i] = 1
		mParent[i] = i
		mSize[i] = 1
		members[i] = []int{i}
	}

	for j := 0; j < m; j++ {
		var op string
		if _, err := fmt.Fscan(in, &op); err != nil {
			return
		}
		switch op {
		case "U":
			var a, b int
			fmt.Fscan(in, &a, &b)
			unionU(a, b)
		case "M":
			var c, d int
			fmt.Fscan(in, &c, &d)
			unionM(c, d)
		case "A":
			var x int
			fmt.Fscan(in, &x)
			rx := findU(x)
			uAdd[rx] += int64(uSize[rx])
		case "Z":
			var y int
			fmt.Fscan(in, &y)
			raid(y)
		case "Q":
			var q int
			fmt.Fscan(in, &q)
			syncDorm(q)
			fmt.Fprintln(out, base[q])
		}
	}
}
