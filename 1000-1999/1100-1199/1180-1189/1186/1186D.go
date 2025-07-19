package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   b := make([]int64, n)
   var sum int64
   for i := 0; i < n; i++ {
       var ai int
       var bi int64
       fmt.Fscan(in, &ai, &bi)
       if ai < 0 {
           bi = -bi
       }
       a[i] = ai
       b[i] = bi
       sum += bi
   }
   const M int64 = 100000
   res := sum / M
   // sum must be divisible by M
   for i := 0; i < n; i++ {
       if res != 0 && b[i] != 0 {
           if res > 0 {
               if b[i] > 0 {
                   fmt.Fprintln(out, a[i]+1)
                   res--
               } else {
                   fmt.Fprintln(out, a[i])
               }
           } else {
               if b[i] < 0 {
                   fmt.Fprintln(out, a[i]-1)
                   res++
               } else {
                   fmt.Fprintln(out, a[i])
               }
           }
       } else {
           fmt.Fprintln(out, a[i])
       }
   }
}
