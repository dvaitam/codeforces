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
   fmt.Fscan(in, &n)
   // prepare grid with bits on odd-sum cells
   a := make([][]int64, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int64, n)
       for j := 0; j < n; j++ {
           if (i+j)%2 == 1 {
               a[i][j] = 1 << uint(i)
           }
       }
   }
   // print grid
   fmt.Fprintln(out, n)
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if j > 0 {
               out.WriteByte(' ')
           }
           fmt.Fprint(out, a[i][j])
       }
       fmt.Fprintln(out)
   }
   out.Flush()

   var q int
   fmt.Fscan(in, &q)
   for qi := 0; qi < q; qi++ {
       var k int64
       fmt.Fscan(in, &k)
       x, y := 0, 0
       // output path
       fmt.Fprintln(out, 1, 1)
       for step := 0; step < 2*n-2; step++ {
           if (x+y)%2 == 0 {
               // at even-sum cell, children have weights
               down := int64(0)
               if x+1 < n {
                   down = a[x+1][y]
               }
               if down > 0 && (k&down) != 0 {
                   x++
               } else {
                   y++
               }
           } else {
               // at odd-sum cell, move to the only remaining to match parity
               if x+1 < n {
                   x++
               } else {
                   y++
               }
           }
           fmt.Fprintln(out, x+1, y+1)
       }
       out.Flush()
   }
