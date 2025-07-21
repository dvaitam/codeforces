package main

import (
   "bufio"
   "fmt"
   "os"
)

// Edge holds endpoints and flow weight
type Edge struct {
   u, v int
   w    int64
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   edges := make([]Edge, m)
   adj := make([][]int, n+1)
   sumFlow := make([]int64, n+1)
   for i := 0; i < m; i++ {
       var a, b int
       var c int64
       fmt.Fscan(in, &a, &b, &c)
       edges[i] = Edge{u: a, v: b, w: c}
       adj[a] = append(adj[a], i)
       adj[b] = append(adj[b], i)
       sumFlow[a] += c
       sumFlow[b] += c
   }
   needOut := make([]int64, n+1)
   for v := 1; v <= n; v++ {
       switch v {
       case 1:
           needOut[v] = sumFlow[v]
       case n:
           needOut[v] = 0
       default:
           needOut[v] = sumFlow[v] / 2
       }
   }
   outSum := make([]int64, n+1)
   processed := make([]bool, m)
   dir := make([]int, m)
   // queue of saturated nodes
   queue := make([]int, 0, n)
   for v := 1; v <= n; v++ {
       if needOut[v] == 0 {
           queue = append(queue, v)
       }
   }
   head := 0
   for head < len(queue) {
       v := queue[head]
       head++
       for _, ei := range adj[v] {
           if processed[ei] {
               continue
           }
           processed[ei] = true
           e := edges[ei]
           // orient edge toward v: u -> v
           var u int
           if e.u == v {
               u = e.v
               // original a->b is e.u->e.v; now e.v->e.u, reversed
               dir[ei] = 1
           } else {
               u = e.u
               // original a->b is e.u->e.v; now e.u->e.v, same
               dir[ei] = 0
           }
           outSum[u] += e.w
           if outSum[u] == needOut[u] {
               queue = append(queue, u)
           }
       }
   }
   // output directions
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < m; i++ {
       // 0 means a->b, 1 means b->a
       fmt.Fprintln(writer, dir[i])
   }
}
