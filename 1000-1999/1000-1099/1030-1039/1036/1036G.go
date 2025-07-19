package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   indegree := make([]int, n+1)
   outdegree := make([]int, n+1)
   // reverse graph: for each edge u->v, store u in preds[v]
   preds := make([][]int, n+1)
   var u, v int
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &u, &v)
       indegree[v]++
       outdegree[u]++
       preds[v] = append(preds[v], u)
   }
   // assign source and sink indices
   inIndex := make([]int, n+1)
   outIndex := make([]int, n+1)
   ii, oo := 0, 0
   for i := 1; i <= n; i++ {
       if indegree[i] == 0 {
           ii++
           inIndex[i] = ii
       }
       if outdegree[i] == 0 {
           oo++
           outIndex[i] = oo
       }
   }
   // DP values
   val := make([]int, n+1)
   // copy outdegree for processing
   remOut := make([]int, n+1)
   copy(remOut, outdegree)
   // hiss stores bitmask for sources
   // size ii+1 to access [1..ii]
   hiss := make([]int, ii+1)
   // queue for reverse topological order
   queue := make([]int, 0, n)
   // initialize sinks
   for i := 1; i <= n; i++ {
       if remOut[i] == 0 {
           // sink
           val[i] = 1 << outIndex[i]
           queue = append(queue, i)
           if indegree[i] == 0 {
               // source that is also sink
               hiss[inIndex[i]] = val[i]
           }
       }
   }
   // process in reverse topological order
   for qi := 0; qi < len(queue); qi++ {
       u = queue[qi]
       for _, p := range preds[u] {
           // propagate masks
           val[p] |= val[u]
           remOut[p]--
           if remOut[p] == 0 {
               queue = append(queue, p)
               if indegree[p] == 0 {
                   hiss[inIndex[p]] = val[p]
               }
           }
       }
   }
   // shift hiss masks: hiss[i] = hiss[i+1] >> 1
   for i := 0; i < ii; i++ {
       hiss[i] = hiss[i+1] >> 1
   }
   // check Hall's condition complements
   full := 1 << ii
   for mask := 1; mask < full; mask++ {
       if mask == full-1 {
           fmt.Println("YES")
           return
       }
       reach := 0
       // union of reachable sinks
       for j := 0; j < ii; j++ {
           if mask&(1<<j) != 0 {
               reach |= hiss[j]
           }
       }
       if bits.OnesCount(uint(reach)) == bits.OnesCount(uint(mask)) {
           fmt.Println("NO")
           return
       }
   }
   // if none failed, it's always possible
   fmt.Println("YES")
}
