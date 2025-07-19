package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // Compute total xor
   x := 0
   for _, v := range a {
       x ^= v
   }
   if x == 0 {
       fmt.Fprintln(writer, -1)
       return
   }
   ans := 0
   // Gaussian elimination over GF(2) to find basis size
   for i := 29; i >= 0; i-- {
       // find vector with bit i set
       id := -1
       mask := 1 << i
       for j, v := range a {
           if v&mask != 0 {
               id = j
               break
           }
       }
       if id == -1 {
           continue
       }
       ans++
       // move pivot to end
       last := len(a) - 1
       a[id], a[last] = a[last], a[id]
       // eliminate bit i from all other vectors
       pivot := a[last]
       for j := 0; j < last; j++ {
           if a[j]&mask != 0 {
               a[j] ^= pivot
           }
       }
       // remove pivot
       a = a[:last]
   }
   fmt.Fprintln(writer, ans)
}
