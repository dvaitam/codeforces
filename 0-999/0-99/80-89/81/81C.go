package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, a, b int
   if _, err := fmt.Fscan(reader, &n, &a, &b); err != nil {
       return
   }
   t := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &t[i])
   }
   f := make([]int, n)
   // Equal group sizes: lexicographically minimal is first a as 1
   if a == b {
       for i := 0; i < n; i++ {
           if i < a {
               f[i] = 1
           } else {
               f[i] = 2
           }
       }
   } else {
       // Determine whether to pick largest or smallest for subject 1
       pickLargest := b > a
       u := make([]int, n)
       copy(u, t)
       sort.Ints(u)
       if pickLargest {
           // pick a largest values
           threshold := u[n-a]
           cntGreater := 0
           cntEqTotal := 0
           for _, v := range t {
               if v > threshold {
                   cntGreater++
               } else if v == threshold {
                   cntEqTotal++
               }
           }
           eqSelect := a - cntGreater
           for i, v := range t {
               if v > threshold {
                   f[i] = 1
               } else if v == threshold && eqSelect > 0 {
                   f[i] = 1
                   eqSelect--
               } else {
                   f[i] = 2
               }
           }
       } else {
           // pick a smallest values
           threshold := u[a-1]
           cntLess := 0
           cntEqTotal := 0
           for _, v := range t {
               if v < threshold {
                   cntLess++
               } else if v == threshold {
                   cntEqTotal++
               }
           }
           eqSelect := a - cntLess
           for i, v := range t {
               if v < threshold {
                   f[i] = 1
               } else if v == threshold && eqSelect > 0 {
                   f[i] = 1
                   eqSelect--
               } else {
                   f[i] = 2
               }
           }
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < n; i++ {
       if i > 0 {
           writer.WriteString(" ")
       }
       fmt.Fprint(writer, f[i])
   }
   writer.WriteByte('\n')
}
