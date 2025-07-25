package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// This program computes E_max for the first corridor: the maximum weight e (<=1e9)
// such that the edge between u1 and v1 may appear in some MST when its weight is set to e.
// This equals the minimax path value between u1 and v1 in the graph without that edge.
func main() {
   reader := bufio.NewReader(os.Stdin)
   // read n, m
   n := readInt(reader)
   m := readInt(reader)
   // read first edge
   u1 := readInt(reader)
   v1 := readInt(reader)
   _ = readInt(reader) // original weight, unused
   // read remaining edges
   edges := make([]edge, 0, m-1)
   for i := 1; i < m; i++ {
       u := readInt(reader)
       v := readInt(reader)
       w := readInt(reader)
       edges = append(edges, edge{u: u, v: v, w: w})
   }
   // sort edges by weight ascending
   sort.Slice(edges, func(i, j int) bool {
       return edges[i].w < edges[j].w
   })
   // DSU initialization
   parent := make([]int, n+1)
   for i := 1; i <= n; i++ {
       parent[i] = i
   }
   // find with path halving
   var find func(int) int
   find = func(x int) int {
       for parent[x] != x {
           parent[x] = parent[parent[x]]
           x = parent[x]
       }
       return x
   }
   // union sets
   unite := func(a, b int) {
       a = find(a)
       b = find(b)
       if a != b {
           parent[b] = a
       }
   }
   // Kruskal until u1 and v1 become connected
   const INF = 1000000001
   ans := INF
   for _, e := range edges {
       ua := find(e.u)
       ub := find(e.v)
       if ua != ub {
           unite(ua, ub)
           if find(u1) == find(v1) {
               ans = e.w
               break
           }
       }
   }
   // if never connected, ans remains INF
   if ans > 1000000000 {
       ans = 1000000000
   }
   fmt.Println(ans)
}

// edge represents an undirected edge with weight w
type edge struct {
   u, v, w int
}

// readInt reads next integer from buffered reader
func readInt(r *bufio.Reader) int {
   neg := false
   c, err := r.ReadByte()
   if err != nil {
       return 0
   }
   // skip non-numeric
   for (c < '0' || c > '9') && c != '-' {
       c, err = r.ReadByte()
       if err != nil {
           return 0
       }
   }
   if c == '-' {
       neg = true
       c, _ = r.ReadByte()
   }
   x := 0
   for c >= '0' && c <= '9' {
       x = x*10 + int(c-'0')
       c, err = r.ReadByte()
       if err != nil {
           break
       }
   }
   if neg {
       return -x
   }
   return x
}
