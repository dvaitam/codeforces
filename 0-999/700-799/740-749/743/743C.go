package main

import (
   "fmt"
   "os"
)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   var n int64
   if _, err := fmt.Fscan(os.Stdin, &n); err != nil {
       return
   }
   if n == 1 {
       fmt.Println(-1)
       return
   }
   // even case
   if n%2 == 0 {
       k := n / 2
       x := k + 2
       y := k * (k + 1)
       z := (k + 1) * (k + 2)
       fmt.Println(x, y, z)
       return
   }
   // odd case: brute search for x,y,z
   half := n/2 + 1
   var found bool
   var rx, ry, rz int64
   maxX := half + 2000
   if maxX > 2*n {
       maxX = 2 * n
   }
   for x := half; x <= maxX; x++ {
       num := 2*x - n
       den := n * x
       g := gcd(num, den)
       a := num / g
       b := den / g
       // find y > x such that 1/y + 1/z = a/b
       // minimal y such that a*y > b
       y0 := b/a + 1
       if y0 <= x {
           y0 = x + 1
       }
       // try a few y candidates
       for y := y0; y < y0+10; y++ {
           D := a*y - b
           if D <= 0 {
               continue
           }
           by := b * y
           if by % D != 0 {
               continue
           }
           z := by / D
           if z <= 0 || z > 1000000000 {
               continue
           }
           if z == y || z == x {
               continue
           }
           // found
           rx, ry, rz = x, y, z
           found = true
           break
       }
       if found {
           break
       }
   }
   if !found {
       fmt.Println(-1)
   } else {
       fmt.Println(rx, ry, rz)
   }
}
