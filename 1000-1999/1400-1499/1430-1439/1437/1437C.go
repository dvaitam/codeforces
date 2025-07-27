package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var q int
   if _, err := fmt.Fscan(in, &q); err != nil {
       return
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for q > 0 {
       q--
       var n int
       fmt.Fscan(in, &n)
       t := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &t[i])
       }
       sort.Ints(t)
       maxT := 2 * n
       const INF = 1_000_000_000
       // dpPrev[j]: min unpleasant for previous i-1 dishes ending at time j
       dpPrev := make([]int, maxT+1)
       for j := 1; j <= maxT; j++ {
           dpPrev[j] = INF
       }
       // dpPrev[0] = 0 implicitly
       for i := 1; i <= n; i++ {
           // prefix minimum of dpPrev
           prefix := make([]int, maxT+1)
           prefix[0] = dpPrev[0]
           for j := 1; j <= maxT; j++ {
               if dpPrev[j] < prefix[j-1] {
                   prefix[j] = dpPrev[j]
               } else {
                   prefix[j] = prefix[j-1]
               }
           }
           dpCur := make([]int, maxT+1)
           for j := 1; j <= maxT; j++ {
               // choose previous time k < j
               best := prefix[j-1]
               dpCur[j] = best + abs(j - t[i-1])
           }
           // no valid assignment ending at time 0 for i>0
           dpCur[0] = INF
           dpPrev = dpCur
       }
       // answer is min over dpPrev[j]
       ans := INF
       for j := 1; j <= 2*n; j++ {
           if dpPrev[j] < ans {
               ans = dpPrev[j]
           }
       }
       fmt.Fprintln(out, ans)
   }
}
