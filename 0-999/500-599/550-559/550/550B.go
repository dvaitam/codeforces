package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var l, r, x int
   if _, err := fmt.Fscan(reader, &n, &l, &r, &x); err != nil {
       return
   }
   c := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &c[i])
   }
   count := 0
   // iterate over all subsets
   for mask := 0; mask < (1 << n); mask++ {
       cnt, sum := 0, 0
       minv, maxv := int(1e9+5), 0
       for i := 0; i < n; i++ {
           if mask&(1<<i) != 0 {
               cnt++
               v := c[i]
               sum += v
               if v < minv {
                   minv = v
               }
               if v > maxv {
                   maxv = v
               }
           }
       }
       if cnt < 2 {
           continue
       }
       if sum >= l && sum <= r && maxv-minv >= x {
           count++
       }
   }
   fmt.Println(count)
}
