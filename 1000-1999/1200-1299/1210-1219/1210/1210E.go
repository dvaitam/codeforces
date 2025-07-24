package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

type GroupKey struct{ lo, hi uint64 }

type Group struct {
	bits  GroupKey
	size  int
	trans []int
}

var (
	k         int
	m         int
	perms     [][]int
	permIndex map[[5]int]int
	mul       [][]int
	groups    []Group
	groupMap  map[GroupKey]int
	idIdx     int
)

func addBit(b *GroupKey, idx int) {
	if idx < 64 {
		b.lo |= 1 << uint(idx)
	} else {
		b.hi |= 1 << uint(idx-64)
	}
}

func hasBit(b GroupKey, idx int) bool {
	if idx < 64 {
		return (b.lo>>uint(idx))&1 == 1
	}
	return (b.hi>>uint(idx-64))&1 == 1
}

func bitCount(b GroupKey) int {
	return bits.OnesCount64(b.lo) + bits.OnesCount64(b.hi)
}

func composeIndex(a, b int) int {
	key := [5]int{}
	pa := perms[a]
	pb := perms[b]
	for i := 0; i < k; i++ {
		key[i] = pa[pb[i]]
	}
	return permIndex[key]
}

func closure(bits GroupKey) GroupKey {
	queue := make([]int, 0, m)
	for i := 0; i < m; i++ {
		if hasBit(bits, i) {
			queue = append(queue, i)
		}
	}
	for idx := 0; idx < len(queue); idx++ {
		x := queue[idx]
		for j := 0; j < len(queue); j++ {
			y := queue[j]
			xy := mul[x][y]
			if !hasBit(bits, xy) {
				addBit(&bits, xy)
				queue = append(queue, xy)
			}
			yx := mul[y][x]
			if !hasBit(bits, yx) {
				addBit(&bits, yx)
				queue = append(queue, yx)
			}
		}
	}
	return bits
}

func getGroup(bits GroupKey) int {
	if id, ok := groupMap[bits]; ok {
		return id
	}
	id := len(groups)
	g := Group{bits: bits, size: bitCount(bits), trans: make([]int, m)}
	for i := 0; i < m; i++ {
		g.trans[i] = -1
	}
	groups = append(groups, g)
	groupMap[bits] = id
	return id
}

func transition(gid, perm int) int {
	g := &groups[gid]
	if g.trans[perm] != -1 {
		return g.trans[perm]
	}
	bits := g.bits
	if !hasBit(bits, perm) {
		addBit(&bits, perm)
	}
	bits = closure(bits)
	id := getGroup(bits)
	g.trans[perm] = id
	return id
}

func genPerms(cur []int, used []bool) {
	if len(cur) == k {
		p := make([]int, k)
		copy(p, cur)
		id := len(perms)
		perms = append(perms, p)
		var key [5]int
		for i := 0; i < k; i++ {
			key[i] = cur[i]
		}
		permIndex[key] = id
		allId := true
		for i := 0; i < k; i++ {
			if cur[i] != i {
				allId = false
				break
			}
		}
		if allId {
			idIdx = id
		}
		return
	}
	for i := 0; i < k; i++ {
		if !used[i] {
			used[i] = true
			cur = append(cur, i)
			genPerms(cur, used)
			cur = cur[:len(cur)-1]
			used[i] = false
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n, &k)

	permIndex = make(map[[5]int]int)
	genPerms([]int{}, make([]bool, k))
	m = len(perms)

	mul = make([][]int, m)
	for i := 0; i < m; i++ {
		mul[i] = make([]int, m)
		for j := 0; j < m; j++ {
			mul[i][j] = composeIndex(i, j)
		}
	}

	groupMap = make(map[GroupKey]int)
	var idBits GroupKey
	addBit(&idBits, idIdx)
	getGroup(idBits) // id group at index 0

	seq := make([]int, n)
	for i := 0; i < n; i++ {
		p := make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &p[j])
			p[j]--
		}
		var key [5]int
		for j := 0; j < k; j++ {
			key[j] = p[j]
		}
		seq[i] = permIndex[key]
	}

	cur := map[int]int{}
	var ans int64
	for r := 1; r <= n; r++ {
		pr := seq[r-1]
		nxt := make(map[int]int)
		g := transition(0, pr)
		nxt[g] = r
		for gid, start := range cur {
			ng := transition(gid, pr)
			if old, ok := nxt[ng]; !ok || start < old {
				nxt[ng] = start
			}
		}
		type pair struct{ g, s int }
		pairs := make([]pair, 0, len(nxt))
		for gid, start := range nxt {
			pairs = append(pairs, pair{gid, start})
		}
		sort.Slice(pairs, func(i, j int) bool { return pairs[i].s < pairs[j].s })
		for i := 0; i < len(pairs); i++ {
			start := pairs[i].s
			end := r
			if i+1 < len(pairs) {
				end = pairs[i+1].s - 1
			}
			if end >= start {
				cnt := end - start + 1
				ans += int64(groups[pairs[i].g].size) * int64(cnt)
			}
		}
		cur = nxt
	}
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, ans)
	out.Flush()
}
