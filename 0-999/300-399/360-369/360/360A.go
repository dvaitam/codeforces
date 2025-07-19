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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   plu := make([]int, n)
   mx := make([]int, n)
   col := make([][]int, n)
   const INF = int(1e9)
   for i := 0; i < n; i++ {
       mx[i] = INF
   }
   s := 0
   for j := 0; j < m; j++ {
       var a, b, c, d int
       fmt.Fscan(reader, &a, &b, &c, &d)
       b--
       c--
       if a == 1 {
           for i := b; i <= c; i++ {
               plu[i] += d
           }
       } else {
           s++
           idx := s - 1
           for i := b; i <= c; i++ {
               val := d - plu[i]
               if mx[i] > val {
                   mx[i] = val
                   col[i] = col[i][:0]
               }
               if mx[i] == val {
                   col[i] = append(col[i], idx)
               }
           }
       }
   }
   used := make([]bool, s)
   cnt := 0
   for i := 0; i < n; i++ {
       for _, idx := range col[i] {
           if !used[idx] {
               used[idx] = true
               cnt++
           }
       }
   }
   if cnt == s {
       fmt.Fprintln(writer, "YES")
       for i := 0; i < n; i++ {
           fmt.Fprint(writer, mx[i])
           if i+1 < n {
               fmt.Fprint(writer, " ")
           }
       }
       fmt.Fprintln(writer)
   } else {
       fmt.Fprintln(writer, "NO")
   }
}
