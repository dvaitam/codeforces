package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for t > 0 {
       t--
       var n, k int
       fmt.Fscan(in, &n, &k)
       a := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &a[i])
       }
       sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
       var sum int64
       for i := 0; i <= k; i++ {
           sum += a[i]
       }
       fmt.Fprintln(out, sum)
   }
}
