package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // sort descending
   sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
   // count negative values
   cntNeg := 0
   for _, v := range a {
       if v < 0 {
           cntNeg++
       }
   }
   // we can isolate at most k negatives
   p := k
   if p > cntNeg {
       p = cntNeg
   }
   mainSize := n - p
   var ans int64
   // in main segment [0..mainSize), compute sum of prefix sums
   for i := 0; i < mainSize; i++ {
       mult := mainSize - i - 1
       ans += int64(a[i]) * int64(mult)
   }
   fmt.Println(ans)
}
