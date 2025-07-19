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

   var n, k, p, x, y int
   fmt.Fscan(reader, &n, &k, &p, &x, &y)
   a := make([]int, k)
   tot := 0
   for i := 0; i < k; i++ {
       fmt.Fscan(reader, &a[i])
       tot += a[i]
   }
   nNew := n - k
   // Minimum sum with all new scores = tot + nNew*1
   if tot+nNew > x {
       fmt.Fprintln(writer, -1)
       return
   }
   // Distribute extra points beyond the baseline of 1
   cur := x - (tot + nNew)
   b := make([]int, 0, nNew)
   for i := 0; i < nNew; i++ {
       if cur >= y-1 {
           b = append(b, y)
           cur -= y - 1
       } else {
           b = append(b, 1)
       }
   }
   // Combine and sort to check median and max
   all := make([]int, 0, n)
   all = append(all, a...)
   all = append(all, b...)
   sort.Ints(all)
   // Median at zero-based index n/2 corresponds to one-based n/2+1
   if all[n/2] < y || all[n-1] > p {
       fmt.Fprintln(writer, -1)
       return
   }
   // Output the added scores
   for i, v := range b {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
