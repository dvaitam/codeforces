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
   // a[x][y][z]
   a := make([][][]int, n)
   for x := 0; x < n; x++ {
      a[x] = make([][]int, n)
      for y := 0; y < n; y++ {
         a[x][y] = make([]int, n)
      }
   }
   cnt := 1
   // Traverse columns (x,y) in 2D snake order, and for each column traverse z-axis alternating direction
   idx := 0
   for x := 0; x < n; x++ {
      if x%2 == 0 {
         for y := 0; y < n; y++ {
            if idx%2 == 0 {
               for z := 0; z < n; z++ {
                  a[x][y][z] = cnt
                  cnt++
               }
            } else {
               for z := n - 1; z >= 0; z-- {
                  a[x][y][z] = cnt
                  cnt++
               }
            }
            idx++
         }
      } else {
         for y := n - 1; y >= 0; y-- {
            if idx%2 == 0 {
               for z := 0; z < n; z++ {
                  a[x][y][z] = cnt
                  cnt++
               }
            } else {
               for z := n - 1; z >= 0; z-- {
                  a[x][y][z] = cnt
                  cnt++
               }
            }
            idx++
         }
      }
   }
   // Output layers z=0..n-1
   for z := 0; z < n; z++ {
      for x := 0; x < n; x++ {
         for y := 0; y < n; y++ {
            if y > 0 {
               out.WriteByte(' ')
            }
            fmt.Fprint(out, a[x][y][z])
         }
         out.WriteByte('\n')
      }
      if z != n-1 {
         out.WriteByte('\n')
      }
   }
}
