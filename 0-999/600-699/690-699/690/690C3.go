package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program computes the brain latency (tree diameter) after each
// new brain is connected. The input gives for each brain k (k >= 2) the
// parent brain pk it is attached to. After each insertion we output the
// current diameter of the tree.

const maxN = 200005
const logN = 20

var up [logN][maxN]int
var depth [maxN]int

// lca returns the lowest common ancestor of u and v using binary lifting.
func lca(u, v int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for i := 0; i < logN; i++ {
		if diff>>i&1 == 1 {
			u = up[i][u]
		}
	}
	if u == v {
		return u
	}
	for i := logN - 1; i >= 0; i-- {
		if up[i][u] != up[i][v] {
			u = up[i][u]
			v = up[i][v]
		}
	}
	return up[0][u]
}

// dist returns the distance between u and v in the tree.
func dist(u, v int) int {
	w := lca(u, v)
	return depth[u] + depth[v] - 2*depth[w]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	a, b, d := 1, 1, 0 // current diameter endpoints and length
	for k := 2; k <= n; k++ {
		var p int
		fmt.Fscan(reader, &p)
		up[0][k] = p
		depth[k] = depth[p] + 1
		for i := 1; i < logN; i++ {
			up[i][k] = up[i-1][up[i-1][k]]
		}

		da := dist(k, a)
		db := dist(k, b)
		if da > d {
			b = k
			d = da
		}
		if db > d {
			a = k
			d = db
		}

		if k > 2 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, d)
	}
	if n > 1 {
		fmt.Fprintln(writer)
	}
}
