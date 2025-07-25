package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var T int
   if _, err := fmt.Fscan(in, &T); err != nil {
       return
   }
   for t := 0; t < T; t++ {
       var n, k int
       fmt.Fscan(in, &n, &k)
       a := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &a[i])
       }
       bestSpan := int64(1<<62 - 1)
       bestI := 0
       for i := 0; i + k < n; i++ {
           span := a[i+k] - a[i]
           if span < bestSpan {
               bestSpan = span
               bestI = i
           }
       }
       x := (a[bestI] + a[bestI+k]) / 2
       fmt.Fprintln(out, x)
   }
}
