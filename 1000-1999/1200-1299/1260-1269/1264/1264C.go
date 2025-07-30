package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// --- Global Variables ---
var (
	n, q    int
	invP    []int64
	segTree []Node
)

// --- Constants and I/O ---
const M = 998244353

var scanner *bufio.Scanner
var writer *bufio.Writer

func init() {
	scanner = bufio.NewScanner(os.Stdin)
	// Increase buffer size for large inputs
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	scanner.Split(bufio.ScanWords)
	writer = bufio.NewWriter(os.Stdout)
}

func nextInt() int {
	scanner.Scan()
	i, _ := strconv.Atoi(scanner.Text())
	return i
}

func writeln(a ...interface{}) {
	fmt.Fprintln(writer, a...)
}

// --- Modular Arithmetic ---
func power(a, b int64) int64 {
	res := int64(1)
	a %= M
	for b > 0 {
		if b%2 == 1 {
			res = (res * a) % M
		}
		a = (a * a) % M
		b /= 2
	}
	return res
}

func modInverse(n int64) int64 {
	return power(n, M-2)
}

// --- Segment Tree ---
type Node struct {
	val  int64
	prod int64
}

// merge combines two adjacent segments [L, M] and [M+1, R]
func merge(left, right Node) Node {
	// The formula for combining segments is:
	// val(L, R+1) = val(L, M+1) * prod(M+1, R+1) + val(M+1, R+1)
	newVal := (left.val*right.prod + right.val) % M
	newProd := (left.prod * right.prod) % M
	return Node{val: newVal, prod: newProd}
}

func buildSegTree(v, tl, tr int) {
	if tl == tr {
		segTree[v] = Node{val: invP[tl], prod: invP[tl]}
	} else {
		tm := (tl + tr) / 2
		buildSegTree(v*2, tl, tm)
		buildSegTree(v*2+1, tm+1, tr)
		segTree[v] = merge(segTree[v*2], segTree[v*2+1])
	}
}

func querySegTree(v, tl, tr, l, r int) Node {
	if l > r {
		return Node{val: 0, prod: 1} // Identity element
	}
	if l == tl && r == tr {
		return segTree[v]
	}
	tm := (tl + tr) / 2
	if r <= tm {
		return querySegTree(v*2, tl, tm, l, r)
	}
	if l > tm {
		return querySegTree(v*2+1, tm+1, tr, l, r)
	}
	leftResult := querySegTree(v*2, tl, tm, l, tm)
	rightResult := querySegTree(v*2+1, tm+1, tr, tm+1, r)
	return merge(leftResult, rightResult)
}

// --- BIT (Fenwick Tree) ---
type BIT struct {
	size int
	tree []int
}

func newBIT(size int) *BIT {
	return &BIT{
		size: size,
		tree: make([]int, size+1),
	}
}

func (b *BIT) add(idx, delta int) {
	for ; idx <= b.size; idx += idx & -idx {
		b.tree[idx] += delta
	}
}

func (b *BIT) prefixSum(idx int) int {
	if idx < 1 {
		return 0
	}
	sum := 0
	for ; idx > 0; idx -= idx & -idx {
		sum += b.tree[idx]
	}
	return sum
}

// findKth finds the index of the k-th checkpoint (k is 1-based)
func (b *BIT) findKth(k int) int {
	if k <= 0 {
		return 0
	}
	pos := 0
	sum := 0
	
	logN := 0
	for (1 << (logN + 1)) <= b.size {
		logN++
	}

	for i := logN; i >= 0; i-- {
		p := 1 << i
		if pos+p <= b.size && sum+b.tree[pos+p] < k {
			sum += b.tree[pos+p]
			pos += p
		}
	}
	return pos + 1
}


func main() {
	defer writer.Flush()

	n = nextInt()
	q = nextInt()

	p := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		p[i] = int64(nextInt())
	}

	invP = make([]int64, n+1)
	for i := 1; i <= n; i++ {
		invP[i] = (100 * modInverse(p[i])) % M
	}

	segTree = make([]Node, 4*(n+1))
	buildSegTree(1, 1, n)

	checkpointBIT := newBIT(n)
	isCheckpoint := make(map[int]bool)

	// Initial state: only mirror 1 is a checkpoint
	isCheckpoint[1] = true
	checkpointBIT.add(1, 1)
	totalE := querySegTree(1, 1, n, 1, n).val

	for i := 0; i < q; i++ {
		u := nextInt()

		var cp, cn int
		isUCheckpoint := isCheckpoint[u]
		
		numCheckpoints := checkpointBIT.prefixSum(n)

		if isUCheckpoint { // Case: removing u from checkpoints
			rankU := checkpointBIT.prefixSum(u)
			
			// Find previous checkpoint
			cpRank := rankU - 1
			cp = checkpointBIT.findKth(cpRank)
			
			// Find next checkpoint
			nextRank := rankU + 1
			if nextRank > numCheckpoints {
				cn = n + 1
			} else {
				cn = checkpointBIT.findKth(nextRank)
			}
		} else { // Case: adding u to checkpoints
			rankBeforeU := checkpointBIT.prefixSum(u - 1)

			// Find previous checkpoint
			cpRank := rankBeforeU
			// Since u > 1, cpRank >= 1 and a checkpoint must exist.
			cp = checkpointBIT.findKth(cpRank)

			// Find next checkpoint
			nextRank := rankBeforeU + 1
			if nextRank > numCheckpoints {
				cn = n + 1
			} else {
				cn = checkpointBIT.findKth(nextRank)
			}
		}

		// Calculate change in total expected value
		// Change is val(cp, u) * (1 - prod(u, cn))
		valCP_U_node := querySegTree(1, 1, n, cp, u-1)
		prodU_CN_node := querySegTree(1, 1, n, u, cn-1)

		delta := (valCP_U_node.val * (1 - prodU_CN_node.prod + M)) % M
		
		if isUCheckpoint {
			totalE = (totalE - delta + M) % M
			checkpointBIT.add(u, -1)
			isCheckpoint[u] = false
		} else {
			totalE = (totalE + delta) % M
			checkpointBIT.add(u, 1)
			isCheckpoint[u] = true
		}
		
		writeln(totalE)
	}
}
