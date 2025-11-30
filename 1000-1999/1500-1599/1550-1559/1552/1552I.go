package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const MOD int64 = 998244353

type Set struct {
	bits     [2]uint64
	size     int
	children []*Set
}

func newSet(nums []int) *Set {
	s := &Set{}
	for _, v := range nums {
		v--
		if v < 64 {
			s.bits[0] |= 1 << uint(v)
		} else {
			s.bits[1] |= 1 << uint(v-64)
		}
	}
	s.size = bits.OnesCount64(s.bits[0]) + bits.OnesCount64(s.bits[1])
	return s
}

func newUniversal(n int) *Set {
	nums := make([]int, n)
	for i := 1; i <= n; i++ {
		nums[i-1] = i
	}
	return newSet(nums)
}

func (s *Set) subsetOf(t *Set) bool {
	return s.bits[0]&^t.bits[0] == 0 && s.bits[1]&^t.bits[1] == 0
}

func (s *Set) intersects(t *Set) bool {
	return (s.bits[0]&t.bits[0]) != 0 || (s.bits[1]&t.bits[1]) != 0
}

var fac []int64

func calc(node *Set) int64 {
	sum := 0
	res := int64(1)
	for _, ch := range node.children {
		res = res * calc(ch) % MOD
		sum += ch.size
	}
	blocks := len(node.children) + node.size - sum
	if blocks < 0 {
		return 0
	}
	res = res * fac[blocks] % MOD
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	mp := make(map[[2]uint64]bool)
	sets := []*Set{}
	for i := 0; i < m; i++ {
		var q int
		fmt.Fscan(reader, &q)
		nums := make([]int, q)
		for j := 0; j < q; j++ {
			fmt.Fscan(reader, &nums[j])
		}
		st := newSet(nums)
		key := [2]uint64{st.bits[0], st.bits[1]}
		if !mp[key] {
			mp[key] = true
			sets = append(sets, st)
		}
	}

	// check laminar property
	for i := 0; i < len(sets); i++ {
		for j := i + 1; j < len(sets); j++ {
			a := sets[i]
			b := sets[j]
			if a.intersects(b) && !a.subsetOf(b) && !b.subsetOf(a) {
				fmt.Println(0)
				return
			}
		}
	}

	root := newUniversal(n)
	all := append([]*Set{root}, sets...)

	// build tree
	for _, child := range sets {
		var parent *Set
		parentSize := n + 1
		for _, p := range all {
			if p.size > child.size && child.subsetOf(p) && p.size < parentSize {
				parent = p
				parentSize = p.size
			}
		}
		if parent == nil {
			parent = root
		}
		parent.children = append(parent.children, child)
	}

	fac = make([]int64, n+1)
	fac[0] = 1
	for i := 1; i <= n; i++ {
		fac[i] = fac[i-1] * int64(i) % MOD
	}

	ans := calc(root)
	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, ans)
	writer.Flush()
}
