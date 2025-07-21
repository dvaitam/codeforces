package main

import (
   "bufio"
   "fmt"
   "os"
)

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int) int {
   if b == 0 {
       return a
   }
   return gcd(b, a%b)
}

// abs returns the absolute value of x.
func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, x0, y0 int
   if _, err := fmt.Fscan(in, &n, &x0, &y0); err != nil {
       return
   }
   // use a map to track unique undirected directions
   dirs := make(map[[2]int]struct{})
   for i := 0; i < n; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       dx := x - x0
       dy := y - y0
       g := gcd(abs(dx), abs(dy))
       dx /= g
       dy /= g
       // normalize direction: ensure a unique representation for the line
       if dx < 0 || (dx == 0 && dy < 0) {
           dx = -dx
           dy = -dy
       }
       dirs[[2]int{dx, dy}] = struct{}{}
   }
   // the number of unique directions equals the number of shots
   fmt.Println(len(dirs))
