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

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   sum := make([]int, n)
   if n > 0 {
       sum[0] = a[0]
       for i := 1; i < n; i++ {
           sum[i] = sum[i-1] + a[i]
       }
   }

   var avg float64
   for L := k; L <= n; L++ {
       mx := sum[L-1]
       for j := L; j < n; j++ {
           s := sum[j] - sum[j-L]
           if s > mx {
               mx = s
           }
       }
       cur := float64(mx) / float64(L)
       if cur > avg {
           avg = cur
       }
   }

   fmt.Fprintf(writer, "%.9f\n", avg)
}
