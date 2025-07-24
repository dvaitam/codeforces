package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Fenwick tree for int counts
type Fenwick struct {
   n    int
   tree []int
}

func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, tree: make([]int, n+1)}
}

func (f *Fenwick) Update(i, v int) {
   for x := i; x <= f.n; x += x & -x {
       f.tree[x] += v
   }
}

// Query prefix sum [1..i]
func (f *Fenwick) Query(i int) int {
   s := 0
   for x := i; x > 0; x -= x & -x {
       s += f.tree[x]
   }
   return s
}

// Range sum [l..r]
func (f *Fenwick) Range(l, r int) int {
   if l > r {
       return 0
   }
   return f.Query(r) - f.Query(l-1)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   adj := make([][]struct{ to int; w int64 }, n+1)
   for i := 2; i <= n; i++ {
       var p int
       var w int64
       fmt.Fscan(in, &p, &w)
       adj[p] = append(adj[p], struct{ to int; w int64 }{i, w})
   }
   // DFS to compute tin, tout, depth sums
   tin := make([]int, n+1)
   tout := make([]int, n+1)
   depth := make([]int64, n+1)
   time := 1
   type stackEntry struct{ v, parent, idx int }
   stack := make([]stackEntry, 0, n)
   stack = append(stack, stackEntry{v: 1, parent: 0, idx: 0})
   for len(stack) > 0 {
       top := &stack[len(stack)-1]
       v := top.v
       if top.idx == 0 {
           tin[v] = time
           time++
       }
       if top.idx < len(adj[v]) {
           e := adj[v][top.idx]
           top.idx++
           if e.to == top.parent {
               continue
           }
           depth[e.to] = depth[v] + e.w
           stack = append(stack, stackEntry{v: e.to, parent: v, idx: 0})
       } else {
           tout[v] = time - 1
           stack = stack[:len(stack)-1]
       }
   }
   // Prepare nodes sorted by depth
   type nodeInfo struct{ sum int64; tin int }
   nodes := make([]nodeInfo, n)
   for v := 1; v <= n; v++ {
       nodes[v-1] = nodeInfo{sum: depth[v], tin: tin[v]}
   }
   sort.Slice(nodes, func(i, j int) bool {
       return nodes[i].sum < nodes[j].sum
   })
   // Prepare queries
   type query struct{ threshold int64; v int }
   qs := make([]query, n)
   for v := 1; v <= n; v++ {
       qs[v-1] = query{threshold: depth[v] + a[v], v: v}
   }
   sort.Slice(qs, func(i, j int) bool {
       return qs[i].threshold < qs[j].threshold
   })
   // Process queries
   bit := NewFenwick(n)
   ans := make([]int, n+1)
   j := 0
   for _, q := range qs {
       for j < n && nodes[j].sum <= q.threshold {
           bit.Update(nodes[j].tin, 1)
           j++
       }
       v := q.v
       cnt := bit.Range(tin[v], tout[v])
       // exclude v itself
       if cnt > 0 {
           cnt--
       }
       ans[v] = cnt
   }
   // Output
   for i := 1; i <= n; i++ {
       if i > 1 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, ans[i])
   }
   out.WriteByte('\n')
}
