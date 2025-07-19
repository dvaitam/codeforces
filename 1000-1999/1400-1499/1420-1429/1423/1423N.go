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
   fmt.Fscan(reader, &n, &m)
   in := make([][]int, n+1)
   vw := make([][]int, n+1)
   su := make([]int, n+1)
   zz := make([]int, n+1)
   A := make([]int, m)
   B := make([]int, m)
   jg := make([]int, m)
   bk := make([]bool, m+n+5)

   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       A[i], B[i] = a, b
       if a > b {
           a, b = b, a
       }
       in[b] = append(in[b], a)
       vw[b] = append(vw[b], i)
       su[a]++
       su[b]++
       jg[i] = 1
   }

   for i := 1; i <= n; i++ {
       s0, s1 := 0, 0
       for _, x := range in[i] {
           bk[su[x]] = true
           if zz[x] != 0 {
               s1++
           } else {
               s0++
           }
       }
       w := su[i] - s0
       for bk[w] {
           w++
       }
       c := w - su[i]
       su[i] = w
       for idx, x := range in[i] {
           y := vw[i][idx]
           bk[su[x]] = false
           if c < 0 {
               if zz[x] == 0 {
                   jg[y]--
                   c++
                   zz[x] = 1
               }
           } else if c > 0 {
               if zz[x] == 1 {
                   jg[y]++
                   c--
                   zz[x] = 0
               }
           }
       }
   }

   cnt := 0
   for i := 1; i <= n; i++ {
       cnt += zz[i]
   }
   fmt.Fprintln(writer, cnt)
   for i := 1; i <= n; i++ {
       if zz[i] != 0 {
           fmt.Fprint(writer, i, " ")
       }
   }
   fmt.Fprintln(writer)
   for i := 0; i < m; i++ {
       fmt.Fprintln(writer, A[i], B[i], jg[i])
   }
}
