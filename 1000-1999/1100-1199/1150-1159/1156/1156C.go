package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, z int
   if _, err := fmt.Fscan(reader, &n, &z); err != nil {
       return
   }
   x := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x[i])
   }
   sort.Ints(x)

   lo, hi := 0, n/2
   for lo < hi {
       mid := (lo + hi + 1) / 2
       ok := true
       for i := 0; i < mid; i++ {
           if x[n-mid+i] - x[i] < z {
               ok = false
               break
           }
       }
       if ok {
           lo = mid
       } else {
           hi = mid - 1
       }
   }
   fmt.Fprintln(writer, lo)
}
