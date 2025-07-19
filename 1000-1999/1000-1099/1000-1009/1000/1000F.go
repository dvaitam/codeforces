package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const INF = int64(1) << 60

// Query holds a range query with original index
type Query struct {
   l, r, i int
}

// Node stores a value and its position in the segment tree
type Node struct {
   v   int64
   idx int
}

// SegTree supports range minimum query with index
type SegTree struct {
   n, size int
   data     []Node
}

// NewSegTree builds a segment tree for n elements, initializing all values to INF
func NewSegTree(n int) *SegTree {
   size := 1
   for size < n {
       size <<= 1
   }
   data := make([]Node, 2*size)
   // initialize leaves
   for i := 0; i < size; i++ {
       data[size+i] = Node{v: INF, idx: i}
   }
   // build internal nodes
   for i := size - 1; i > 0; i-- {
       data[i] = minNode(data[2*i], data[2*i+1])
   }
   return &SegTree{n: n, size: size, data: data}
}

func minNode(a, b Node) Node {
   if a.v < b.v || (a.v == b.v && a.idx < b.idx) {
       return a
   }
   return b
}

// Update sets position pos to value v
func (st *SegTree) Update(pos int, v int64) {
   i := pos + st.size
   st.data[i].v = v
   // propagate up
   for i >>= 1; i > 0; i >>= 1 {
       st.data[i] = minNode(st.data[2*i], st.data[2*i+1])
   }
}

// Query returns the node with minimum v in [l, r)
func (st *SegTree) Query(l, r int) Node {
   var resl = Node{v: INF, idx: -1}
   var resr = Node{v: INF, idx: -1}
   l += st.size
   r += st.size
   for l < r {
       if l&1 == 1 {
           resl = minNode(resl, st.data[l])
           l++
       }
       if r&1 == 1 {
           r--
           resr = minNode(st.data[r], resr)
       }
       l >>= 1
       r >>= 1
   }
   return minNode(resl, resr)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var N int
   fmt.Fscan(reader, &N)
   A := make([]int, N)
   for i := 0; i < N; i++ {
       fmt.Fscan(reader, &A[i])
   }
   var Q int
   fmt.Fscan(reader, &Q)
   queries := make([]Query, Q)
   for i := 0; i < Q; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       l-- // zero-based inclusive
       queries[i] = Query{l: l, r: r, i: i}
   }
   sort.Slice(queries, func(i, j int) bool {
       return queries[i].r < queries[j].r
   })
   ans := make([]int, Q)
   lastPos := make(map[int]int)

   seg := NewSegTree(N)
   qi := 0
   // iterate over positions up to N to process queries at r == i
   for i := 0; i <= N; i++ {
       // process all queries with r == i
       for qi < Q && queries[qi].r <= i {
           q := queries[qi]
           node := seg.Query(q.l, q.r)
           if node.v < int64(q.l) {
               ans[q.i] = A[node.idx]
           }
           qi++
       }
       if i == N {
           break
       }
       val := A[i]
       prev, ok := lastPos[val]
       if !ok {
           prev = -1
       }
       if prev >= 0 {
           seg.Update(prev, INF)
       }
       seg.Update(i, int64(prev))
       lastPos[val] = i
   }
   for i := 0; i < Q; i++ {
       fmt.Fprintln(writer, ans[i])
   }
}
