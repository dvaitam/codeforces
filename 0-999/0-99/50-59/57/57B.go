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

   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   a := make([]int, m)
   b := make([]int, m)
   c := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &a[i], &b[i], &c[i])
   }
   queries := make([]int, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(reader, &queries[i])
   }
   var ans int64
   for _, q := range queries {
       for i := 0; i < m; i++ {
           if a[i] <= q && q <= b[i] {
               ans += int64(c[i]) + int64(q - a[i])
           }
       }
   }
   fmt.Fprint(writer, ans)
}
