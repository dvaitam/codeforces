package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // positions of each value
   pos := make([][]int, n+1)
   for i, v := range a {
       if v >= 1 && v <= n {
           pos[v] = append(pos[v], i)
       }
   }
   // compute maximum frequency
   maxf := 0
   for v := 1; v <= n; v++ {
       if len(pos[v]) > maxf {
           maxf = len(pos[v])
       }
   }
   if maxf < 2 {
       fmt.Fprintln(out, 0)
       return
   }
   cnt := make([]int, maxf+2)
   L := make([]int, maxf+2)
   R := make([]int, maxf+2)
   A := make([]int, maxf+2)
   for k := 1; k <= maxf; k++ {
       L[k] = n
       R[k] = -1
       A[k] = n
   }
   // fill arrays
   for v := 1; v <= n; v++ {
       pv := pos[v]
       m := len(pv)
       if m == 0 {
           continue
       }
       first := pv[0]
       // cnt, L, R
       for k := 1; k <= m; k++ {
           cnt[k]++
           if first < L[k] {
               L[k] = first
           }
           idx := pv[k-1]
           if idx > R[k] {
               R[k] = idx
           }
       }
       // next occurrence for A
       for k := 1; k < m; k++ {
           nxt := pv[k]
           if nxt < A[k] {
               A[k] = nxt
           }
       }
   }
   // compute answer
   ans := 0
   for k := 1; k <= maxf; k++ {
       if cnt[k] < 2 {
           continue
       }
       // ensure no value occurs >k in subarray [L[k],R[k]]
       if A[k] <= R[k] {
           continue
       }
       length := R[k] - L[k] + 1
       if length > ans {
           ans = length
       }
   }
   fmt.Fprintln(out, ans)
}
