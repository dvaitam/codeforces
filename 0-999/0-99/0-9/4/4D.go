package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Item struct {
   w, h, id int
}

type Node struct {
   length, id int
}

// Segment tree for range max query of Node
type SegTree struct {
   n    int
   tree []Node
}

func NewSegTree(n int) *SegTree {
   size := 4 * n
   return &SegTree{n: n, tree: make([]Node, size)}
}

// update position pos with node if node.length is greater
func (st *SegTree) update(pos int, node Node, idx, l, r int) {
   if l == r {
       if node.length > st.tree[idx].length {
           st.tree[idx] = node
       }
       return
   }
   mid := (l + r) >> 1
   if pos <= mid {
       st.update(pos, node, idx<<1, l, mid)
   } else {
       st.update(pos, node, idx<<1|1, mid+1, r)
   }
   // pull up
   left := st.tree[idx<<1]
   right := st.tree[idx<<1|1]
   if left.length >= right.length {
       st.tree[idx] = left
   } else {
       st.tree[idx] = right
   }
}

// query max Node in range [ql, qr]
func (st *SegTree) query(ql, qr, idx, l, r int) Node {
   if ql > r || qr < l {
       return Node{length: -1, id: 0}
   }
   if ql <= l && r <= qr {
       return st.tree[idx]
   }
   mid := (l + r) >> 1
   left := st.query(ql, qr, idx<<1, l, mid)
   right := st.query(ql, qr, idx<<1|1, mid+1, r)
   if left.length >= right.length {
       return left
   }
   return right
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, w0, h0 int
   fmt.Fscan(reader, &n, &w0, &h0)
   items := make([]Item, n+1)
   widths := make([]int, 0, n+1)
   widths = append(widths, w0)
   for i := 0; i < n; i++ {
       var w, h int
       fmt.Fscan(reader, &w, &h)
       items[i] = Item{w: w, h: h, id: i + 1}
       widths = append(widths, w)
   }
   // special item
   items[n] = Item{w: w0, h: h0, id: n + 1}
   // compress widths
   sort.Ints(widths)
   uniqW := widths[:0]
   for i, x := range widths {
       if i == 0 || x != widths[i-1] {
           uniqW = append(uniqW, x)
       }
   }
   m := len(uniqW)
   // sort items by descending h
   sort.Slice(items, func(i, j int) bool {
       return items[i].h > items[j].h
   })
   // parent pointers
   parent := make([]int, n+2)
   st := NewSegTree(m)
   type Buf struct{ fid, length, id int }
   buf := make([]Buf, 0, n+1)
   // helper to get compressed id
   getId := func(x int) int {
       // 1-based
       pos := sort.SearchInts(uniqW, x)
       return pos + 1
   }

   // process items
   for i, it := range items {
       fid := getId(it.w)
       // flush when height group changes
       if i > 0 && it.h != items[i-1].h {
           for _, b := range buf {
               st.update(b.fid, Node{length: b.length, id: b.id}, 1, 1, m)
           }
           buf = buf[:0]
       }
       // query range fid+1..m
       var best Node
       if fid < m {
           best = st.query(fid+1, m, 1, 1, m)
       } else {
           best = Node{length: -1, id: 0}
       }
       length := best.length + 1
       parent[it.id] = best.id
       if it.id == n+1 {
           // reached special
           fmt.Fprintln(writer, length)
           // backtrack
           tp := it.id
           for parent[tp] != 0 {
               fmt.Fprintf(writer, "%d ", parent[tp])
               tp = parent[tp]
           }
           fmt.Fprintln(writer)
           return
       }
       buf = append(buf, Buf{fid: fid, length: length, id: it.id})
   }
}
