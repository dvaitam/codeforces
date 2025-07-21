package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // directed edge map: from->to = cost to reverse
   dirCost := make(map[[2]int]int)
   // undirected adjacency
   adj := make(map[int][]int)
   for i := 0; i < n; i++ {
       var a, b, c int
       fmt.Fscan(reader, &a, &b, &c)
       dirCost[[2]int{a, b}] = c
       // build undirected graph
       adj[a] = append(adj[a], b)
       adj[b] = append(adj[b], a)
   }
   // reconstruct cycle order
   order := make([]int, 0, n)
   // start from any node, say the first a from input, but use 1
   // find a node present in adj
   start := 0
   for u := range adj {
       start = u
       break
   }
   cur := start
   prev := -1
   for len(order) < n {
       order = append(order, cur)
       // next neighbor not equal prev
       for _, nei := range adj[cur] {
           if nei != prev {
               prev, cur = cur, nei
               break
           }
       }
   }
   // compute cost for two orientations
   costCW := 0
   costCCW := 0
   for i := 0; i < n; i++ {
       u := order[i]
       v := order[(i+1)%n]
       // cost to orient u->v
       if _, ok := dirCost[[2]int{u, v}]; ok {
           // already u->v, no cost for CW
       } else if c, ok2 := dirCost[[2]int{v, u}]; ok2 {
           costCW += c
       }
       // cost to orient v->u (CCW)
       if _, ok := dirCost[[2]int{v, u}]; ok {
           // already v->u, no cost for CCW
       } else if c, ok2 := dirCost[[2]int{u, v}]; ok2 {
           costCCW += c
       }
   }
   if costCW < costCCW {
       fmt.Println(costCW)
   } else {
       fmt.Println(costCCW)
   }
}
