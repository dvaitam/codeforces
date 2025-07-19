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
   for {
       var n int
       if _, err := fmt.Fscan(reader, &n); err != nil {
           break
       }
       a := make([]int, n)
       b := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &b[i])
       }
       x := make([]int, 0, 2)
       y := make([]int, 0, 2)
       num := make([]bool, n+1)
       for i := 0; i < n; i++ {
           if a[i] != b[i] {
               x = append(x, a[i])
               y = append(y, b[i])
               a[i] = -1
           } else {
               num[a[i]] = true
           }
       }
       f := make([]int, 0, 2)
       for v := 1; v <= n; v++ {
           if !num[v] {
               f = append(f, v)
           }
       }
       if len(f) == 1 {
           // only one missing
           miss := f[0]
           for i := 0; i < n; i++ {
               if a[i] != -1 {
                   fmt.Fprint(writer, a[i])
               } else {
                   fmt.Fprint(writer, miss)
               }
               if i < n-1 {
                   writer.WriteByte(' ')
               }
           }
           writer.WriteByte('\n')
       } else {
           // two missing, decide order
           var s int
           // x[0],y[0] and x[1],y[1]
           if (f[0] == x[0] && f[0] != y[0] && f[1] == y[1] && f[1] != x[1]) ||
               (f[0] != x[0] && f[0] == y[0] && f[1] != y[1] && f[1] == x[1]) {
               s = 0
           } else {
               s = 1
           }
           for i := 0; i < n; i++ {
               if a[i] != -1 {
                   fmt.Fprint(writer, a[i])
               } else {
                   fmt.Fprint(writer, f[s])
                   // toggle s between 0 and 1
                   if s == 0 {
                       s = 1
                   } else {
                       s = 0
                   }
               }
               if i < n-1 {
                   writer.WriteByte(' ')
               }
           }
           writer.WriteByte('\n')
       }
   }
}
