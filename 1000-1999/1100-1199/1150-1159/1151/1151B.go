package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   rdr := bufio.NewReader(os.Stdin)
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()

   var n, m int
   if _, err := fmt.Fscan(rdr, &n, &m); err != nil {
       return
   }
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(rdr, &a[i][j])
       }
   }
   // Try each bit
   for k := 0; k < 10; k++ {
       s0, s1, s01 := 0, 0, 0
       pos0 := make([]int, n)
       pos1 := make([]int, n)
       for i := 0; i < n; i++ {
           // find last positions
           p0, p1 := 0, 0
           for j := 0; j < m; j++ {
               if (a[i][j]>>k)&1 == 1 {
                   p1 = j + 1
               } else {
                   p0 = j + 1
               }
           }
           pos0[i], pos1[i] = p0, p1
           if p0 == 0 {
               s1++
           } else if p1 == 0 {
               s0++
           } else {
               s01++
           }
       }
       // check possibility
       if s1%2 == 1 || s01 > 0 {
           fmt.Fprintln(w, "TAK")
           ans := make([]int, n)
           if s1%2 == 1 {
               for i := 0; i < n; i++ {
                   if pos0[i] == 0 {
                       ans[i] = pos1[i]
                   } else {
                       ans[i] = pos0[i]
                   }
               }
           } else {
               // use one flexible to flip
               flip := true
               for i := 0; i < n; i++ {
                   if pos0[i] == 0 {
                       ans[i] = pos1[i]
                   } else if pos1[i] == 0 {
                       ans[i] = pos0[i]
                   } else {
                       if flip {
                           ans[i] = pos1[i]
                           flip = false
                       } else {
                           ans[i] = pos0[i]
                       }
                   }
               }
           }
           for i, v := range ans {
               if i > 0 {
                   w.WriteByte(' ')
               }
               fmt.Fprint(w, v)
           }
           w.WriteByte('\n')
           return
       }
   }
   fmt.Fprintln(w, "NIE")
}
