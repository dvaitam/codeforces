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
   var T int
   fmt.Fscan(in, &T)
   for T > 0 {
       T--
       var n int
       fmt.Fscan(in, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &a[i])
       }
       // check all equal and any adjacent equal (including wrap)
       tag := true
       flag := false
       for i := 1; i < n; i++ {
           if a[i] != a[i-1] {
               tag = false
           } else {
               flag = true
           }
       }
       if a[0] == a[n-1] {
           flag = true
       }
       // result colors
       var k int
       colors := make([]int, n)
       if tag {
           // all same
           k = 1
           for i := 0; i < n; i++ {
               colors[i] = 1
           }
       } else if flag {
           // can do with 2 colors
           k = 2
           if n%2 == 0 {
               // even n, simple alt
               for i := 0; i < n; i++ {
                   if i%2 == 0 {
                       colors[i] = 2
                   } else {
                       colors[i] = 1
                   }
               }
           } else {
               // odd n, find first equal adjacency
               op := 1
               vis := false
               colors[0] = op + 1
               for i := 1; i < n; i++ {
                   if !vis && a[i] == a[i-1] {
                       vis = true
                   } else {
                       op ^= 1
                   }
                   colors[i] = op + 1
               }
           }
       } else {
           // no adjacent equal, need 2 or 3
           if n%2 == 0 {
               k = 2
               // even n, simple alt
               for i := 0; i < n; i++ {
                   if i%2 == 0 {
                       colors[i] = 2
                   } else {
                       colors[i] = 1
                   }
               }
           } else {
               k = 3
               // odd n, use 3 colors
               op := 1
               colors[0] = 2
               for i := 1; i < n-1; i++ {
                   op ^= 1
                   colors[i] = op + 1
               }
               colors[n-1] = 3
           }
       }
       // output
       fmt.Fprintln(out, k)
       for i, c := range colors {
           if i > 0 {
               out.WriteByte(' ')
           }
           fmt.Fprint(out, c)
       }
       fmt.Fprintln(out)
   }
}
