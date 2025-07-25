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

   var n int
   fmt.Fscan(reader, &n)
   total := 2 * n
   a := make([]int, total)
   for i := 0; i < total; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Ints(a)
   // Case 1: first n to x, last n to y
   // area1 = (a[n-1]-a[0]) * (a[2n-1]-a[n])
   var area1 int64
   dx1 := int64(a[n-1] - a[0])
   dy1 := int64(a[total-1] - a[n])
   area1 = dx1 * dy1
   // Case 2: x is any window of size n, y gets extremes
   // area2 = min_{i=0..n} (a[i+n-1]-a[i]) * (a[total-1]-a[0])
   var minDx int64 = dx1
   for i := 1; i <= n; i++ {
       j := i + n - 1
       if j >= total {
           break
       }
       d := int64(a[j] - a[i])
       if d < minDx {
           minDx = d
       }
   }
   var area2 int64 = minDx * int64(a[total-1]-a[0])
   // result is minimal area
   res := area1
   if area2 < res {
       res = area2
   }
   fmt.Fprint(writer, res)
}
