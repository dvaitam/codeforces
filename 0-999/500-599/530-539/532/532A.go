package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   hi := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &hi[i])
   }
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       adj[a] = append(adj[a], b)
       adj[b] = append(adj[b], a)
   }
   var k int
   fmt.Fscan(in, &k)
   s := make([]int64, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(in, &s[i])
   }
   // DFS to compute mi, minNode, pmin, cnt
   const INF = int64(1e18)
   mi := make([]int64, n+1)
   pmin := make([]int64, n+1)
   cnt := make([]int, n+1)
   type stkElem struct{ u, p int; pminIn int64; minNode int }
   stk := make([]stkElem, 0, n)
   stk = append(stk, stkElem{1, 0, INF, 1})
   for len(stk) > 0 {
       e := stk[len(stk)-1]
       stk = stk[:len(stk)-1]
       u, p := e.u, e.p
       pminIn, mn := e.pminIn, e.minNode
       // compute mi[u]
       if hi[u] < pminIn {
           mi[u] = hi[u]
           mn = u
       } else {
           mi[u] = pminIn
       }
       pmin[u] = pminIn
       cnt[mn]++
       // traverse children
       for _, v := range adj[u] {
           if v == p {
               continue
           }
           stk = append(stk, stkElem{v, u, mi[u], mn})
       }
   }
   // Build capacities
   M := make([]int64, n)
   for i := 1; i <= n; i++ {
       M[i-1] = mi[i]
   }
   sort.Slice(M, func(i, j int) bool { return M[i] > M[j] })
   sort.Slice(s, func(i, j int) bool { return s[i] > s[j] })
   // initial match
   possible := true
   if k <= n {
       for i := 0; i < k; i++ {
           if s[i] > M[i] {
               possible = false
               break
           }
       }
   } else {
       possible = false
   }
   if possible {
       fmt.Fprintln(out, 0)
       return
   }
   // compute deficits
   D := 0
   var T_max int64
   // two pointers
   // uniq s in descending order
   uniq := make([]int64, 0, k)
   for i, v := range s {
       if i == 0 || v != s[i-1] {
           uniq = append(uniq, v)
       }
   }
   i_s, i_M := 0, 0
   for _, t := range uniq {
       for i_s < k && s[i_s] >= t {
           i_s++
       }
       for i_M < n && M[i_M] >= t {
           i_M++
       }
       d := i_s - i_M
       if d > D {
           D = d
       }
       if d > 0 && t > T_max {
           T_max = t
       }
   }
   if D <= 0 {
       fmt.Fprintln(out, 0)
       return
   }
   // find minimal x
   ans := INF
   for v := 1; v <= n; v++ {
       h := hi[v]
       if h >= T_max {
           continue
       }
       if pmin[v] < T_max {
           continue
       }
       if cnt[v] < D {
           continue
       }
       x := T_max - h
       if int64(x) < ans {
           ans = int64(x)
       }
   }
   if ans == INF {
       fmt.Fprintln(out, -1)
   } else {
       fmt.Fprintln(out, ans)
   }
}
