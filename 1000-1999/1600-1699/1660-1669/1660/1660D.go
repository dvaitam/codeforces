package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       p, l, r := 0, 0, 0
       j := -1
       for i := 0; i <= n; i++ {
           if i == n || a[i] == 0 {
               mn0x, mn0y := 0, j+1
               mn1x, mn1y := n, -1
               pw, sign := 0, 0
               for k := j+1; k < i; k++ {
                   if a[k] < 0 {
                       sign ^= 1
                   }
                   if abs(a[k]) == 2 {
                       pw++
                   }
                   // update best
                   if sign == 0 {
                       if pw-mn0x > p {
                           p = pw - mn0x
                           l = mn0y
                           r = k + 1
                       }
                   } else {
                       if pw-mn1x > p {
                           p = pw - mn1x
                           l = mn1y
                           r = k + 1
                       }
                   }
                   // update mn for sign
                   if sign == 0 {
                       if pw < mn0x {
                           mn0x = pw
                           mn0y = k + 1
                       }
                   } else {
                       if pw < mn1x {
                           mn1x = pw
                           mn1y = k + 1
                       }
                   }
               }
               j = i
           }
       }
       fmt.Fprintln(writer, l, n-r)
   }
}
