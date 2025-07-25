package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int
   fmt.Fscan(reader, &n, &m, &k)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Ints(a)
   if k > n {
       k = n
   }
   // take k cheapest
   b := a[:k]
   // best free count for each purchase size x
   bestY := make([]int, k+1)
   for j := 0; j < m; j++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       if x <= k && y > bestY[x] {
           bestY[x] = y
       }
   }
   // prefix sums of b, 0..k
   prefix := make([]int64, k+1)
   for i := 1; i <= k; i++ {
       prefix[i] = prefix[i-1] + int64(b[i-1])
   }
   const INF = int64(4e18)
   dp := make([]int64, k+1)
   for i := 1; i <= k; i++ {
       dp[i] = INF
   }
   dp[0] = 0
   for i := 1; i <= k; i++ {
       // try all purchase sizes x
       for x := 1; x <= i; x++ {
           y := bestY[x]
           // cost to pay for segment of last x items
           // sum of largest (x-y) = prefix[i] - prefix[i-x+y]
           cost := dp[i-x] + (prefix[i] - prefix[i-x+y])
           if cost < dp[i] {
               dp[i] = cost
           }
       }
   }
   fmt.Println(dp[k])
}
