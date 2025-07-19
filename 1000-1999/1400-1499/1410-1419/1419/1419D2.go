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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Ints(a)

   res := make([]int, n)
   // Place smallest elements at positions 1,3,5... (0-based)
   i, j := 0, 1
   for j < n {
       res[j] = a[i]
       i++
       j += 2
   }
   // Fill remaining positions
   j = 0
   for i < n {
       if res[j] != 0 {
           j++
       }
       res[j] = a[i]
       i++
       j++
   }
   // Count cheap spheres (local minima)
   cheap := 0
   for k := 1; k+1 < n; k++ {
       if res[k] < res[k-1] && res[k] < res[k+1] {
           cheap++
       }
   }
   // Output
   fmt.Fprintln(writer, cheap)
   for idx, v := range res {
       if idx > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
