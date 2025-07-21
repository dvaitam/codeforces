package main

import (
   "bufio"
   "fmt"
   "os"
)

type Edge struct {
   to    int
   wF    int // cost to orient from parent->child
   wB    int // cost to orient from child->parent
}

var (
   n    int
   adj  [][]Edge
   cost []int
)

// initial DFS to compute cost[0]
func dfs1(u, p int) int {
   sum := 0
   for _, e := range adj[u] {
       v := e.to
       if v == p {
           continue
       }
       sum += e.wF
       sum += dfs1(v, u)
   }
   return sum
}

// reroot DP to compute cost for all nodes
func dfs2(u, p int) {
   for _, e := range adj[u] {
       v := e.to
       if v == p {
           continue
       }
       // move root from u to v: remove edge u->v, add v->u
       cost[v] = cost[u] - e.wF + e.wB
       dfs2(v, u)
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n)
   adj = make([][]Edge, n)
   for i := 0; i < n-1; i++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       a--
       b--
       // original a->b
       adj[a] = append(adj[a], Edge{to: b, wF: 0, wB: 1})
       adj[b] = append(adj[b], Edge{to: a, wF: 1, wB: 0})
   }
   if n == 0 {
       fmt.Fprintln(out, 0)
       return
   }
   cost = make([]int, n)
   cost[0] = dfs1(0, -1)
   dfs2(0, -1)
   // prepare arrays for delta DFS
   prefix := make([]int, n+1)
   var answer = cost[0] // initialize
   // for each source candidate s1
   var dfsDelta func(u, p, depth, s1 int, best *int)
   dfsDelta = func(u, p, depth, s1 int, best *int) {
       // compute delta for u
       d := depth
       mid := (d + 1) / 2
       delta := prefix[d] - prefix[mid]
       cur := cost[s1] + delta
       if cur < *best {
           *best = cur
       }
       for _, e := range adj[u] {
           v := e.to
           if v == p {
               continue
           }
           // diff for this edge in path from u->v
           diff := -1
           if e.wF == 0 {
               diff = 1
           }
           prefix[d+1] = prefix[d] + diff
           dfsDelta(v, u, d+1, s1, best)
       }
   }
   // find global minimal
   for s1 := 0; s1 < n; s1++ {
       best := cost[s1] // delta >=0 at root yields cost[s1]
       // init prefix
       prefix[0] = 0
       dfsDelta(s1, -1, 0, s1, &best)
       if best < answer {
           answer = best
       }
   }
   fmt.Fprintln(out, answer)
}
