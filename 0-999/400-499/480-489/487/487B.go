package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, l int
   var s int64
   if _, err := fmt.Fscan(reader, &n, &s, &l); err != nil {
       return
   }
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   pos := make([]int, n+1)
   maxD := make([]int, 0, n)
   minD := make([]int, 0, n)
   left := 1
   for i := 1; i <= n; i++ {
       for len(maxD) > 0 && a[maxD[len(maxD)-1]] <= a[i] {
           maxD = maxD[:len(maxD)-1]
       }
       maxD = append(maxD, i)
       for len(minD) > 0 && a[minD[len(minD)-1]] >= a[i] {
           minD = minD[:len(minD)-1]
       }
       minD = append(minD, i)
       for a[maxD[0]]-a[minD[0]] > s {
           if maxD[0] == left {
               maxD = maxD[1:]
           }
           if minD[0] == left {
               minD = minD[1:]
           }
           left++
       }
       pos[i] = left
   }
   const INF = 1000000000
   dp := make([]int, n+1)
   for i := 1; i <= n; i++ {
       dp[i] = INF
   }
   dq := make([]int, 0, n+1)
   for i := 1; i <= n; i++ {
       if i >= l {
           k := i - l
           for len(dq) > 0 && dp[dq[len(dq)-1]] >= dp[k] {
               dq = dq[:len(dq)-1]
           }
           dq = append(dq, k)
       }
       leftBound := pos[i] - 1
       for len(dq) > 0 && dq[0] < leftBound {
           dq = dq[1:]
       }
       if len(dq) > 0 {
           dp[i] = dp[dq[0]] + 1
       }
   }
   if dp[n] >= INF {
       fmt.Println(-1)
   } else {
       fmt.Println(dp[n])
   }
}
