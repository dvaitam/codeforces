package main

import (
   "bufio"
   "fmt"
   "os"
)

// Segment tree for range minimum query and point update
type SegTree struct {
   n int
   t []int
}

func NewSegTree(n int) *SegTree {
   N := 1
   for N < n {
       N <<= 1
   }
   t := make([]int, 2*N)
   const INF = 1<<60
   for i := range t {
       t[i] = INF
   }
   return &SegTree{n: N, t: t}
}

func (st *SegTree) Update(pos, val int) {
   i := pos + st.n
   st.t[i] = val
   for i >>= 1; i > 0; i >>= 1 {
       a, b := st.t[2*i], st.t[2*i+1]
       if a < b {
           st.t[i] = a
       } else {
           st.t[i] = b
       }
   }
}

func (st *SegTree) Query(l, r int) int {
   if l > r {
       return 1<<60
   }
   l += st.n; r += st.n
   res := 1<<60
   for l <= r {
       if l&1 == 1 {
           if st.t[l] < res {
               res = st.t[l]
           }
           l++
       }
       if r&1 == 0 {
           if st.t[r] < res {
               res = st.t[r]
           }
           r--
       }
       l >>= 1; r >>= 1
   }
   return res
}

var (
   adj [][]int
   parent, tin, tout []int
   children [][]int
   timer int
)

func dfs(u, p int) {
   parent[u] = p
   tin[u] = timer
   timer++
   for _, v := range adj[u] {
       if v == p {
           continue
       }
       children[u] = append(children[u], v)
       dfs(v, u)
   }
   tout[u] = timer - 1
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, q int
   fmt.Fscan(in, &n, &q)
   adj = make([][]int, n)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--; v--
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   parent = make([]int, n)
   tin = make([]int, n)
   tout = make([]int, n)
   children = make([][]int, n)
   timer = 0
   dfs(0, -1)

   // build segtree
   st := NewSegTree(n)
   p := make([]int, n)
   for i := 0; i < n; i++ {
       p[i] = i + 1
       st.Update(tin[i], p[i])
   }
   maxP := n

   for i := 0; i < q; i++ {
       var typ string
       fmt.Fscan(in, &typ)
       if typ == "up" {
           var v int
           fmt.Fscan(in, &v)
           v--
           maxP++
           p[v] = maxP
           st.Update(tin[v], p[v])
       } else if typ == "when" {
           var v int
           fmt.Fscan(in, &v)
           v--
           cnt := 0
           pv := p[v]
           // children
           for _, u := range children[v] {
               if st.Query(tin[u], tout[u]) < pv {
                   cnt++
               }
           }
           // parent side
           if parent[v] != -1 {
               m1 := st.Query(0, tin[v]-1)
               m2 := st.Query(tout[v]+1, n-1)
               if m1 < m2 {
                   if m1 < pv {
                       cnt++
                   }
               } else {
                   if m2 < pv {
                       cnt++
                   }
               }
           }
           fmt.Fprintln(out, cnt+1)
       } else if typ == "compare" {
           var v, u int
           fmt.Fscan(in, &v, &u)
           v--; u--
           // compute times
           calc := func(v int) int {
               cnt := 0
               pv := p[v]
               for _, u := range children[v] {
                   if st.Query(tin[u], tout[u]) < pv {
                       cnt++
                   }
               }
               if parent[v] != -1 {
                   m1 := st.Query(0, tin[v]-1)
                   m2 := st.Query(tout[v]+1, n-1)
                   mn := m1
                   if m2 < mn {
                       mn = m2
                   }
                   if mn < pv {
                       cnt++
                   }
               }
               return cnt + 1
           }
           tv := calc(v)
           tu := calc(u)
           if tv < tu {
               fmt.Fprintln(out, v+1)
           } else {
               fmt.Fprintln(out, u+1)
           }
       }
   }
}
