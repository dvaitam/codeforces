package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var y, k, n int64
   if _, err := fmt.Fscan(reader, &y, &k, &n); err != nil {
       return
   }
   // Find the smallest multiple of k strictly greater than y
   start := ((y + k) / k) * k
   if start <= y {
       start += k
   }
   if start > n {
       fmt.Println(-1)
       return
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   first := true
   for m := start; m <= n; m += k {
       x := m - y
       if !first {
           writer.WriteByte(' ')
       }
       first = false
       writer.WriteString(fmt.Sprint(x))
   }
   writer.WriteByte('\n')
}
