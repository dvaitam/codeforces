package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var h int64
   if _, err := fmt.Fscan(reader, &n, &h); err != nil {
       return
   }
   x := make([]int64, n)
   y := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x[i], &y[i])
   }

   // two pointers with budget to cover gaps
   i, j := 0, 0
   sum := y[0] - x[0]
   budget := h
   // initial expand
   for j+1 < n {
       gap := x[j+1] - y[j]
       if budget-gap > 0 {
           budget -= gap
       } else {
           break
       }
       j++
       sum += gap + (y[j] - x[j])
   }
   ans := sum + budget
   if j == n-1 {
       fmt.Fprintln(writer, ans)
       return
   }
   // slide window
   for i < n {
       // move left pointer
       i++
       // adjust covered span
       sum -= x[i] - x[i-1]
       // free budget from the gap between segments
       budget += x[i] - y[i-1]
       // expand right pointer as much as possible
       for j+1 < n {
           gap := x[j+1] - y[j]
           if budget-gap > 0 {
               budget -= gap
           } else {
               break
           }
           j++
           sum += gap + (y[j] - x[j])
       }
       // update answer
       ans = max(ans, sum+budget)
       if j == n-1 {
           break
       }
   }
   fmt.Fprintln(writer, ans)
}
