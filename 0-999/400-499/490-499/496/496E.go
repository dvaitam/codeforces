package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type part struct {
   a, b int64
   idx  int
}
type actor struct {
   c, d, k int64
   idx      int
}

// Segment tree for range sum and find-first query
type SegTree struct {
   n    int
   tree []int64
}

func NewSegTree(n int) *SegTree {
   size := 1
   for size < n {
       size <<= 1
   }
   return &SegTree{n: size, tree: make([]int64, 2*size)}
}

// add delta at position i
func (st *SegTree) Update(i int, delta int64) {
   i += st.n
   st.tree[i] += delta
   for i >>= 1; i > 0; i >>= 1 {
       st.tree[i] = st.tree[2*i] + st.tree[2*i+1]
   }
}

// find first position >= l with positive sum, or -1
func (st *SegTree) FindFirst(l int) int {
   return st.findFirst(1, 0, st.n-1, l)
}

func (st *SegTree) findFirst(node, l, r, ql int) int {
   if r < ql || st.tree[node] == 0 {
       return -1
   }
   if l == r {
       return l
   }
   mid := (l + r) >> 1
   if ql <= mid {
       if res := st.findFirst(2*node, l, mid, ql); res != -1 {
           return res
       }
       return st.findFirst(2*node+1, mid+1, r, ql)
   }
   return st.findFirst(2*node+1, mid+1, r, ql)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   parts := make([]part, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &parts[i].a, &parts[i].b)
       parts[i].idx = i
   }
   var m int
   fmt.Fscan(in, &m)
   actors := make([]actor, m)
   diVals := make([]int64, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &actors[i].c, &actors[i].d, &actors[i].k)
       actors[i].idx = i + 1
       diVals[i] = actors[i].d
   }
   // compress di values
   sort.Slice(diVals, func(i, j int) bool { return diVals[i] < diVals[j] })
   diVals = uniqueInt64(diVals)
   // sort parts by a
   sort.Slice(parts, func(i, j int) bool { return parts[i].a < parts[j].a })
   // sort actors by c
   sort.Slice(actors, func(i, j int) bool { return actors[i].c < actors[j].c })
   // prepare segment tree and actor lists per di index
   st := NewSegTree(len(diVals))
   type actorEntry struct { idx int; k int64 }
   actorLists := make([][]*actorEntry, len(diVals))
   ans := make([]int, n)
   ai := 0
   for _, p := range parts {
       // insert actors with c <= a
       for ai < m && actors[ai].c <= p.a {
           d := actors[ai].d
           pos := sort.Search(len(diVals), func(i int) bool { return diVals[i] >= d })
           // add capacity
           st.Update(pos, actors[ai].k)
           // add to list
           entry := &actorEntry{idx: actors[ai].idx, k: actors[ai].k}
           actorLists[pos] = append(actorLists[pos], entry)
           ai++
       }
       // find actor with d >= b
       j := sort.Search(len(diVals), func(i int) bool { return diVals[i] >= p.b })
       pos := st.FindFirst(j)
       if pos == -1 {
           fmt.Println("NO")
           return
       }
       // pick first actor at pos
       e := actorLists[pos][0]
       ans[p.idx] = e.idx
       // decrement capacity
       st.Update(pos, -1)
       e.k--
       if e.k == 0 {
           actorLists[pos] = actorLists[pos][1:]
       }
   }
   // success
   fmt.Println("YES")
   // print ans in original order
   for i, v := range ans {
       if i > 0 {
           fmt.Print(" ")
       }
       fmt.Print(v)
   }
   fmt.Println()
}

func uniqueInt64(a []int64) []int64 {
   j := 0
   for i := 0; i < len(a); i++ {
       if j == 0 || a[i] != a[j-1] {
           a[j] = a[i]
           j++
       }
   }
   return a[:j]
}
