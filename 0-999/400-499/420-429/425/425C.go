package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m, s, e int
   fmt.Fscan(reader, &n, &m, &s, &e)
   a := make([]int, n)
   b := make([]int, m)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   maxVal := 0
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &b[i])
       if b[i] > maxVal {
           maxVal = b[i]
       }
   }
   // positions of values in b (1-based)
   posList := make([][]int, maxVal+1)
   for i, v := range b {
       posList[v] = append(posList[v], i+1)
   }
   maxK := s / e
   if maxK > n {
       maxK = n
   }
   if maxK > m {
       maxK = m
   }
   // dp[t] = minimal b-position for matching t elements
   INFJ := m + 1
   dp := make([]int, maxK+2)
   for i := range dp {
       dp[i] = INFJ
   }
   dp[0] = 0
   // minSum[t] = minimal i+j sum when dp[t] updated
   INF_SUM := n + m + 5
   minSum := make([]int, maxK+2)
   for i := range minSum {
       minSum[i] = INF_SUM
   }
   curMax := 0
   // iterate a
   for i := 1; i <= n; i++ {
       v := a[i-1]
       if v <= maxVal {
           lst := posList[v]
           if len(lst) == 0 {
               continue
           }
           // try to extend matches
           ub := curMax
           if ub > maxK-1 {
               ub = maxK - 1
           }
           for t := ub; t >= 0; t-- {
               prevJ := dp[t]
               if prevJ >= INFJ {
                   continue
               }
               // find pos > prevJ
               idx := sort.Search(len(lst), func(k int) bool { return lst[k] > prevJ })
               if idx < len(lst) {
                   j := lst[idx]
                   if j < dp[t+1] {
                       dp[t+1] = j
                       sum := i + j
                       if sum < minSum[t+1] {
                           minSum[t+1] = sum
                       }
                       if t+1 > curMax {
                           curMax = t + 1
                       }
                   }
               }
           }
       }
   }
   // compute answer
   ans := 0
   for t := 1; t <= curMax; t++ {
       if t*e > s {
           break
       }
       if minSum[t] <= s - t*e {
           ans = t
       }
   }
   fmt.Fprintln(writer, ans)
}
