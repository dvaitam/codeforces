package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"
)

const mod = 998244353
const maxN = 300000

var (
	fac     [2*maxN + 5]int
	ifac    [2*maxN + 5]int
	cat     [maxN + 5]int
	nGlobal int
)

func modPow(a, e int) int {
	res := 1
	for e > 0 {
		if e&1 == 1 {
			res = int(int64(res) * int64(a) % mod)
		}
		a = int(int64(a) * int64(a) % mod)
		e >>= 1
	}
	return res
}

func initComb() {
	fac[0] = 1
	for i := 1; i < len(fac); i++ {
		fac[i] = int(int64(fac[i-1]) * int64(i) % mod)
	}
	ifac[len(fac)-1] = modPow(fac[len(fac)-1], mod-2)
	for i := len(fac) - 2; i >= 0; i-- {
		ifac[i] = int(int64(ifac[i+1]) * int64(i+1) % mod)
	}
	cat[0] = 1
	for k := 1; k < len(cat); k++ {
		c := int(int64(fac[2*k]) * int64(ifac[k]) % mod)
		c = int(int64(c) * int64(ifac[k]) % mod)
		c = int(int64(c) * int64(modPow(k+1, mod-2)) % mod)
		cat[k] = c
	}
}

func catLen(length int) int {
	if length < 0 || length&1 == 1 {
		return 0
	}
	return cat[length/2]
}

// treap storing active children of a node, ordered by l.
var (
	treapLeft   []int
	treapRight  []int
	treapPri    []uint32
	treapFirstL []int
	treapLastR  []int
	treapProd   []int
	treapID     []int // corresponding bracket pair id
	treapCnt    int
)

func treapNew(id int) int {
	treapCnt++
	if treapCnt >= len(treapLeft) {
		size := len(treapLeft) * 2
		if size == 0 {
			size = 4
		}
		treapLeft = append(treapLeft, make([]int, size-len(treapLeft))...)
		treapRight = append(treapRight, make([]int, size-len(treapRight))...)
		treapPri = append(treapPri, make([]uint32, size-len(treapPri))...)
		treapFirstL = append(treapFirstL, make([]int, size-len(treapFirstL))...)
		treapLastR = append(treapLastR, make([]int, size-len(treapLastR))...)
		treapProd = append(treapProd, make([]int, size-len(treapProd))...)
		treapID = append(treapID, make([]int, size-len(treapID))...)
	}
	treapPri[treapCnt] = rand.Uint32()
	treapID[treapCnt] = id
	treapLeft[treapCnt], treapRight[treapCnt] = 0, 0
	treapFirstL[treapCnt], treapLastR[treapCnt] = L[id], R[id]
	treapProd[treapCnt] = 1
	return treapCnt
}

func treapPushUp(x int) {
	id := treapID[x]
	first := L[id]
	last := R[id]
	prod := 1
	if treapLeft[x] != 0 {
		prod = int(int64(prod) * int64(treapProd[treapLeft[x]]) % mod)
		gap := L[id] - treapLastR[treapLeft[x]] - 1
		prod = int(int64(prod) * int64(catLen(gap)) % mod)
		first = treapFirstL[treapLeft[x]]
	}
	if treapRight[x] != 0 {
		gap := treapFirstL[treapRight[x]] - R[id] - 1
		prod = int(int64(prod) * int64(catLen(gap)) % mod)
		prod = int(int64(prod) * int64(treapProd[treapRight[x]]) % mod)
		last = treapLastR[treapRight[x]]
	}
	treapFirstL[x] = first
	treapLastR[x] = last
	treapProd[x] = prod
}

func treapSplit(root int, key int) (int, int) {
	if root == 0 {
		return 0, 0
	}
	id := treapID[root]
	if L[id] < key {
		a, b := treapSplit(treapRight[root], key)
		treapRight[root] = a
		treapPushUp(root)
		return root, b
	}
	a, b := treapSplit(treapLeft[root], key)
	treapLeft[root] = b
	treapPushUp(root)
	return a, root
}

func treapInsert(root, node int) int {
	if root == 0 {
		return node
	}
	if treapPri[node] < treapPri[root] {
		a, b := treapSplit(root, L[treapID[node]])
		treapLeft[node] = a
		treapRight[node] = b
		treapPushUp(node)
		return node
	}
	if L[treapID[node]] < L[treapID[root]] {
		treapLeft[root] = treapInsert(treapLeft[root], node)
	} else {
		treapRight[root] = treapInsert(treapRight[root], node)
	}
	treapPushUp(root)
	return root
}

var (
	L       []int
	R       []int
	parent  []int
	childRT []int // treap root of active children
	gapProd []int // gap product for node with active children
	active  []bool
)

func computeGapProd(id int) int {
	start := 1
	end := 2 * nGlobal
	if id != 0 {
		start = L[id] + 1
		end = R[id] - 1
	}
	if childRT[id] == 0 {
		return catLen(end - start + 1)
	}
	t := childRT[id]
	res := int(int64(catLen(treapFirstL[t]-start)) * int64(treapProd[t]) % mod)
	res = int(int64(res) * int64(catLen(end-treapLastR[t])) % mod)
	return res
}

func main() {
	rand.Seed(time.Now().UnixNano())
	initComb()
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		nGlobal = n
		L = make([]int, n+1)
		R = make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &L[i], &R[i])
		}
		L[0] = 0
		R[0] = 2*n + 1

		idx := make([]int, n)
		for i := 0; i < n; i++ {
			idx[i] = i + 1
		}
		sort.Slice(idx, func(i, j int) bool {
			if L[idx[i]] == L[idx[j]] {
				return R[idx[i]] > R[idx[j]]
			}
			return L[idx[i]] < L[idx[j]]
		})

		parent = make([]int, n+1)
		stack := make([]int, 0)
		stack = append(stack, 0)
		for _, id := range idx {
			for len(stack) > 0 && R[stack[len(stack)-1]] < L[id] {
				stack = stack[:len(stack)-1]
			}
			parent[id] = stack[len(stack)-1]
			stack = append(stack, id)
		}

		childRT = make([]int, n+1)
		gapProd = make([]int, n+1)
		active = make([]bool, n+1)
		for i := 0; i <= n; i++ {
			gapProd[i] = computeGapProd(i)
		}

		answer := gapProd[0]
		fmt.Fprint(out, answer)
		for i := 1; i <= n; i++ {
			id := i
			active[id] = true
			answer = int(int64(answer) * int64(gapProd[id]) % mod)
			p := parent[id]
			oldP := gapProd[p]
			node := treapNew(id)
			childRT[p] = treapInsert(childRT[p], node)
			gapProd[p] = computeGapProd(p)
			answer = int(int64(answer) * int64(modPow(oldP, mod-2)) % mod)
			answer = int(int64(answer) * int64(gapProd[p]) % mod)
			fmt.Fprint(out, " ", answer)
		}
		fmt.Fprintln(out)
	}
}
