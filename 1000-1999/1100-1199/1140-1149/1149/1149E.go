package main

import (
   "bufio"
   "fmt"
   "os"
   "sync"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   h := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &h[i])
   }
   adjRev := make([][]int, n)
   outdeg := make([]int, n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--
       v--
       // edge u->v
       adjRev[v] = append(adjRev[v], u)
       outdeg[u]++
   }
   // dp: longest path to leaf
   dp := make([]int, n)
   // queue of nodes with outdeg 0
   q := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if outdeg[i] == 0 {
           q = append(q, i)
       }
   }
   // process
   for qi := 0; qi < len(q); qi++ {
       u := q[qi]
       for _, p := range adjRev[u] {
           if dp[p] < dp[u]+1 {
               dp[p] = dp[u] + 1
           }
           outdeg[p]--
           if outdeg[p] == 0 {
               q = append(q, p)
           }
       }
   }
   // compute xor over nodes at even distance parity (dp%2==0)
   var x int64
   par0 := make([]bool, n)
   for i := 0; i < n; i++ {
       if dp[i]&1 == 0 {
           par0[i] = true
           x ^= h[i]
       }
   }
   if x == 0 {
       fmt.Fprintln(out, "LOSE")
       return
   }
   // winning move exists
   fmt.Fprintln(out, "WIN")
   newH := make([]int64, n)
   copy(newH, h)
   for i := 0; i < n; i++ {
       if par0[i] {
           want := h[i] ^ x
           if want < h[i] {
               newH[i] = want
               // children can be left unchanged
               break
           }
       }
   }
   // output newH
   for i, v := range newH {
       if i > 0 {
           out.WriteString(" ")
       }
       fmt.Fprint(out, v)
   }
   out.WriteByte('\n')
}
// ensure stack underflow does not occur
var _ = sync.Once{}
