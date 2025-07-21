package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, k int
   fmt.Fscan(in, &n, &k)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }

   // Flip negatives greedily
   for i := 0; i < n && k > 0; i++ {
       if a[i] < 0 {
           a[i] = -a[i]
           k--
       }
   }

   sum := int64(0)
   minAbs := int(1e18)
   zeroExists := false
   for _, v := range a {
       if v == 0 {
           zeroExists = true
       }
       if abs(v) < minAbs {
           minAbs = abs(v)
       }
       sum += int64(v)
   }

   // If remaining flips is odd and no zero, subtract twice the smallest abs value
   if k > 0 && k%2 == 1 && !zeroExists {
       sum -= int64(2 * minAbs)
   }

   fmt.Fprintln(out, sum)
}
