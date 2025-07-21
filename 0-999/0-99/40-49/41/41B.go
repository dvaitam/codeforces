package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, b int
   if _, err := fmt.Fscan(reader, &n, &b); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   ans := b
   for i := 0; i < n; i++ {
       priceBuy := a[i]
       if priceBuy <= 0 {
           continue
       }
       k := b / priceBuy
       if k <= 0 {
           continue
       }
       for j := i + 1; j < n; j++ {
           delta := a[j] - priceBuy
           if delta > 0 {
               val := b + k*delta
               if val > ans {
                   ans = val
               }
           }
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
}
