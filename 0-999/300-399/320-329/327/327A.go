package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   arr := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &arr[i])
   }
   // count original ones and transform array for max subarray
   ones := 0
   b := make([]int, n)
   for i, x := range arr {
       if x == 1 {
           ones++
           b[i] = -1
       } else {
           b[i] = 1
       }
   }
   // find maximum subarray sum (at least one element)
   best := b[0]
   curr := b[0]
   for i := 1; i < n; i++ {
       if curr < 0 {
           curr = b[i]
       } else {
           curr += b[i]
       }
       if curr > best {
           best = curr
       }
   }
   result := ones + best
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, result)
}
