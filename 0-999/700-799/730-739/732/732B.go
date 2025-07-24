package main

import (
   "bufio"
   "fmt"
   "os"
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
   b := make([]int, n)
   totalExtra := 0
   if n > 0 {
       // first day has no lower bound from previous day (pre-day has k walks)
       b[0] = a[0]
       // subsequent days
       for i := 1; i < n; i++ {
           need := k - b[i-1]
           if need < 0 {
               need = 0
           }
           if a[i] >= need {
               b[i] = a[i]
           } else {
               b[i] = need
           }
           totalExtra += b[i] - a[i]
       }
   }
   // no extra for first day since b[0] == a[0]
   fmt.Println(totalExtra)
   for i := 0; i < n; i++ {
       if i > 0 {
           fmt.Print(" ")
       }
       fmt.Print(b[i])
   }
   fmt.Println()
}
