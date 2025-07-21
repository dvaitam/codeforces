package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   pies := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &pies[i])
   }
   sort.Ints(pies)

   // Compute maximum number of free pies (k)
   k := 0
   i, j := 0, 0
   for j < n && i < n {
       if pies[j] > pies[i] {
           k++
           i++
           j++
       } else {
           j++
       }
   }
   // Ensure free count does not exceed paid count
   if k > n-k {
       k = n - k
   }

   var total, freeSum int64
   for _, v := range pies {
       total += int64(v)
   }
   for idx := 0; idx < k; idx++ {
       freeSum += int64(pies[idx])
   }

   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, total-freeSum)
}
