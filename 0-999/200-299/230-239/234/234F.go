package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, a, b int
   fmt.Fscan(reader, &n)
   fmt.Fscan(reader, &a, &b)
   h := make([]int, n)
   total := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &h[i])
       total += h[i]
   }
   // impossible if total area exceeds available paint
   if total > a+b {
       fmt.Fprintln(writer, -1)
       return
   }
   const INF = 1000000000
   // dp arrays: prev0[r] = min cost ending with red using r red area
   prev0 := make([]int, a+1)
   prev1 := make([]int, a+1)
   for i := 0; i <= a; i++ {
       prev0[i] = INF
       prev1[i] = INF
   }
   // first board
   if h[0] <= a {
       prev0[h[0]] = 0
   }
   prev1[0] = 0

   // DP over boards
   for i := 1; i < n; i++ {
       hi := h[i]
       hip := h[i-1]
       // next dp
       next0 := make([]int, a+1)
       next1 := make([]int, a+1)
       for j := 0; j <= a; j++ {
           next0[j] = INF
           next1[j] = INF
       }
       // transitions
       cost := min(hip, hi)
       for r := 0; r <= a; r++ {
           v0 := prev0[r]
           if v0 < INF {
               // paint red
               nr := r + hi
               if nr <= a && v0 < next0[nr] {
                   next0[nr] = v0
               }
               // paint green
               if v0+cost < next1[r] {
                   next1[r] = v0 + cost
               }
           }
           v1 := prev1[r]
           if v1 < INF {
               // paint red
               nr := r + hi
               if nr <= a && v1+cost < next0[nr] {
                   next0[nr] = v1 + cost
               }
               // paint green
               if v1 < next1[r] {
                   next1[r] = v1
               }
           }
       }
       prev0, prev1 = next0, next1
   }
   // compute answer
   ans := INF
   // red used r must satisfy green = total - r <= b => r >= total-b
   start := 0
   if total > b {
       start = total - b
   }
   for r := start; r <= a; r++ {
       if prev0[r] < ans {
           ans = prev0[r]
       }
       if prev1[r] < ans {
           ans = prev1[r]
       }
   }
   if ans == INF {
       fmt.Fprintln(writer, -1)
   } else {
       fmt.Fprintln(writer, ans)
   }
}
