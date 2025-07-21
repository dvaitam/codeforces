package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       if err != io.EOF {
           fmt.Fprintln(os.Stderr, "failed to read n and k:", err)
       }
       return
   }
   prices := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &prices[i])
   }
   // if k >= n/2, can make unlimited transactions
   if k >= n/2 {
       var profit int64
       for i := 1; i < n; i++ {
           if diff := prices[i] - prices[i-1]; diff > 0 {
               profit += diff
           }
       }
       fmt.Println(profit)
       return
   }
   // dp for at most k transactions
   // prev and curr dp arrays: profit until day i
   prev := make([]int64, n)
   curr := make([]int64, n)
   for t := 1; t <= k; t++ {
       // maxDiff = max(prev[j] - prices[j]) for j before i
       maxDiff := prev[0] - prices[0]
       curr[0] = 0
       for i := 1; i < n; i++ {
           // either no transaction at day i, or sell at i
           // sell profit = prices[i] + maxDiff
           sell := prices[i] + maxDiff
           if sell > curr[i-1] {
               curr[i] = sell
           } else {
               curr[i] = curr[i-1]
           }
           // update maxDiff for next i
           if candidate := prev[i] - prices[i]; candidate > maxDiff {
               maxDiff = candidate
           }
       }
       // swap prev and curr
       prev, curr = curr, prev
   }
   // result in prev[n-1]
   fmt.Println(prev[n-1])
}
