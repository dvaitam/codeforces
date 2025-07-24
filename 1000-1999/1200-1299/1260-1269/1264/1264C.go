package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

const MOD int64 = 998244353

// treap for ordered set
type node struct {
	key   int
	prio  int
	left  *node
	right *node
}

func split(root *node, key int) (l, r *node) {
	if root == nil {
		return nil, nil
	}
	if root.key < key {
		root.right, r = split(root.right, key)
		l = root
	} else {
		l, root.left = split(root.left, key)
		r = root
	}
	return
}

func merge(l, r *node) *node {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	if l.prio > r.prio {
		l.right = merge(l.right, r)
		return l
	}
	r.left = merge(l, r.left)
	return r
}

func insertNode(root *node, nd *node) *node {
	if root == nil {
		return nd
	}
	if nd.prio > root.prio {
		nd.left, nd.right = split(root, nd.key)
		return nd
	}
	if nd.key < root.key {
		root.left = insertNode(root.left, nd)
	} else {
		root.right = insertNode(root.right, nd)
	}
	return root
}

func deleteNode(root *node, key int) *node {
	if root == nil {
		return nil
	}
	if root.key == key {
		return merge(root.left, root.right)
	}
	if key < root.key {
		root.left = deleteNode(root.left, key)
	} else {
		root.right = deleteNode(root.right, key)
	}
	return root
}

func predecessor(root *node, key int) int {
	res := -1 << 60
	for root != nil {
		if root.key < key {
			if root.key > res {
				res = root.key
			}
			root = root.right
		} else {
			root = root.left
		}
	}
	return res
}

func successor(root *node, key int) int {
	res := 1<<60 - 1
	for root != nil {
		if root.key > key {
			if root.key < res {
				res = root.key
			}
			root = root.left
		} else {
			root = root.right
		}
	}
	return res
}

type orderedSet struct{ root *node }

func (s *orderedSet) insert(key int)   { s.root = insertNode(s.root, &node{key: key, prio: rand.Int()}) }
func (s *orderedSet) remove(key int)   { s.root = deleteNode(s.root, key) }
func (s *orderedSet) prev(key int) int { return predecessor(s.root, key) }
func (s *orderedSet) next(key int) int { return successor(s.root, key) }

// modular exponentiation
func powmod(a, e int64) int64 {
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}

	p := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		p[i] = x
	}

	inv100 := powmod(100, MOD-2)

	prob := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prob[i] = p[i] * inv100 % MOD
	}

	pref := make([]int64, n+1)
	pref[0] = 1
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] * prob[i] % MOD
	}

	invPref := make([]int64, n+1)
	for i := 0; i <= n; i++ {
		invPref[i] = powmod(pref[i], MOD-2)
	}

	sumInv := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		sumInv[i] = (sumInv[i-1] + invPref[i]) % MOD
	}

	segExp := func(l, r int) int64 {
		if l >= r {
			return 0
		}
		t := sumInv[r-1] - sumInv[l-1]
		t %= MOD
		if t < 0 {
			t += MOD
		}
		return pref[l-1] * t % MOD
	}

	st := &orderedSet{}
	st.insert(1)
	st.insert(n + 1)
	inSet := make(map[int]bool)
	inSet[1] = true
	inSet[n+1] = true

	ans := segExp(1, n+1)

	for i := 0; i < q; i++ {
		var u int
		fmt.Fscan(reader, &u)
		if !inSet[u] {
			// insert u
			l := st.prev(u)
			r := st.next(u)
			ans -= segExp(l, r)
			ans %= MOD
			if ans < 0 {
				ans += MOD
			}
			ans += segExp(l, u)
			ans %= MOD
			ans += segExp(u, r)
			ans %= MOD
			st.insert(u)
			inSet[u] = true
		} else {
			// remove u
			l := st.prev(u)
			r := st.next(u)
			ans -= segExp(l, u)
			ans %= MOD
			if ans < 0 {
				ans += MOD
			}
			ans -= segExp(u, r)
			ans %= MOD
			if ans < 0 {
				ans += MOD
			}
			ans += segExp(l, r)
			ans %= MOD
			st.remove(u)
			delete(inSet, u)
		}
		fmt.Fprintln(writer, ans)
	}
}
