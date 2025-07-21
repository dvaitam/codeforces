package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const INF = 1e9

// segment tree for range minimum query and point update
type SegTree struct {
   n, size int
   tree     []int
}

func NewSegTree(n int) *SegTree {
   size := 1
   for size < n {
       size <<= 1
   }
   tree := make([]int, 2*size)
   // initialize all values to negative infinity (no rook)
   for i := range tree {
       tree[i] = -INF
   }
   return &SegTree{n: n, size: size, tree: tree}
}

func (st *SegTree) Update(pos, val int) {
   p := pos + st.size
   st.tree[p] = val
   for p >>= 1; p > 0; p >>= 1 {
       left, right := st.tree[2*p], st.tree[2*p+1]
       if left < right {
           st.tree[p] = left
       } else {
           st.tree[p] = right
       }
   }
}

// query minimum on [l, r]
func (st *SegTree) Query(l, r int) int {
   res := INF
   l += st.size
   r += st.size
   for l <= r {
       if (l & 1) == 1 {
           if st.tree[l] < res {
               res = st.tree[l]
           }
           l++
       }
       if (r & 1) == 0 {
           if st.tree[r] < res {
               res = st.tree[r]
           }
           r--
       }
       l >>= 1
       r >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, k, q int
   fmt.Fscan(reader, &n, &m, &k, &q)
   rooks := make([][2]int, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(reader, &rooks[i][0], &rooks[i][1])
       rooks[i][0]--
       rooks[i][1]--
   }
   x1 := make([]int, q)
   y1 := make([]int, q)
   x2 := make([]int, q)
   y2 := make([]int, q)
   for i := 0; i < q; i++ {
       fmt.Fscan(reader, &x1[i], &y1[i], &x2[i], &y2[i])
       x1[i]--
       y1[i]--
       x2[i]--
       y2[i]--
   }

   // Aempty: rows condition
   type RowQ struct{ y2, x1, x2, y1, idx int }
   rowQs := make([]RowQ, q)
   for i := 0; i < q; i++ {
       rowQs[i] = RowQ{y2[i], x1[i], x2[i], y1[i], i}
   }
   sort.Slice(rowQs, func(i, j int) bool {
       return rowQs[i].y2 < rowQs[j].y2
   })
   sort.Slice(rooks, func(i, j int) bool {
       return rooks[i][1] < rooks[j][1]
   })
   Aempty := make([]bool, q)
   stR := NewSegTree(n)
   ri := 0
   for _, rq := range rowQs {
       for ri < k && rooks[ri][1] <= rq.y2 {
           x := rooks[ri][0]
           y := rooks[ri][1]
           stR.Update(x, y)
           ri++
       }
       // check min last_y in rows [x1,x2]
       mn := stR.Query(rq.x1, rq.x2)
       Aempty[rq.idx] = (mn >= rq.y1)
   }

   // Bempty: columns condition
   type ColQ struct{ x2, y1, y2, x1, idx int }
   colQs := make([]ColQ, q)
   for i := 0; i < q; i++ {
       colQs[i] = ColQ{x2[i], y1[i], y2[i], x1[i], i}
   }
   sort.Slice(colQs, func(i, j int) bool {
       return colQs[i].x2 < colQs[j].x2
   })
   // resort rooks by x
   sort.Slice(rooks, func(i, j int) bool {
       return rooks[i][0] < rooks[j][0]
   })
   Bempty := make([]bool, q)
   stC := NewSegTree(m)
   ci := 0
   for _, cq := range colQs {
       for ci < k && rooks[ci][0] <= cq.x2 {
           x := rooks[ci][0]
           y := rooks[ci][1]
           stC.Update(y, x)
           ci++
       }
       mn := stC.Query(cq.y1, cq.y2)
       Bempty[cq.idx] = (mn >= cq.x1)
   }

   for i := 0; i < q; i++ {
       if Aempty[i] || Bempty[i] {
           fmt.Fprintln(writer, "YES")
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}
