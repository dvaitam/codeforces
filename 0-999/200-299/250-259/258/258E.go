package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m int
   adj  [][]int
   tin, tout []int
   rev  []int
   timer int
   l1, r1, l2, r2 []int
   cntOp []int
   events [][]OpEvent
   ans   []int
)

// OpEvent marks start(+1) or end(-1) of operation i
type OpEvent struct { i, d int }

// Segment tree for union-length of intervals with dynamic add/remove
type SegTree struct {
   n int
   cov []int
   len []int
}

func NewSegTree(n int) *SegTree {
   size := 4 * n
   return &SegTree{n: n, cov: make([]int, size), len: make([]int, size)}
}

func (st *SegTree) update(L, R, v, l, r, idx int) {
   if R < l || r < L {
       return
   }
   if L <= l && r <= R {
       st.cov[idx] += v
   } else {
       mid := (l + r) >> 1
       st.update(L, R, v, l, mid, idx<<1)
       st.update(L, R, v, mid+1, r, idx<<1|1)
   }
   if st.cov[idx] > 0 {
       st.len[idx] = r - l + 1
   } else {
       if l == r {
           st.len[idx] = 0
       } else {
           st.len[idx] = st.len[idx<<1] + st.len[idx<<1|1]
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // read n, m
   fmt.Fscan(reader, &n, &m)
   adj = make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   tin = make([]int, n+1)
   tout = make([]int, n+1)
   rev = make([]int, n+1)
   // iterative DFS for tin/tout
   type Item struct{ u, p, idx int }
   stack := make([]Item, 0, n*2)
   stack = append(stack, Item{u:1, p:0, idx:-1})
   for len(stack) > 0 {
       it := &stack[len(stack)-1]
       if it.idx < 0 {
           timer++
           tin[it.u] = timer
           rev[timer] = it.u
           it.idx = 0
       } else if it.idx < len(adj[it.u]) {
           v := adj[it.u][it.idx]
           it.idx++
           if v == it.p {
               continue
           }
           stack = append(stack, Item{u: v, p: it.u, idx: -1})
       } else {
           tout[it.u] = timer
           stack = stack[:len(stack)-1]
       }
   }
   // read operations
   l1 = make([]int, m)
   r1 = make([]int, m)
   l2 = make([]int, m)
   r2 = make([]int, m)
   events = make([][]OpEvent, n+2)
   cntOp = make([]int, m)
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       l1[i], r1[i] = tin[a], tout[a]
       l2[i], r2[i] = tin[b], tout[b]
       // event start/end for op i
       events[l1[i]] = append(events[l1[i]], OpEvent{i, +1})
       events[r1[i]+1] = append(events[r1[i]+1], OpEvent{i, -1})
       events[l2[i]] = append(events[l2[i]], OpEvent{i, +1})
       events[r2[i]+1] = append(events[r2[i]+1], OpEvent{i, -1})
   }
   // sweep
   st := NewSegTree(n)
   ans = make([]int, n+1)
   for p := 1; p <= n; p++ {
       for _, e := range events[p] {
           i, d := e.i, e.d
           old := cntOp[i]
           cntOp[i] = old + d
           if old == 0 && cntOp[i] == 1 {
               st.update(l1[i], r1[i], +1, 1, n, 1)
               st.update(l2[i], r2[i], +1, 1, n, 1)
           } else if old == 1 && cntOp[i] == 0 {
               st.update(l1[i], r1[i], -1, 1, n, 1)
               st.update(l2[i], r2[i], -1, 1, n, 1)
           }
       }
       u := rev[p]
       total := st.len[1]
       if total > 0 {
           ans[u] = total - 1
       } else {
           ans[u] = 0
       }
   }
   // output
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, ans[i])
   }
   writer.WriteByte('\n')
}
