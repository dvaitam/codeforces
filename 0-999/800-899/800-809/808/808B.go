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
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   var s int64
   for i := 0; i < k; i++ {
       s += a[i]
   }
   var res int64 = s
   for i := k; i < n; i++ {
       s -= a[i-k]
       s += a[i]
       res += s
   }

   windows := int64(n - k + 1)
   avg := float64(res) / float64(windows)
   fmt.Fprintf(writer, "%.6f", avg)
}
